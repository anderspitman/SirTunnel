package main

import (
        "log"
)


func main() {
        log.Println("Starting up")

        proxy := NewBoringProxy()
        proxy.Run()
}
