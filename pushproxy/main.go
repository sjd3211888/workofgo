package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	sccsql "golearn/gomysql"
	sccredis "golearn/goredis"
	. "golearn/sccprotobuf"
	"net"
	"strconv"

	. "golearn/pushproxy/pushporxygateway"

	"github.com/BurntSushi/toml"
	"github.com/golang/protobuf/proto"
)

type coreinfo struct {
	tmpsql   sccsql.Mysqlconnectpool
	tmpredis sccredis.Redisconnectpool
}
type SCChead struct {
	cmdid    int16
	seq      int16
	bodyleng int32
}

var buffer bytes.Buffer //Buffer是一个实现了读写方法的可变大小的字节缓冲
var sccinfo coreinfo

func init() {
	var conf map[string]map[string]string
	if _, err := toml.DecodeFile("./sccconfig.toml", &conf); err != nil {
		// handle error
	}
	Host := conf["pusherproxy"]["Host"]
	Username := conf["pusherproxy"]["Username"]
	Password := conf["pusherproxy"]["Password"]
	Dbname := conf["pusherproxy"]["Dbname"]
	Port := conf["pusherproxy"]["Port"]
	iport, _ := strconv.Atoi(Port)
	Serhost := conf["pusherproxy"]["ToHttphost"]
	Redisip := conf["coreservice"]["Redisip"]
	go func(Host string, Username string, Password string, Dbname string, Serhost string, Redisip string, iport int) {
		sccinfo.tmpsql.Initmysql(Host, Username, Password, Dbname, iport)
		sccinfo.tmpredis.Redisip = (Redisip)
		sccinfo.tmpredis.ConnectRedis()
		/*gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		setrouter(r)
		if err := r.Run(Serhost); err != nil {
			fmt.Println("startup service failed, err:\n", err)
		}*/
	}(Host, Username, Password, Dbname, Serhost, Redisip, iport)

}

