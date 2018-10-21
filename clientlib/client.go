package clientlib

import (
	"log"
	"time"

	"github.com/shiguanghuxian/micro-base/internal/config"
	"github.com/shiguanghuxian/micro-base/internal/register"
	"github.com/shiguanghuxian/micro-base/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	svcName = config.GetSvcName("grpc")
)

// NewGRPCClient 创建grpc客户端
func NewGRPCClient(etcdAddr string) (pb.HelloClient, error) {
	if etcdAddr == "" {
		etcdAddr = config.GetETCDAddr()
	}
	// 获取grpc服务
	r := register.NewResolver(etcdAddr)
	resolver.Register(r)

	log.Println(r.Scheme() + "://" + svcName)

	conn, err := grpc.Dial(r.Scheme()+"://author/"+svcName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure(), grpc.WithTimeout(15*time.Second))
	if err != nil {
		return nil, err
	}
	client := pb.NewHelloClient(conn)
	return client, nil
}
