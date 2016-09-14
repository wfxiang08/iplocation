package zk

/*
TODO:
* make sure a ping response comes back in a reasonable time

Possible watcher events:
* Event{Type: EventNotWatching, State: StateDisconnected, Path: path, Err: err}
*/

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var ErrNoServer = errors.New("zk: could not connect to a server")

const (
	bufferSize      = 10 * 1024 * 1024
	eventChanSize   = 6
	sendChanSize    = 16
	protectedPrefix = "_c_"
)

type watchType int

const (
	watchTypeData  = iota
	watchTypeExist = iota
	watchTypeChild = iota
)

type watchPathType struct {
	path  string
	wType watchType
}

type Dialer func(network, address string, timeout time.Duration) (net.Conn, error)

// Zxid: 为zookeeper维持的id
// xid: 为本地的id
type Conn struct {
	lastZxid  int64
	sessionID int64
	state     State // must be 32-bit aligned
	xid       int32
	timeout   int32 // session timeout in seconds
	passwd    []byte

	dialer         Dialer
	servers        []string
	serverIndex    int
	conn           net.Conn
	eventChan      chan Event
	shouldQuit     chan bool
	pingInterval   time.Duration
	recvTimeout    time.Duration
	connectTimeout time.Duration

	// 待发送的请求
	sendChan chan *request

	// 正在处理的请求
	requests     map[int32]*request // Xid -> pending request
	requestsLock sync.Mutex

	// 正在处理的watchers
	watchers     map[watchPathType][]chan Event
	watchersLock sync.Mutex

	// Debug (used by unit tests)
	reconnectDelay time.Duration
}

type request struct {
	xid        int32
	opcode     int32
	pkt        interface{}
	recvStruct interface{}
	recvChan   chan response

	// Because sending and receiving happen in separate go routines, there's
	// a possible race condition when creating watches from outside the read
	// loop. We must ensure that a watcher gets added to the list synchronously
	// with the response from the server on any request that creates a watch.
	// In order to not hard code the watch logic for each opcode in the recv
	// loop the caller can use recvFunc to insert some synchronously code
	// after a response.
	recvFunc func(*request, *responseHeader, error)
}

type response struct {
	zxid int64
	err  error
}

type Event struct {
	Type  EventType
	State State
	Path  string // For non-session events, the path of the watched node.
	Err   error
}

func Connect(servers []string, recvTimeout time.Duration, zkSessionTimeout int) (*Conn, <-chan Event, error) {
	return ConnectWithDialer(servers, recvTimeout, nil, zkSessionTimeout)
}

func ConnectWithDialer(servers []string, recvTimeout time.Duration, dialer Dialer, zkSessionTimeout int) (*Conn, <-chan Event, error) {
	// 1. 预处理zk servers
	// Randomize the order of the servers to avoid creating hotspots
	stringShuffle(servers)

	for i, addr := range servers {
		if !strings.Contains(addr, ":") {
			servers[i] = addr + ":" + strconv.Itoa(DefaultPort)
		}
	}

	// 2. 创建connection&dialer方法
	ec := make(chan Event, eventChanSize)
	if dialer == nil {
		dialer = net.DialTimeout
	}
	conn := Conn{
		dialer:         dialer,
		servers:        servers,
		serverIndex:    0,
		conn:           nil, // 同一个时间点只连接到一个Server
		state:          StateDisconnected,
		eventChan:      ec,
		shouldQuit:     make(chan bool),
		recvTimeout:    recvTimeout,
		pingInterval:   time.Duration((int64(recvTimeout) / 2)),
		connectTimeout: 1 * time.Second,

		// 待发送的请求
		sendChan: make(chan *request, sendChanSize),

		//
		requests: make(map[int32]*request),
		watchers: make(map[watchPathType][]chan Event),

		passwd:  emptyPassword,
		timeout: int32(zkSessionTimeout),

		// Debug
		reconnectDelay: time.Second,
	}

	go func() {
		conn.loop()

		// 将已经发送到zk的请求处理掉，直接报错: ErrClosing(Close之后, Conn就完事，程序也该关闭了)
		conn.flushRequests(ErrClosing)
		// 所有的Watcher也会收到 ErrClosing 消息
		conn.invalidateWatches(ErrClosing)
		close(conn.eventChan)
	}()
	return &conn, ec, nil
}

