package coreservice

import (
	"fmt"
	sccsql "golearn/gomysql"
	sccredis "golearn/goredis"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "golearn/sccprotobuf"

	_ "golearn/docs"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
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
func reverse(s []map[string]string) []map[string]string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// @查询部门信息
// @Description querydepartment 查询部门信息
// @Accept  json
// @Produce json
// @Param article body coreservice.Querysccdeparment true "查询部门信息"
// @Router /querydepartment [post]
func querysccdepartment(c *gin.Context) {
	var json Querysccdeparment
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

// @查询部门成员信息
// @Description querydepartmentuser 查询部门成员信息
// @Accept  json
// @Produce json
// @Param article body coreservice.Querysccdeparmentuser true "查询部门成员信息"
// @Router /querydepartmentuser [post]
func querydepartmentuser(c *gin.Context) {

	onlydispatcher := 0
	var json Querysccdeparmentuser
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

// @查询群组信息
// @Description querygroup 查询群组信息
// @Accept  json
// @Produce json
// @Param article body coreservice.Querygroupinfo true "查询群组信息"
// @Router /querygroup [post]
func querygroup(c *gin.Context) {
	var json Querygroupinfo
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

// @查询个人成员详细信息
// @Description queryuser 查询个人成员详细信息
// @Accept  json
// @Produce json
// @Param article body coreservice.Queryuserinfo true "查询个人成员详细信息"
// @Router /queryuser [post]
func queryuser(c *gin.Context) {
	var json Queryuserinfo
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

// @查询群组成员
// @Description querygroupuser 查询群组成员
// @Accept  json
// @Produce json
// @Param article body coreservice.Quserygroupuserinfo true "查询群组成员"
// @Router /querygroupuser [post]
func querygroupuser(c *gin.Context) {

	var json Quserygroupuserinfo
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

// @查询离线消息
// @Description queryofflinemsg 查询离线消息
// @Accept  json
// @Produce json
// @Param article body coreservice.Querypersonofflinemsg true "查询离线消息"
// @Router /queryofflinemsg [post]
func queryofflinemsg(c *gin.Context) {
	var json Querypersonofflinemsg
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

// @查询最近会话
// @Description queryRecnetSession 查询最近会话
// @Accept  json
// @Produce json
// @Param article body coreservice.QueryRecntSession true "查询最近会话"
// @Router /queryRecnetSession [post]
func queryRecnetSession(c *gin.Context) {
	var json QueryRecntSession
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

// @上报轨迹
// @Description reportgps 上报轨迹
// @Accept  json
// @Produce json
// @Param article body coreservice.Reportgps true "上报轨迹"
// @Router /reportgps [post]
func reportgps(c *gin.Context) {
	var json Reportgps
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

// @查询个人的历史轨迹
// @Description querygps 查询个人的历史轨迹
// @Accept  json
// @Produce json
// @Param article body coreservice.Querygps true "sccid和时间查询历史轨迹"
// @Router /querygps [post]
func querygps(c *gin.Context) {
	var json Querygps
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
				sqlcmd1 := fmt.Sprintf("Select longitude,latitude,reporttime,gps,angle,speed,description from gpsinfo where sccid= '%v' order by id DESC limit 1", json.Sccid)
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

// @查询个人的历史信息
// @Description querypersonhistoryim 查询个人消息
// @Accept  json
// @Produce json
// @Param article body coreservice.Querrypersonimhistory true "查询历史信息"
// @Router /querypersonhistoryim [post]
func querypersonhistoryim(c *gin.Context) {
	var json Querrypersonimhistory
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

// @查询群组的历史信息
// @Description querygrouphistoryim 根据群组查询历史消息
// @Accept  json
// @Produce json
// @Param article body coreservice.Querygroupimhistory true "根据群组查询历史消息"
// @Router /querygrouphistoryim [post]
func querygrouphistoryim(c *gin.Context) {
	var json Querygroupimhistory
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

// @修改用户信息
// @Description moduserdetail 修改用户详细信息
// @Accept  json
// @Produce json
// @Param article body coreservice.Moduserdetail true "修改用户信息"
// @Router /moduserdetail [post]
func moduserdetail(c *gin.Context) {
	var json Moduserdetail
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(json)
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

// @查询附近的人
// @Description querynearbyscc 查询附近的人
// @Accept  json
// @Produce json
// @Param article body coreservice.Querynearby true "查询附近的人"
// @Router /querynearbyscc [post]
func querynearbyscc(c *gin.Context) {
	var json Querynearby
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

// @查询用户的详细信息
// @Description querysccuserdetail 查询用户详细信息
// @Accept  json
// @Produce json
// @Param article body coreservice.Queryuserdetail true "查询用户详细信息"
// @Router /querysccuserdetail [post]
func querysccdetail(c *gin.Context) {

	var json Queryuserdetail
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

// @查询必达消息 根据sccid查询和我有关的群组必达 //0 是和我相关的  1 是我发送的 2 是我接收的
// @Description querypersondingbysccid
// @Accept  json
// @Produce json
// @Param article body coreservice.Relationding true "0 是和我相关的  1 是我发送的 2 是我接收的"
// @Router /querypersondingbysccid [post]
func querydingbysccidbyperson(c *gin.Context) {
	var json Relationding
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var exrencmd string
	if json.Sccidstatus == 0 {
		exrencmd = fmt.Sprintf("sccfromding= '%v'  or scctoding= '%v'", json.Sccid, json.Sccid)
	} else if json.Sccidstatus == 1 {
		exrencmd = fmt.Sprintf("sccfromding= '%v'", json.Sccid)
	} else if json.Sccidstatus == 2 {
		exrencmd = fmt.Sprintf("scctoding= '%v'", json.Sccid)
	}
	var sqlcmd1 string
	if "3" == json.Dingstatus {
		sqlcmd1 = fmt.Sprintf("select count(*) from scc_ding where messagtype = 0 and (%v)", exrencmd)
	} else {
		sqlcmd1 = fmt.Sprintf("select count(*) from scc_ding where dingstatus='%v' and messagtype = 0 and (%v) ", json.Dingstatus, exrencmd)
	}
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var sqlret1 []map[string]string
	var icount int
	if _, ok := sqlresult[0]["count(*)"]; ok {
		count, _ := strconv.Atoi(sqlresult[0]["count(*)"])
		icount = count
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

			if "3" == json.Dingstatus {
				sqlcmd1 = fmt.Sprintf("select id,messgaeid,messagtype,dingtype,info,sccfromding,scctoding,dingstatus,createtime,dingtime,groupid,filepath,replyfilepath,replymsg from scc_ding where  messagtype = 0 and (%v)  order by id  limit %v,%v", exrencmd, fromcount, numberperpage)
				sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)
			} else {
				sqlcmd1 = fmt.Sprintf("select id,messgaeid,messagtype,dingtype,info,sccfromding,scctoding,dingstatus,createtime,dingtime,groupid,filepath,replyfilepath,replymsg from scc_ding where  messagtype = 0 and (%v) and dingstatus = '%v'  order by id   limit %v,%v", exrencmd, json.Dingstatus, fromcount, numberperpage)
				sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)
			}

		}
	}
	tmpresult := reverse(sqlret1)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": tmpresult, "pagenum": pagenum, "total": icount})
}

// @查询必达消息 根据sccid查询和我有关的群组必达 //0 是和我相关的  1 是我发送的 2 是我接收的
// @Description querygroupdingbysccid
// @Accept  json
// @Produce json
// @Param article body coreservice.Relationding true "0 是和我相关的  1 是我发送的 2 是我接收的"
// @Router /querygroupdingbysccid [post]
func querydingbysccidbygroup(c *gin.Context) {

	var json Relationding
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var exrencmd string
	if json.Sccidstatus == 0 {
		exrencmd = fmt.Sprintf("sccfromding= '%v' or scctoding= '%v'", json.Sccid, json.Sccid)
	} else if json.Sccidstatus == 1 {
		exrencmd = fmt.Sprintf("sccfromding= '%v'", json.Sccid)
	} else if json.Sccidstatus == 2 {
		exrencmd = fmt.Sprintf("scctoding= '%v'", json.Sccid)
	}
	var sqlcmd1 string
	if "3" == json.Dingstatus {
		sqlcmd1 = fmt.Sprintf("select count(distinct messgaeid,groupid) as total from scc_ding where  messagtype = 1 and (%v) ", exrencmd)
	} else {
		sqlcmd1 = fmt.Sprintf("select count(distinct messgaeid,groupid) as total from scc_ding where dingstatus='%v' and messagtype = 1 and (%v) ", json.Dingstatus, exrencmd)
	}
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var sqlret1 []map[string]string
	var icount int
	if _, ok := sqlresult[0]["total"]; ok {
		count, _ := strconv.Atoi(sqlresult[0]["total"])
		icount = count
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
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}

			if "3" == json.Dingstatus {
				sqlcmd1 = fmt.Sprintf("select  messgaeid,groupid,info,sccfromding,createtime,filepath,replyfilepath,replymsg from scc_ding where   messagtype = 1 and (%v) group by messgaeid , groupid order by id limit %v,%v", exrencmd, fromcount, numberperpage)
				sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)
			} else {
				sqlcmd1 = fmt.Sprintf("select  messgaeid,groupid,info,sccfromding,createtime,filepath,replyfilepath,replymsg from scc_ding where  messagtype = 1 and (%v) and dingstatus = '%v' group by messgaeid , groupid order by id limit %v,%v", exrencmd, json.Dingstatus, fromcount, numberperpage)
				sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)
			}
			sqlret1 = reverse(sqlret1)
			// 再此基础上把 每个确认的个数查一下  性能十分差劲 要优化
			sqllen := len(sqlret1)
			for i := 0; i < sqllen; i++ {
				sqlcmd1 = fmt.Sprintf("select count(*) as total from scc_ding where messagtype = 1 and messgaeid='%v' and groupid = '%v'", sqlret1[i]["messgaeid"], sqlret1[i]["groupid"])
				sqlret2 := sccinfo.tmpsql.SelectData(sqlcmd1)
				if _, ok := sqlret2[0]["total"]; ok {
					//totalsend, _ := strconv.Atoi(sqlret2[0]["total"])
					sqlret1[i]["totalsend"] = sqlret2[0]["total"]
					sqlcmd1 = fmt.Sprintf("select count(*) as total from scc_ding where messagtype = 1 and messgaeid='%v' and groupid = '%v' and dingstatus=1", sqlret1[i]["messgaeid"], sqlret1[i]["groupid"])
					sqlret3 := sccinfo.tmpsql.SelectData(sqlcmd1)
					if _, ok1 := sqlret3[0]["total"]; ok1 {
						sqlret1[i]["totalding"] = sqlret3[0]["total"]
					}
				}
			}

		}
	}
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlret1, "pagenum": pagenum, "total": icount})
}

// @查询必达消息 如果个人必达 messagetype是0  groupid是0  群组必达 messagetype是1 groupid是群组id
// @Description querydingbymsgid
// @Accept  json
// @Produce json
// @Param article body coreservice.Dingbyid true "如果个人必达 messagetype是0  groupid是0  群组必达 messagetype是1 groupid是群组id"
// @Router /querydingbymsgid [post]
func querydingbymsgid(c *gin.Context) {

	var json Dingbyid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd := fmt.Sprintf("Select id,messgaeid,messagtype,dingtype,info,sccfromding,scctoding,dingstatus,createtime,dingtime,filepath,replyfilepath,replymsg from scc_ding where messgaeid = '%v' and messagtype = '%v' and groupid = '%v'", json.Messageid, json.Messagetype, json.Groupid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})

}

// @通过groupid和messageid查询必达的详情
// @Description querydingbysccidandgroupid
// @Accept  json
// @Produce json
// @Param article body coreservice.Dingfrommsgidindgroupid true "根据群组id和msessageid查询群组必达的必达情况"
// @Router /querydingbysccidandgroupid [post]
func querydingbysccidandgroupid(c *gin.Context) {

	var json Dingfrommsgidindgroupid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//var totalding int
	var sqlcmd1 string

	sqlcmd1 = fmt.Sprintf("select count(*) as total,info,sccfromding,filepath from scc_ding where messagtype = 1 and messgaeid='%v' and groupid = '%v'", json.Messageid, json.Groupid)

	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd1)

	//totalding, _ = strconv.Atoi(sqlresult[0]["total"])

	if json.Dingstatus == "0" {
		sqlcmd1 = fmt.Sprintf("select count(*)  from scc_ding where messagtype = 1 and messgaeid='%v' and groupid = '%v' and dingstatus = 0", json.Messageid, json.Groupid)
	} else {
		sqlcmd1 = fmt.Sprintf("select count(*)  from scc_ding where messagtype = 1 and messgaeid='%v' and groupid = '%v' and dingstatus = 1", json.Messageid, json.Groupid)
	}
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var sqlret1 []map[string]string
	var icount int
	if _, ok := sqlresult[0]["count(*)"]; ok {
		count, _ := strconv.Atoi(sqlresult[0]["count(*)"])
		icount = count
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
			if json.Dingstatus == "0" {
				sqlcmd1 = fmt.Sprintf("select id,messgaeid,messagtype,dingtype,info,sccfromding,scctoding,dingstatus,createtime,dingtime,groupid,filepath,replyfilepath,replymsg from scc_ding where  messagtype = 1 and messgaeid='%v' and groupid = '%v'  and dingstatus = 0  order by id desc limit %v,%v", json.Messageid, json.Groupid, fromcount, numberperpage)
			} else {
				sqlcmd1 = fmt.Sprintf("select id,messgaeid,messagtype,dingtype,info,sccfromding,scctoding,dingstatus,createtime,dingtime,groupid,filepath,replyfilepath,replymsg from scc_ding where  messagtype = 1 and messgaeid='%v' and groupid = '%v'  and dingstatus = 1 order by id desc limit %v,%v", json.Messageid, json.Groupid, fromcount, numberperpage)
			}

			sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)

		}
	}
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlret1, "pagenum": pagenum, "totalInfo": sqlresult3, "currenttypecount": icount})
}

