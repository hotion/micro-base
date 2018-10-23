package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/shiguanghuxian/micro-base/clientlib"
	"github.com/shiguanghuxian/micro-base/pb"
)

// 测试代码

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	c, err := clientlib.NewGRPCClient("")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	login(c)

	// tcp测试 go run main.go tcp.go
	// tcpHello()

	select {}
}

// 登录测试
func login(c pb.AccountClient) {
	u, err := c.Login(context.Background(), &pb.LoginRequest{
		Username: "admin",
		Password: "111111",
	})
	if err != nil {
		log.Println(err)
		return
	}
	js, _ := json.Marshal(u)
	log.Println(string(js))
}
