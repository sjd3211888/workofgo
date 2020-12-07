package workflow

import (
	"fmt"
	sccsql "golearn/gomysql"
	sccredis "golearn/goredis"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

type coreinfo struct {
	tmpsql    sccsql.Mysqlconnectpool
	scctmpsql sccsql.Mysqlconnectpool
	tmpredis  sccredis.Redisconnectpool
}

var sccinfo coreinfo

func init() {

	/*	year, month, _ := time.Now().Date()
		thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		start := thisMonth.AddDate(0, -1, 0).Format("2006-01-02 15:04:05")
		end := thisMonth.AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
		timeRange := fmt.Sprintf("%s~%s", start, end)
		fmt.Println(timeRange)*/
	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation

	var conf map[string]map[string]string
	if _, err := toml.DecodeFile("./sccconfig.toml", &conf); err != nil {
		// handle error
	}
	Host := conf["workdaping"]["Host"]
	Username := conf["workdaping"]["Username"]
	Password := conf["workdaping"]["Password"]
	Dbname := conf["workdaping"]["Dbname"]
	Port := conf["workdaping"]["Port"]
	iport, _ := strconv.Atoi(Port)
	Serhost := conf["workdaping"]["Httpserverhost"]
	Redisip := conf["workdaping"]["Redisip"]
	fmt.Println(Host)
	go func(Host string, Username string, Password string, Dbname string, Serhost string, Redisip string, iport int) {
		sccinfo.tmpsql.Initmysql(Host, Username, Password, Dbname, iport)
		sccinfo.scctmpsql.Initmysql(Host, Username, Password, "SCC", iport)
		sccinfo.tmpredis.Redisip = (Redisip)
		sccinfo.tmpredis.ConnectRedis()
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		r.Use(cors())
		setrouter(r)

		if err := r.Run(Serhost); err != nil {
			fmt.Println("startup service failed, err:\n", err)
		}
	}(Host, Username, Password, Dbname, Serhost, Redisip, iport)

}
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
func querylefttopquestion(c *gin.Context) {

	var totalworknum int
	var resolveed int
	var yesterdayproblem int

	var lightunCompleted int
	var normalunCompleted int
	var seriousunCompleted int

	var lightCompleted int
	var normalCompleted int
	var seriousCompleted int
	sqlcmd := fmt.Sprintf("select count(distinct appid ) as total from scc_workflow") //总问题数
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	if tmpworknum, ok := sqlresult[0]["total"]; ok {
		sqlcmd1 := fmt.Sprintf("select count(distinct appid ) as total from scc_workflow where appnextnode=200;") //总解决问题数
		sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
		totalworknum, _ = strconv.Atoi(tmpworknum)
		if tmpresolveed, ok := sqlresult1[0]["total"]; ok {
			resolveed, _ = strconv.Atoi(tmpresolveed)
		}
	}

	timeStr := time.Now().Format("2006-01-02")
	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
	yesterdattime := (t.Unix() - 3600*24*2)

	sqlcmd2 := fmt.Sprintf("select count(*) as total from scc_apply where createtime>'%v' and createtime< '%v';", yesterdattime, yesterdattime+3600*24) //总解决问题数
	fmt.Println(sqlcmd2)
	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult3[0]["total"]; ok {
		yesterdayproblem, _ = strconv.Atoi(tmpyesterday)
	}

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply where grade =0 and status=1;") //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmplightunCompleted, ok := sqlresult3[0]["total"]; ok {
		lightunCompleted, _ = strconv.Atoi(tmplightunCompleted)
	}

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply where grade =1 and status=1;") //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpnormalunCompleted, ok := sqlresult3[0]["total"]; ok {
		normalunCompleted, _ = strconv.Atoi(tmpnormalunCompleted)
	}

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply where grade =2 and status=1;") //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpseriousunCompleted, ok := sqlresult3[0]["total"]; ok {
		seriousunCompleted, _ = strconv.Atoi(tmpseriousunCompleted)
	}

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply where grade =0;") //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmplightCompleted, ok := sqlresult3[0]["total"]; ok {
		lightCompleted, _ = strconv.Atoi(tmplightCompleted)
	}

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply where grade =1;") //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpnormalCompleted, ok := sqlresult3[0]["total"]; ok {
		normalCompleted, _ = strconv.Atoi(tmpnormalCompleted)
	}

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply where grade =2 ;") //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpseriousCompleted, ok := sqlresult3[0]["total"]; ok {
		seriousCompleted, _ = strconv.Atoi(tmpseriousCompleted)
	}
	tmprestlt := make(map[string]interface{})
	tmprestlt["resolved"] = resolveed
	tmprestlt["workYesterday"] = yesterdayproblem
	tmprestlt["remain"] = totalworknum - resolveed
	tmprestlt["lightunCompleted"] = lightunCompleted
	tmprestlt["normalunCompleted"] = normalunCompleted
	tmprestlt["seriousunCompleted"] = seriousunCompleted
	tmprestlt["lightCompleted"] = lightCompleted
	tmprestlt["normalCompleted"] = normalCompleted
	tmprestlt["seriousCompleted"] = seriousCompleted
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tmprestlt})
}
func querylefttopmidquestion(c *gin.Context) {

	var lightunCompleted int
	var normalunCompleted int
	var seriousunCompleted int
	sqlcmd2 := fmt.Sprintf("select count(*) as total from scc_apply where grade =0 and status=1;") //总解决问题数
	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmplightunCompleted, ok := sqlresult3[0]["total"]; ok {
		lightunCompleted, _ = strconv.Atoi(tmplightunCompleted)
	}

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply where grade =1 and status=1;") //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpnormalunCompleted, ok := sqlresult3[0]["total"]; ok {
		normalunCompleted, _ = strconv.Atoi(tmpnormalunCompleted)
	}

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply where grade =2 and status=1;") //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpseriousunCompleted, ok := sqlresult3[0]["total"]; ok {
		seriousunCompleted, _ = strconv.Atoi(tmpseriousunCompleted)
	}

	tmprestlt := make(map[string]interface{})
	tmprestlt["lightunCompleted"] = lightunCompleted
	tmprestlt["normalunCompleted"] = normalunCompleted
	tmprestlt["seriousunCompleted"] = seriousunCompleted

	var Dishquality int
	var serverquality int
	var environmentalquality int

	sqlcmd2 = fmt.Sprintf("select t1.templateid,t2.p_type from scc_apply t1 inner join scc_worktempplate t2 on t1.templateid = t2.workid;") //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	for _, v := range sqlresult3 {
		if v["p_type"] == "1" {
			environmentalquality++
		} else if v["p_type"] == "10" {
			serverquality++
		} else if v["p_type"] == "100" {
			Dishquality++
		}
	}
	tmprestlt2 := make([]map[string]interface{}, 0)
	Dishqualitymao := make(map[string]interface{})
	Dishqualitymao["number"] = Dishquality
	Dishqualitymao["companyName"] = "菜品问题"
	tmprestlt2 = append(tmprestlt2, Dishqualitymao)

	serverqualitymao := make(map[string]interface{})
	serverqualitymao["number"] = serverquality
	serverqualitymao["companyName"] = "服务问题"
	tmprestlt2 = append(tmprestlt2, serverqualitymao)

	environmentalqualitmao := make(map[string]interface{})
	environmentalqualitmao["number"] = environmentalquality
	environmentalqualitmao["companyName"] = "环境问题"
	tmprestlt2 = append(tmprestlt2, environmentalqualitmao)
	tmprestlt["unAnswerTopTen"] = tmprestlt2

	c.JSON(http.StatusOK, gin.H{"success": true, "data": tmprestlt})
}
func queryleftobottom(c *gin.Context) {

	var totalnum int
	var totalnumresove int
	sqlcmd2 := fmt.Sprintf("select count(*) as total from scc_apply;") //总问题数
	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmplightCompleted, ok := sqlresult3[0]["total"]; ok {
		totalnum, _ = strconv.Atoi(tmplightCompleted)
	}

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply where  status=1;") //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmplightunCompleted, ok := sqlresult3[0]["total"]; ok {
		totalnumresove, _ = strconv.Atoi(tmplightunCompleted)
	}

	tmprestlt2 := make(map[string]interface{})
	tmprestlt2["totalNumber"] = totalnum
	tmprestlt2["totalnumresove"] = totalnumresove
	tmprestlt2["totalNumbercollection"] = totalnum
	tmprestlt2["totalnumresovetotalNumbercollection"] = totalnumresove
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tmprestlt2})
}
func querycenterbottom(c *gin.Context) {
	tmprestlt2 := make([]map[string]interface{}, 0)
	sqlcmd2 := fmt.Sprintf("select t1.grade,t1.textinfo,t2.p_type,t1.createtime,t1.status from scc_apply t1 inner join scc_worktempplate t2 where t2.workid=t1.templateid;") //总解决问题数
	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd2)
	for _, v := range sqlresult3 {
		tmpmap := make(map[string]interface{})
		tmpgrade, _ := strconv.Atoi(v["grade"])
		tmpmap["number"] = tmpgrade
		tmpmap["title"] = v["textinfo"]
		tmptype, _ := strconv.Atoi(v["p_type"])
		tmpmap["type"] = tmptype
		tmptime, _ := strconv.ParseInt(v["createtime"], 10, 64)
		t := time.Unix(int64(tmptime), 0)

		//返回string
		dateStr := t.Format("2006/01/02 15:04:05")

		tmpmap["reportTime"] = dateStr
		tmpstatus, _ := strconv.Atoi(v["status"])
		tmpmap["status"] = tmpstatus
		tmprestlt2 = append(tmprestlt2, tmpmap)
	}
	//fmt.Println(tmprestlt2)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tmprestlt2})
}

