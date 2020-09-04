package livewebsocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	zlmeidahook "golearn/liveroom/httpser"

	"github.com/gorilla/websocket"
)

const (
	// 允许等待的写入时间
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second * 2

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

// 客户端读写消息
type wsMessage struct {
	// websocket.TextMessage 消息类型
	messageType int
	data        []byte
}

type roominfo struct {
	singlelock sync.Mutex
	Conn       map[string]*wsConnection //登陆名字  连接端口
}

func (roomlive *roominfo) sendmsg(username string, dat map[string]string) {
	for k, v := range roomlive.Conn {
		//fmt.Println("zxnmcbzxmcbzx", k, v)
		if username != k {
			fmt.Println("sending")
			v.mutex.Lock()
			if v.isClosed {
				roomlive.singlelock.Lock()
				delete(roomlive.Conn, k)
				roomlive.singlelock.Unlock()
			}
			v.mutex.Unlock()
			v.wsWrite(1, v.maptobyte(dat))
		}
	}
}
func intrefacetostring(in interface{}) string {
	var tmpschema string
	switch in.(type) {
	case string:
		{
			tmpschema = in.(string)
		}
	}
	return tmpschema
}
func zlreportroom(msg map[string]interface{}) {
	//暂时就处理消息直播间上报
	if value, ok := msg["schema"]; ok {
		var tmpschema string
		tmpschema = intrefacetostring(value)
		if "rtmp" == tmpschema {
			tmpapp := intrefacetostring(msg["app"])
			stream := intrefacetostring(msg["stream"])
			liveroom := tmpapp + stream
			liveroominfo.rommlock.Lock()

			zlroom := &roominfo{
				Conn: make(map[string]*wsConnection),
			}
			liveroominfo.roommap[liveroom] = zlroom
			liveroominfo.rommlock.Unlock()
			return
		}
	}
	fmt.Println("xxxsjd ", liveroominfo.roommap)
}

type sccliveroom struct {
	rommlock  sync.Mutex
	roommap   map[string]*roominfo //直播间名字  直播间
	connlock  sync.Mutex
	wsConnAll map[string]*wsConnection
}

// ws 的所有连接
// 用于广播
var liveroominfo sccliveroom
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有的CORS 跨域请求，正式环境可以关闭
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端连接
type wsConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan   chan *wsMessage // 读队列
	outChan  chan *wsMessage // 写队列

	mutex     sync.Mutex // 避免重复关闭管道,加锁处理
	isClosed  bool
	closeChan chan byte            // 关闭通知
	username  string               //客户端登陆的名字 默认肯定是“” 加个定时器 过了多久还没用用户名登陆就让它滚蛋
	sccticker *time.Ticker         //定时器，如果连接没登陆就用来做校验的定时器，如果登陆了就用来做心跳的定时器
	sccroom   map[string]*roominfo //ws  所连接的直播间 可能不止一个
}

func wsHandler(resp http.ResponseWriter, req *http.Request) {
	// 应答客户端告知升级连接为websocket
	wsSocket, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println("升级为websocket失败", err.Error())
		return
	}
	// TODO 如果要控制连接数可以计算，wsConnAll长度
	// 连接数保持一定数量，超过的部分不提供服务
	wsConn := &wsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *wsMessage, 1000),
		outChan:   make(chan *wsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
		sccticker: time.NewTicker(time.Second * 30),
		sccroom:   make(map[string]*roominfo),
	}
	//wsConnAll[maxConnId] = wsConn
	//	log.Println("当前在线人数", len(wsConnAll))

	// 处理器,发送定时信息，避免意外关闭

	go wsConn.processLoop() //处理携程
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()
}
func (wsConn *wsConnection) chenklogin() bool {
	if wsConn.username == "" {
		return false
	}
	return true
}

