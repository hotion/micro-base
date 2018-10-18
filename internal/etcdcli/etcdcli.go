package etcdcli

import (
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/shiguanghuxian/micro-base/internal/config"
)

/* etcd3 连接对象 */

var (
	// EtcdCli etcd连接对象
	EtcdCli *clientv3.Client
)

func init() {
	etcdAddr := config.GetETCDAddr()
	var err error
	EtcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ";"),
		DialTimeout: 15 * time.Second,
	})
	if err != nil {
		log.Panicln(err)
	}
}