type userinfo struct {
	researcher string
	count      int
}
type UserList []userinfo

func (p UserList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p UserList) Len() int           { return len(p) }
func (p UserList) Less(i, j int) bool { return p[i].count > p[j].count }

func queryrightcenter(c *gin.Context) {
	tmprestlt2 := make([]map[string]interface{}, 0)
	sqlcmd2 := fmt.Sprintf("select t2.researcher from scc_apply t1 inner join scc_worktempplate t2 on t2.workid = t1.templateid where t1.status=1;") //总解决问题数
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)

	tmpusermap := make(map[string]int)
	for _, v := range sqlresult2 {

		if _, ok := tmpusermap[v["researcher"]]; ok {
			//存在
			tmpusermap[v["researcher"]] = tmpusermap[v["researcher"]] + 1
		} else {
			tmpusermap[v["researcher"]] = 0
		}
	}
	tmplist := func(m map[string]int) UserList {
		p := make(UserList, len(m))
		i := 0
		for k, v := range m {
			p[i] = userinfo{k, v}
			i++
		}
		sort.Sort(p)
		return p
	}(tmpusermap)
	fmt.Println("ssssssssssss", tmplist)
	for i := 0; i < len(tmplist); i++ {
		if i > 9 {
			break
		}
		tmpsaveList := make(map[string]interface{})

		sqlcmd3 := fmt.Sprintf("select s_displayname from scc_user  where s_user= %v;", tmplist[i].researcher) //总解决问题数

		sqlresult3 := sccinfo.scctmpsql.SelectData(sqlcmd3)
		tmpsaveList["securityName"] = sqlresult3[0]["s_displayname"]
		tmpsaveList["securityIntegral"] = tmplist[i].count
		tmprestlt2 = append(tmprestlt2, tmpsaveList)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tmprestlt2})
}
func queryrightcenterunsolve(c *gin.Context) {
	tmprestlt2 := make([]map[string]interface{}, 0)
	sqlcmd2 := fmt.Sprintf("select t2.researcher from scc_apply t1 inner join scc_worktempplate t2 on t2.workid = t1.templateid where t1.status=0;") //总解决问题数
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	tmpusermap := make(map[string]int)
	for _, v := range sqlresult2 {

		if _, ok := tmpusermap[v["researcher"]]; ok {
			//存在
			tmpusermap[v["researcher"]] = tmpusermap[v["researcher"]] + 1
		} else {
			tmpusermap[v["researcher"]] = 1
		}
	}
	tmplist := func(m map[string]int) UserList {
		p := make(UserList, len(m))
		i := 0
		for k, v := range m {
			p[i] = userinfo{k, v}
			i++
		}
		sort.Sort(p)
		return p
	}(tmpusermap)
	fmt.Println("ssssssssssss", tmplist)
	for i := 0; i < len(tmplist); i++ {
		if i > 9 {
			break
		}
		tmpsaveList := make(map[string]interface{})
		sqlcmd3 := fmt.Sprintf("select s_displayname from scc_user  where s_user= %v;", tmplist[i].researcher) //总解决问题数
		sqlresult3 := sccinfo.scctmpsql.SelectData(sqlcmd3)
		tmpsaveList["secfirmName"] = sqlresult3[0]["s_displayname"]
		tmpsaveList["securityIntegral"] = tmplist[i].count
		tmprestlt2 = append(tmprestlt2, tmpsaveList)
	}
	//fmt.Println(tmprestlt2)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tmprestlt2})
}
func Decimal(value float64) string {
	return (fmt.Sprintf("%.2f", value))
}
func queryrighttop(c *gin.Context) {
	var manyicount int
	var bumanyicount int
	tmprestlt2 := make([]map[string]interface{}, 12)
	sqlcmd2 := fmt.Sprintf("select count(*) as total from scc_apply;") //不满意个数
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		bumanyicount, _ = strconv.Atoi(tmpyesterday)
	}

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_statistica;") //不满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		manyicount, _ = strconv.Atoi(tmpyesterday)
	}
	tmptimestr := make([]string, 24)
	tmptimestr[0] = "2020-01-01 00:00:00"
	tmptimestr[1] = "2020-01-31 23:59:59"
	tmptimestr[2] = "2020-02-01 00:00:00"
	tmptimestr[3] = "2020-02-29 23:59:59"
	tmptimestr[4] = "2020-03-01 00:00:00"
	tmptimestr[5] = "2020-03-31 23:59:59"
	tmptimestr[6] = "2020-04-01 00:00:00"
	tmptimestr[7] = "2020-04-30 23:59:59"
	tmptimestr[8] = "2020-05-01 00:00:00"
	tmptimestr[9] = "2020-05-31 23:59:59"
	tmptimestr[10] = "2020-06-01 00:00:00"
	tmptimestr[11] = "2020-06-30 23:59:59"
	tmptimestr[12] = "2020-07-01 00:00:00"
	tmptimestr[13] = "2020-07-31 23:59:59"
	tmptimestr[14] = "2020-08-01 00:00:00"
	tmptimestr[15] = "2020-08-31 23:59:59"
	tmptimestr[16] = "2020-09-01 00:00:00"
	tmptimestr[17] = "2020-09-30 23:59:59"
	tmptimestr[18] = "2020-10-01 00:00:00"
	tmptimestr[19] = "2020-10-31 23:59:59"
	tmptimestr[20] = "2020-11-01 00:00:00"
	tmptimestr[21] = "2020-11-30 23:59:59"
	tmptimestr[22] = "2020-12-01 00:00:00"
	tmptimestr[23] = "2020-12-31 23:59:59"
	for i := 0; i < 12; i++ {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", tmptimestr[i*2], time.Local)
		start := t.Unix()

		t1, _ := time.ParseInLocation("2006-01-02 15:04:05", tmptimestr[i*2+1], time.Local)
		end := t1.Unix()
		sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply where createtime > %v and createtime < %v;", start, end) //不满意个数
		sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
		var tmpmanyi int
		var tmpbumanyi int
		if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
			tmpbumanyi, _ = strconv.Atoi(tmpyesterday)
		}
		sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_statistica where createtime > %v and createtime < %v;", start, end) //满意个数
		sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
		if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
			/*tmpmanyi := make(map[string]interface{})
			sss, _ := strconv.Atoi(tmpyesterday)
			tmpsjd := strconv.Itoa(i + 1)
			tmpsjd = tmpsjd + "月"
			tmpmanyi["date"] = tmpsjd
			tmpmanyi["videoTimes"] = sss
			tmprestlt2[i] = tmpmanyi*/
			tmpmanyi, _ = strconv.Atoi(tmpyesterday)
		}
		tmpmanyimap := make(map[string]interface{})
		tmpsjd := strconv.Itoa(i + 1)
		tmpsjd = tmpsjd + "月"
		tmpmanyimap["date"] = tmpsjd
		if tmpmanyi != 0 {

			tmpmanyimap["videoTimes"] = Decimal((float64(tmpmanyi) / (float64(tmpmanyi) + float64(tmpbumanyi))) * 100)

			fmt.Println(Decimal(float64(1135) / (float64(1135) + float64(243)) * 100))
		} else {
			tmpmanyimap["videoTimes"] = "0"
		}
		tmprestlt2[i] = tmpmanyimap
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"averageVideoTime": manyicount, "totalVideoTime": bumanyicount + manyicount, "details": tmprestlt2}})
}

