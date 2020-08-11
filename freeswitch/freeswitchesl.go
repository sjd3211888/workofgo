package main

import (
	"fmt"
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
	fslock   sync.Mutex
	userinfo map[string]*Freeswitchuser
}

var sccfsinfo Fseslinfo

func handlemsg(msg map[string]string) {
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
								//push_userinfo(tmpuid,itfreuser->second.iuserstatus);推送到rabbitmq
							}
						} else {
							//不存在就添加
							var tmpuserinfo Freeswitchuser
							tmpuserinfo.iuserstatus = 0x00001
							if host, bexist3 := msg["from-host"]; bexist3 {
								tmpuserinfo.userip = host
							}
							//fmt.Println("host ip is ", tmpuserinfo)
							sccfsinfo.userinfo[existuser] = &tmpuserinfo

						}
						sccfsinfo.fslock.Unlock()
					} else {
						fmt.Println("from-user not exist do nothing")
					}

				} else if subevent == "sofia::unregister" {
					if existuser, bexist := msg["from-user"]; bexist {
						sccfsinfo.fslock.Lock()

						delete(sccfsinfo.userinfo, existuser)
						sccfsinfo.fslock.Unlock()
					} else {
						fmt.Println("from-user not exist do nothing")
					}
				} else if subevent == "sofia::sip_user_state" {
					if existuser, bexist := msg["from-user"]; bexist {
						sccfsinfo.fslock.Lock()

						delete(sccfsinfo.userinfo, existuser)
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
						existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					}
					if !bassin {
						fmt.Println("ss", tmpcallinfo)
						existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					}
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
						existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					}
					if !bassin {
						existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					}
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
				}
				if !bassign {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
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
						existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					}
					if !bassign {
						existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					}
				} else {
					fmt.Println("callee  not exist", calleename)
				}
			} else {
				fmt.Println("caller  not exist", callername)
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
				for _, v := range existuserinfo.usercallinfo {
					fmt.Println("zzzzzzzzzzzzzzzzzz", v.uuid)
				}
				if 1 == len(existuserinfo.usercallinfo) && index == 0 {
					existuserinfo.usercallinfo = existuserinfo.usercallinfo[:0]

				} else {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo[:index], existuserinfo.usercallinfo[index+1:]...)
				}

				fmt.Println("caller is ", existuserinfo.usercallinfo, callername, index, len(existuserinfo.usercallinfo))
				if len(existuserinfo.usercallinfo) >= 0 {
					existuserinfo.iuserstatus = existuserinfo.iuserstatus & 0x0001
				}
			}

			if existuserinfo, bexist1 := sccfsinfo.userinfo[calleename]; bexist1 {
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
				if 1 == len(existuserinfo.usercallinfo) && index == 0 {
					existuserinfo.usercallinfo = existuserinfo.usercallinfo[:0]
				} else {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo[:index], existuserinfo.usercallinfo[index+1:]...)
				}
				fmt.Println("callee is ", existuserinfo.usercallinfo)
				if len(existuserinfo.usercallinfo) >= 0 {
					existuserinfo.iuserstatus = existuserinfo.iuserstatus & 0x0001
				}
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
					fmt.Println("sssssssssssssssssssss", tmpcallinfo)
				}
				if !bassign {
					existuserinfo.usercallinfo = append(existuserinfo.usercallinfo, &tmpcallinfo)
					fmt.Println(tmpcallinfo)
				}
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
func prinfstatus() {
	tiker := time.NewTicker(time.Second * 24)
	for i := 0; ; i++ {

		fmt.Println(<-tiker.C)
		sccfsinfo.fslock.Lock()

		for k, v := range sccfsinfo.userinfo {
			fmt.Println(k, v.userip)
			if 0 != len(v.usercallinfo) {
				fmt.Println(v.usercallinfo[0].callernum, v.usercallinfo[0].applicationdata, v.usercallinfo[0].calleenum, v.usercallinfo[0].callstarttime, v.usercallinfo[0].application)
			}
		}
		sccfsinfo.fslock.Unlock()
	}
}
func main() {

	// Boost it as much as it can go ...
	// We don't need this since Go 1.5
	// runtime.GOMAXPROCS(runtime.NumCPU())

	client, err := NewClient("192.168.1.200", 8021, "ClueCon", 10)
	sccfsinfo.userinfo = make(map[string]*Freeswitchuser, 0)
	if err != nil {
		Error("Error while creating new client: %s", err)
		return
	}

	// Apparently all is good... Let us now handle connection :)
	// We don't want this to be inside of new connection as who knows where it my lead us.
	// Remember that this is crutial part in handling incoming messages. This is a must!

	go client.Handle()

	client.Send("events json PLAYBACK_START CHANNEL_ANSWER CHANNEL_ORIGINATE CHANNEL_HANGUP DTMF CUSTOM conference::maintenance sofia::register sofia::unregister sofia::sip_user_state")

	//client.BgApi(fmt.Sprintf("originate %s %s", "sofia/internal/1001@127.0.0.1", "&socket(192.168.1.2:8084 async full)"))
	go prinfstatus()
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
		handlemsg(msg.Headers)
	}
}
