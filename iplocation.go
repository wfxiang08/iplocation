package main

import (
	"fmt"
	thrift "git.apache.org/thrift.git/lib/go/thrift"
	ips "git.chunyu.me/infra/ip_utils/gen-go/ip_service"
	ip_query "git.chunyu.me/infra/ip_utils/ip_query"
	"github.com/docopt/docopt-go"
	color "github.com/fatih/color"
	zmq "github.com/pebbe/zmq4"
	topozk "github.com/wandoulabs/go-zookeeper/zk"
	utils "github.com/wfxiang08/rpc_proxy/utils"
	"github.com/wfxiang08/rpc_proxy/utils/bytesize"
	"github.com/wfxiang08/rpc_proxy/utils/log"
	zk "github.com/wfxiang08/rpc_proxy/zk"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	HEARTBEAT_INTERVAL = 1000
)

var magenta = color.New(color.FgMagenta).SprintFunc()

var verbose = false
var usage = `usage: ip_service -c <config_file>[-L <log_file>] [--log-level=<loglevel>] [--log-filesize=<filesize>] 

options:
   -c <config_file>
   -L	set output log file, default is stdout
   --log-level=<loglevel>	set log level: info, warn, error, debug [default: info]
   --log-filesize=<maxsize>  set max log file size, suffixes "KB", "MB", "GB" are allowed, 1KB=1024 bytes, etc. Default is 1GB.
`

func main() {
	args, err := docopt.Parse(usage, nil, true, "Chunyu Ip Service v0.1", true)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var maxFileFrag = 2
	var maxFragSize int64 = bytesize.GB * 1
	if s, ok := args["--log-filesize"].(string); ok && s != "" {
		v, err := bytesize.Parse(s)
		if err != nil {
			log.PanicErrorf(err, "invalid max log file size = %s", s)
		}
		maxFragSize = v
	}

	// set output log file
	if s, ok := args["-L"].(string); ok && s != "" {
		f, err := log.NewRollingFile(s, maxFileFrag, maxFragSize)
		if err != nil {
			log.PanicErrorf(err, "open rolling log file failed: %s", s)
		} else {
			defer f.Close()
			log.StdLog = log.New(f, "")
		}
	}
	log.SetLevel(log.LEVEL_INFO)
	log.SetFlags(log.Flags() | log.Lshortfile)

	// set log level
	if s, ok := args["--log-level"].(string); ok && s != "" {
		setLogLevel(s)
	}
	var frontendAddr, zkAddr, productName, serviceName string

	// set config file

	configFile := args["-c"].(string)
	conf, err := utils.LoadConf(configFile)
	if err != nil {
		log.PanicErrorf(err, "load config failed")
	}
	productName = conf.ProductName
	if productName == "" {
		// 既没有config指定，也没有命令行指定，则报错
		log.PanicErrorf(err, "Invalid ProductName")
	}

	if conf.FrontHost == "" {
		fmt.Println("FrontHost: ", conf.FrontHost, ", Prefix: ", conf.IpPrefix)
		if conf.IpPrefix != "" {
			conf.FrontHost = utils.GetIpWithPrefix(conf.IpPrefix)
		}
	}
	if conf.FrontPort != "" && conf.FrontHost != "" {
		frontendAddr = fmt.Sprintf("tcp://%s:%s", conf.FrontHost, conf.FrontPort)
	}
	if frontendAddr == "" {
		log.PanicErrorf(err, "Invalid frontend address")
	}

	serviceName = conf.Service
	if serviceName == "" {
		log.PanicErrorf(err, "Invalid ServiceName")
	}

	zkAddr = conf.ZkAddr
	if zkAddr == "" {
		log.PanicErrorf(err, "Invalid zookeeper address")
	}
	verbose = conf.Verbose

	// 正式的服务
	mainBody(zkAddr, productName, serviceName, frontendAddr)
}

// tcp://127.0.0.1:5555 --> tcp://127_0_0_1:5555
func GetServiceIdentity(frontendAddr string) string {
	fid := strings.Replace(frontendAddr, ".", "_", -1)
	fid = strings.Replace(fid, ":", "_", -1)
	fid = strings.Replace(fid, "//", "", -1)
	return fid
}

