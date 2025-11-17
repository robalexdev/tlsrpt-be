package main

import (
	"fmt"
	"net/http"
	"os"
)

func healthcheck() {
	resp, err := http.Get("http://127.0.0.1:8080/health")
	if err != nil {
		fmt.Printf("Bad GET: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	// Just look for 200 OK status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad status code: %d\n", resp.StatusCode)
		os.Exit(1)
	}
	fmt.Println("UP")
}

func main() {
  healthcheck()
}
