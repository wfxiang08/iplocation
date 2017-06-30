package main

import (
	"fmt"
	"git.chunyu.me/infra/go_thrift/thrift"
	ips "git.chunyu.me/infra/iplocation/gen-go/ip_service"
	"git.chunyu.me/golang/rpc_proxy_base/src/rpc_utils"
	"git.chunyu.me/infra/rpc_proxy/proxy"
	log "git.chunyu.me/golang/cyutils/utils/rolling_log"
	"math/rand"
	"sync"
	"time"
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
		socket   thrift.TTransport
		protocol thrift.TProtocol
	)
	wait := &sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wait.Add(1)
		go func() {
			var client *ips.IpServiceClient
			if useProxy {
				sockFile = "/usr/local/rpc_proxy/proxy.sock"
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
			client.SeqId = int32(i * 100000)

			t1 := time.Now().UnixNano()
			for k := 0; k < 10000; k++ {
				interval := rand.Int63n(100)
				time.Sleep(time.Duration(time.Microsecond * time.Duration(interval)))

				ip := rand.Uint32()

				ipStr := fmt.Sprintf("%d.%d.%d.%d", ip&0xFF, (ip>>8)&0xFF, (ip>>16)&0xFF, (ip>>24)&0xFF)
				location, err := client.IpToLocation(ipStr)
				if location != nil {
					//					log.Printf("%s ==> %s %s", ipStr, location.City, location.Province)
				} else {
					log.ErrorErrorf(err, proxy.Red("%s ==> Error: %s, Index: %d-[%d]"), ipStr, err, k, i)
					break
				}
			}
			t2 := time.Now().UnixNano()
			fmt.Printf("T: %.3fms", float64(t2-t1)*0.000001)

			wait.Done()
		}()
	}

	wait.Wait()
	fmt.Println("================ DONE ====================")

}
