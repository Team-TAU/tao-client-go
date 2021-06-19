package main

import (
	"fmt"
	tau "github.com/Team-TAU/tao-client-go"
	"os"
	"strconv"
	"time"
)

// This application gets all messages from TAU and prints them to stdout
func main() {
	host := os.Getenv("TAU_HOST")
	port := os.Getenv("TAU_PORT")
	portInt, _ := strconv.Atoi(port)
	token := os.Getenv("TAU_TOKEN")
	ssl := os.Getenv("TAU_SSL")
	hasSSL, _ := strconv.ParseBool(ssl)

	client, err := tau.NewClient(host, portInt, token, hasSSL)
	if err != nil {
		panic(err)
	}

	client.SetRawCallback(rawCallback)

	for {
		time.Sleep(time.Minute)
	}
}

func rawCallback(msg []byte) {
	println(fmt.Sprintf("%s", msg))
}
