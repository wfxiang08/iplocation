package main

//import (
//	ips "git.chunyu.me/infra/iplocation/gen-go/ip_service"
//	"git.chunyu.me/infra/rpc_proxy/proxy"
//	"git.chunyu.me/infra/iplocation/ip_query"
//)
//
//const (
//	BINARY_NAME = "iplocation"
//	SERVICE_DESC = "Chunyu Ip Service v0.1"
//	IP_DATA = "/usr/local/ip/qqwry.dat"
//)
//
//var (
//	buildDate string
//	gitVersion string
//)
//
//func main() {
//
//	proxy.RpcMain(BINARY_NAME, SERVICE_DESC,
//		// 默认的ThriftServer的配置checker
//		proxy.ConfigCheckThriftService,
//
//		// 可以根据配置config来创建processor
//		func(config *proxy.Config) proxy.Server {
//			handler := ip_query.NewHandler(IP_DATA)
//			processor := ips.NewIpServiceProcessor(handler)
//			return proxy.NewThriftRpcServer(config, processor)
//		}, buildDate, gitVersion)
//}
