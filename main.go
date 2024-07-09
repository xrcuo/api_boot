package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/xrcuo/api_boot/motd"
)

func main() {
	Host := "103.40.13.47:20002"
	data, err := motd.MotdBE(Host)
	if err != nil {
		fmt.Println(err)
	}
	logrus.Info(data)
}