func mainBody(zkAddr string, productName string, serviceName string, frontendAddr string) {
	// 1. 创建到zk的连接
	var topo *zk.Topology
	topo = zk.NewTopology(productName, zkAddr)

	// 2. 启动服务
	frontend, _ := zmq.NewSocket(zmq.ROUTER)
	defer frontend.Close()

	// ROUTER/ROUTER绑定到指定的端口

	// tcp://127.0.0.1:5555 --> tcp://127_0_0_1:5555
	lbServiceName := GetServiceIdentity(frontendAddr)

	frontend.SetIdentity(lbServiceName)
	frontend.Bind(frontendAddr)

	log.Printf("FrontAddr: %s\n", magenta(frontendAddr))

	poller1 := zmq.NewPoller()
	poller1.Add(frontend, zmq.POLLIN)

	// 3. 注册zk
	var endpointInfo map[string]interface{} = make(map[string]interface{})
	endpointInfo["frontend"] = frontendAddr

	topo.AddServiceEndPoint(serviceName, lbServiceName, endpointInfo)

	isAlive := true
	isAliveLock := &sync.RWMutex{}

	go func() {
		servicePath := topo.ProductServicePath(serviceName)
		evtbus := make(chan interface{})
		for true {
			// 只是为了监控状态
			_, err := topo.WatchNode(servicePath, evtbus)

			if err == nil {
				// 等待事件
				e := (<-evtbus).(topozk.Event)
				if e.State == topozk.StateExpired || e.Type == topozk.EventNotWatching {
					// Session过期了，则需要删除之前的数据，因为这个数据的Owner不是当前的Session
					topo.DeleteServiceEndPoint(serviceName, lbServiceName)
					topo.AddServiceEndPoint(serviceName, lbServiceName, endpointInfo)
				}
			} else {
				time.Sleep(time.Second)
			}

			isAliveLock.RLock()
			isAlive1 := isAlive
			isAliveLock.RUnlock()
			if !isAlive1 {
				break
			}

		}
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	// syscall.SIGKILL
	// kill -9 pid
	// kill -s SIGKILL pid 还是留给运维吧
	//
	handler := ip_query.NewHandler("qqwry.dat")
	processor := ips.NewIpServiceProcessor(handler)

	// 自动退出条件:
	//

	var suideTime time.Time

	for {
		var sockets []zmq.Polled
		var err error

		sockets, err = poller1.Poll(time.Second)
		if err != nil {
			log.Errorf("Error When Pollling: %v\n", err)
			continue
		}

		hasValidMsg := false
		log.Printf("Sockets: %d\n", len(sockets))
		for _, socket := range sockets {
			switch socket.Socket {
			case frontend:
				hasValidMsg = true
				if verbose {
					log.Println("----->Message from front: ")
				}
				msgs, err := frontend.RecvMessage(0)
				if err != nil {
					log.Errorf("Error when reading from frontend: %v\n", err)
					continue
				}

				// msgs:
				// <proxy_id, "", client_id, "", rpc_data>
				if verbose {
					utils.PrintZeromqMsgs(msgs, "frontend")
				}
				msgs = utils.TrimLeftEmptyMsg(msgs)

				bufferIn := thrift.NewTMemoryBufferLen(0)
				bufferIn.WriteString(msgs[len(msgs)-1])
				procIn := thrift.NewTBinaryProtocolTransport(bufferIn)
				bufferOut := thrift.NewTMemoryBufferLen(0)
				procOut := thrift.NewTBinaryProtocolTransport(bufferOut)
				if verbose {
					log.Println("----->Message from front Process: ")
				}
				processor.Process(procIn, procOut)

				result := bufferOut.Bytes()
				// <proxy_id, "", client_id, "", rpc_data>
				frontend.SendMessage(msgs[0:(len(msgs)-1)], result)
			}
		}

		// 如果安排的suiside, 则需要处理 suiside的时间
		isAliveLock.RLock()
		isAlive1 := isAlive
		isAliveLock.RUnlock()

		if !isAlive1 {
			if hasValidMsg {
				suideTime = time.Now().Add(time.Second * 3)
			} else {
				if time.Now().After(suideTime) {
					log.Println(utils.Green("Load Balance Suiside Gracefully"))
					break
				}
			}
		}

		// 心跳同步
		select {
		case sig := <-ch:
			isAliveLock.Lock()
			isAlive1 := isAlive
			isAlive = false
			isAliveLock.Unlock()

			if isAlive1 {
				// 准备退出(但是需要处理完毕手上的活)

				// 需要退出:
				topo.DeleteServiceEndPoint(serviceName, lbServiceName)

				if sig == syscall.SIGKILL {
					log.Println(utils.Red("Got Kill Signal, Return Directly"))
					break
				} else {
					suideTime = time.Now().Add(time.Second * 3)
					log.Println(utils.Red("Schedule to suicide at: "), suideTime.Format("@2006-01-02 15:04:05"))
				}
			}
		default:
		}
	}
}

func init() {
	log.SetLevel(log.LEVEL_INFO)
}

func setLogLevel(level string) {
	var lv = log.LEVEL_INFO
	switch strings.ToLower(level) {
	case "error":
		lv = log.LEVEL_ERROR
	case "warn", "warning":
		lv = log.LEVEL_WARN
	case "debug":
		lv = log.LEVEL_DEBUG
	case "info":
		fallthrough
	default:
		lv = log.LEVEL_INFO
	}
	log.SetLevel(lv)
	log.Infof("set log level to %s", lv)
}

func setCrashLog(file string) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.InfoErrorf(err, "cannot open crash log file: %s", file)
	} else {
		syscall.Dup2(int(f.Fd()), 2)
	}
}
