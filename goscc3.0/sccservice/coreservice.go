package coreservice

import (
	"fmt"
	sccsql "golearn/gomysql"
	sccredis "golearn/goredis"
	"net/http"
	"strconv"
	"time"

	. "golearn/sccprotobuf"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
)

type coreinfo struct {
	tmpsql   sccsql.Mysqlconnectpool
	tmpredis sccredis.Redisconnectpool
}

var sccinfo coreinfo

func init() {
	var conf map[string]map[string]string
	if _, err := toml.DecodeFile("./sccconfig.toml", &conf); err != nil {
		// handle error
	}
	Host := conf["coreservice"]["Host"]
	Username := conf["coreservice"]["Username"]
	Password := conf["coreservice"]["Password"]
	Dbname := conf["coreservice"]["Dbname"]
	Port := conf["coreservice"]["Port"]
	iport, _ := strconv.Atoi(Port)
	Serhost := conf["coreservice"]["Httpserverhost"]
	Redisip := conf["coreservice"]["Redisip"]
	//fmt.Println("Hostxxxxxxxxxxxxxx", Host)
	go func(Host string, Username string, Password string, Dbname string, Serhost string, Redisip string, iport int) {
		sccinfo.tmpsql.Initmysql(Host, Username, Password, Dbname, iport)
		sccinfo.tmpredis.Redisip = (Redisip)
		sccinfo.tmpredis.ConnectRedis()
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		setrouter(r)
		if err := r.Run(Serhost); err != nil {
			fmt.Println("startup service failed, err:\n", err)
		}
	}(Host, Username, Password, Dbname, Serhost, Redisip, iport)

}
func querysccdepartment(c *gin.Context) {
	type sccdeparment struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Departmentid string `json:"departmentid" binding:"required"`
	}
	var json sccdeparment
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlcmd := fmt.Sprintf("Select s_departmentname,s_departmentid,s_grade,s_path,s_createtime,s_updatetime  from scc_department  where s_path like '%%/%v/%%'", json.Departmentid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})

}
func querydepartmentuser(c *gin.Context) {

	onlydispatcher := 0

	type sccdeparmentuser struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Departmentid   string `json:"departmentid" binding:"required"`
		Onlydispatcher string `json:"onlydispatcher" binding:"required"`
	}
	var json sccdeparmentuser
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if "yes" == json.Onlydispatcher {
		onlydispatcher = 1
	} else {
		onlydispatcher = 0
	}
	sqlcmd := fmt.Sprintf("Select s_user,s_grade,s_usertype,s_createtime,s_updatetime,s_alias,s_displayname from scc_user  where s_departmentid = '%v' and s_usertype>='%v'", json.Departmentid, onlydispatcher)

	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})

}
func querygroup(c *gin.Context) {
	type groupinfo struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Sccid string `json:"sccid" binding:"required"`
	}
	var json groupinfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd := fmt.Sprintf("Select s_groupid from scc_groupuser where s_user = '%v'", json.Sccid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})
}
func queryuser(c *gin.Context) {

	type userinfo struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Sccid string `json:"sccid" binding:"required"`
	}
	var json userinfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd := fmt.Sprintf("Select s_grade,s_usertype,s_createtime,s_updatetime,s_alias,s_displayname from scc_user  where s_user = '%v'", json.Sccid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})
}
func convemapstringtointerface(data map[string]string) map[string]interface{} {
	interfacedata := make(map[string]interface{})
	for k, v := range data {
		interfacedata[k] = v
	}
	return interfacedata
}
func querygroupuser(c *gin.Context) {

	type groupuserinfo struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Groupid string `json:"groupid" binding:"required"`
	}
	var json groupuserinfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlcmd := fmt.Sprintf("Select s_groupname,s_grade,s_grouptype,s_createtime,s_updatetime,s_creater from scc_group where s_groupid = %v", json.Groupid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
	sqlcmd1 := fmt.Sprintf("Select scc_groupuser.s_user,scc_groupuser.s_displayname,scc_groupuser.s_grade,scc_groupuser.s_usertype,scc_groupuser.s_createtime,scc_groupuser.s_updatetime,scc_user.s_alias from scc_groupuser inner join scc_user on scc_groupuser.s_user =scc_user.s_user where scc_groupuser.s_groupid= %v", json.Groupid)
	sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
	tmggroupid := fmt.Sprintf("group_%s", json.Groupid)
	newmsg, _ := sccinfo.tmpredis.SccredisBHget(tmggroupid, "newestmsg")
	data := &SccIMPush{}
	proto.Unmarshal(newmsg, data)
	//("反序列化之后的信息为：", data)

	//	sqlresult = append(sqlresult, tmpmsg)
	lastmsg := make(map[string]interface{})
	lastmsg["messageid"] = data.GetMessageid()
	lastmsg["fromsccid"] = data.GetFromsccid()
	lastmsg["tosccid"] = data.GetTosccid()
	lastmsg["sendtype"] = data.GetSendtype()
	lastmsg["imtype"] = data.GetImtype()
	lastmsg["iminfo"] = data.GetIminfo()
	lastmsg["filetpath"] = data.GetFiletpath()
	lastmsg["createtime"] = data.GetCreatetime()
	lastmsg["groupmessageid"] = data.GetGroupmessageid()

	sqlresultinterface := convemapstringtointerface(sqlresult[0]) //把查询的结果转换为[string]interface{}
	sqlresultinterface["lastmsg"] = lastmsg
	for i := range sqlresult1 {

		//(sqlresult1[i]["s_user"])
		userstatus, _ := sccinfo.tmpredis.SccredisHget(tmggroupid, sqlresult1[i]["s_user"])
		//statusmap := map[string]string{"status": userstatus}
		if "" == userstatus {
			sqlresult1[i]["status"] = "-1"
		} else {
			sqlresult1[i]["status"] = userstatus
		}

	}

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"groupinfo": sqlresultinterface, "userinfo": sqlresult1}})

}
func queryofflinemsg(c *gin.Context) {
	type personofflinemsg struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Sccid string `json:"sccid" binding:"required"`
	}
	var json personofflinemsg
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scckey := fmt.Sprintf("offlinemsg_%s", json.Sccid)
	myredisresult, _ := sccinfo.tmpredis.SccredisGetAll(scckey)
	sccinfo.tmpredis.SccredisDel(scckey)

	Offlinemsg := make([]map[string]interface{}, 0)
	for _, v := range myredisresult {
		data := &SccIMPush{}
		b := []byte(v)
		proto.Unmarshal(b, data)
		//fmt.Println(data)
		onemsg := make(map[string]interface{})
		onemsg["messageid"] = data.GetMessageid()
		onemsg["fromsccid"] = data.GetFromsccid()
		onemsg["tosccid"] = data.GetTosccid()
		onemsg["sendtype"] = data.GetSendtype()
		onemsg["imtype"] = data.GetImtype()
		onemsg["iminfo"] = data.GetIminfo()
		onemsg["filetpath"] = data.GetFiletpath()
		onemsg["createtime"] = data.GetCreatetime()
		Offlinemsg = append(Offlinemsg, onemsg)
	}

	//sqlresultinterface := convemapstringtointerface(sqlresult[0]) //把查询的结果转换为[string]interface{}
	//sqlresultinterface["lastmsg"] = lastmsg

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": Offlinemsg})
}
func queryRecnetSession(c *gin.Context) {
	type personofflinemsg struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Sccid string `json:"sccid" binding:"required"`
	}
	var json personofflinemsg
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scckey := fmt.Sprintf("session_%s", json.Sccid)
	myredisresult, _ := sccinfo.tmpredis.SccredisGetAll(scckey)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": myredisresult})
}
func reportgps(c *gin.Context) {
	type gps struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Sccid       string `json:"sccid" binding:"required"`
		Longitude   string `json:"longitude" binding:"required"`
		Latitude    string `json:"latitude" binding:"required"`
		Gps         string `json:"gps" binding:"required"`
		Speed       int    `json:"speed"`
		Angle       string `json:"angle"`
		Description string `json:"description"`
	}

	var json gps
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//fmt.Println(json.Sccid)
	sqlcmd := fmt.Sprintf("insert into gpsinfo(sccid,longitude,latitude,reporttime,gps,angle,speed,description) values('%v','%v','%v','%v','%v','%v','%v','%v');", json.Sccid, json.Longitude, json.Latitude, time.Now().Unix(), json.Gps, json.Angle, json.Speed, json.Description)
	//fmt.Println(sqlcmd)
	sccinfo.tmpsql.Execsqlcmd(sqlcmd, true)
	sccinfo.tmpredis.SccredisAddgps(json.Longitude, json.Latitude, json.Sccid)
	c.JSON(http.StatusOK, gin.H{"status": "200"})
}
func querygps(c *gin.Context) {
	type querygps struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Sccid           string `json:"sccid" binding:"required"`
		Starttime       int    `json:"starttime" binding:"required"`
		Endtime         int    `json:"endtime" binding:"required"`
		Pagenum         int    `json:"pagenum" binding:"required"`
		Needdescription string `json:"needdescription"`
	}

	var json querygps
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd1 := fmt.Sprintf("select count(*) from gpsinfo where sccid= '%v' and reporttime > '%v' and reporttime < '%v'", json.Sccid, json.Starttime, json.Endtime)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	if _, ok := sqlresult[0]["count(*)"]; ok {
		count, _ := strconv.Atoi(sqlresult[0]["count(*)"])
		//(sqlresult) //存在
		if 0 != count {
			if 0 == json.Pagenum {
				pagenum = 1 //不存在第0页
			}
			if count < numberperpage {
				if 1 == pagenum {
					fromcount = 0
				} else {
					fromcount = numberperpage
				}

			} else {
				fromcount = count - (numberperpage)*pagenum
				if fromcount < 0 {
					//int tmpnumberperpage = numberperpage;
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}
			if 0 == json.Starttime && 0 == json.Endtime {
				sqlcmd1 := fmt.Sprintf("Select longitude,latitude,reporttime,gps,angle,speed,description from gpsinfo where uid= '%v' order by id DESC limit 1", json.Sccid)
				sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})
				return
			}
			if "yes" == json.Needdescription {
				sqlcmd1 := fmt.Sprintf("Select longitude,latitude,reporttime,gps,angle,speed,description from gpsinfo where sccid = '%v' and reporttime > '%v' and reporttime < '%v' order by id  limit %v,%v", json.Sccid, json.Starttime, json.Endtime, fromcount, numberperpage)
				sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"pagenum": pagenum, "total": count, "gpsinfo": sqlresult}})
				return
			} else {
				sqlcmd1 := fmt.Sprintf("Select longitude,latitude,reporttime,gps,angle,speed from gpsinfo where sccid = '%v' and reporttime > '%v' and reporttime < '%v' order by id  limit %v,%v", json.Sccid, json.Starttime, json.Endtime, fromcount, numberperpage)
				sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"pagenum": pagenum, "totoal": count, "gpsinfo": sqlresult}})
				return
			}

		} else {
			c.JSON(http.StatusOK, gin.H{"result": "success"})
		}
	}
}
func querypersonhistoryim(c *gin.Context) {
	type personimhistory struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Sccid   string `json:"sccid" binding:"required"`
		Peerid  string `json:"peerid" binding:"required"`
		Pagenum int    `json:"pagenum" binding:"required"`
	}
	var json personimhistory
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd1 := fmt.Sprintf("select count(*) from scc_IMMessage where fromsccid= '%v' and tosccid = '%v' or fromsccid= '%v' and tosccid = '%v'", json.Sccid, json.Peerid, json.Peerid, json.Sccid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	if _, ok := sqlresult[0]["count(*)"]; ok {
		count, _ := strconv.Atoi(sqlresult[0]["count(*)"])
		//fmt.Println(sqlresult) //存在
		if 0 != count {
			if 0 == json.Pagenum {
				pagenum = 1 //不存在第0页
			}
			if count < numberperpage {
				if 1 == pagenum {
					fromcount = 0
				} else {
					fromcount = numberperpage
				}

			} else {
				fromcount = count - (numberperpage)*pagenum
				if fromcount < 0 {
					//int tmpnumberperpage = numberperpage;
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}
			sqlcmd1 := fmt.Sprintf("select id,fromsccid,tosccid,iminfo,imtype,filepath,created from scc_IMMessage where fromsccid= '%v' and tosccid = '%v' or fromsccid= '%v' and tosccid = '%v' order by id  limit %v,%v", json.Sccid, json.Peerid, json.Peerid, json.Sccid, fromcount, numberperpage)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"pagenum": pagenum, "total": count, "msginfo": sqlresult}})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "success"})
		}
	}

}
func querygrouphistoryim(c *gin.Context) {
	type personimhistory struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Groupid int `json:"groupid" binding:"required"`
		Pagenum int `json:"pagenum" binding:"required"`
	}
	var json personimhistory
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd1 := fmt.Sprintf("select count(*) from scc_IMGROUPMessage where groupid= '%v'", json.Groupid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	if _, ok := sqlresult[0]["count(*)"]; ok {
		count, _ := strconv.Atoi(sqlresult[0]["count(*)"])
		//(sqlresult) //存在
		if 0 != count {
			if 0 == json.Pagenum {
				pagenum = 1 //不存在第0页
			}
			if count < numberperpage {
				if 1 == pagenum {
					fromcount = 0
				} else {
					fromcount = numberperpage
				}

			} else {
				fromcount = count - (numberperpage)*pagenum
				if fromcount < 0 {
					//int tmpnumberperpage = numberperpage;
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}
			sqlcmd1 := fmt.Sprintf("select id,fromsccid,groupid,iminfo,imtype,filepath,created,msgid from scc_IMGROUPMessage where groupid = '%v' order by id  limit %v,%v", json.Groupid, fromcount, numberperpage)
			//(sqlcmd1)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"pagenum": pagenum, "total": count, "msginfo": sqlresult}})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "success"})
		}
	}
}
func moduserdetail(c *gin.Context) {
	type userdetail struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Sccid       string `json:"sccid" binding:"required"`
		Post        string `json:"post"`
		Mailbox     string `json:"mailbox"`
		Addr        string `json:"addr"`
		Phone       string `json:"phone"`
		Mobilephone string `json:"mobilephone"`
	}
	var json userdetail
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlcmd1 := fmt.Sprintf("Select count(*) from scc_userdetailed where sccid = '%v'", json.Sccid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)

	if _, ok := sqlresult[0]["count(*)"]; ok {

		//fmt.Println(sqlresult) //存在
		if "0" != sqlresult[0]["count(*)"] {
			sqlcmd := fmt.Sprintf("update scc_userdetailed  set post  = '%v' ,mailbox = '%v',address='%v',phone='%v',mobliephone='%v' where sccid = '%v'", json.Post, json.Mailbox, json.Addr, json.Phone, json.Mobilephone, json.Sccid)
			//(sqlcmd)
			sccinfo.tmpsql.Execsqlcmd(sqlcmd, true)
		} else {
			sqlcmd := fmt.Sprintf("insert into scc_userdetailed(sccid,post,mailbox,address,phone,mobliephone) values('%v','%v','%v','%v','%v','%v')", json.Sccid, json.Post, json.Mailbox, json.Addr, json.Phone, json.Mobilephone)
			//fmt.Println(sqlcmd)
			sccinfo.tmpsql.Execsqlcmd(sqlcmd, true)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "200"})

}
func querynearbyscc(c *gin.Context) {
	type nearby struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Distance  int    `json:"distance" binding:"required"`
		Longitude string `json:"longitude" binding:"required"`
		Latitude  string `json:"latitude" binding:"required"`
	}
	var json nearby
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	myredisnearby, _ := sccinfo.tmpredis.SccredisGetmembernearby(json.Longitude, json.Latitude, json.Distance)
	myredisnearbyinfo, _ := sccinfo.tmpredis.SccredisGetmembergpsinf(myredisnearby)
	var sccresult []map[string]string
	for i := range myredisnearby {
		long := strconv.FormatFloat(myredisnearbyinfo[i][0], 'f', -1, 64)
		lat := strconv.FormatFloat(myredisnearbyinfo[i][1], 'f', -1, 64)
		tmpgps := map[string]string{"longitude": long, "latitude": lat, "sccid": myredisnearby[i]}
		sccresult = append(sccresult, tmpgps)
	}
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sccresult})
}
func querysccdetail(c *gin.Context) {

	type userdetail struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Sccid string `json:"sccid" binding:"required"`
	}
	var json userdetail
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd := fmt.Sprintf("select  sccid,post,mailbox,address,phone,mobliephone from scc_userdetailed where sccid='%v'", json.Sccid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})

}
func querydingbysccid(c *gin.Context) {

	type dinginfo struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Sccid      string `json:"fromsccid" binding:"required"`
		Pagenum    int    `json:"pagenum" binding:"required"`
		Dingstatus string `json:"dingstatus" binding:"required"`
	}
	var json dinginfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sqlcmd1 string
	if "3" == json.Dingstatus {
		sqlcmd1 = fmt.Sprintf("select count(*) from scc_ding where sccfromding= '%v'  or scctoding= '%v'", json.Sccid, json.Sccid)
	} else {
		sqlcmd1 = fmt.Sprintf("select count(*) from scc_ding where sccfromding= '%v'  or scctoding= '%v'  and dingstatus='%v'", json.Sccid, json.Sccid, json.Dingstatus)
	}
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var sqlret1 []map[string]string
	if _, ok := sqlresult[0]["count(*)"]; ok {
		count, _ := strconv.Atoi(sqlresult[0]["count(*)"])
		//(sqlresult) //存在
		if 0 != count {
			if 0 == json.Pagenum {
				pagenum = 1 //不存在第0页
			}
			if count < numberperpage {
				if 1 == pagenum {
					fromcount = 0
				} else {
					fromcount = numberperpage
				}

			} else {
				fromcount = count - (numberperpage)*pagenum
				if fromcount < 0 {
					//int tmpnumberperpage = numberperpage;
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}

			if "" == json.Dingstatus {
				sqlcmd1 = fmt.Sprintf("select id,messgaeid,messagtype,dingtype,info,sccfromding,scctoding,dingstatus from scc_ding where sccfromding= '%v'  or scctoding= '%v'  order by id  limit %v,%v", json.Sccid, json.Sccid, fromcount, numberperpage)
				sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)
			} else {
				sqlcmd1 = fmt.Sprintf("select id,messgaeid,messagtype,dingtype,info,sccfromding,scctoding,dingstatus from scc_ding where sccfromding= '%v'  or scctoding= '%v' and dingstatus = '%v'  order by id  limit %v,%v", json.Sccid, json.Sccid, json.Dingstatus, fromcount, numberperpage)
				sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)
			}

		}
	}
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlret1})
}
func querydingbymsgid(c *gin.Context) {

	type dingbyid struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Messagetype string `json:"messagetype" binding:"required"`
		Messageid   string `json:"messageid" binding:"required"`
	}
	var json dingbyid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd := fmt.Sprintf("Select id,messgaeid,messagtype,dingtype,info,sccfromding,scctoding,dingstatus from scc_ding where messgaeid = '%v' and messagtype = '%v'", json.Messageid, json.Messagetype)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})

}
func setrouter(r *gin.Engine) {
	r.POST("/querydepartment", querysccdepartment)
	r.POST("/querydepartmentuser", querydepartmentuser)
	r.POST("/querygroup", querygroup)
	r.POST("/queryuser", queryuser)
	r.POST("/querygroupuser", querygroupuser)
	r.POST("/queryofflinemsg", queryofflinemsg)
	r.POST("/queryRecnetSession", queryRecnetSession)
	r.POST("/reportgps", reportgps)
	r.POST("/querygps", querygps)
	r.POST("/querypersonhistoryim", querypersonhistoryim)
	r.POST("/querygrouphistoryim", querygrouphistoryim)
	r.POST("/moduserdetail", moduserdetail)
	r.POST("/querysccuserdetail", querysccdetail)
	r.POST("/querynearbyscc", querynearbyscc)
	r.POST("/querydingbysccid", querydingbysccid)
	r.POST("/querydingbymsgid", querydingbymsgid)
}
