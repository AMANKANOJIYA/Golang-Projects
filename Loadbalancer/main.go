package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	serverList = []*server{
		newServer("Server-1", "http://localhost:5000"),
		newServer("Server-2", "http://localhost:5001"),
		newServer("Server-3", "http://localhost:5002"),
		newServer("Server-4", "http://localhost:5003"),
		newServer("Server-5", "http://localhost:5004"),
	}
	lastSerevrIndex = 0
)

func main() {
	http.HandleFunc("/", forwardRequest)
	go startHealthCheck()
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func forwardRequest(res http.ResponseWriter, req *http.Request) {
	server, err := getHealthyServer()
	if err != nil {
		http.Error(res, "Couldn't process Request: "+err.Error(), http.StatusServiceUnavailable)
		// fmt.Fprintf(res, "Couldn't process Request: %s", err.Error())
	}
	server.ReverseProxy.ServeHTTP(res, req)
}

func getHealthyServer() (*server, error) {
	for i := 0; i < len(serverList); i++ {
		server := getServer()
		if server.Health {
			return server, nil
		}
	}
	return nil, fmt.Errorf("No healthy Hosts")
}

func getServer() *server {
	nextIndex := (lastSerevrIndex + 1) % len(serverList)
	server := serverList[nextIndex]
	lastSerevrIndex = nextIndex
	return server
}
