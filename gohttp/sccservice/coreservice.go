package coreservice

import (
	"fmt"
	sccsql "golearn/gomysql"
	sccredis "golearn/goredis"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type coreinfo struct {
	tmpsql   sccsql.Mysqlconnectpool
	tmpredis sccredis.Redisconnectpool
}

var sccinfo coreinfo

func init() {
	sccinfo.tmpsql.Initmysql("127.0.0.1", "root", "root", "SCC", 3306)
	sccinfo.tmpredis.Redisip = ("127.0.0.1:6379")
	sccinfo.tmpredis.ConnectRedis()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	setrouter(r)
	if err := r.Run(":9888"); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
func querysccdepartment(c *gin.Context) {
	bjson := c.DefaultQuery("json", "yes")
	departmentid := c.DefaultQuery("departmentid", "1")
	sqlcmd := fmt.Sprintf("Select s_departmentname,s_departmentid,s_grade,s_path,s_createtime,s_updatetime  from scc_department  where s_path like '%%/%v/%%'", departmentid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	if "yes" == bjson {
		c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})
	} else {

	}
}
func querydepartmentuser(c *gin.Context) {
	bjson := c.DefaultQuery("json", "yes")
	departmentid := c.DefaultQuery("departmentid", "1")
	onlydispatcher := c.DefaultQuery("onlydispatcher", "1")
	sqlcmd := fmt.Sprintf("Select s_user,s_grade,s_usertype,s_createtime,s_updatetime,s_alias,s_displayname from scc_user  where s_departmentid = '%v' and s_usertype>='%v'", departmentid, onlydispatcher)

	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	if "yes" == bjson {
		c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})
	} else {

	}
}
func querygroup(c *gin.Context) {
	bjson := c.DefaultQuery("json", "yes")
	sccid := c.DefaultQuery("sccid", "1")
	sqlcmd := fmt.Sprintf("Select s_groupid from scc_groupuser where s_user = '%v'", sccid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	if "yes" == bjson {
		c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})
	} else {

	}
}
func queryuser(c *gin.Context) {
	bjson := c.DefaultQuery("json", "yes")
	sccid := c.DefaultQuery("sccid", "1")
	sqlcmd := fmt.Sprintf("Select s_grade,s_usertype,s_createtime,s_updatetime,s_alias,s_displayname from scc_user  where s_user = '%v'", sccid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	if "yes" == bjson {
		c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})
	} else {

	}
}
func querygroupuser(c *gin.Context) {
	bjson := c.DefaultQuery("json", "yes")
	groupid := c.DefaultQuery("groupid", "1")
	sqlcmd := fmt.Sprintf("Select s_groupname,s_grade,s_grouptype,s_createtime,s_updatetime,s_creater from scc_group where s_groupid = %v", groupid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
	sqlcmd1 := fmt.Sprintf("Select scc_groupuser.s_user,scc_groupuser.s_displayname,scc_groupuser.s_grade,scc_groupuser.s_usertype,scc_groupuser.s_createtime,scc_groupuser.s_updatetime,scc_user.s_alias from scc_groupuser inner join scc_user on scc_groupuser.s_user =scc_user.s_user where scc_groupuser.s_groupid= %v", groupid)
	sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
	tmggroupid := fmt.Sprintf("group_%s", groupid)
	newmsg, _ := sccinfo.tmpredis.SccredisHget(tmggroupid, "newestmsg")
	tmpmsg := map[string]string{"lastmsg": newmsg}
	sqlresult = append(sqlresult, tmpmsg)
	for i := range sqlresult1 {

		fmt.Println(sqlresult1[i]["s_user"])
		userstatus, _ := sccinfo.tmpredis.SccredisHget(tmggroupid, sqlresult1[i]["s_user"])
		//statusmap := map[string]string{"status": userstatus}
		sqlresult1[i]["status"] = userstatus
	}
	if "yes" == bjson {
		c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"groupinfo": sqlresult, "userinfo": sqlresult1}})
	} else {

	}
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
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": myredisresult})
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
	fmt.Println(json.Sccid)
	sqlcmd := fmt.Sprintf("insert into gpsinfo(sccid,longitude,latitude,reporttime,gps,angle,speed,description) values('%v','%v','%v','%v','%v','%v','%v','%v');", json.Sccid, json.Longitude, json.Latitude, time.Now().Unix(), json.Gps, json.Angle, json.Speed, json.Description)
	fmt.Println(sqlcmd)
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
		fmt.Println(sqlresult) //存在
		if 0 != count {
			fmt.Println("存在 updata")
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
				sqlcmd1 := fmt.Sprintf("Select longitude,latitude,reporttime,gps,angle,speed,description from gpsinfo where sccid = '%v' and reporttime > '%v' and reporttime < '%v' limit %v,%v", json.Sccid, json.Starttime, json.Endtime, fromcount, numberperpage)
				sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"pagenum": pagenum, "total": count, "gpsinfo": sqlresult}})
				return
			} else {
				sqlcmd1 := fmt.Sprintf("Select longitude,latitude,reporttime,gps,angle,speed from gpsinfo where sccid = '%v' and reporttime > '%v' and reporttime < '%v' limit %v,%v", json.Sccid, json.Starttime, json.Endtime, fromcount, numberperpage)
				sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"pagenum": pagenum, "totoal": count, "gpsinfo": sqlresult}})
				return
			}

		} else {
			fmt.Println("不存在 insert")
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
		fmt.Println(sqlresult) //存在
		if 0 != count {
			fmt.Println("存在 updata")
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
			sqlcmd1 := fmt.Sprintf("select id,fromsccid,tosccid,iminfo,imtype,filepath,created from scc_IMMessage where fromsccid= '%v' and tosccid = '%v' or fromsccid= '%v' and tosccid = '%v' limit %v,%v", json.Sccid, json.Peerid, json.Peerid, json.Sccid, fromcount, numberperpage)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"pagenum": pagenum, "totoal": count, "msginfo": sqlresult}})
		} else {
			fmt.Println("不存在 insert")
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
		fmt.Println(sqlresult) //存在
		if 0 != count {
			fmt.Println("存在 updata", count)
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
			sqlcmd1 := fmt.Sprintf("select id,fromsccid,groupid,iminfo,imtype,filepath,created,msgid from scc_IMGROUPMessage where groupid = '%v' limit %v,%v", json.Groupid, fromcount, numberperpage)
			fmt.Println(sqlcmd1)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"pagenum": pagenum, "total": count, "msginfo": sqlresult}})
		} else {
			fmt.Println("不存在 insert")
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

		fmt.Println(sqlresult) //存在
		if "0" != sqlresult[0]["count(*)"] {
			fmt.Println("存在 updata")
			sqlcmd := fmt.Sprintf("update scc_userdetailed  set post  = '%v' ,mailbox = '%v',address='%v',phone='%v',mobliephone='%v' where sccid = '%v'", json.Post, json.Mailbox, json.Addr, json.Phone, json.Mobilephone, json.Sccid)
			fmt.Println(sqlcmd)
			sccinfo.tmpsql.Execsqlcmd(sqlcmd, true)
		} else {
			sqlcmd := fmt.Sprintf("insert into scc_userdetailed(sccid,post,mailbox,address,phone,mobliephone) values('%v','%v','%v','%v','%v','%v')", json.Sccid, json.Post, json.Mailbox, json.Addr, json.Phone, json.Mobilephone)
			fmt.Println(sqlcmd)
			sccinfo.tmpsql.Execsqlcmd(sqlcmd, true)
			fmt.Println("不存在 insert")
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
	bjson := c.DefaultQuery("json", "yes")
	sccid := c.DefaultQuery("sccid", "0")
	sqlcmd := fmt.Sprintf("select  sccid,post,mailbox,address,phone,mobliephone from scc_userdetailed where sccid='%v'", sccid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
	if "yes" == bjson {
		c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult})
	} else {

	}
}
func setrouter(r *gin.Engine) {
	r.GET("/querydepartment", querysccdepartment)
	r.GET("/querydepartmentuser", querydepartmentuser)
	r.GET("/querygroup", querygroup)
	r.GET("/queryuser", queryuser)
	r.GET("/querygroupuser", querygroupuser)
	r.POST("/queryofflinemsg", queryofflinemsg)
	r.POST("/queryRecnetSession", queryRecnetSession)
	r.POST("/reportgps", reportgps)
	r.GET("/querygps", querygps)
	r.POST("/querypersonhistoryim", querypersonhistoryim)
	r.POST("/querygrouphistoryim", querygrouphistoryim)
	r.POST("/moduserdetail", moduserdetail)
	r.GET("/querysccuserdetail", querysccdetail)
	r.POST("/querynearbyscc", querynearbyscc)
}
