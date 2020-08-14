package fstoesl

import (
	"encoding/json"
	"fmt"
	RabbitMQ "golearn/Rabbitmq"
	"strings"
	"sync"
	"time"

	. "github.com/0x19/goesl"
)

type Callinfo struct {
	callernum       string
	calleenum       string
	brocastuuid     string
	uuid            string
	buuid           string //被叫腿的uuid
	application     string
	applicationdata string
	callstarttime   int64
}
type Freeswitchuser struct {
	usertype int
	//用户的固定组ID//
	gid int
	//用户的注册的IP 貌似没啥卵用 现在换成用户对讲UDP端口所用的IP
	userip string
	//fs注册的ip 用来发短信
	hostip string
	/********************
	  0x00000
	  从右往左 sip在线离线 振铃中  广播与否 通话中 对讲在线离线
	 ********************/
	iuserstatus  int
	sccid        string
	usercallinfo []*Callinfo
}
type Fseslinfo struct {
	fslock         sync.Mutex
	userinfo       map[string]*Freeswitchuser
	fsclient       *Client
	rabbitmqpusher *RabbitMQ.RabbitMQ
}

//   这个函数在在线程里有锁的清空下调用 不用特意加锁
func getuserinfo(userinfo *Freeswitchuser) (info string) {
	tmpinfo := make(map[string]interface{}, 0)
	tmpinfo["uid"] = userinfo.sccid
	tmpinfo["status"] = userinfo.iuserstatus
	len := len(userinfo.usercallinfo)
	callinfo := make([]map[string]interface{}, 0)
	for i := 0; i < len; i++ {
		callinfosub := make(map[string]interface{})
		callinfosub["callernum"] = userinfo.usercallinfo[i].callernum
		callinfosub["calleenum"] = userinfo.usercallinfo[i].calleenum
		callinfosub["application"] = userinfo.usercallinfo[i].application
		callinfosub["applicationdata"] = userinfo.usercallinfo[i].applicationdata
		callinfosub["buuid"] = userinfo.usercallinfo[i].buuid
		callinfosub["uuid"] = userinfo.usercallinfo[i].uuid
		callinfosub["callstarttime"] = userinfo.usercallinfo[i].callstarttime
		callinfo = append(callinfo, callinfosub)
	}
	tmpinfo["callinfo"] = callinfo

	mjson, _ := json.Marshal(tmpinfo)
	mString := string(mjson)
	//fmt.Println("print mString:", mString)
	return mString
}
func (sccfsinfo *Fseslinfo) handlemsg(msg map[string]string) {
	if value, ok := msg["Event-Name"]; ok {
		if "CUSTOM" == value {
			if subevent, ok1 := msg["Event-Subclass"]; ok1 {
				if "sofia::register" == subevent {
					if existuser, bexist := msg["from-user"]; bexist {
						sccfsinfo.fslock.Lock()
						if existuserinfo, bexist1 := sccfsinfo.userinfo[existuser]; bexist1 {
							//存在就更新一下下
							if 0x00000 == (existuserinfo.iuserstatus & 0x00001) {
								existuserinfo.iuserstatus = existuserinfo.iuserstatus | 0x00001

								if host, bexist3 := msg["from-host"]; bexist3 {
									existuserinfo.userip = host
								}
								existuserinfo.sccid = existuser
								tmptopic := fmt.Sprintf("SCC.onlinestatus.%v", existuser)
								sccfsinfo.rabbitmqpusher.PublishTopic(getuserinfo(existuserinfo), "SCCSIPSTATUS", tmptopic)
							}
						} else {
							//不存在就添加
							var tmpuserinfo Freeswitchuser
							tmpuserinfo.iuserstatus = 0x00001
							if host, bexist3 := msg["from-host"]; bexist3 {
								tmpuserinfo.userip = host
							}
							tmpuserinfo.sccid = existuser
							//fmt.Println("host ip is ", tmpuserinfo)
							sccfsinfo.userinfo[existuser] = &tmpuserinfo
							tmptopic := fmt.Sprintf("SCC.onlinestatus.%v", existuser)
							sccfsinfo.rabbitmqpusher.PublishTopic(getuserinfo(&tmpuserinfo), "SCCSIPSTATUS", tmptopic)
						}
						sccfsinfo.fslock.Unlock()
					} else {
						fmt.Println("from-user not exist do nothing")
					}

				} else if subevent == "sofia::unregister" {
					if existuser, bexist := msg["from-user"]; bexist {
						sccfsinfo.fslock.Lock()

						delete(sccfsinfo.userinfo, existuser)
						tmptopic := fmt.Sprintf("SCC.offlinestatus.%v", existuser)
						tmpjson := fmt.Sprintf("{\"callinfo\":[],\"status\":0,\"uid\":\"%v\"}", existuser)
						sccfsinfo.rabbitmqpusher.PublishTopic(tmpjson, "SCCSIPSTATUS", tmptopic)
						sccfsinfo.fslock.Unlock()
					} else {
						fmt.Println("from-user not exist do nothing")
					}
				} else if subevent == "sofia::sip_user_state" {
					if existuser, bexist := msg["from-user"]; bexist {
						sccfsinfo.fslock.Lock()

						delete(sccfsinfo.userinfo, existuser)
						tmptopic := fmt.Sprintf("SCC.offlinestatus.%v", existuser)
						tmpjson := fmt.Sprintf("{\"callinfo\":[],\"status\":0,\"uid\":\"%v\"}", existuser)
						sccfsinfo.rabbitmqpusher.PublishTopic(tmpjson, "SCCSIPSTATUS", tmptopic)
						sccfsinfo.fslock.Unlock()
					} else {
						fmt.Println("from-user not exist do nothing")
					}
				} else {

				}
			} else {
				fmt.Println("Event-Subclass no exist")
			}
		} else if "CHANNEL_ANSWER" == value {
			var callername string
			var calleename string
			callstarttime := time.Now().Unix()
			uuid := msg["variable_call_uuid"]
			buuid := msg["Unique-ID"]
			out := msg["Channel-State"]
			if out == "CS_CONSUME_MEDIA" {
				var tmpcallinfo Callinfo
				callername = msg["Caller-Caller-ID-Number"]
				calleename = msg["Caller-Callee-ID-Number"]
				tmpcallinfo.callernum = callername
				tmpcallinfo.calleenum = calleename
				tmpcallinfo.uuid = uuid
				tmpcallinfo.buuid = buuid
				tmpcallinfo.application = ""
				tmpcallinfo.applicationdata = ""
				tmpcallinfo.callstarttime = callstarttime
				sccfsinfo.fslock.Lock()
				if existuserinfo, bexist1 := sccfsinfo.userinfo[callername]; bexist1 {
					existuserinfo.iuserstatus = existuserinfo.iuserstatus | 0x01000
					bassin := false
					if 0 != len(existuserinfo.usercallinfo) {
						for _, v := range existuserinfo.usercallinfo {
							if v.buuid == buuid || v.uuid == uuid {
								v.callstarttime = callstarttime
								bassin = true
								break
							} //这里玩意uuid都没匹配 还需要把对象加上去，目前感觉这种情况不会出现所以每家
						}
					} else {
						bassin = true
						existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
						fmt.Println("add1")
					}
					if !bassin {
						existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
						fmt.Println("add2")
					}
					tmptopic := fmt.Sprintf("SCC.sipcallinfo.%v", callername)
					sccfsinfo.rabbitmqpusher.PublishTopic(getuserinfo(existuserinfo), "SCCSIPSTATUS", tmptopic)
				} else {
					fmt.Println("caller not exist", callername)
				}
				if existuserinfo, bexist1 := sccfsinfo.userinfo[calleename]; bexist1 {
					existuserinfo.iuserstatus = existuserinfo.iuserstatus | 0x01000
					bassin := false
					fmt.Println("ssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss", len(existuserinfo.usercallinfo), calleename)
					if 0 != len(existuserinfo.usercallinfo) {
						for _, v := range existuserinfo.usercallinfo {
							if v.buuid == buuid || v.uuid == uuid {
								v.callstarttime = callstarttime
								bassin = true
								break
							} //这里玩意uuid都没匹配 还需要把对象加上去，目前感觉这种情况不会出现所以每家
						}
					} else {
						bassin = true
						existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
						fmt.Println("add3")
					}
					if !bassin {
						existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
						fmt.Println("add4")
					}
					tmptopic := fmt.Sprintf("SCC.sipcallinfo.%v", calleename)
					sccfsinfo.rabbitmqpusher.PublishTopic(getuserinfo(existuserinfo), "SCCSIPSTATUS", tmptopic)
				} else {
					fmt.Println("callee  not exist", calleename)
				}
				sccfsinfo.fslock.Unlock()
			}

		} else if "CHANNEL_ORIGINATE" == value {
			var callername string
			var calleename string
			uuid := msg["variable_call_uuid"]
			buuid := msg["Unique-ID"]
			callername = msg["Caller-Caller-ID-Number"]
			calleename = msg["Caller-Callee-ID-Number"]
			var tmpcallinfo Callinfo
			tmpcallinfo.callernum = callername
			tmpcallinfo.calleenum = calleename
			tmpcallinfo.uuid = uuid
			tmpcallinfo.buuid = buuid
			tmpcallinfo.application = ""
			tmpcallinfo.applicationdata = ""
			tmpcallinfo.callstarttime = 0
			sccfsinfo.fslock.Lock()
			if existuserinfo, bexist1 := sccfsinfo.userinfo[callername]; bexist1 {
				existuserinfo.iuserstatus = existuserinfo.iuserstatus | 0x01000
				bassign := false
				if 0 != len(existuserinfo.usercallinfo) {
					for _, v := range existuserinfo.usercallinfo {
						if v.buuid == buuid || v.uuid == uuid {
							v.uuid = uuid
							bassign = true
							break
						} //这里玩意uuid都没匹配 还需要把对象加上去，目前感觉这种情况不会出现所以每家
					}
				} else {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					bassign = true
					fmt.Println("add5")
				}
				if !bassign {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					fmt.Println("add6")
				}
				tmptopic := fmt.Sprintf("SCC.sipcallinfo.%v", callername)
				sccfsinfo.rabbitmqpusher.PublishTopic(getuserinfo(existuserinfo), "SCCSIPSTATUS", tmptopic)
			} else {
				fmt.Println("caller  not exist", callername)
			}
			if existuserinfo, bexist1 := sccfsinfo.userinfo[calleename]; bexist1 {
				existuserinfo.iuserstatus = existuserinfo.iuserstatus | 0x01000
				bassign := false
				if 0 != len(existuserinfo.usercallinfo) {
					for _, v := range existuserinfo.usercallinfo {
						if v.buuid == buuid || v.uuid == uuid {
							v.uuid = uuid
							bassign = true
							break
						} //这里玩意uuid都没匹配 还需要把对象加上去，目前感觉这种情况不会出现所以每家
					}
				} else {
					bassign = true
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					fmt.Println("add7")
				}
				if !bassign {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					fmt.Println("add8")
				}
				tmptopic := fmt.Sprintf("SCC.sipcallinfo.%v", calleename)
				sccfsinfo.rabbitmqpusher.PublishTopic(getuserinfo(existuserinfo), "SCCSIPSTATUS", tmptopic)
			} else {
				fmt.Println("callee  not exist", calleename)
			}

			sccfsinfo.fslock.Unlock()
		} else if "CHANNEL_HANGUP" == value {
			var callername string
			var calleename string
			uuid := msg["variable_call_uuid"]
			buuid := msg["Unique-ID"]
			callername = msg["Caller-Caller-ID-Number"]
			calleename = msg["Caller-Callee-ID-Number"]
			sccfsinfo.fslock.Lock()
			if existuserinfo, bexist1 := sccfsinfo.userinfo[callername]; bexist1 {
				existuserinfo.iuserstatus = existuserinfo.iuserstatus & 0x0001
				var index int
				if 0 != len(existuserinfo.usercallinfo) {
					for k, v := range existuserinfo.usercallinfo {
						if v.buuid == buuid || v.uuid == uuid {
							//delete()
							index = k
							break
						} //这里玩意uuid都没匹配 还需要把对象加上去，目前感觉这种情况不会出现所以每家
					}
				}
				//
				println("bbbbbbbbbbbbb is ", len(existuserinfo.usercallinfo), index)
				if 1 == len(existuserinfo.usercallinfo) && index == 0 {
					existuserinfo.usercallinfo = existuserinfo.usercallinfo[:0]

				} else if 0 == len(existuserinfo.usercallinfo) {
					existuserinfo.usercallinfo = existuserinfo.usercallinfo[:0]
				} else {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo[:index], existuserinfo.usercallinfo[index+1:]...)
				}

				fmt.Println("caller is ", existuserinfo.usercallinfo, callername, index, len(existuserinfo.usercallinfo))
				if len(existuserinfo.usercallinfo) >= 0 {
					existuserinfo.iuserstatus = existuserinfo.iuserstatus & 0x0001
				}
				tmptopic := fmt.Sprintf("SCC.sipcallinfo.%v", callername)
				sccfsinfo.rabbitmqpusher.PublishTopic(getuserinfo(existuserinfo), "SCCSIPSTATUS", tmptopic)
			}

			if existuserinfo, bexist1 := sccfsinfo.userinfo[calleename]; bexist1 {
				existuserinfo.iuserstatus = existuserinfo.iuserstatus & 0x0001
				var index int
				if 0 != len(existuserinfo.usercallinfo) {
					for k, v := range existuserinfo.usercallinfo {
						if v.buuid == buuid || v.uuid == uuid {
							index = k
							break
						} //这里玩意uuid都没匹配 还需要把对象加上去，目前感觉这种情况不会出现所以每家
					}
				}
				println("xxxxxxxxx is ", len(existuserinfo.usercallinfo), index)
				if 1 <= len(existuserinfo.usercallinfo) && index == 0 {
					existuserinfo.usercallinfo = existuserinfo.usercallinfo[:0]
				} else if 0 == len(existuserinfo.usercallinfo) {
					existuserinfo.usercallinfo = existuserinfo.usercallinfo[:0]
				} else {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo[:index], existuserinfo.usercallinfo[index+1:]...)
				}
				fmt.Println("callee is ", existuserinfo.usercallinfo)
				if len(existuserinfo.usercallinfo) >= 0 {
					existuserinfo.iuserstatus = existuserinfo.iuserstatus & 0x0001
				}
				tmptopic := fmt.Sprintf("SCC.sipcallinfo.%v", calleename)
				sccfsinfo.rabbitmqpusher.PublishTopic(getuserinfo(existuserinfo), "SCCSIPSTATUS", tmptopic)
			}
			sccfsinfo.fslock.Unlock()
		} else if "PLAYBACK_START" == value {
			fmt.Println("ddddddddddddddddddddddddddd")
			var callername string
			callstarttime := time.Now().Unix()
			uuid := msg["variable_call_uuid"]
			buuid := msg["Unique-ID"]
			callername = msg["Caller-Caller-ID-Number"]
			appinfo := msg["variable_current_application"]
			appdata := msg["variable_current_application_data"]
			var tmpcallinfo Callinfo
			tmpcallinfo.callernum = "0"
			tmpcallinfo.calleenum = callername
			tmpcallinfo.brocastuuid = uuid
			tmpcallinfo.uuid = ""
			tmpcallinfo.buuid = buuid
			tmpcallinfo.application = appinfo
			tmpcallinfo.applicationdata = appdata
			tmpcallinfo.callstarttime = callstarttime
			sccfsinfo.fslock.Lock()
			if existuserinfo, bexist1 := sccfsinfo.userinfo[callername]; bexist1 {
				existuserinfo.iuserstatus = existuserinfo.iuserstatus | 0x10100
				bassign := false
				fmt.Println("sjdjsdjsjdasjdjsajd", len(existuserinfo.usercallinfo))
				if 0 != len(existuserinfo.usercallinfo) {
					for _, v := range existuserinfo.usercallinfo {
						if v.buuid == buuid || v.uuid == uuid {
							v.application = appinfo
							v.applicationdata = appdata
							v.uuid = uuid
							tmpcallinfo.callstarttime = callstarttime
							bassign = true
							break
						} //这里玩意uuid都没匹配 还需要把对象加上去，目前感觉这种情况不会出现所以每家
					}
				} else {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					fmt.Println("add9")
				}
				if !bassign {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					fmt.Println("add10")
				}
				tmptopic := fmt.Sprintf("SCC.sipcallinfo.%v", callername)
				sccfsinfo.rabbitmqpusher.PublishTopic(getuserinfo(existuserinfo), "SCCSIPSTATUS", tmptopic)
			}
			sccfsinfo.fslock.Unlock()
		} else {
			fmt.Println("unknow Event-Name")
		}
	} else {
		fmt.Println(msg)
		fmt.Println("Event-Name not exist")
	}
}

