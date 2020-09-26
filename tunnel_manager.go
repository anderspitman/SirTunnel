package main

import (
        "errors"
        "sync"
)


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

func (m *TunnelManager) DeleteTunnel(host string) {
        m.mutex.Lock()
        delete(m.tunnels, host)
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
