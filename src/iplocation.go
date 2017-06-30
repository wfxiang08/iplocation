package main

import (
	ips "gen-go/ip_service"
	"github.com/wfxiang08/rpc_proxy/src/proxy"
	"ip_query"
)

const (
	BINARY_NAME = "iplocation"
	SERVICE_DESC = "Chunyu Ip Service v0.1"
	IP_DATA = "/usr/local/ip/qqwry.dat"
)

var (
	buildDate string
	gitVersion string
)

func main() {

	proxy.RpcMain(BINARY_NAME, SERVICE_DESC,
		// 默认的ThriftServer的配置checker
		proxy.ConfigCheckThriftService,

		// 可以根据配置config来创建processor
		func(config *proxy.Config) proxy.Server {
			handler := ip_query.NewHandler(IP_DATA)
			processor := ips.NewIpServiceProcessor(handler)
			return proxy.NewThriftRpcServer(config, processor)
		}, buildDate, gitVersion)
}