func (sccfsinfo *Fseslinfo) Hangupuser(uuid string) {
	sendcmd := fmt.Sprintf("bgapi uuid_kill %s\r\n", uuid)
	sccfsinfo.fsclient.BgApi(sendcmd)
}
func (sccfsinfo *Fseslinfo) Servercall(callerl string, callerr string, calleeid string) {
	sendcmd := fmt.Sprintf("bgapi originate user/%v,user/%v %v XML default\r\n", callerl, callerr, calleeid)
	sccfsinfo.fsclient.BgApi(sendcmd)
}
func (sccfsinfo *Fseslinfo) Monitoruser(callerl string, callerr string, calleeid string) {
	sendcmd := fmt.Sprintf("bgapi originate {absolute_codec_string=PCMA}{sip_h_Call-info=<uri>;audiomode=2}user/%s,user/%s %s XML default\r\n", callerl, callerr, calleeid)
	sccfsinfo.fsclient.BgApi(sendcmd)
}
func (sccfsinfo *Fseslinfo) Yellinguser(callerl string, callerr string, calleeid string) {
	sendcmd := fmt.Sprintf("bgapi originate {absolute_codec_string=PCMA}user/%s,user/%s &bridge({sip_h_Call-info=<uri>;audiomode=1}user/%s)\r\n", callerl, callerr, calleeid)
	sccfsinfo.fsclient.BgApi(sendcmd)
}
func (sccfsinfo *Fseslinfo) prinfstatus() {
	tiker := time.NewTicker(time.Second * 24)
	for i := 0; ; i++ {

		fmt.Println(<-tiker.C)
		sccfsinfo.fslock.Lock()

		for k, v := range sccfsinfo.userinfo {
			fmt.Println(k, v.userip)
			for sjd := 0; sjd < len(v.usercallinfo); sjd++ {
				fmt.Println(v.usercallinfo[sjd].callernum, v.usercallinfo[sjd].applicationdata, v.usercallinfo[sjd].calleenum, v.usercallinfo[sjd].callstarttime, v.usercallinfo[sjd].application)
			}
		}
		sccfsinfo.fslock.Unlock()
	}
}
func (sccfsinfo *Fseslinfo) Fseslclientrun() {

	// Boost it as much as it can go ...
	// We don't need this since Go 1.5
	// runtime.GOMAXPROCS(runtime.NumCPU())

	client, err := NewClient("192.168.1.200", 8021, "ClueCon", 10)
	sccfsinfo.userinfo = make(map[string]*Freeswitchuser, 0)
	if err != nil {
		Error("Error while creating new client: %s", err)
		return
	}
	sccfsinfo.fsclient = client
	// Apparently all is good... Let us now handle connection :)
	// We don't want this to be inside of new connection as who knows where it my lead us.
	// Remember that this is crutial part in handling incoming messages. This is a must!
	sccfsinfo.rabbitmqpusher = RabbitMQ.NewRabbitMQTopic("amqp://sjd:sjd@*.*.*.*:5672/admin")
	go client.Handle()

	client.Send("events json PLAYBACK_START CHANNEL_ANSWER CHANNEL_ORIGINATE CHANNEL_HANGUP DTMF CUSTOM conference::maintenance sofia::register sofia::unregister sofia::sip_user_state")

	//client.BgApi(fmt.Sprintf("originate %s %s", "sofia/internal/1001@127.0.0.1", "&socket(192.168.1.2:8084 async full)"))
	go sccfsinfo.prinfstatus()
	for {
		msg, err := client.ReadMessage()

		if err != nil {

			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				Error("Error while reading Freeswitch message: %s", err)
			}

			break
		}

		//Debug("Got new message: %s", msg)
		sccfsinfo.handlemsg(msg.Headers)
	}
}
