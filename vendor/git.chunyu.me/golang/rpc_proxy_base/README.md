# rpc_proxy_gobase 文档

* 如何创建一个Thrift RPC Client Pool呢?
* 在src/Godeps目录下添加引用

```bash
# Thrift RPC
git.chunyu.me/infra/go_thrift master
git.chunyu.me/golang/rpc_proxy_base master
```
* 在代码中添加如下的实现 

```go
package rpc_pool
import (
	"fmt"
	"git.chunyu.me/golang/chunyu_sms_client/src/chunyu_sms_service"
	"git.chunyu.me/infra/go_thrift/thrift"
	"git.chunyu.me/golang/rpc_proxy_base/src/rpc_utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//
// go test rpc_pool -v -run "TestSmsRpc"
//
func TestSmsRpc(t *testing.T) {

	var ThriftPool *rpc_utils.Pool

	//	proxyAddress := "/usr/local/rpc_proxy/proxy.sock"

	proxyAddress := "60.29.255.199:5550" // 60.29.255.199
	//	proxyAddress := "/usr/local/rpc_proxy/online_proxy.sock"

	ThriftPool = &rpc_utils.Pool{
		Dial: func() (thrift.TTransport, error) {
			// 如何创建一个Transport
			t, err := thrift.NewTSocketTimeout(proxyAddress, time.Second*5)
			trans := rpc_utils.NewTFramedTransport(t)
			return trans, err
		},
		MaxActive:   30,
		MaxIdle:     30,
		IdleTimeout: time.Second * 3600 * 24,
		Wait:        true,
	}

	transport := ThriftPool.Get()
	defer transport.Close()

	ip, op := rpc_utils.GetProtocolFromTransport(transport, "sms")
	client := chunyu_sms_service.NewSmsServiceClientProtocol(transport, ip, op)


	r, err := client.AddSendSmsRequest("18611730934", "测试信息")
	fmt.Printf("R: %v\n", r)
	fmt.Printf("R: %v\n", err)
	fmt.Printf("Response: %d, %s --> %v \n", r.Code, r.ErrorMsg, err)

	assert.True(t, true)
}
```