// @通过主叫的被必达的SCC查询
// @Description querydingbyfromsccid
// @Accept  json
// @Produce json
// @Param article body coreservice.Fromdinginfo true "根据被叫SCCid查询必达的情况"
// @Router /querydingbyfromsccid [post]
func querydingbyfromsccid(c *gin.Context) {
	var json Fromdinginfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sqlcmd1 string
	if "3" == json.Dingstatus {
		sqlcmd1 = fmt.Sprintf("select count(distinct messgaeid,groupid) as total from scc_ding where   sccfromding= '%v' ", json.Sccid)
	} else {
		sqlcmd1 = fmt.Sprintf("select count(distinct messgaeid,groupid) as total from scc_ding where dingstatus='%v'  and sccfromding= '%v' ", json.Dingstatus, json.Sccid)
	}
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var sqlret1 []map[string]string
	var icount int
	if _, ok := sqlresult[0]["total"]; ok {
		count, _ := strconv.Atoi(sqlresult[0]["total"])
		icount = count
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
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}

			if "3" == json.Dingstatus {
				sqlcmd1 = fmt.Sprintf("select distinct messgaeid,groupid,info,sccfromding,filepath,replyfilepath,replymsg from scc_ding where   sccfromding= '%v' order by id  desc limit %v,%v", json.Sccid, fromcount, numberperpage)
				sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)
			} else {
				sqlcmd1 = fmt.Sprintf("select distinct messgaeid,groupid,info,sccfromding,filepath,replyfilepath,replymsg from scc_ding where   sccfromding= '%v' and dingstatus = '%v' order by id desc limit %v,%v", json.Sccid, json.Dingstatus, fromcount, numberperpage)
				sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)
			}

		}
	}
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlret1, "pagenum": pagenum, "total": icount})
}

