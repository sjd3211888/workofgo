package fsmysql

import (
	"fmt"
	sccsql "golearn/gomysql"
	"strconv"

	"github.com/BurntSushi/toml"
)

var tmpsql sccsql.Mysqlconnectpool

func init() {
	var conf map[string]map[string]string
	if _, err := toml.DecodeFile("./sccconfig.toml", &conf); err != nil {
		// handle error
	}
	Host := conf["sccfs"]["Host"]
	Username := conf["sccfs"]["Username"]
	Password := conf["sccfs"]["Password"]
	Dbname := conf["sccfs"]["Dbname"]
	Port := conf["sccfs"]["Port"]
	iport, _ := strconv.Atoi(Port)
	//fmt.Println("Hostxxxxxxxxxxxxxx", Host)
	go func(Host string, Username string, Password string, Dbname string, iport int) {
		tmpsql.Initmysql(Host, Username, Password, Dbname, iport)
	}(Host, Username, Password, Dbname, iport)

}

//
func Getuseronline() []map[string]string {
	sqlcmd := fmt.Sprintf("Select sip_user,status,network_ip,sip_host from sip_registrations")

	return tmpsql.SelectData(sqlcmd)
}
