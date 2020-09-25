package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"net/http"
	"bytes"
)


func main() {

	args := os.Args[1:]
	tunnelSpecParts := strings.Split(args[0], ":")
	tunnelId := args[0]
	host := tunnelSpecParts[0]
	port := tunnelSpecParts[1]

	client := &http.Client{}

	caddyAddRouteStr := fmt.Sprintf("{\"@id\":\"%s\",\"match\":[{\"host\":[\"%s\"]}],\"handle\":[{\"handler\":\"reverse_proxy\",\"upstreams\":[{\"dial\":\":%s\"}]}]}", tunnelId, host, port);

	resp, err := http.Post("http://127.0.0.1:2019/config/apps/http/servers/sirtunnel/routes", "application/json", bytes.NewBuffer([]byte(caddyAddRouteStr)))

	if err != nil {
		fmt.Println("Tunnel creation failed")
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Tunnel created successfully")

	// wait for CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Cleaning up tunnel")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://127.0.0.1:2019/id/%s", tunnelId), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	_, err = client.Do(req)

	fmt.Println("Exiting")
}