// 处理队列中的消息
func (wsConn *wsConnection) processLoop() {
	// 处理消息队列中的消息
	// 获取到消息队列中的消息，处理完成后，发送消息给客户端
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	for {
		msg, err := wsConn.wsRead()
		if err != nil {
			log.Println("获取消息出现错误222", err.Error())
			break
		}
		if nil == msg {
			log.Println("channel closed end loop")
			break
		}
		var dat map[string]string
		if err := json.Unmarshal([]byte(string(msg.data)), &dat); err == nil {
			//  证明是json，处理信令
			if value, ok := dat["msgtype"]; ok {
				//存在
				switch value {
				case "login":
					{
						reply := make(map[string]string)
						reason := wsConn.handlelogin(dat)
						if "" == reason {
							reply["result"] = "success"
							wsConn.wsWrite(msg.messageType, wsConn.maptobyte(reply))
						} else {
							reply["result"] = "failed"
							reply["reason"] = reason
							wsConn.wsWrite(msg.messageType, wsConn.maptobyte(reply))
							wsConn.close()
						}

						break
					}
				case "sendmsg":
					{
						reply := make(map[string]string)
						if !wsConn.chenklogin() {
							reply["result"] = "failed"
							reply["reason"] = "not login"
						} else {
							ret := wsConn.handlesendmsg(dat)

							reply["result"] = "success"
							if "" != ret {
								reply["reason"] = ret
							}
						}

						wsConn.wsWrite(msg.messageType, wsConn.maptobyte(reply))
						break
					}
				case "enterroom":
					{

						reply := make(map[string]string)
						if !wsConn.chenklogin() {
							reply["result"] = "failed"
							reply["reason"] = "not login"
						} else {
							ret := wsConn.enterroom(dat)

							if "" == ret {
								reply["result"] = "success"
							} else {
								reply["result"] = "failed"
								reply["data"] = ret
							}
						}

						wsConn.wsWrite(msg.messageType, wsConn.maptobyte(reply))
						break
					}
				case "leaveroom":
					{
						reply := make(map[string]string)
						if !wsConn.chenklogin() {
							reply["result"] = "failed"
							reply["reason"] = "not login"
						} else {
							wsConn.leaveroom(dat)

							reply["result"] = "success"
							wsConn.wsWrite(msg.messageType, wsConn.maptobyte(reply))
						}

						break
					}
				default:
					{
						break
					}
				}
			}
		} else {
			fmt.Println(err)
			wsConn.close()
		}
	}

}

// 处理消息队列中的消息
func (wsConn *wsConnection) wsReadLoop() {
	// 设置消息的最大长度
	//wsConn.wsSocket.SetReadLimit(maxMessageSize)//暂时屏蔽
	//wsConn.wsSocket.SetReadDeadline(time.Now().Add(pongWait))
	for {
		// 读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
			log.Println("消息读取出现错误", err.Error())
			wsConn.close()
			return
		}
		req := &wsMessage{
			msgType,
			data,
		}
		// 放入请求队列,消息入栈
		select {
		case wsConn.inChan <- req:
			{
				break
			}
		case <-wsConn.closeChan:
			{
				fmt.Println("closing the socket  3")
				return
			}

		}
	}
}

// 发送消息给客户端
func (wsConn *wsConnection) wsWriteLoop() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("xxxxxxxxxxx", err) // 将 interface{} 转型为具体类型。
		}
	}()
	for {
		select {
		// 取一个应答
		case msg := <-wsConn.outChan:
			// 写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				log.Println("发送消息给客户端发生错误", err.Error())
				// 切断服务
				wsConn.close()
				return
			}
		case <-wsConn.closeChan:
			fmt.Println("closing the socket4")
			// 获取到关闭通知
			return
		case <-wsConn.sccticker.C: //定时器到时间了 如果username还是空 证明没登陆 关闭，如果不为空 发送心跳
			if "" == wsConn.username {
				wsConn.close()
				log.Println("user has not login yet")
				return
			}
			if err := wsConn.wsSocket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 写入消息到队列中
func (wsConn *wsConnection) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &wsMessage{messageType, data}:
	case <-wsConn.closeChan:
		fmt.Println("closing the socket1")
		return errors.New("连接已经关闭")
	}
	return nil
}

// 读取消息队列中的消息
func (wsConn *wsConnection) wsRead() (*wsMessage, error) {
	select {

	case msg := <-wsConn.inChan:
		// 获取到消息队列中的消息
		return msg, nil
	case <-wsConn.closeChan:
		fmt.Println("closing the socket")

	}
	return nil, errors.New("连接已经关闭")
}

