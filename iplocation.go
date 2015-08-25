package main

import (
	ips "git.chunyu.me/infra/iplocation/gen-go/ip_service"
	ip_query "git.chunyu.me/infra/iplocation/ip_query"
	rpc_commons "git.chunyu.me/infra/rpc_commons"
	utils "git.chunyu.me/infra/rpc_proxy/utils"
)

const (
	BINARY_NAME  = "iplocation"
	SERVICE_DESC = "Chunyu Ip Service v0.1"
	IP_DATA      = "qqwry.dat"
)

func main() {

	rpc_commons.RpcMain(BINARY_NAME, SERVICE_DESC,
		// 默认的ThriftServer的配置checker
		rpc_commons.ConfigCheckThriftService,

		// 可以根据配置config来创建processor
		func(config *utils.Config) rpc_commons.Server {

			handler := ip_query.NewHandler(IP_DATA)
			processor := ips.NewIpServiceProcessor(handler)
			return rpc_commons.NewThriftRpcServer(config, processor)
		})
}
