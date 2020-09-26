package main

import (
        "fmt"
	"log"
	"net"
	"net/http"
	"crypto/tls"
        "io"
        "sync"
        "errors"
        "strconv"
        "encoding/json"
)


type BoringProxy struct {
        tunMan *TunnelManager
        adminListener *AdminListener
}

func NewBoringProxy() *BoringProxy {

        tunMan := NewTunnelManager()
        adminListener := NewAdminListener()

        p := &BoringProxy{tunMan, adminListener}

	http.HandleFunc("/", p.handleAdminRequest)
        go http.Serve(adminListener, nil)

        return p
}

func (p *BoringProxy) Run() {

        listener, err := net.Listen("tcp", ":443")
        if err != nil {
                log.Fatal(err)
        }

        for {
                conn, err := listener.Accept()
                if err != nil {
                        log.Print(err)
                        continue
                }
                go p.handleConnection(conn)
        }
}

func (p *BoringProxy) handleAdminRequest(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/tunnels" {
                p.handleTunnels(w, r)
        }
}

func (p *BoringProxy) handleTunnels(w http.ResponseWriter, r *http.Request) {
        fmt.Println("handleTunnels")

        if r.Method == "GET" {
                body, err := json.Marshal(p.tunMan.tunnels)
                if err != nil {
                w.WriteHeader(500)
                w.Write([]byte("Error encoding tunnels"))
                return
                }
                w.Write([]byte(body))
        } else if r.Method == "POST" {
                p.handleCreateTunnel(w, r)
        }
}

func (p *BoringProxy) handleCreateTunnel(w http.ResponseWriter, r *http.Request) {
        fmt.Println("handleCreateTunnel")

        query := r.URL.Query()

        if len(query["host"]) != 1 {
                w.WriteHeader(400)
                w.Write([]byte("Invalid host parameter"))
                return
        }
        host := query["host"][0]

        if len(query["port"]) != 1 {
                w.WriteHeader(400)
                w.Write([]byte("Invalid port parameter"))
                return
        }

        port, err := strconv.Atoi(query["port"][0])
        if err != nil {
                w.WriteHeader(400)
                w.Write([]byte("Invalid port parameter"))
                return
        }

        p.tunMan.SetTunnel(host, port)
}

func (p *BoringProxy) handleConnection(clientConn net.Conn) {
        // TODO: does this need to be closed manually, or is it handled when decryptedConn is closed?
        //defer clientConn.Close()

        certBaseDir := "/home/anders/.local/share/caddy/certificates/acme-v02.api.letsencrypt.org-directory/"

        var serverName string

        decryptedConn := tls.Server(clientConn, &tls.Config{
                GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {

                        serverName = clientHello.ServerName

                        certPath := certBaseDir + clientHello.ServerName + "/" + clientHello.ServerName + ".crt"
                        keyPath := certBaseDir + clientHello.ServerName + "/" + clientHello.ServerName + ".key"

                        cert, err := tls.LoadX509KeyPair(certPath, keyPath)
                        if err != nil {
                                log.Println("getting cert failed")
                                return nil, err
                        }

                        return &cert, nil
                },
        })
        //defer decryptedConn.Close()

        // Need to manually do handshake to ensure serverName is populated by this point. Usually Handshake()
        // is automatically called on first read/write
        decryptedConn.Handshake()

        adminDomain := "anders.webstreams.io"
        if serverName == adminDomain {
                p.handleAdminConnection(decryptedConn)
        } else {
                p.handleTunnelConnection(decryptedConn, serverName)
        }
}

func (p *BoringProxy) handleAdminConnection(decryptedConn net.Conn) {
        p.adminListener.connChan <- decryptedConn
}

func (p *BoringProxy) handleTunnelConnection(decryptedConn net.Conn, serverName string) {

        defer decryptedConn.Close()

        port, err := p.tunMan.GetPort(serverName)
        if err != nil {
                log.Print(err)
                errMessage := fmt.Sprintf("HTTP/1.1 500 Internal server error\n\nNo tunnel attached to %s", serverName)
                decryptedConn.Write([]byte(errMessage))
                return
        }

        upstreamAddr := fmt.Sprintf("127.0.0.1:%d", port)

        upstreamConn, err := net.Dial("tcp", upstreamAddr)
        if err != nil {
                log.Print(err)
                return
        }
        defer upstreamConn.Close()

        var wg sync.WaitGroup
        wg.Add(2)

        go func() {
                io.Copy(decryptedConn, upstreamConn)
                //decryptedConn.(*net.TCPConn).CloseWrite()
                wg.Done()
        }()
        go func() {
                io.Copy(upstreamConn, decryptedConn)
                //upstreamConn.(*net.TCPConn).CloseWrite()
                wg.Done()
        }()

        wg.Wait()
}


type TunnelManager struct {
        tunnels map[string]int
        mutex *sync.Mutex
}

func NewTunnelManager() *TunnelManager {
        tunnels := make(map[string]int)
        mutex := &sync.Mutex{}
        return &TunnelManager{tunnels, mutex}
}

func (m *TunnelManager) SetTunnel(host string, port int) {
        m.mutex.Lock()
        m.tunnels[host] = port
        m.mutex.Unlock()
}

func (m *TunnelManager) GetPort(serverName string) (int, error) {
        m.mutex.Lock()
        port, exists := m.tunnels[serverName]
        m.mutex.Unlock()

        if !exists {
                return 0, errors.New("Doesn't exist")
        }

        return port, nil
}


type AdminListener struct {
        connChan chan(net.Conn)
}

func NewAdminListener() *AdminListener {
        connChan := make(chan(net.Conn))
        return &AdminListener{connChan}
}

// implement net.Listener
func (l *AdminListener) Accept() (net.Conn, error) {
        // TODO: error conditions?
        conn := <-l.connChan
        return conn, nil
}
func (l *AdminListener) Close() error {
        // TODO
        fmt.Println("AdminListener Close")
        return nil
}
func (l *AdminListener) Addr() net.Addr {
        // TODO
        fmt.Println("AdminListener Addr")
        return nil
}

func main() {
        log.Println("Starting up")

        proxy := NewBoringProxy()
        proxy.Run()
}