func (c *Conn) Close() {
	close(c.shouldQuit)

	select {
	case <-c.queueRequest(opClose, &closeRequest{}, &closeResponse{}, nil):
	case <-time.After(time.Second):
	}
}

func (c *Conn) State() State {
	return State(atomic.LoadInt32((*int32)(&c.state)))
}

func (c *Conn) setState(state State) {
	atomic.StoreInt32((*int32)(&c.state), int32(state))
	select {
	case c.eventChan <- Event{Type: EventSession, State: state}:
	default:
		// panic("zk: event channel full - it must be monitored and never allowed to be full")
	}
}

// 选择一个Server
// 如果总是不成功，则等待，并且告知新请求， ErrNoServer
func (c *Conn) connect() {
	c.serverIndex = (c.serverIndex + 1) % len(c.servers)
	startIndex := c.serverIndex
	c.setState(StateConnecting)
	for {
		// 尝试连接， 成功则返回；失败则打印warning, 切换到下一个
		zkConn, err := c.dialer("tcp", c.servers[c.serverIndex], c.connectTimeout)
		if err == nil {
			c.conn = zkConn
			c.setState(StateConnected)
			return
		}

		log.Printf("Failed to connect to %s: %+v", c.servers[c.serverIndex], err)

		c.serverIndex = (c.serverIndex + 1) % len(c.servers)

		// 避免死循环
		if c.serverIndex == startIndex {
			c.flushUnsentRequests(ErrNoServer)
			time.Sleep(time.Second)
		}
	}
}

//
// zk的主循环体在做什么呢?
//
func (c *Conn) loop() {
	for {
		// 1. 确保能选择一个可用的Connection
		c.connect()

		// Session准备开始， 或者出现问题
		err := c.authenticate()

		switch {
		case err == ErrSessionExpired:
			// 1. 可能需要处理Session过期的问题
			c.invalidateWatches(err)
		case err != nil && c.conn != nil:
			// 2. 其他的err, 则关闭conn
			c.conn.Close()
		case err == nil:
			// 正常开启connection, 则开启两个loop
			closeChan := make(chan bool) // channel to tell send loop stop
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				// 开启sendLoop
				c.sendLoop(c.conn, closeChan)
				c.conn.Close() // causes recv loop to EOF/exit
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				// 开启recvLoop
				err = c.recvLoop(c.conn)
				if err == nil {
					panic("zk: recvLoop should never return nil error")
				}
				close(closeChan) // tell send loop to exit
				wg.Done()
			}()

			wg.Wait()
		}

		c.setState(StateDisconnected)

		// Yeesh
		if err != io.EOF && err != ErrSessionExpired && !strings.Contains(err.Error(), "use of closed network connection") {
			log.Println(err)
		}

		// 如果收到 shouldQuit, 则准备退出
		select {
		case <-c.shouldQuit:
			c.flushRequests(ErrClosing)
			return
		default:
		}

		if err != ErrSessionExpired {
			err = ErrConnectionClosed
		}
		c.flushRequests(err)

		// 等待重连
		if c.reconnectDelay > 0 {
			select {
			case <-c.shouldQuit:
				return
			case <-time.After(c.reconnectDelay):
			}
		}
	}
}

// 处理为发送的Request
// 和 flushRequests 的区别在于: request来自于 c.sendChan的buffer
func (c *Conn) flushUnsentRequests(err error) {
	for {
		select {
		default:
			return
		case req := <-c.sendChan:
			req.recvChan <- response{-1, err}
		}
	}
}

// Send error to all pending requests and clear request map
// 将当前的Request全部都处理掉(直接报错), 让clients有一个后续的处理
// 这些都是已经发送的Request, 保存在: c.requests中
//
func (c *Conn) flushRequests(err error) {

	c.requestsLock.Lock()
	for _, req := range c.requests {
		req.recvChan <- response{-1, err}
	}
	c.requests = make(map[int32]*request) // 直接reset
	c.requestsLock.Unlock()
}

