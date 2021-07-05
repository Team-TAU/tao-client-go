package main

import (
	"fmt"
	tau "github.com/Team-TAU/tau-client-go"
	"os"
	"strconv"
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

	streamer, err := client.FollowStreamerOnTau("fiercekittenz")
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("%s - %s", streamer.TwitchUsername, streamer.TwitchID))
}
