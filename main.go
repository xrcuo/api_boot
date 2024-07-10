package main

import (
	"fmt"

	"github.com/xrcuo/api_boot/motd"
)

func main() {
	Host := ""
	data, err := motd.MotdBE(Host)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