// Send error to all watchers and clear watchers map
// 通知所有的Watcher, 连接断开
//
func (c *Conn) invalidateWatches(err error) {
	c.watchersLock.Lock()
	defer c.watchersLock.Unlock()

	// 出错之后，watcher有什么消息呢?
	if len(c.watchers) >= 0 {
		for pathType, watchers := range c.watchers {
			// 通过channel, 以 StateDisconnected 告知 Client
			ev := Event{Type: EventNotWatching, State: StateDisconnected, Path: pathType.path, Err: err}
			for _, ch := range watchers {
				ch <- ev
				close(ch)
			}
		}

		// 然后清空所有的watchers
		c.watchers = make(map[watchPathType][]chan Event)
	}
}

func (c *Conn) sendSetWatches() {
	c.watchersLock.Lock()
	defer c.watchersLock.Unlock()

	if len(c.watchers) == 0 {
		return
	}

	// 如何设置Watchers呢?
	req := &setWatchesRequest{
		RelativeZxid: c.lastZxid,
		DataWatches:  make([]string, 0),
		ExistWatches: make([]string, 0),
		ChildWatches: make([]string, 0),
	}

	// 1. 将watch的数据发送给zk
	n := 0
	for pathType, watchers := range c.watchers {
		if len(watchers) == 0 {
			continue
		}
		switch pathType.wType {
		case watchTypeData:
			req.DataWatches = append(req.DataWatches, pathType.path)
		case watchTypeExist:
			req.ExistWatches = append(req.ExistWatches, pathType.path)
		case watchTypeChild:
			req.ChildWatches = append(req.ChildWatches, pathType.path)
		}
		n++
	}
	if n == 0 {
		return
	}

	// 2. 异步地将watch发送给zk
	go func() {
		res := &setWatchesResponse{}
		_, err := c.request(opSetWatches, req, res, nil)
		if err != nil {
			log.Printf("Failed to set previous watches: %s", err.Error())
		}
	}()
}

func (c *Conn) authenticate() error {
	buf := make([]byte, 256)

	// connect request
	// 继续使用: sessionID等信息继续之前的Session
	n, err := encodePacket(buf[4:], &connectRequest{
		ProtocolVersion: protocolVersion,
		LastZxidSeen:    c.lastZxid,
		TimeOut:         c.timeout,
		SessionID:       c.sessionID,
		Passwd:          c.passwd,
	})
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(buf[:4], uint32(n))

	// 1. 写入建立连接的参数信息(SessionID, LastZxidSeen)等, 以维持状态的统一
	_, err = c.conn.Write(buf[:n+4])
	if err != nil {
		return err
	}

	// 2. 将watchers添加到 zk中
	c.sendSetWatches()

	// connect response

	// package length
	_, err = io.ReadFull(c.conn, buf[:4])
	if err != nil {
		return err
	}

	blen := int(binary.BigEndian.Uint32(buf[:4]))
	if cap(buf) < blen {
		buf = make([]byte, blen)
	}

	_, err = io.ReadFull(c.conn, buf[:blen])
	if err != nil {
		return err
	}

	r := connectResponse{}
	_, err = decodePacket(buf[:blen], &r)
	if err != nil {
		return err
	}

	// Session过期了，就放弃
	if r.SessionID == 0 {
		c.sessionID = 0
		c.passwd = emptyPassword
		c.lastZxid = 0
		c.setState(StateExpired)
		return ErrSessionExpired
	}

	// SessionID变化了，则重置 xid
	// 为什么为有SessionID不一致的情况呢?
	if c.sessionID != r.SessionID {
		atomic.StoreInt32(&c.xid, 0)
	}
	c.timeout = r.TimeOut
	c.sessionID = r.SessionID
	c.passwd = r.Passwd
	c.setState(StateHasSession)

	return nil
}

