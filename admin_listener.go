package main

import(
        "fmt"
	"net"
)


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
