package main

import (
	"fmt"

	"github.com/xrcuo/api_boot/YiYan"
)

func main() {
	W, _ := YiYan.GetYiYan()
	fmt.Println(W.Text)
}