// 关闭连接
func (wsConn *wsConnection) close() {
	log.Println("关闭连接被调用了")
	wsConn.wsSocket.Close()
	wsConn.mutex.Lock()
	wsConn.sccticker.Stop()

	defer wsConn.mutex.Unlock()
	if wsConn.isClosed == false {
		wsConn.isClosed = true
		for k, v := range wsConn.sccroom {
			v.singlelock.Lock()
			leavemsg := make(map[string]string)
			leavemsg["username"] = "sccsystem"
			leavemsg["roomname"] = k
			leavemsg["msgtype"] = "sendmsg"
			leavemsg["msg"] = wsConn.username + " leave the liveroom"
			v.sendmsg(wsConn.username, leavemsg)
			delete(v.Conn, wsConn.username)
			v.singlelock.Unlock()
		}
		wsConn.sccroom = make(map[string]*roominfo)
		close(wsConn.closeChan)
		close(wsConn.inChan)
		close(wsConn.outChan)
	}
	delete(liveroominfo.wsConnAll, wsConn.username)
}
func (wsConn *wsConnection) handlelogin(dat map[string]string) (reason string) {
	if wsConn.username == "" {
		if value, ok := dat["username"]; ok {
			liveroominfo.connlock.Lock()
			if conn, connok := liveroominfo.wsConnAll[value]; connok {
				//如果有证明这个哥们已经登陆 了 先T
				conn.close()
			}
			wsConn.username = value
			liveroominfo.wsConnAll[value] = wsConn
			liveroominfo.connlock.Unlock()
		} else {
			return "no username"
		}
	} else {
		return "double login"
	}
	return ""
}
func (wsConn *wsConnection) handlesendmsg(dat map[string]string) (reason string) {
	username, _ := dat["username"]
	if username != wsConn.username {
		return "usename not mathch login"
	}
	if value, ok := dat["roomname"]; ok {
		tmproom := wsConn.sccroom[value]
		tmproom.singlelock.Lock()
		tmproom.sendmsg(wsConn.username, dat)
		tmproom.singlelock.Unlock()
	} else {
		return "no such room"
	}

	return ""
}
func (wsConn *wsConnection) enterroom(dat map[string]string) (reason string) {
	username, _ := dat["username"]
	if username != wsConn.username {
		return "usename not mathch login"
	}
	if value, ok := dat["roomname"]; ok {
		liveroominfo.rommlock.Lock()
		if room, roomok := liveroominfo.roommap[value]; roomok {
			room.singlelock.Lock()
			entermsg := make(map[string]string)
			entermsg["username"] = "sccsystem"
			entermsg["roomname"] = value
			entermsg["msgtype"] = "sendmsg"
			entermsg["msg"] = username + " enter the liveroom"
			room.sendmsg(wsConn.username, entermsg)
			room.Conn[username] = wsConn
			room.singlelock.Unlock()
			wsConn.sccroom[value] = room
		} else {
			liveroominfo.rommlock.Unlock()
			return "no such room"
		}
		liveroominfo.rommlock.Unlock()
		return ""
	}
	return "roomname failed"
}
func (wsConn *wsConnection) leaveroom(dat map[string]string) (reason string) {
	username, _ := dat["username"]
	if username != wsConn.username {
		return "usename not mathch login"
	}
	if value, ok := dat["roomname"]; ok {

		tmproom := wsConn.sccroom[value]
		tmproom.singlelock.Lock()
		delete(tmproom.Conn, username)
		leavemsg := make(map[string]string)
		leavemsg["username"] = "sccsystem"
		leavemsg["roomname"] = value
		leavemsg["msgtype"] = "sendmsg"
		leavemsg["msg"] = username + " leave the liveroom"
		tmproom.sendmsg(wsConn.username, leavemsg)
		tmproom.singlelock.Unlock()
		delete(wsConn.sccroom, value)
	} else {
		return "no such room"
	}
	return ""
}
func (wsConn *wsConnection) maptobyte(reply map[string]string) (data2 []byte) {
	dataType, _ := json.Marshal(reply)
	dataString := string(dataType)
	data2 = []byte(dataString)
	return data2
}

// 启动程序
func StartWebsocket() {
	zlmeidahook.Setcallback(zlreportroom)
	liveroominfo.roommap = make(map[string]*roominfo)
	liveroominfo.wsConnAll = make(map[string]*wsConnection)
	http.HandleFunc("/", wsHandler)
	http.ListenAndServe(":20080", nil)

}