// @通过被叫的被必达的SCC查询
// @Description querydingbytosccid
// @Accept  json
// @Produce json
// @Param article body coreservice.Todinginfo true "根据被叫SCCid查询必达的情况"
// @Router /querydingbytosccid [post]
func querydingbytosccid(c *gin.Context) {
	var json Todinginfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sqlcmd1 string
	if "3" == json.Dingstatus {
		sqlcmd1 = fmt.Sprintf("select count(distinct messgaeid,groupid) as total from scc_ding where scctoding= '%v' ", json.Sccid)
	} else {
		sqlcmd1 = fmt.Sprintf("select count(distinct messgaeid,groupid) as total from scc_ding where dingstatus='%v'  and scctoding= '%v' ", json.Dingstatus, json.Sccid)
	}
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var sqlret1 []map[string]string
	var icount int
	if _, ok := sqlresult[0]["total"]; ok {
		count, _ := strconv.Atoi(sqlresult[0]["total"])
		icount = count
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

			if "3" == json.Dingstatus {
				sqlcmd1 = fmt.Sprintf("select distinct messgaeid,groupid,info,sccfromding,filepath,replyfilepath,replymsg from scc_ding where  scctoding= '%v' order by id  desc limit %v,%v", json.Sccid, fromcount, numberperpage)
				sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)
			} else {
				sqlcmd1 = fmt.Sprintf("select distinct messgaeid,groupid,info,sccfromding,filepath,replyfilepath,replymsg from scc_ding where scctoding= '%v' and dingstatus = '%v' order by id desc limit %v,%v", json.Sccid, json.Dingstatus, fromcount, numberperpage)
				sqlret1 = sccinfo.tmpsql.SelectData(sqlcmd1)
			}

		}
	}
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlret1, "pagenum": pagenum, "total": icount})
}

