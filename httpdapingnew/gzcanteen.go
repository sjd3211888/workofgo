package gzcanteen

import (
	"fmt"
	sccsql "golearn/gomysql"
	sccredis "golearn/goredis"
	"net/http"
	"strconv"
	"strings"
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

func handlemouth(mounth string) int {
	switch mounth {
	case "January":
		{
			return 1
		}
	case "February":
		{
			return 2
		}
	case "March":
		{
			return 3
		}
	case "April":
		{
			return 4
		}
	case "May":
		{
			return 5
		}
	case "June":
		{
			return 6
		}
	case "July":
		{
			return 7
		}
	case "August":
		{
			return 8
		}
	case "September":
		{
			return 9
		}
	case "October":
		{
			return 10
		}
	case "November":
		{
			return 11
		}
	case "December":
		{
			return 12
		}
	}
	return 0
}
func Decimal(value float64) string {
	return (fmt.Sprintf("%.2f", value))
}

//不满意个数=轻微问题+一般问题+严重问题，但是[敏感词]]让轻微问题变成满意
//2.0
//【敏感词】又改了 现在0默认未定型 但是做问题统计，不做严重性统计 grade 0 未定性问题  1定性为轻微问题但是TM的竟然是基本满意 2定性为一般问题，统计到不满意中 3 定性为严重问题统计到不满意
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
	go func(Host string, Username string, Password string, Dbname string, Serhost string, Redisip string, iport int) {
		sccinfo.tmpsql.Initmysql(Host, Username, Password, Dbname, iport)
		sccinfo.scctmpsql.Initmysql(Host, Username, Password, "SCC", iport)
		sccinfo.tmpredis.Redisip = (Redisip)
		sccinfo.tmpredis.ConnectRedis()
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		r.Use(cors())
		setrouter(r)
		fmt.Println(Serhost)
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
func gzlefttop(c *gin.Context) {
	var json GZworkplace
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var totalworknum int     //总问题数
	var resolveed int        //已解决问题数
	var yesterdayproblem int //昨日新增问题数

	var lightunCompleted int   //轻微问题未解决数
	var normalunCompleted int  //一般问题未解决数
	var seriousunCompleted int //严重问题未解决数

	var lightCompleted int   //轻微问题数                                                                                                                                                   //轻微问题解决数
	var normalCompleted int  //一般问题数
	var seriousCompleted int //严重问题数

	var foodproblem int
	var serviceproblem int
	var environment int

	sqlcmd := fmt.Sprintf("select count(distinct appid ) as total from scc_workflow A inner join scc_worktempplate B on A.templateid = B.workid where B.department= %v;", json.Workplace) //总问题数
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	if tmpworknum, ok := sqlresult[0]["total"]; ok {
		sqlcmd1 := fmt.Sprintf("select count(distinct appid ) as total from scc_workflow A inner join scc_worktempplate B on A.templateid = B.workid where A.appnextnode=200 and B.department= %v", json.Workplace) //总解决问题数
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

	sqlcmd2 := fmt.Sprintf("select count(*) as total from scc_apply A inner join scc_worktempplate B on A.templateid=B.workid where A.createtime>'%v'  and A.createtime< '%v' and department = '%v';", yesterdattime, yesterdattime+3600*24, json.Workplace) //总解决问题数
	fmt.Println(sqlcmd2)

	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult3[0]["total"]; ok {
		yesterdayproblem, _ = strconv.Atoi(tmpyesterday)
	}
	//  0未定型问题  1轻微问题   2 一般问题  3严重问题
	//lightunCompleted normalunCompleted seriousunCompleted
	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid  where A.grade =1 and A.status=0 and B.department=%v;", json.Workplace) //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)

	for i := 0; i < len(sqlresult3); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult3[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			lightunCompleted++
		}
	}

	/*if tmplightunCompleted, ok := sqlresult3[0]["total"]; ok {
		normalunCompleted, _ = strconv.Atoi(tmplightunCompleted)
	}*/

	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid  where A.grade =2 and A.status=0 and B.department='%v';", json.Workplace) //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)

	for i := 0; i < len(sqlresult3); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult3[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			normalunCompleted++
		}
	}

	/*if tmpnormalunCompleted, ok := sqlresult3[0]["total"]; ok {
		seriousunCompleted, _ = strconv.Atoi(tmpnormalunCompleted)
	}*/

	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid  where A.grade =3 and A.status=0 and B.department='%v';", json.Workplace) //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)

	for i := 0; i < len(sqlresult3); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult3[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			seriousunCompleted++
		}
	}
	/*if tmpseriousunCompleted, ok := sqlresult3[0]["total"]; ok {
		lightunCompleted, _ = strconv.Atoi(tmpseriousunCompleted)
	}*/

	//  0未定型问题  1轻微问题   2 一般问题  3严重问题
	//lightCompleted normalCompleted seriousCompleted
	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid  where A.grade =1 and B.department ='%v';", json.Workplace) //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)

	for i := 0; i < len(sqlresult3); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult3[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			lightCompleted++
		}
	}

	/*if tmplightCompleted, ok := sqlresult3[0]["total"]; ok {
		normalCompleted, _ = strconv.Atoi(tmplightCompleted)
	}*/

	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid  where A.grade =2 and B.department ='%v'", json.Workplace) //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)

	for i := 0; i < len(sqlresult3); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult3[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			normalCompleted++
		}
	}
	/*if tmpenvironment, ok := sqlresult3[0]["total"]; ok {
		seriousCompleted, _ = strconv.Atoi(tmpenvironment)
	}*/

	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid  where A.grade =3 and B.department ='%v'", json.Workplace) //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)

	for i := 0; i < len(sqlresult3); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult3[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			seriousCompleted++
		}
	}
	/*if tmpseriousCompleted, ok := sqlresult3[0]["total"]; ok {
		lightCompleted, _ = strconv.Atoi(tmpseriousCompleted)
	}*/

	sqlcmd8 := fmt.Sprintf("select count(distinct appid ) as total from scc_workflow A inner join scc_worktempplate B on A.templateid = B.workid where B.department= %v and B.p_type = 10;", json.Workplace)
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd8)
	if tmpserviceproblem, ok := sqlresult3[0]["total"]; ok {
		serviceproblem, _ = strconv.Atoi(tmpserviceproblem)
	}

	sqlcmd8 = fmt.Sprintf("select count(distinct appid ) as total from scc_workflow A inner join scc_worktempplate B on A.templateid = B.workid where B.department= %v and B.p_type = 100;", json.Workplace)
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd8)
	if tmpfoodproblem, ok := sqlresult3[0]["total"]; ok {
		foodproblem, _ = strconv.Atoi(tmpfoodproblem)
	}

	sqlcmd8 = fmt.Sprintf("select count(distinct appid ) as total from scc_workflow A inner join scc_worktempplate B on A.templateid = B.workid where B.department= %v and B.p_type = 1;", json.Workplace)
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd8)
	if tmpserviceproblem, ok := sqlresult3[0]["total"]; ok {
		environment, _ = strconv.Atoi(tmpserviceproblem)
	}

	tmpleft1 := make(map[string]interface{})
	//tmpleft2 := make(map[string]interface{})
	tmpleft3 := make(map[string]interface{})
	//tmpleft4 := make(map[string]interface{})

	tmpleft1["currentproblem"] = totalworknum - resolveed
	tmpleft1["lastinsertproblem"] = yesterdayproblem
	tmpleft1["totalproblemsolved"] = resolveed

	lenged := make([]string, 2)
	lenged[0] = "问题总数"
	lenged[1] = "已解决数"
	question := make([]string, 3)
	question[0] = "轻微问题"
	question[1] = "一般问题"
	question[2] = "严重问题"

	problemcount := make([]int, 3)

	problemcount[0] = lightCompleted
	problemcount[1] = normalCompleted
	problemcount[2] = seriousCompleted

	problemsolved := make([]int, 3)
	problemsolved[0] = lightCompleted - lightunCompleted
	problemsolved[1] = normalCompleted - normalunCompleted
	problemsolved[2] = seriousCompleted - seriousunCompleted

	/*tmpleft2["minorproblem"] = lightCompleted
	tmpleft2["minorproblemsolved"] = lightCompleted - lightunCompleted
	tmpleft2["generalprobelm"] = normalCompleted
	tmpleft2["generalproblemsolved"] = normalCompleted - normalunCompleted
	tmpleft2["seriousproblem"] = seriousCompleted
	tmpleft2["seriousproblemsolved"] = seriousCompleted - seriousunCompleted*/

	tmpleft3["minorproblemsolved"] = lightCompleted - lightunCompleted
	tmpleft3["generalproblemsolved"] = normalCompleted - normalunCompleted
	tmpleft3["seriousproblemsolved"] = seriousCompleted - seriousunCompleted

	tmptype := make([]string, 3)
	tmptype[0] = "菜品质量"
	tmptype[1] = "服务态度"
	tmptype[2] = "环境卫生"

	tmpcount := make([]int, 3)
	tmpcount[0] = foodproblem
	tmpcount[1] = serviceproblem
	tmpcount[2] = environment

	/*tmpleft4["foodproblem"] = foodproblem
	tmpleft4["serviceproblem"] = serviceproblem
	tmpleft4["environmentalproblem"] = environment*/

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"left1": tmpleft1, "left2": gin.H{"lenged": lenged, "question": question, "problemcount": problemcount, "problemsolved": problemsolved}, "left3": tmpleft3, "left4": gin.H{"type": tmptype, "tmpcount": tmpcount}}})
}
func gzleftbottom(c *gin.Context) {
	var json GZworkplace
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var totalworknumbycustomer int
	var totalworknumbyself int
	var resolveedcustomer int
	var resolveedself int
	sqlcmd := fmt.Sprintf("select count(distinct appid ) as total from scc_workflow A inner join scc_worktempplate B on A.templateid = B.workid where B.department= %v and B.trade =1;", json.Workplace)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)

	if tmptotalworknumbycustomer, ok := sqlresult[0]["total"]; ok {
		totalworknumbycustomer, _ = strconv.Atoi(tmptotalworknumbycustomer)
	}

	sqlcmd = fmt.Sprintf("select count(distinct appid ) as total from scc_workflow A inner join scc_worktempplate B on A.templateid = B.workid where B.department= %v and B.trade =2;", json.Workplace)
	sqlresult = sccinfo.tmpsql.SelectData(sqlcmd)

	if tmptotalworknumbyself, ok := sqlresult[0]["total"]; ok {
		totalworknumbyself, _ = strconv.Atoi(tmptotalworknumbyself)
	}

	sqlcmd1 := fmt.Sprintf("select count(distinct appid ) as total from scc_workflow A inner join scc_worktempplate B on A.templateid = B.workid where A.appnextnode=200 and B.department= %v and B.trade =1;", json.Workplace) //总解决问题数
	sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
	if tmpresolveed, ok := sqlresult1[0]["total"]; ok {
		resolveedcustomer, _ = strconv.Atoi(tmpresolveed)
	}

	sqlcmd2 := fmt.Sprintf("select count(distinct appid ) as total from scc_workflow A inner join scc_worktempplate B on A.templateid = B.workid where A.appnextnode=200 and B.department= %v and B.trade =2;", json.Workplace) //总解决问题数
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpresolveed, ok := sqlresult2[0]["total"]; ok {
		resolveedself, _ = strconv.Atoi(tmpresolveed)
	}
	title := make([]interface{}, 4)
	/*title[0] = "投诉总数"                 //用户提问题的总数
	title[1] = totalworknumbycustomer //用户提问题的总数*/
	totalnumbycu := make(map[string]interface{})
	totalnumbycu["name"] = "投诉总数"
	totalnumbycu["value"] = totalworknumbycustomer
	/*title[2] = "已整改总数"            //用户提问题解决的总数
	title[3] = resolveedcustomer  //用户提问题解决的总数*/
	totalnumbycuso := make(map[string]interface{})
	totalnumbycuso["name"] = "已整改投诉总数"
	totalnumbycuso["value"] = resolveedcustomer
	/*title[4] = totalworknumbyself //自己人提问题的总数
	title[5] = "自查总数"             //自己人提问题的总数*/
	totalnumbyself := make(map[string]interface{})
	totalnumbyself["name"] = "自查总数"
	totalnumbyself["value"] = totalworknumbyself

	/*title[6] = "自查已整改总数"     //自己人提问题已解决的总数*
	title[7] = resolveedself //自己人提问题已解决的总数*/
	totalnumbyselfso := make(map[string]interface{})
	totalnumbyselfso["name"] = "已整改自查总数"
	totalnumbyselfso["value"] = resolveedself
	title[0] = totalnumbycu
	title[1] = totalnumbycuso
	title[2] = totalnumbyself
	title[3] = totalnumbyselfso
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": title})
}
func gzmiddle(c *gin.Context) {
	var json GZworkplace
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tmpresult := make([]map[string]interface{}, 0)
	sqlcmd2 := fmt.Sprintf("select A.appid,A.createtime from scc_workflow A inner join scc_worktempplate B on A.templateid = B.workid where A.appnextnode = 200 and B.department= '%v' order by A.createtime desc limit 0,8;", json.Workplace) //最近解决的8个
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	for i := 0; i < len(sqlresult2); i++ {
		tmpproblem := make(map[string]interface{})
		sqlcmd3 := fmt.Sprintf("select textinfo,filepath from scc_apply where appid = %v;", sqlresult2[i]["appid"]) //最近解决的8个
		sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd3)
		tmpproblem["probleminfo"] = sqlresult3[0]["textinfo"]
		problempathinfo := strings.Split(sqlresult3[0]["filepath"], ",")

		if len(problempathinfo) == 1 {
			if 5 > len(problempathinfo[0]) {
				tmppro := make([]string, 0)
				tmpproblem["problempath"] = tmppro
				fmt.Println("yyyyyyyyyyyyyyyyyyyy", problempathinfo, len(problempathinfo))
			} else {
				tmpproblem["problempath"] = problempathinfo
				fmt.Println("ssssssssssssssssssssss", problempathinfo, len(problempathinfo))
			}
		} else {
			tmpproblem["problempath"] = problempathinfo
			fmt.Println("zzzzzzzzzzzzzzzzz", problempathinfo, len(problempathinfo))
		}
		//tmpproblem["problempath"] = problempathinfo
		sqlcmd3 = fmt.Sprintf("select max(id),advise,filepath from scc_workflow where appid=%v and appnextnode=2;", sqlresult2[i]["appid"]) //最近解决的8个
		sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd3)
		tmpproblem["rectificationresult"] = sqlresult3[0]["advise"]
		solvepathinfo := strings.Split(sqlresult3[0]["filepath"], ",")

		if len(solvepathinfo) == 1 {
			if 5 > len(solvepathinfo[0]) {
				tmppro := make([]string, 0)
				tmpproblem["rectificationpath"] = tmppro
				fmt.Println("yyyyyyyyyyyyyyyyyyyy", solvepathinfo, len(solvepathinfo))
			} else {
				tmpproblem["rectificationpath"] = solvepathinfo
				fmt.Println("ssssssssssssssssssssss", solvepathinfo, len(solvepathinfo))
			}
		} else {
			tmpproblem["rectificationpath"] = solvepathinfo
			fmt.Println("zzzzzzzzzzzzzzzzz", solvepathinfo, len(solvepathinfo))
		}
		/*solutioninfo
		solutionpath
		rectificationresult
		rectificationpath*/
		//tmpproblem["solutionpath"] = solvepathinfo
		sqlcmd3 = fmt.Sprintf("select max(id),advise,filepath from scc_workflow where appid= %v and appnextnode=102;", sqlresult2[i]["appid"]) //最近解决的8个
		sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd3)
		tmpproblem["solutioninfo"] = sqlresult3[0]["advise"]
		rectificationpathinfo := strings.Split(sqlresult3[0]["filepath"], ",")
		//fmt.Println(sqlresult3[0]["filepath"])
		if len(rectificationpathinfo) == 1 {
			if 5 > len(rectificationpathinfo[0]) {
				tmppro := make([]string, 0)
				tmpproblem["solutionpath"] = tmppro
				//fmt.Println("yyyyyyyyyyyyyyyyyyyy", rectificationpathinfo, len(rectificationpathinfo))
			} else {
				tmpproblem["solutionpath"] = rectificationpathinfo
				//fmt.Println("ssssssssssssssssssssss", rectificationpathinfo, len(rectificationpathinfo))
			}
		} else {
			tmpproblem["solutionpath"] = rectificationpathinfo
			//fmt.Println("zzzzzzzzzzzzzzzzz", rectificationpathinfo, len(rectificationpathinfo))
		}
		/*tmppro := make([]string, 0)
		tmpproblem["rectificationpath"] = tmppro
		fmt.Println(tmpproblem)*/
		tmpresult = append(tmpresult, tmpproblem)
	}
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": tmpresult})
}
func gzrightop(c *gin.Context) {
	var json GZworkplace
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tmpmounthsjd := make([]int, 6)
	percent := make([]string, 6)
	tmpleft1 := make(map[string]interface{})
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	fmt.Println(currentMonth)
	currentmonth := handlemouth(currentMonth.String())
	currentstart := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	currentlast := currentstart.AddDate(0, 1, -1)
	//fmt.Println(currentstart.Unix())
	//所有问题，但是又可能又未定性的  求出定性的问题 也就是被审批过的
	sqlcmd2 := fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid where A.createtime > %v and A.createtime < %v and B.department = '%v' and B.trade = 1 and A.grade!=1;", currentstart.Unix(), currentlast.Unix(), json.Workplace) //不满意个数=轻微问题+一般问题+严重问题，但是傻逼让轻微问题变成满意
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)

	sqlcmd8 := fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid where A.createtime > %v and A.createtime < %v and B.department = '%v' and B.trade = 1  and A.grade=1;", currentstart.Unix(), currentlast.Unix(), json.Workplace) //不满意个数=轻微问题+一般问题+严重问题，但是傻逼让轻微问题变成满意
	sqlresult8 := sccinfo.tmpsql.SelectData(sqlcmd8)

	var tmpmanyi int
	var tmpbumanyi int
	var sbproblem int

	for i := 0; i < len(sqlresult2); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult2[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			tmpbumanyi++
		}
	}
	for i := 0; i < len(sqlresult8); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult8[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			sbproblem++
		}
	}
	/*if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpbumanyi, _ = strconv.Atoi(tmpyesterday)
	}*/

	/*if tmpsbproblem, ok := sqlresult8[0]["total"]; ok {
		sbproblem, _ = strconv.Atoi(tmpsbproblem)
	}*/
	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_statistica where createtime > %v and createtime < %v;", currentstart.Unix(), currentlast.Unix()) //满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)

	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpmanyi, _ = strconv.Atoi(tmpyesterday)
	}
	tmpmounth := strconv.Itoa(currentmonth)
	tmpmounthsjd[4] = currentmonth
	if tmpmanyi != 0 {

		tmpleft1[tmpmounth] = Decimal(((float64(tmpmanyi) + float64(sbproblem)) / (float64(tmpmanyi) + float64(tmpbumanyi) + float64(sbproblem))) * 100)
		percent[4] = Decimal(((float64(tmpmanyi) + float64(sbproblem)) / (float64(tmpmanyi) + float64(tmpbumanyi) + float64(sbproblem))) * 100)
		//fmt.Println(Decimal(float64(1135) / (float64(1135) + float64(243)) * 100))
	} else {
		tmpleft1[tmpmounth] = "0"
		percent[4] = "0"
	}

	last1start := time.Date(currentYear, currentMonth-1, 1, 0, 0, 0, 0, currentLocation)
	last1last := last1start.AddDate(0, 1, -1)

	sqlcmd2 = fmt.Sprintf("select A.appid  from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid where (A.createtime > %v and A.createtime < %v) and B.department = '%v' and B.trade = 1 and A.grade!=1;", last1start.Unix(), last1last.Unix(), json.Workplace) //不满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)

	for i := 0; i < len(sqlresult2); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult2[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			tmpbumanyi++
		}
	}
	/*if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpbumanyi, _ = strconv.Atoi(tmpyesterday)
	}*/

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_statistica where (createtime > %v and createtime < %v);", last1start.Unix(), last1last.Unix()) //满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)

	sqlcmd8 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid where A.createtime > %v and A.createtime < %v and B.department = '%v' and B.trade = 1  and A.grade=1;", currentstart.Unix(), currentlast.Unix(), json.Workplace) //不满意个数=轻微问题+一般问题+严重问题，但是傻逼让轻微问题变成满意
	sqlresult8 = sccinfo.tmpsql.SelectData(sqlcmd8)
	for i := 0; i < len(sqlresult8); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult8[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			sbproblem++
		}
	}
	/*if tmpsbproblem, ok := sqlresult8[0]["total"]; ok {
		sbproblem, _ = strconv.Atoi(tmpsbproblem)
	}*/
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpmanyi, _ = strconv.Atoi(tmpyesterday)
	}
	tmpsjd := currentmonth - 1
	if tmpsjd <= 0 {
		tmpsjd = tmpsjd + 12
	}
	tmpmounthlast1 := strconv.Itoa(tmpsjd)
	tmpmounthsjd[3] = tmpsjd
	fmt.Println(tmpmanyi, tmpmounth, tmpbumanyi)
	if tmpmanyi != 0 {

		tmpleft1[tmpmounthlast1] = Decimal(((float64(tmpmanyi) + float64(sbproblem)) / (float64(tmpmanyi) + float64(tmpbumanyi) + float64(sbproblem))) * 100)
		percent[3] = Decimal(((float64(tmpmanyi) + float64(sbproblem)) / (float64(tmpmanyi) + float64(tmpbumanyi) + float64(sbproblem))) * 100)
		fmt.Println(Decimal(float64(1135) / (float64(1135) + float64(243)) * 100))
	} else {
		tmpleft1[tmpmounthlast1] = "0"
		percent[3] = "0"
	}

	last2start := time.Date(currentYear, currentMonth-2, 1, 0, 0, 0, 0, currentLocation)
	last2last := last2start.AddDate(0, 1, -1)

	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid where A.createtime > %v and A.createtime < %v and B.department = '%v' and B.trade = 1 and A.grade!=1;", last2start.Unix(), last2last.Unix(), json.Workplace) //不满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	//var tmpmanyi int
	//var tmpbumanyi int
	for i := 0; i < len(sqlresult2); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult2[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			tmpbumanyi++
		}
	}
	/*if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpbumanyi, _ = strconv.Atoi(tmpyesterday)
	}*/
	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_statistica where createtime > %v and createtime < %v;", last2start.Unix(), last2last.Unix()) //满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpmanyi, _ = strconv.Atoi(tmpyesterday)
	}
	sqlcmd8 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid where A.createtime > %v and A.createtime < %v and B.department = '%v' and B.trade = 1  and A.grade=1;", currentstart.Unix(), currentlast.Unix(), json.Workplace) //不满意个数=轻微问题+一般问题+严重问题，但是傻逼让轻微问题变成满意
	sqlresult8 = sccinfo.tmpsql.SelectData(sqlcmd8)
	for i := 0; i < len(sqlresult8); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult8[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			sbproblem++
		}
	}
	/*if tmpsbproblem, ok := sqlresult8[0]["total"]; ok {
		sbproblem, _ = strconv.Atoi(tmpsbproblem)
	}*/
	tmpsjd = currentmonth - 2
	if tmpsjd <= 0 {
		tmpsjd = tmpsjd + 12
	}
	tmpmounthlast2 := strconv.Itoa(tmpsjd)
	tmpmounthsjd[2] = tmpsjd
	if tmpmanyi != 0 {

		tmpleft1[tmpmounthlast2] = Decimal(((float64(tmpmanyi) + float64(sbproblem)) / (float64(tmpmanyi) + float64(tmpbumanyi) + float64(sbproblem))) * 100)
		percent[2] = Decimal(((float64(tmpmanyi) + float64(sbproblem)) / (float64(tmpmanyi) + float64(tmpbumanyi) + float64(sbproblem))) * 100)
		//fmt.Println(Decimal(float64(1135) / (float64(1135) + float64(243)) * 100))
	} else {
		tmpleft1[tmpmounthlast2] = "0"
		percent[2] = "0"
	}

	last3start := time.Date(currentYear, currentMonth-3, 1, 0, 0, 0, 0, currentLocation)
	last3last := last3start.AddDate(0, 1, -1)

	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid where A.createtime > %v and A.createtime < %v and B.department = '%v' and B.trade = 1 and A.grade!=1;", last3start.Unix(), last3last.Unix(), json.Workplace) //不满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	//var tmpmanyi int
	//var tmpbumanyi int
	for i := 0; i < len(sqlresult2); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult2[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			tmpbumanyi++
		}
	}
	/*if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpbumanyi, _ = strconv.Atoi(tmpyesterday)
	}*/
	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_statistica where createtime > %v and createtime < %v;", last3start.Unix(), last3last.Unix()) //满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpmanyi, _ = strconv.Atoi(tmpyesterday)
	}
	sqlcmd8 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid where A.createtime > %v and A.createtime < %v and B.department = '%v' and B.trade = 1  and A.grade=1;", currentstart.Unix(), currentlast.Unix(), json.Workplace) //不满意个数=轻微问题+一般问题+严重问题，但是傻逼让轻微问题变成满意
	sqlresult8 = sccinfo.tmpsql.SelectData(sqlcmd8)
	for i := 0; i < len(sqlresult8); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult8[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			sbproblem++
		}
	}
	/*if tmpsbproblem, ok := sqlresult8[0]["total"]; ok {
		sbproblem, _ = strconv.Atoi(tmpsbproblem)
	}*/
	tmpsjd = currentmonth - 3
	if tmpsjd <= 0 {
		tmpsjd = tmpsjd + 12
	}
	tmpmounthlast3 := strconv.Itoa(tmpsjd)
	tmpmounthsjd[1] = tmpsjd
	if tmpmanyi != 0 {

		tmpleft1[tmpmounthlast3] = Decimal(((float64(tmpmanyi) + float64(sbproblem)) / (float64(tmpmanyi) + float64(tmpbumanyi) + float64(sbproblem))) * 100)
		percent[1] = Decimal(((float64(tmpmanyi) + float64(sbproblem)) / (float64(tmpmanyi) + float64(tmpbumanyi) + float64(sbproblem))) * 100)
		//fmt.Println(Decimal(float64(1135) / (float64(1135) + float64(243)) * 100))
	} else {
		tmpleft1[tmpmounthlast3] = "0"

		percent[1] = "0"
	}

	last4start := time.Date(currentYear, currentMonth-4, 1, 0, 0, 0, 0, currentLocation)
	last4last := last4start.AddDate(0, 1, -1)

	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid where A.createtime > %v and A.createtime < %v and B.department = '%v' and B.trade = 1 and A.grade!=1;", last4start.Unix(), last4last.Unix(), json.Workplace) //不满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	//var tmpmanyi int
	//var tmpbumanyi int

	for i := 0; i < len(sqlresult2); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult2[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			tmpbumanyi++
		}
	}
	/*if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpbumanyi, _ = strconv.Atoi(tmpyesterday)
	}*/
	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_statistica where createtime > %v and createtime < %v;", last4start.Unix(), last4last.Unix()) //满意个数
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpmanyi, _ = strconv.Atoi(tmpyesterday)
	}
	sqlcmd8 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid where A.createtime > %v and A.createtime < %v and B.department = '%v' and B.trade = 1  and A.grade=1;", currentstart.Unix(), currentlast.Unix(), json.Workplace) //不满意个数=轻微问题+一般问题+严重问题，但是傻逼让轻微问题变成满意
	sqlresult8 = sccinfo.tmpsql.SelectData(sqlcmd8)

	for i := 0; i < len(sqlresult8); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult8[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			sbproblem++
		}
	}
	/*if tmpsbproblem, ok := sqlresult8[0]["total"]; ok {
		sbproblem, _ = strconv.Atoi(tmpsbproblem)
	}*/
	tmpsjd = currentmonth - 4
	if tmpsjd <= 0 {
		tmpsjd = tmpsjd + 12
	}
	tmpmounthlast4 := strconv.Itoa(tmpsjd)
	tmpmounthsjd[0] = tmpsjd
	if tmpmanyi != 0 {

		tmpleft1[tmpmounthlast4] = Decimal(((float64(tmpmanyi) + float64(sbproblem)) / (float64(tmpmanyi) + float64(tmpbumanyi) + float64(sbproblem))) * 100)
		percent[0] = Decimal(((float64(tmpmanyi) + float64(sbproblem)) / (float64(tmpmanyi) + float64(tmpbumanyi) + float64(sbproblem))) * 100)
		//fmt.Println(Decimal(float64(1135) / (float64(1135) + float64(243)) * 100))
	} else {
		tmpleft1[tmpmounthlast4] = "0"

		percent[0] = "0"
	}
	tmpsjd = currentmonth + 1
	if tmpsjd > 12 {
		tmpsjd = tmpsjd - 12
	}
	tmpmounthnext := strconv.Itoa(tmpsjd)
	tmpmounthsjd[5] = tmpsjd
	//percent[5] = "0"
	tmpleft1[tmpmounthnext] = "0"
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"mounth": tmpmounthsjd, "percent": percent}})
}
func gzrightmid(c *gin.Context) {
	var json GZworkplace
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tmplevel := make([]string, 4)
	tmplevel[0] = "满意"
	tmplevel[1] = "基本满意"
	tmplevel[2] = "不满意"
	tmplevel[3] = "非常不满意"

	amount := make([]int, 4)
	var best int
	var better int
	var bad int
	var verybad int
	var undefined int

	//  0是未定型问题  1轻微问题   2 一般问题  3 严重问题
	// bad verybad better
	sqlcmd2 := fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid  where A.grade =2 and B.department ='%v';", json.Workplace)
	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd2)

	for i := 0; i < len(sqlresult3); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult3[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			bad++
		}
	}
	/*if tmplightCompleted, ok := sqlresult3[0]["total"]; ok {
		bad, _ = strconv.Atoi(tmplightCompleted)
	}*/

	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid  where A.grade =3 and B.department ='%v'", json.Workplace)
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)

	for i := 0; i < len(sqlresult3); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult3[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			verybad++
		}
	}
	/*if tmpenvironment, ok := sqlresult3[0]["total"]; ok {
		verybad, _ = strconv.Atoi(tmpenvironment)
	}*/

	sqlcmd2 = fmt.Sprintf("select A.appid from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid  where A.grade =1 and B.department ='%v'", json.Workplace)
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)

	for i := 0; i < len(sqlresult3); i++ {
		sqlcmd9 := fmt.Sprintf("select appnextnode from scc_workflow where appid= '%v' order by id desc limit 1", sqlresult3[i]["appid"]) //总解决问题数
		sqlresult9 := sccinfo.tmpsql.SelectData(sqlcmd9)
		tmpappnextnode, _ := strconv.Atoi(sqlresult9[0]["appnextnode"])
		if tmpappnextnode > 2 {
			better++
		}
	}

	/*if tmpseriousCompleted, ok := sqlresult3[0]["total"]; ok {
		better, _ = strconv.Atoi(tmpseriousCompleted)
	}*/

	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_statistica") //满意个数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult3[0]["total"]; ok {
		best, _ = strconv.Atoi(tmpyesterday)
	}

	amount[0] = best
	amount[1] = better
	amount[2] = bad
	amount[3] = verybad

	//未定性问题
	sqlcmd2 = fmt.Sprintf("select count(*) as total from scc_apply A inner join scc_worktempplate B on A.templateid = B.workid  where A.grade =0 and B.department ='%v'", json.Workplace) //总解决问题数
	sqlresult3 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpundefined, ok := sqlresult3[0]["total"]; ok {
		undefined, _ = strconv.Atoi(tmpundefined)
	}

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"level": tmplevel, "amount": amount, "undefined": undefined}})
}
func gzrightbuttom(c *gin.Context) {
	var json GZworkplace
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tmpmounthsjd := make([]int, 6)
	percent := make([]string, 6)
	var tmpday int
	var tmpmonth int
	var tmpweek int
	sqlcmd2 := fmt.Sprintf(" select count(*) as total from scc_workflow t1 inner join scc_apply t2 on t2.appid = t1.appid inner join scc_worktempplate t3 on t3.workid=t2.templateid where t2.status=1 and t1.appnextnode=200 and t3.department = %v and (t1.createtime-t2.createtime<3600*24);", json.Workplace) //一天内解决问题
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpday, _ = strconv.Atoi(tmpyesterday)
	}
	sqlcmd2 = fmt.Sprintf(" select count(*) as total from scc_workflow t1 inner join scc_apply t2 on t2.appid = t1.appid inner join scc_worktempplate t3 on t3.workid=t2.templateid where t2.status=1 and t1.appnextnode=200 and t3.department = %v and (t1.createtime-t2.createtime<3600*24*7);", json.Workplace) //一周内解决问题
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpweek, _ = strconv.Atoi(tmpyesterday)
	}
	sqlcmd2 = fmt.Sprintf(" select count(*)  as total from scc_workflow t1 inner join scc_apply t2 on t2.appid = t1.appid inner join scc_worktempplate t3 on t3.workid=t2.templateid where t2.status=1 and t1.appnextnode=200 and t3.department = %v and (t1.createtime-t2.createtime<3600*24*30);", json.Workplace) //一月内解决问题
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	if tmpyesterday, ok := sqlresult2[0]["total"]; ok {
		tmpmonth, _ = strconv.Atoi(tmpyesterday)
	}

	tmpleft1 := make(map[string]interface{})
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	fmt.Println(currentMonth)
	currentmonth := handlemouth(currentMonth.String())
	currentstart := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	currentlast := currentstart.AddDate(0, 1, -1)
	//fmt.Println(currentstart.Unix())
	sqlcmd2 = fmt.Sprintf("select A.createtime-B.createtime as solvetime from scc_workflow A inner join scc_apply B on A.appid=B.appid inner join scc_worktempplate C on C.workid=B.templateid where C.department = %v and (A.createtime>%v and A.createtime<%v) and A.appnextnode=200;", json.Workplace, currentstart.Unix(), currentlast.Unix())
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)

	//currentmonth := handlemouth(currentMonth.String())

	tmpmounth := strconv.Itoa(currentmonth)
	tmpmounthsjd[4] = currentmonth
	tmplen := len(sqlresult2)
	alltime := 0
	if 0 == tmplen {
		tmpleft1[tmpmounth] = "0"
		percent[4] = "0"
	} else {
		//tmpleft1[tmpmounth] = sqlresult2[0]["solvetime"]
		//percent[4] = sqlresult2[0]["solvetime"]

		for i := 0; i < tmplen; i++ {
			tmponetime, _ := strconv.Atoi(sqlresult2[i]["solvetime"])
			alltime = alltime + tmponetime
		}
		//tmpleft1[tmpmounthlast1] = strconv.Itoa(alltime / 3600 / tmplen)
		percent[4] = strconv.Itoa(alltime / 3600 / tmplen)
	}

	last1start := time.Date(currentYear, currentMonth-1, 1, 0, 0, 0, 0, currentLocation) //上个月
	last1last := last1start.AddDate(0, 1, -1)

	sqlcmd2 = fmt.Sprintf("select A.createtime-B.createtime as solvetime from scc_workflow A inner join scc_apply B on A.appid=B.appid inner join scc_worktempplate C on C.workid=B.templateid where C.department = %v and (A.createtime>%v and A.createtime<%v) and A.appnextnode=200;", json.Workplace, last1start.Unix(), last1last.Unix())
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)
	fmt.Println(sqlcmd2)
	tmpsjd := currentmonth - 1
	if tmpsjd <= 0 {
		tmpsjd = tmpsjd + 12
	}
	tmpmounthlast1 := strconv.Itoa(tmpsjd)
	tmpmounthsjd[3] = tmpsjd
	tmplen = len(sqlresult2)
	alltime = 0
	if 0 == tmplen {
		tmpleft1[tmpmounthlast1] = "0"
		percent[3] = "0"
	} else {
		for i := 0; i < tmplen; i++ {
			tmponetime, _ := strconv.Atoi(sqlresult2[i]["solvetime"])
			alltime = alltime + tmponetime
		}
		tmpleft1[tmpmounthlast1] = strconv.Itoa(alltime / 3600 / tmplen)
		percent[3] = strconv.Itoa(alltime / 3600 / tmplen)
	}

	last2start := time.Date(currentYear, currentMonth-2, 1, 0, 0, 0, 0, currentLocation) //上2个月
	last2last := last2start.AddDate(0, 1, -1)

	sqlcmd2 = fmt.Sprintf("select A.createtime-B.createtime as solvetime from scc_workflow A inner join scc_apply B on A.appid=B.appid inner join scc_worktempplate C on C.workid=B.templateid where C.department = %v and (A.createtime>%v and A.createtime<%v) and A.appnextnode=200;", json.Workplace, last2start.Unix(), last2last.Unix())
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)

	tmpsjd = currentmonth - 2
	if tmpsjd <= 0 {
		tmpsjd = tmpsjd + 12
	}
	tmpmounthlast2 := strconv.Itoa(tmpsjd)
	tmpmounthsjd[2] = tmpsjd
	tmplen = len(sqlresult2)
	if 0 == tmplen {
		tmpleft1[tmpmounthlast2] = "0"
		percent[2] = "0"
	} else {
		for i := 0; i < tmplen; i++ {
			tmponetime, _ := strconv.Atoi(sqlresult2[i]["solvetime"])
			alltime = alltime + tmponetime
		}
		tmpleft1[tmpmounthlast2] = strconv.Itoa(alltime / 3600 / tmplen)
		percent[2] = strconv.Itoa(alltime / 3600 / tmplen)
	}

	last3start := time.Date(currentYear, currentMonth-3, 1, 0, 0, 0, 0, currentLocation) //上3个月
	last3last := last3start.AddDate(0, 1, -1)

	sqlcmd2 = fmt.Sprintf("select A.createtime-B.createtime as solvetime from scc_workflow A inner join scc_apply B on A.appid=B.appid inner join scc_worktempplate C on C.workid=B.templateid where C.department = %v and (A.createtime>%v and A.createtime<%v) and A.appnextnode=200;", json.Workplace, last3start.Unix(), last3last.Unix())
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)

	tmpsjd = currentmonth - 3
	if tmpsjd <= 0 {
		tmpsjd = tmpsjd + 12
	}
	tmpmounthlast3 := strconv.Itoa(tmpsjd)
	tmpmounthsjd[1] = tmpsjd
	tmplen = len(sqlresult2)
	if 0 == tmplen {
		tmpleft1[tmpmounthlast3] = "0"
		percent[1] = "0"
	} else {
		for i := 0; i < tmplen; i++ {
			tmponetime, _ := strconv.Atoi(sqlresult2[i]["solvetime"])
			alltime = alltime + tmponetime
		}
		tmpleft1[tmpmounthlast3] = strconv.Itoa(alltime / 3600 / tmplen)
		percent[1] = strconv.Itoa(alltime / 3600 / tmplen)
	}
	last4start := time.Date(currentYear, currentMonth-4, 1, 0, 0, 0, 0, currentLocation) //上4个月
	last4last := last4start.AddDate(0, 1, -1)

	sqlcmd2 = fmt.Sprintf("select A.createtime-B.createtime as solvetime from scc_workflow A inner join scc_apply B on A.appid=B.appid inner join scc_worktempplate C on C.workid=B.templateid where C.department = %v and (A.createtime>%v and A.createtime<%v) and A.appnextnode=200;", json.Workplace, last4start.Unix(), last4last.Unix())
	sqlresult2 = sccinfo.tmpsql.SelectData(sqlcmd2)

	tmpsjd = currentmonth - 4
	if tmpsjd <= 0 {
		tmpsjd = tmpsjd + 12
	}
	tmpmounthlast4 := strconv.Itoa(tmpsjd)
	tmpmounthsjd[0] = tmpsjd
	tmplen = len(sqlresult2)
	if 0 == tmplen {
		tmpleft1[tmpmounthlast4] = "0"
		percent[0] = "0"
	} else {
		for i := 0; i < tmplen; i++ {
			tmponetime, _ := strconv.Atoi(sqlresult2[i]["solvetime"])
			alltime = alltime + tmponetime
		}
		tmpleft1[tmpmounthlast4] = strconv.Itoa(alltime / 3600 / tmplen)
		percent[0] = strconv.Itoa(alltime / 3600 / tmplen)
	}

	tmpsjd = currentmonth + 1
	if tmpsjd > 12 {
		tmpsjd = tmpsjd - 12
	}
	tmpmounthnext := strconv.Itoa(tmpsjd)
	tmpmounthsjd[5] = tmpsjd
	tmpleft1[tmpmounthnext] = "0"
	//percent[5] = "0"

	solveproblem := make([]string, 3)
	solveproblem[0] = "一天内解决"
	solveproblem[1] = "一周内解决"
	solveproblem[2] = "一月内解决"

	amount := make([]int, 3)
	amount[0] = tmpday
	amount[1] = tmpweek
	amount[2] = tmpmonth
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"mounth": tmpmounthsjd, "solvetime": percent, "solvedproblem": solveproblem, "amount": amount}})
}
func setrouter(r *gin.Engine) {
	r.POST("/lefttop", gzlefttop)
	r.POST("/leftbottom", gzleftbottom)
	r.POST("/middle", gzmiddle)
	r.POST("/rightop", gzrightop)
	r.POST("/rightmid", gzrightmid)
	r.POST("/rightbottom", gzrightbuttom) //未开发完成  先查这个月解决的问题，在for循环查询这个问题创建的时间 一减

}
