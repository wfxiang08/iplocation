package main

import (
	"fmt"
	ips "git.chunyu.me/infra/iplocation/gen-go/ip_service"
	"git.chunyu.me/golang/rpc_proxy_base/src/rpc_utils"
	"git.chunyu.me/infra/rpc_proxy/proxy"
	"sync"
	"time"
	"git.chunyu.me/infra/go_thrift/thrift"
)

const (
	IP_DATA = "/usr/local/ip/qqwry.dat"
)

func main() {
	// 假定proxy, iplocation都启动改了
	// 做啥呢?
	useProxy := true
	var (
		sockFile string
		socket thrift.TTransport
		protocol thrift.TProtocol
	)
	wait := &sync.WaitGroup{}

	wait.Add(1)
	go func() {
		var client *ips.IpServiceClient
		if useProxy {
			sockFile = "/usr/local/rpc_proxy/online_proxy.sock"
			sk, _ := rpc_utils.NewTUnixDomain(sockFile)
			//				sockFile = "127.0.0.1:5550"
			//				sk, _ := thrift.NewTSocket(sockFile)

			sk.Open()
			defer sk.Close()
			sk.SetTimeout(time.Second * 5)
			socket = sk
		} else {
			sockFile = "127.0.0.1:5563"
			sk, _ := thrift.NewTSocket(sockFile)

			sk.Open()
			defer sk.Close()
			sk.SetTimeout(time.Second * 5)
			socket = sk
		}

		transport := proxy.NewTBufferedFramedTransport(socket, 0, 0)
		framedTransport := thrift.NewTBinaryProtocol(transport, false, true)
		if useProxy {
			protocol = thrift.NewTMultiplexedProtocol(framedTransport, "iplocation")
		} else {
			protocol = framedTransport
		}

		client = ips.NewIpServiceClientProtocol(transport, protocol, protocol)
		client.SeqId = 0

		t1 := time.Now().UnixNano()
		iteration := 100000
		for k := 0; k < iteration; k++ {

			err := client.Ping()
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
		t2 := time.Now().UnixNano()
		fmt.Printf("T: %.3fms", float64(t2 - t1) * 0.000001 / float64(iteration))

		wait.Done()
	}()

	wait.Wait()
	fmt.Println("================ DONE ====================")

}