func main() {
	fmt.Println("Starting the server ...")
	// 创建 listener
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":4456")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	//server, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接
	//listener.SetReadBuffer(1024 * 1024)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		conn.SetReadBuffer(int(1024 * 1024 * 64))
		go doServerStuff(conn)
	}
}
func doServerStuff(conn net.Conn) {
	for {
		buf := make([]byte, 4096)
		readlen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return //终止程序
		}
		fmt.Println("len is ", len(buf), "read len is ", readlen)
		buffer.Write(buf[:readlen]) //先添加包
		//b3 := buffer.Bytes()
		for {
			fmt.Println("buf1len is ", buffer.Len())
			if buffer.Len() <= 10 {
				break
			}
			var sss int32
			buf1 := buffer.Next(8)
			buflen := bytes.NewReader(buf1[4:8])
			err = binary.Read(buflen, binary.BigEndian, &sss)
			fmt.Println("ssss354", sss, buffer.Len())
			if buffer.Len() < int(sss) {
				//把读的还原
				var buffer1 bytes.Buffer
				buffer1.Write(buf1)
				var tmpbuf []byte
				tmpbuf = buffer.Next(buffer.Len())
				fmt.Println("bbbbbbbbbbbbbb", len(tmpbuf))
				buffer1.Write(tmpbuf)
				buffer.Reset()
				buffer = buffer1
				fmt.Println("zzzzzzzzzzzzzzzz", buffer1.Len())
				break
			} else {
				tmpmsg := make([]byte, int(sss))
				_, _ = buffer.Read(tmpmsg)
				go func(tmpmsg []byte) {
					data := &SccIMPush{}
					proto.Unmarshal(tmpmsg, data)
					fmt.Println("反序列化之后的信息为：", data)
					pushmsg := make(map[string]string)
					if 0 == data.GetSendtype() {
						//个人短信

						sqlcmd := fmt.Sprintf("Select s_displayname from scc_user where s_user = '%v'", data.GetFromsccid())
						sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
						pushmsg["title"] = sqlresult[0]["s_displayname"]

						switch data.GetImtype() {
						case 0:
							pushmsg["info"] = data.GetIminfo()
						case 1:
							pushmsg["info"] = "[图片]"
						case 2:
							pushmsg["info"] = "[视频]"
						case 3:
							pushmsg["info"] = "[老音频]"
						case 4:
							pushmsg["info"] = "[老文件]"
						case 5:
							pushmsg["info"] = "[语音]"
						case 6:
							pushmsg["info"] = "[文件]"
						case 7:
							pushmsg["info"] = "[位置]"
						case 8:
							pushmsg["info"] = "[链接]"
						case 9:
							pushmsg["info"] = "[系统消息]"
						case 10:
							pushmsg["info"] = "[其他]"
						case 11:
							pushmsg["info"] = "[直播消息]"
						case 12:
							pushmsg["info"] = "[报警消息]"
						case 13:
							pushmsg["info"] = "[必达消息]"
						case 14:
							pushmsg["info"] = "[调度台呼叫]"
						case 15:
							pushmsg["info"] = "[上报gps]"
						case 16:
							pushmsg["info"] = "[调度台结束直播]"
						case 17:
							pushmsg["info"] = "[拒绝直播]"
						case 18:
							pushmsg["info"] = "[音频通话]"
						case 19:
							pushmsg["info"] = "[视频通话]"
						case 20:
							pushmsg["info"] = "[创建组]"
						case 21:
							pushmsg["info"] = "[修改组]"
						case 22:
							pushmsg["info"] = "[删除组]"
						case 23:
							pushmsg["info"] = "[转发直播]"
						case 24:
							pushmsg["info"] = "[挂断音频通话]"
						case 25:
							pushmsg["info"] = "[挂断视频通话]"
						case 26:
							pushmsg["info"] += "[评价采集]"
						case 27:
							pushmsg["info"] += "[确认必达]"
						case 28:
							pushmsg["info"] += "[挂断视频通话]"
						case 29:
							pushmsg["info"] += "[必达]"
						case 30:
							pushmsg["info"] += "[审批通知]"
						}
						tmptype := "unknow"
						token := ""
						newmsg, _ := sccinfo.tmpredis.SccredisBGetAll(strconv.Itoa(int(data.GetTosccid())))
						tmplen := len(newmsg)
						for i := 0; i < tmplen; i = i + 2 {
							if string(newmsg[i][:]) == "type" && i+1 < tmplen {
								tmptype = string(newmsg[i+1][:])
							}
							if string(newmsg[i][:]) == "token" && i+1 < tmplen {
								token = string(newmsg[i+1][:])
							}
						}
						pushmsg["token"] = token
						if tmptype == "HUAWEI" || tmptype == "HONOR" {

							mjson, _ := json.Marshal(pushmsg)
							Push2huawei(mjson)
							//走华为推送
						} else if tmptype == "XIAOMI" || tmptype == "REDMI" {
							//走小米推送
							mjson, _ := json.Marshal(pushmsg)
							Push2xiaomi(mjson)
						} else if tmptype == "OPPO" {
							mjson, _ := json.Marshal(pushmsg)
							Push2oppo(mjson)
							//oppo推送
						} else if tmptype == "VIVO" {
							//vivo 推送
							mjson, _ := json.Marshal(pushmsg)
							Push2vivo(mjson)
						} else {

						}
					} else {
						sqlcmd := fmt.Sprintf("Select s_groupname from scc_group where s_groupid = '%v'", data.GetTosccid())
						sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
						pushmsg["title"] = sqlresult[0]["s_groupname"]

						sqlcmd = fmt.Sprintf("Select s_displayname from scc_user where s_user = '%v'", data.GetFromsccid())
						sqlresult = sccinfo.tmpsql.SelectData(sqlcmd)
						pushmsg["info"] = sqlresult[0]["s_displayname"] + ":"
						switch data.GetImtype() {
						case 0:
							pushmsg["info"] += data.GetIminfo()
						case 1:
							pushmsg["info"] += "[图片]"
						case 2:
							pushmsg["info"] += "[视频]"
						case 3:
							pushmsg["info"] += "[老音频]"
						case 4:
							pushmsg["info"] += "[老文件]"
						case 5:
							pushmsg["info"] += "[语音]"
						case 6:
							pushmsg["info"] += "[文件]"
						case 7:
							pushmsg["info"] += "[位置]"
						case 8:
							pushmsg["info"] += "[链接]"
						case 9:
							pushmsg["info"] += "[系统消息]"
						case 10:
							pushmsg["info"] += "[其他]"
						case 11:
							pushmsg["info"] += "[直播消息]"
						case 12:
							pushmsg["info"] += "[报警消息]"
						case 13:
							pushmsg["info"] += "[必达消息]"
						case 14:
							pushmsg["info"] += "[调度台呼叫]"
						case 15:
							pushmsg["info"] += "[上报gps]"
						case 16:
							pushmsg["info"] += "[调度台结束直播]"
						case 17:
							pushmsg["info"] += "[拒绝直播]"
						case 18:
							pushmsg["info"] += "[音频通话]"
						case 19:
							pushmsg["info"] += "[视频通话]"
						case 20:
							pushmsg["info"] += "[创建组]"
						case 21:
							pushmsg["info"] += "[修改组]"
						case 22:
							pushmsg["info"] += "[删除组]"
						case 23:
							pushmsg["info"] += "[转发直播]"
						case 24:
							pushmsg["info"] += "[挂断音频通话]"
						case 25:
							pushmsg["info"] += "[挂断视频通话]"
						case 26:
							pushmsg["info"] += "[评价采集]"
						case 27:
							pushmsg["info"] += "[确认必达]"
						case 28:
							pushmsg["info"] += "[挂断视频通话]"
						case 29:
							pushmsg["info"] += "[必达]"
						case 30:
							pushmsg["info"] += "[审批通知]"
						}
						tmptype := "unknow"
						token := ""
						newmsg, _ := sccinfo.tmpredis.SccredisBGetAll(strconv.Itoa(int(data.GetSccidrevived())))
						tmplen := len(newmsg)
						for i := 0; i < tmplen; i = i + 2 {
							if string(newmsg[i][:]) == "type" && i+1 < tmplen {
								tmptype = string(newmsg[i+1][:])
							}
							if string(newmsg[i][:]) == "token" && i+1 < tmplen {
								token = string(newmsg[i+1][:])
							}
						}
						pushmsg["token"] = token
						if tmptype == "HUAWEI" || tmptype == "HONOR" {
							mjson, _ := json.Marshal(pushmsg)
							Push2huawei(mjson)
							//走华为推送
						} else if tmptype == "XIAOMI" || tmptype == "REDMI" {
							//走小米推送
							mjson, _ := json.Marshal(pushmsg)
							Push2xiaomi(mjson)
						} else if tmptype == "OPPO" {
							//oppo推送
							mjson, _ := json.Marshal(pushmsg)
							Push2oppo(mjson)
						} else if tmptype == "VIVO" {
							//vivo 推送
							mjson, _ := json.Marshal(pushmsg)
							Push2vivo(mjson)
						} else {

						}
					}
					fmt.Println(pushmsg)
				}(tmpmsg)

			}
		}
	}
}