func (c *Conn) sendLoop(conn net.Conn, closeChan <-chan bool) error {
	pingTicker := time.NewTicker(c.pingInterval)
	defer pingTicker.Stop()

	buf := make([]byte, bufferSize)
	for {
		// 如何开始 sendLoop呢?
		// 1. sendChan: 记录了待发送到zk的请求（为一个Requests Buffer)
		select {
		case req := <-c.sendChan:
			header := &requestHeader{req.xid, req.opcode}
			n, err := encodePacket(buf[4:], header)
			// 出错返回
			if err != nil {
				req.recvChan <- response{-1, err}
				continue
			}

			// 编码pkt数据
			n2, err := encodePacket(buf[4+n:], req.pkt)
			if err != nil {
				req.recvChan <- response{-1, err}
				continue
			}

			n += n2

			// 最后将n写入buf, 注意放在buf的最前面
			binary.BigEndian.PutUint32(buf[:4], uint32(n))

			c.requestsLock.Lock()
			select {
			case <-closeChan:
				req.recvChan <- response{-1, ErrConnectionClosed}
				c.requestsLock.Unlock()
				return ErrConnectionClosed
			default:
			}
			// 记录了: requests的信息
			c.requests[req.xid] = req
			c.requestsLock.Unlock()

			conn.SetWriteDeadline(time.Now().Add(c.recvTimeout))
			_, err = conn.Write(buf[:n+4])
			conn.SetWriteDeadline(time.Time{})

			if err != nil {
				req.recvChan <- response{-1, err}
				conn.Close()
				return err
			}
		case <-pingTicker.C:
			// 定期和zk ping
			n, err := encodePacket(buf[4:], &requestHeader{Xid: -2, Opcode: opPing})
			if err != nil {
				panic("zk: opPing should never fail to serialize")
			}

			binary.BigEndian.PutUint32(buf[:4], uint32(n))

			conn.SetWriteDeadline(time.Now().Add(c.recvTimeout))
			_, err = conn.Write(buf[:n+4])
			conn.SetWriteDeadline(time.Time{})
			if err != nil {
				conn.Close()
				return err
			}
		case <-closeChan:
			return nil
		}
	}
}

func (c *Conn) recvLoop(conn net.Conn) error {
	buf := make([]byte, bufferSize)
	for {
		// package length
		conn.SetReadDeadline(time.Now().Add(c.recvTimeout))

		// 1. 读取block length
		_, err := io.ReadFull(conn, buf[:4])
		if err != nil {
			return err
		}
		blen := int(binary.BigEndian.Uint32(buf[:4]))
		if cap(buf) < blen {
			buf = make([]byte, blen)
		}

		// 2. 读取一个完整的block数据
		_, err = io.ReadFull(conn, buf[:blen])
		conn.SetReadDeadline(time.Time{})
		if err != nil {
			return err
		}

		// 16个字节的Header
		res := responseHeader{}
		_, err = decodePacket(buf[:16], &res)
		if err != nil {
			return err
		}

		if res.Xid == -1 {
			res := &watcherEvent{}

			// blen这个地方是否存在bug呢?
			_, err := decodePacket(buf[16:16+blen], res)
			if err != nil {
				return err
			}
			ev := Event{
				Type:  res.Type,
				State: res.State,
				Path:  res.Path,
				Err:   nil,
			}
			// 如果c.eventChan有效，将数据输出
			select {
			case c.eventChan <- ev:
			default:
			}

			// 来自zk的Node
			wTypes := make([]watchType, 0, 2)
			switch res.Type {
			case EventNodeCreated:
				wTypes = append(wTypes, watchTypeExist)
			case EventNodeDeleted, EventNodeDataChanged:
				wTypes = append(wTypes, watchTypeExist, watchTypeData, watchTypeChild)
			case EventNodeChildrenChanged:
				wTypes = append(wTypes, watchTypeChild)
			}

			// 根据不同的eventType, 将事件推送到不同的watcher上
			c.watchersLock.Lock()
			for _, t := range wTypes {
				wpt := watchPathType{res.Path, t}
				if watchers := c.watchers[wpt]; watchers != nil && len(watchers) > 0 {
					for _, ch := range watchers {
						ch <- ev
						close(ch)
					}
					delete(c.watchers, wpt)
				}
			}
			c.watchersLock.Unlock()
		} else if res.Xid == -2 {
			// Ping response. Ignore.
		} else if res.Xid < 0 {
			log.Printf("Xid < 0 (%d) but not ping or watcher event", res.Xid)
		} else {
			// Zxid
			if res.Zxid > 0 {
				c.lastZxid = res.Zxid
			}

			// Xid --> 本地的request
			c.requestsLock.Lock()
			req, ok := c.requests[res.Xid]
			if ok {
				delete(c.requests, res.Xid)
			}
			c.requestsLock.Unlock()

			if !ok {
				log.Printf("Response for unknown request with xid %d", res.Xid)
			} else {
				if res.Err != 0 {
					err = res.Err.toError()
				} else {
					// 本地的req获取response
					_, err = decodePacket(buf[16:16+blen], req.recvStruct)
				}
				if req.recvFunc != nil {
					req.recvFunc(req, &res, err)
				}
				req.recvChan <- response{res.Zxid, err}
				if req.opcode == opClose {
					return io.EOF
				}
			}
		}
	}
}

