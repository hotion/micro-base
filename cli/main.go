package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shiguanghuxian/micro-base/clientlib"
	"github.com/shiguanghuxian/micro-base/pb"
	"google.golang.org/grpc"
)

// 测试代码

func main() {
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
