package clientlib

import (
	"fmt"

	"github.com/shiguanghuxian/micro-base/internal/config"
	"github.com/shiguanghuxian/micro-base/internal/register"
	"github.com/shiguanghuxian/micro-base/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	svcName = config.GetSvcName()
)

// NewGRPCClient 创建grpc客户端
func NewGRPCClient(etcdAddr string) (pb.HelloClient, error) {
	if etcdAddr == "" {
		etcdAddr = config.GetETCDAddr()
	}
	// 获取grpc服务
	r := register.NewResolver(etcdAddr)
	resolver.Register(r)

	fmt.Println(r.Scheme() + "://" + svcName)

	conn, err := grpc.Dial(r.Scheme()+"://author/"+svcName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewHelloClient(conn)
	return client, nil
}
