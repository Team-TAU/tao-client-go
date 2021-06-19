package main

import (
	"github.com/Team-TAU/tau-client-go/helix"
	"os"
	"strconv"
)

func main() {
	host := os.Getenv("TAU_HOST")
	port := os.Getenv("TAU_PORT")
	portInt, _ := strconv.Atoi(port)
	token := os.Getenv("TAU_TOKEN")
	ssl := os.Getenv("TAU_SSL")
	hasSSL, _ := strconv.ParseBool(ssl)

	client, err := helix.NewClient(host, portInt, token, hasSSL)
	if err != nil {
		panic(err)
	}
	users, err := client.GetTwitchUsers([]string{"wwsean08", "finitesingularity"}, nil)
	if err != nil {
		println(err.Error())
	}
	for _, user := range users.Data {
		println(user.DisplayName)
	}
}
