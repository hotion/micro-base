package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/shiguanghuxian/micro-base/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/kit/middleware/logging"
	"github.com/shiguanghuxian/micro-base/kit/service"
	"github.com/shiguanghuxian/micro-base/kit/transport"
	"github.com/shiguanghuxian/micro-base/pb"
	"github.com/shiguanghuxian/micro-common/config"
	"github.com/shiguanghuxian/micro-common/log"
	"github.com/shiguanghuxian/micro-common/register"
	"google.golang.org/grpc"
)

var (
	grpcSvcName = config.GetSvcName("grpc")
	httpSvcName = config.GetSvcName("http")
	tcpSvcName  = config.GetSvcName("tcp")
	VERSION     string // 程序版本
	GIT_HASH    string // git hash
)

func main() {
	// 日志
	logger := log.Logger

	// 服务对象
	var s service.Service
	{
		s = service.NewBasicService()
		// 注册日志中间件
		s = logging.LoggingMiddleware(logger)(s)
	}

	// 创建所有端点
	endpoints := endpoint.MakeServerEndpoints(s)

	/* grpc */

	// transport grpc方式
	var grpcServer pb.HelloServer
	{
		grpcServer = transport.NewGRPCServer(endpoints, s, logger)
	}

	// 启动grpc监听
	go func() {
		grpcListener, err := net.Listen("tcp", config.GetGRPCAddr())
		if err != nil {
			logger.Panicf("TCP监听错误:%v", err)
			return
		}
		logger.Infof("TCP监听成功:%s", config.GetGRPCAddr())

		baseServer := grpc.NewServer()
		pb.RegisterHelloServer(baseServer, grpcServer)

		err = baseServer.Serve(grpcListener)
		if err != nil {
			logger.Panicf("TCP启动错误:%v", err)
			return
		}
	}()

	// 注册服务 grpc
	go register.Register(config.GetETCDAddr(), grpcSvcName, config.GetGRPCAddr(), 5)

	/* http */

	// transport http方式
	var httpHandler http.Handler
	{
		httpHandler = transport.NewHTTPHandler(endpoints, s, logger)
	}

	// 启动http监听
	go func() {
		httpListener, err := net.Listen("tcp", config.GetHTTPAddr())
		if err != nil {
			logger.Panicf("http监听错误:%v", err)
			return
		}
		logger.Infof("HTTP监听成功:%s", config.GetHTTPAddr())

		err = http.Serve(httpListener, httpHandler)
		if err != nil {
			logger.Panicf("HTTP启动错误:%v", err)
			return
		}
	}()

	// 注册服务 http
	go register.Register(config.GetETCDAddr(), httpSvcName, config.GetHTTPAddr(), 5)

	/* tcp */
	go func() {
		tcpServer, err := transport.NewTCPHandler(endpoints, s, logger)
		if err != nil {
			logger.Panicf("tcp服务创建错误错误:%v", err)
			return
		}
		err = tcpServer.ListenAndServe(config.GetTCPAddr())
		if err != nil {
			logger.Panicf("TCP启动错误:%v", err)
			return
		}
	}()

	// 注册服务 tcp
	go register.Register(config.GetETCDAddr(), httpSvcName, config.GetTCPAddr(), 5)

	// 监听退出信号
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	sig := <-ch
	// 注销服务
	register.UnRegister(grpcSvcName, config.GetGRPCAddr())
	register.UnRegister(httpSvcName, config.GetHTTPAddr())
	register.UnRegister(tcpSvcName, config.GetTCPAddr())

	if i, ok := sig.(syscall.Signal); ok {
		os.Exit(int(i))
	} else {
		os.Exit(0)
	}
}