func queryrightbottom(c *gin.Context) {
	var tmpday int
	var tmpmonth int
	var tmpweek int
	tmprestlt2 := make([]map[string]interface{}, 12)
	sqlcmd2 := fmt.Sprintf(" select count(*) as total from scc_workflow t1 inner join scc_apply t2 on t2.appid = t1.appid where t2.status=1 and t1.appnextnode=200 and (t1.createtime-t2.createtime<3600*24);") //不满意个数
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpday, _ = strconv.Atoi(tmpyesterday)
	}
	sqlcmd2 = fmt.Sprintf(" select count(*) as total from scc_workflow t1 inner join scc_apply t2 on t2.appid = t1.appid where t2.status=1 and t1.appnextnode=200 and (t1.createtime-t2.createtime<3600*24*7);") //不满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpweek, _ = strconv.Atoi(tmpyesterday)
	}
	sqlcmd2 = fmt.Sprintf(" select count(*)  as total from scc_workflow t1 inner join scc_apply t2 on t2.appid = t1.appid where t2.status=1 and t1.appnextnode=200 and (t1.createtime-t2.createtime<3600*24*30);") //不满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpmonth, _ = strconv.Atoi(tmpyesterday)
	}

	tmptimestr := make([]string, 24)
	tmptimestr[0] = "2020-01-01 00:00:00"
	tmptimestr[1] = "2020-01-31 23:59:59"
	tmptimestr[2] = "2020-02-01 00:00:00"
	tmptimestr[3] = "2020-02-29 23:59:59"
	tmptimestr[4] = "2020-03-01 00:00:00"
	tmptimestr[5] = "2020-03-31 23:59:59"
	tmptimestr[6] = "2020-04-01 00:00:00"
	tmptimestr[7] = "2020-04-30 23:59:59"
	tmptimestr[8] = "2020-05-01 00:00:00"
	tmptimestr[9] = "2020-05-31 23:59:59"
	tmptimestr[10] = "2020-06-01 00:00:00"
	tmptimestr[11] = "2020-06-30 23:59:59"
	tmptimestr[12] = "2020-07-01 00:00:00"
	tmptimestr[13] = "2020-07-31 23:59:59"
	tmptimestr[14] = "2020-08-01 00:00:00"
	tmptimestr[15] = "2020-08-31 23:59:59"
	tmptimestr[16] = "2020-09-01 00:00:00"
	tmptimestr[17] = "2020-09-30 23:59:59"
	tmptimestr[18] = "2020-10-01 00:00:00"
	tmptimestr[19] = "2020-10-31 23:59:59"
	tmptimestr[20] = "2020-11-01 00:00:00"
	tmptimestr[21] = "2020-11-30 23:59:59"
	tmptimestr[22] = "2020-12-01 00:00:00"
	tmptimestr[23] = "2020-12-31 23:59:59"
	for i := 0; i < 12; i++ {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", tmptimestr[i*2], time.Local)
		start := t.Unix()

		t1, _ := time.ParseInLocation("2006-01-02 15:04:05", tmptimestr[i*2+1], time.Local)
		end := t1.Unix()
		sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_workflow where  appnextnode =200 and (createtime>%v and createtime <%v);", start, end) //当月解决的个数
		sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
		if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
			tmpmanyi := make(map[string]interface{})
			sss, _ := strconv.Atoi(tmpyesterday)
			tmpsjd := strconv.Itoa(i + 1)
			tmpsjd = tmpsjd + "月"
			tmpmanyi["date"] = tmpsjd
			tmpmanyi["number"] = sss
			tmprestlt2[i] = tmpmanyi
		}

	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"totalNumber": tmpday, "grabDangerNumber": tmpweek, "covertNumber": tmpmonth, "statisticsList": tmprestlt2}})
}
func setrouter(r *gin.Engine) {
	r.POST("api/sign-call/realtime-online", querylefttopquestion)
	r.POST("api/sign-call/call-online", querylefttopmidquestion)
	r.POST("/api/pw-statistics/list", queryleftobottom)
	r.POST("/api/police-work/list", querycenterbottom)
	r.POST("/api/integral-rank/security", queryrightcenter)
	r.POST("/api/integral-rank/secfirm", queryrightcenterunsolve)
	r.POST("/api/video/statistics", queryrighttop)
	r.POST("/api/pw-statistics/graph", queryrightbottom)
}
