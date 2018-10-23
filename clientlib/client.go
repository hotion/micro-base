package clientlib

import (
	"log"
	"time"

	"github.com/shiguanghuxian/micro-base/pb"
	"github.com/shiguanghuxian/micro-common/config"
	"github.com/shiguanghuxian/micro-common/register"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	svcName = config.GetSvcName("grpc")
)

// NewGRPCClient 创建grpc客户端
func NewGRPCClient(etcdAddr string) (pb.AccountClient, error) {
	if etcdAddr == "" {
		etcdAddr = config.GetETCDAddr()
	}
	// 获取grpc服务
	r := register.NewResolver(etcdAddr)
	resolver.Register(r)

	log.Println(r.Scheme() + "://" + svcName)

	conn, err := grpc.Dial(r.Scheme()+"://author/"+svcName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		return nil, err
	}
	client := pb.NewAccountClient(conn)
	return client, nil
}