func (c *Conn) nextXid() int32 {
	return atomic.AddInt32(&c.xid, 1)
}

func (c *Conn) addWatcher(path string, watchType watchType) <-chan Event {
	c.watchersLock.Lock()
	defer c.watchersLock.Unlock()

	ch := make(chan Event, 1)
	wpt := watchPathType{path, watchType}

	// wpt作为key保存到: watchers中了
	c.watchers[wpt] = append(c.watchers[wpt], ch)
	return ch
}

func (c *Conn) queueRequest(opcode int32, req interface{}, res interface{}, recvFunc func(*request, *responseHeader, error)) <-chan response {
	rq := &request{
		xid:        c.nextXid(),
		opcode:     opcode,
		pkt:        req,
		recvStruct: res,
		recvChan:   make(chan response, 1),
		recvFunc:   recvFunc,
	}
	c.sendChan <- rq
	return rq.recvChan
}

// 创建到zk的Request, 并且缓存到Conn的sendChan中
func (c *Conn) request(opcode int32, req interface{}, res interface{}, recvFunc func(*request, *responseHeader, error)) (int64, error) {
	r := <-c.queueRequest(opcode, req, res, recvFunc)
	return r.zxid, r.err
}

func (c *Conn) AddAuth(scheme string, auth []byte) error {
	_, err := c.request(opSetAuth, &setAuthRequest{Type: 0, Scheme: scheme, Auth: auth}, &setAuthResponse{}, nil)
	return err
}

func (c *Conn) Children(path string) ([]string, Stat, error) {
	res := &getChildren2Response{}
	_, err := c.request(opGetChildren2, &getChildren2Request{Path: path, Watch: false}, res, nil)
	return res.Children, &res.Stat, err
}

// 如何监听一个Path呢?
func (c *Conn) ChildrenW(path string) ([]string, Stat, <-chan Event, error) {
	var ech <-chan Event
	res := &getChildren2Response{}
	req := &getChildren2Request{Path: path, Watch: true} // 获取Children数据的同时，也注册了一个Watcher
	callback := func(req *request, res *responseHeader, err error) {
		if err == nil {
			// 添加了ChildWatcher
			ech = c.addWatcher(path, watchTypeChild)
		}
	}
	_, err := c.request(opGetChildren2, req, res, callback)
	if err != nil {
		return nil, nil, nil, err
	}
	return res.Children, &res.Stat, ech, err
}

func (c *Conn) Get(path string) ([]byte, Stat, error) {
	res := &getDataResponse{}
	_, err := c.request(opGetData, &getDataRequest{Path: path, Watch: false}, res, nil)
	return res.Data, &res.Stat, err
}

// GetW returns the contents of a znode and sets a watch
func (c *Conn) GetW(path string) ([]byte, Stat, <-chan Event, error) {
	var ech <-chan Event
	res := &getDataResponse{}
	req := &getDataRequest{Path: path, Watch: true}

	// 首先获取数据，然后再添加Watch
	callback := func(req *request, res *responseHeader, err error) {
		if err == nil {
			ech = c.addWatcher(path, watchTypeData)
		}
	}
	_, err := c.request(opGetData, req, res, callback)
	if err != nil {
		return nil, nil, nil, err
	}
	return res.Data, &res.Stat, ech, err
}

func (c *Conn) Set(path string, data []byte, version int32) (Stat, error) {
	res := &setDataResponse{}
	_, err := c.request(opSetData, &SetDataRequest{path, data, version}, res, nil)
	return &res.Stat, err
}

func (c *Conn) Create(path string, data []byte, flags int32, acl []ACL) (string, error) {
	res := &createResponse{}
	_, err := c.request(opCreate, &CreateRequest{path, data, acl, flags}, res, nil)
	return res.Path, err
}

