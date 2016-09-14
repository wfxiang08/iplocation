//
//  Paranoid Pirate queue. 参考: http://zguide.zeromq.org/php:chapter4
//
package main

import (
	proxy "git.chunyu.me/infra/rpc_proxy/proxy"
)

const (
	BINARY_NAME  = "rpc_lb"
	SERVICE_DESC = "Chunyu RPC Load Balance Service"
)

var (
	buildDate  string
	gitVersion string
)

func main() {
	proxy.RpcMain(BINARY_NAME, SERVICE_DESC,
		// 验证LB的配置
		proxy.ConfigCheckRpcLB,
		// 根据配置创建一个Server
		func(config *proxy.Config) proxy.Server {
			return proxy.NewThriftLoadBalanceServer(config)
		}, buildDate, gitVersion)

}
