package main

import (
	tau "gitlab.com/wwsean08/go-tau"
	"os"
	"strconv"
	"time"
)

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

	client.SetFollowCallback(onFollow)
	client.SetSubscriptionCallback(onSub)
	client.SetCheerCallback(onCheer)

	for {
		time.Sleep(time.Minute)
	}
}

func onFollow(msg *tau.FollowMsg) {
	f, err := os.OpenFile("newestFollower.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		println(err.Error())
	}
	defer f.Close()
	f.Truncate(0)
	f.WriteString(msg.EventData.UserName)
}

func onSub(msg *tau.SubscriptionMsg) {
	f, err := os.OpenFile("newestSub.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		println(err.Error())
	}
	defer f.Close()
	f.Truncate(0)
	f.WriteString(msg.EventData.Data.Message.DisplayName)
}

func onCheer(msg *tau.CheerMsg) {
	f, err := os.OpenFile("newestCheer.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		println(err.Error())
	}
	defer f.Close()
	f.Truncate(0)
	f.WriteString(msg.EventData.UserName)
}