// CreateProtectedEphemeralSequential fixes a race condition if the server crashes
// after it creates the node. On reconnect the session may still be valid so the
// ephemeral node still exists. Therefore, on reconnect we need to check if a node
// with a GUID generated on create exists.
func (c *Conn) CreateProtectedEphemeralSequential(path string, data []byte, acl []ACL) (string, error) {
	var guid [16]byte
	_, err := io.ReadFull(rand.Reader, guid[:16])
	if err != nil {
		return "", err
	}
	guidStr := fmt.Sprintf("%x", guid)

	parts := strings.Split(path, "/")
	parts[len(parts)-1] = fmt.Sprintf("%s%s-%s", protectedPrefix, guidStr, parts[len(parts)-1])
	rootPath := strings.Join(parts[:len(parts)-1], "/")
	protectedPath := strings.Join(parts, "/")

	var newPath string
	for i := 0; i < 3; i++ {
		newPath, err = c.Create(protectedPath, data, FlagEphemeral|FlagSequence, acl)
		switch err {
		case ErrSessionExpired:
			// No need to search for the node since it can't exist. Just try again.
		case ErrConnectionClosed:
			children, _, err := c.Children(rootPath)
			if err != nil {
				return "", err
			}
			for _, p := range children {
				parts := strings.Split(p, "/")
				if pth := parts[len(parts)-1]; strings.HasPrefix(pth, protectedPrefix) {
					if g := pth[len(protectedPrefix) : len(protectedPrefix)+32]; g == guidStr {
						return rootPath + "/" + p, nil
					}
				}
			}
		case nil:
			return newPath, nil
		default:
			return "", err
		}
	}
	return "", err
}

func (c *Conn) Delete(path string, version int32) error {
	_, err := c.request(opDelete, &DeleteRequest{path, version}, &deleteResponse{}, nil)
	return err
}

func (c *Conn) Exists(path string) (bool, Stat, error) {
	res := &existsResponse{}
	_, err := c.request(opExists, &existsRequest{Path: path, Watch: false}, res, nil)
	exists := true
	if err == ErrNoNode {
		exists = false
		err = nil
	}
	return exists, &res.Stat, err
}

func (c *Conn) ExistsW(path string) (bool, Stat, <-chan Event, error) {
	var ech <-chan Event
	res := &existsResponse{}
	_, err := c.request(opExists, &existsRequest{Path: path, Watch: true}, res, func(req *request, res *responseHeader, err error) {
		if err == nil {
			ech = c.addWatcher(path, watchTypeData)
		} else if err == ErrNoNode {
			ech = c.addWatcher(path, watchTypeExist)
		}
	})
	exists := true
	if err == ErrNoNode {
		exists = false
		err = nil
	}
	if err != nil {
		return false, nil, nil, err
	}
	return exists, &res.Stat, ech, err
}

func (c *Conn) GetACL(path string) ([]ACL, Stat, error) {
	res := &getAclResponse{}
	_, err := c.request(opGetAcl, &getAclRequest{Path: path}, res, nil)
	return res.Acl, &res.Stat, err
}

func (c *Conn) SetACL(path string, acl []ACL, version int32) (Stat, error) {
	res := &setAclResponse{}
	_, err := c.request(opSetAcl, &setAclRequest{Path: path, Acl: acl, Version: version}, res, nil)
	return &res.Stat, err
}

func (c *Conn) Sync(path string) (string, error) {
	res := &syncResponse{}
	_, err := c.request(opSync, &syncRequest{Path: path}, res, nil)
	return res.Path, err
}

type MultiOps struct {
	Create  []CreateRequest
	Delete  []DeleteRequest
	SetData []SetDataRequest
	Check   []CheckVersionRequest
}

func (c *Conn) Multi(ops MultiOps) error {
	req := &multiRequest{
		Ops:        make([]multiRequestOp, 0, len(ops.Create)+len(ops.Delete)+len(ops.SetData)+len(ops.Check)),
		DoneHeader: multiHeader{Type: -1, Done: true, Err: -1},
	}
	for _, r := range ops.Create {
		req.Ops = append(req.Ops, multiRequestOp{multiHeader{opCreate, false, -1}, r})
	}
	for _, r := range ops.SetData {
		req.Ops = append(req.Ops, multiRequestOp{multiHeader{opSetData, false, -1}, r})
	}
	for _, r := range ops.Delete {
		req.Ops = append(req.Ops, multiRequestOp{multiHeader{opDelete, false, -1}, r})
	}
	for _, r := range ops.Check {
		req.Ops = append(req.Ops, multiRequestOp{multiHeader{opCheck, false, -1}, r})
	}
	res := &multiResponse{}
	_, err := c.request(opMulti, req, res, nil)
	return err
}
