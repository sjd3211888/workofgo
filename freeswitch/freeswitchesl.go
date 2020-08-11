package main

import (
	"strings"

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
	callstarttime   uint
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
	usercallinfo []Callinfo
}

func main() {

	// Boost it as much as it can go ...
	// We don't need this since Go 1.5
	// runtime.GOMAXPROCS(runtime.NumCPU())

	client, err := NewClient("192.168.1.200", 8021, "ClueCon", 10)

	if err != nil {
		Error("Error while creating new client: %s", err)
		return
	}

	// Apparently all is good... Let us now handle connection :)
	// We don't want this to be inside of new connection as who knows where it my lead us.
	// Remember that this is crutial part in handling incoming messages. This is a must!

	go client.Handle()

	client.Send("events json CHANNEL_ANSWER CHANNEL_ORIGINATE CHANNEL_HANGUP DTMF CUSTOM conference::maintenance sofia::register sofia::unregister sofia::sip_user_state")

	//client.BgApi(fmt.Sprintf("originate %s %s", "sofia/internal/1001@127.0.0.1", "&socket(192.168.1.2:8084 async full)"))

	for {
		msg, err := client.ReadMessage()

		if err != nil {

			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				Error("Error while reading Freeswitch message: %s", err)
			}

			break
		}

		Debug("Got new message: %s", msg)
	}
}
