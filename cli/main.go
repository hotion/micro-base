package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/shiguanghuxian/micro-base/clientlib"
	"github.com/shiguanghuxian/micro-base/pb"
	"google.golang.org/grpc"
)

// 测试代码

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// c, err := clientlib.NewGRPCClient("")
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }

	// login(c)

	// grpc 测试 go run main.go
	// grpcHello()
	// tcp测试 go run main.go tcp.go
	tcpHello()

	select {}
}

func grpcHello() {
	c, err := clientlib.NewGRPCClient("")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	for i := 0; i < 100; i++ {
		resp, err := c.PostHello(context.Background(), &pb.HelloRequest{Name: fmt.Sprintf("vivi:%d", i)}, grpc.FailFast(true))
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("响应:%s\n", resp.GetWord())
		}
	}
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