// @通过msgid查询msg的 详细信息
// @Description querymsgbymsgid
// @Accept  json
// @Produce json
// @Param article body coreservice.Todinginfo true "根据被叫SCCid查询必达的情况"
// @Router /querymsgbymsgid [post]
func querymsgbymsgid(c *gin.Context) {
	var json Querymsgbyid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sqlcmd1 string
	if "group" == json.Msgtype {
		sqlcmd1 = fmt.Sprintf("select fromsccid,tosccid,iminfo,imtype,filepath,created from scc_IMMessage where id= '%v' ", json.Msgid)
	} else if "person" == json.Msgtype {
		sqlcmd1 = fmt.Sprintf("select fromsccid,groupid,iminfo,imtype,filepath,created,msgid from scc_IMGROUPMessage where id= '%v' ", json.Msgid)
	}
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})
}

// @通过msgid查询msg的 详细信息
// @Description sccupdatetoken
// @Accept  json
// @Produce json
// @Param article body coreservice.Todinginfo true "根据被叫SCCid查询必达的情况"
// @Router /sccupdatetoken [post]
func sccupdatetoken(c *gin.Context) {
	var json MobilephoneInfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//tmpkey := json.Sccid
	sccinfo.tmpredis.SccredisHSet(json.Sccid, "token", json.Token)
	sccinfo.tmpredis.SccredisHSet(json.Sccid, "type", strings.ToUpper(json.Mobilephonetype))
	c.JSON(http.StatusOK, gin.H{"result": "success"})
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
	r.POST("/querymsgbymsgid", querymsgbymsgid)
	r.POST("/querysccuserdetail", querysccdetail)
	r.POST("/querynearbyscc", querynearbyscc)
	r.POST("/querypersondingbysccid", querydingbysccidbyperson)
	r.POST("/querygroupdingbysccid", querydingbysccidbygroup)
	r.POST("/querydingbysccidandgroupid", querydingbysccidandgroupid)
	r.POST("/querydingbymsgid", querydingbymsgid)
	r.POST("/querydingbyfromsccid", querydingbyfromsccid)
	r.POST("/querydingbytosccid", querydingbytosccid)
	r.POST("/updatetoken", sccupdatetoken)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
