package commtask

import (
	"fmt"
	sccsql "golearn/gomysql"
	sccredis "golearn/goredis"
	"net/http"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

type coreinfo struct {
	tmpsql sccsql.Mysqlconnectpool
	//scctmpsql sccsql.Mysqlconnectpool
	tmpredis sccredis.Redisconnectpool
}

var sccinfo coreinfo

func reverse(s []map[string]string) []map[string]string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
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
	Host := conf["commwork"]["Host"]
	Username := conf["commwork"]["Username"]
	Password := conf["commwork"]["Password"]
	Dbname := conf["commwork"]["Dbname"]
	Port := conf["commwork"]["Port"]
	iport, _ := strconv.Atoi(Port)
	Serhost := conf["commwork"]["Httpserverhost"]
	Redisip := conf["commwork"]["Redisip"]
	go func(Host string, Username string, Password string, Dbname string, Serhost string, Redisip string, iport int) {
		sccinfo.tmpsql.Initmysql(Host, Username, Password, Dbname, iport)
		//sccinfo.scctmpsql.Initmysql(Host, Username, Password, "SCC", iport)
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
func scccreatecommtask(c *gin.Context) {
	var json CreateCommtask
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//拼接数据库语句
	taskstatus := 0
	if json.Advanced.Taskapprover == "" {
		taskstatus = 1
	}
	fmt.Println("1111")
	sqlcmd1 := fmt.Sprintf("insert into scc_commontask (taskname,createtime,executor,creater,textinfo,filepath,taskstatus)values('%v','%v','%v','%v','%v','%v','%v');", json.Taskname, time.Now().Unix(), json.Executor, json.Creater, json.Textinfo, json.Filepath, taskstatus)
	taskid := sccinfo.tmpsql.Insertmql(sqlcmd1)
	//fmt.Println(sqlcmd1)
	//添加抄送人
	for i := 0; i < len(json.Cc); i++ {
		sqlcmd1 = fmt.Sprintf("insert into scc_commontaskcc (taskid,cc,cctime)values('%v','%v','%v');", taskid, json.Cc[i].CcUsers, 0)
		_ = sccinfo.tmpsql.Insertmql(sqlcmd1)
	}
	fmt.Println("3333")
	//添加更多配置表
	if "" != json.Advanced.Taskbegintime || "" != json.Advanced.Taskendtime || "" != json.Advanced.Tasktype || "" != json.Advanced.Taskegrade || "" != json.Advanced.Taskapprover {

		taskbegintime, _ := strconv.ParseInt(json.Advanced.Taskbegintime, 10, 64)
		taskendtime, _ := strconv.ParseInt(json.Advanced.Taskendtime, 10, 64)
		tasktpye, _ := strconv.ParseInt(json.Advanced.Tasktype, 10, 64)
		taskgrade, _ := strconv.ParseInt(json.Advanced.Taskegrade, 10, 64)
		taskapprover, _ := strconv.ParseInt(json.Advanced.Taskapprover, 10, 64)
		sqlcmd1 = fmt.Sprintf("insert into scc_taskadvanced (taskid,taskbegintime,taskendtime,tasktype,taskgrade,taskapprover)values('%v','%v','%v','%v','%v','%v');", taskid, taskbegintime, taskendtime, tasktpye, taskgrade, taskapprover)
		_ = sccinfo.tmpsql.Insertmql(sqlcmd1)
	}
	fmt.Println(sqlcmd1)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"taskid": taskid}})
}

func scccdocommtask(c *gin.Context) {
	var json Docommtask
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tmpnow := time.Now().Unix()
	//拼接数据库语句
	sqlcmd1 := fmt.Sprintf("insert into scc_commontaskdoinfo (taskid,createtime,executor,textinfo,filepath)values('%v','%v','%v','%v','%v');", json.Taskid, tmpnow, json.Executor, json.Textinfo, json.Filepath)
	_ = sccinfo.tmpsql.Insertmql(sqlcmd1)
	fmt.Println(sqlcmd1)
	//在查下创建的任务表，对比下时间，然后更新下看看任务超时 还是任务完成 还是完成后需要审批
	taskstatus := 3 //先治状态为结束
	sqlcmd1 = fmt.Sprintf("select taskendtime,taskapprover from scc_taskadvanced where taskid = '%v'", json.Taskid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)
	if 0 != len(sqlresult) {
		//if "0" != sqlresult[0]["taskendtime"] {
		tmpendtime, _ := strconv.ParseInt(sqlresult[0]["taskendtime"], 10, 64)
		if tmpendtime <= tmpnow && "0" != sqlresult[0]["taskendtime"] {
			taskstatus = taskstatus + 10 //更新状态为超时
		}
		//}
	}

	//更新task表
	sqlcmd1 = fmt.Sprintf("update scc_commontask set taskstatus = '%v' where taskid = '%v'", taskstatus, json.Taskid)
	sccinfo.tmpsql.Execsqlcmd(sqlcmd1, false)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"taskstaus": taskstatus}})
}

func sccquerycommtaskto(c *gin.Context) {
	var json Tasker
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commontask A left join scc_taskadvanced B on A.taskid=B.taskid left join scc_commontaskdoinfo C on A.taskid=C.taskid where A.executor ='%v' and A.taskstatus = '%v'", json.Taskuser, 1)
	sqlresult4 := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var icount int
	if _, ok := sqlresult4[0]["total"]; ok {
		count, _ := strconv.Atoi(sqlresult4[0]["total"])
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
					//int tmpnumberperpage = numberperpage;
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}
			sqlcmd1 = fmt.Sprintf("select A.taskid,A.taskname,A.creater,A.createtime,A.executor,A.taskstatus,B.taskbegintime,B.taskendtime,B.taskgrade,B.tasktype,C.createtime as taskexectime from scc_commontask A left join scc_taskadvanced B on A.taskid=B.taskid left join scc_commontaskdoinfo C on A.taskid=C.taskid where A.executor ='%v' and A.taskstatus = '%v' order by A.taskid limit %v,%v", json.Taskuser, 1, fromcount, numberperpage)
			sqlresult = sccinfo.tmpsql.SelectData(sqlcmd1)
			fmt.Println(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)

	/*sqlcmd1 := fmt.Sprintf("select A.taskname,A.createtime,A.executor,A.taskstatus,B.taskbegintime,B.taskendtime,C.createtime from scc_commontask A inner join scc_taskadvanced B on A.taskid=B.taskid inner join scc_commontaskdoinfo C on A.taskid=C.taskid where A.executor ='%v' and A.taskstatus = '%v'", json.Taskuser, 1)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)*/
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"taskttodo": tmpresult, "total": icount}})
}
func sccquerycommtaskcc(c *gin.Context) {
	var json Tasker
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commontask A left join scc_taskadvanced B on A.taskid=B.taskid left join scc_commontaskcc D on A.taskid=D.taskid  where D.cc= '%v'", json.Taskuser)
	sqlresult4 := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var icount int
	if _, ok := sqlresult4[0]["total"]; ok {
		count, _ := strconv.Atoi(sqlresult4[0]["total"])
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
					//int tmpnumberperpage = numberperpage;
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}
			sqlcmd1 = fmt.Sprintf("select A.taskid,A.taskname,A.creater,A.createtime,A.executor,A.creater,B.taskbegintime,B.taskendtime from scc_commontask A left join scc_taskadvanced B on A.taskid=B.taskid left join scc_commontaskcc C on A.taskid=C.taskid  where C.cc= '%v' order by A.taskid limit %v,%v", json.Taskuser, fromcount, numberperpage)
			sqlresult = sccinfo.tmpsql.SelectData(sqlcmd1)
			fmt.Println(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)

	/*sqlcmd1 := fmt.Sprintf("select A.taskname,A.createtime,A.executor,A.creater,B.taskbegintime,B.taskendtime from scc_commontask A inner join scc_taskadvanced B on A.taskid=B.taskid inner join scc_commontaskdoinfo C A.taskid=C.taskid on where C.cc= '%v'", json.Taskuser)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)*/
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"taskcc": tmpresult, "total": icount}})
}
func sccquerycommtaskdone(c *gin.Context) {
	var json Tasker
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commontask A left join scc_taskadvanced B on A.taskid=B.taskid where (A.taskstatus = '%v' or A.taskstatus = '%v') and A.executor='%v'", 3, 13, json.Taskuser)
	sqlresult4 := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var icount int
	if _, ok := sqlresult4[0]["total"]; ok {
		count, _ := strconv.Atoi(sqlresult4[0]["total"])
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
					//int tmpnumberperpage = numberperpage;
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}
			sqlcmd1 = fmt.Sprintf("select A.taskid,A.taskname,A.creater,A.createtime,A.taskstatus,B.taskbegintime,B.taskendtime from scc_commontask A left join scc_taskadvanced B on A.taskid=B.taskid where (A.taskstatus = '%v' or A.taskstatus = '%v') and A.executor='%v' order by A.taskid limit %v,%v", 3, 13, json.Taskuser, fromcount, numberperpage)
			sqlresult = sccinfo.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)

	/*sqlcmd1 := fmt.Sprintf("select A.taskname,A.createtime,A.taskstatus,B.taskbegintime,B.taskendtime from scc_commontask A inner join scc_taskadvanced B on A.taskid=B.taskid where (A.taskstatus = '%v' or A.taskstatus = '%v') and A.executor='%v'", 3, 13, json.Taskuser)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)*/
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"taskdone": tmpresult, "total": icount}})
}
func sccquerycommtaskcreated(c *gin.Context) {
	var json Tasker
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commontask A left join scc_taskadvanced B on A.taskid=B.taskid left join scc_commontaskdoinfo C on A.taskid=C.taskid where A.creater='%v'", json.Taskuser)
	sqlresult4 := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var icount int
	if _, ok := sqlresult4[0]["total"]; ok {
		count, _ := strconv.Atoi(sqlresult4[0]["total"])
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
					//int tmpnumberperpage = numberperpage;
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}
			sqlcmd1 = fmt.Sprintf("select A.taskid,A.taskname,A.createtime,A.executor,A.taskstatus,B.taskbegintime,B.taskendtime,B.tasktype,B.taskgrade,C.createtime as taskexectime from scc_commontask A left join scc_taskadvanced B on A.taskid=B.taskid left join scc_commontaskdoinfo C on A.taskid=C.taskid where A.creater='%v' order by A.taskid limit %v,%v", json.Taskuser, fromcount, numberperpage)
			sqlresult = sccinfo.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"taskcreated": tmpresult, "total": icount}})
}
func sccquerycommtaskneedapprove(c *gin.Context) {
	var json Tasker
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commontask A left join scc_taskadvanced B on A.taskid=B.taskid where A.taskstatus = '%v'  and B.taskapprover='%v'", 0, json.Taskuser)
	sqlresult4 := sccinfo.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	var icount int
	if _, ok := sqlresult4[0]["total"]; ok {
		count, _ := strconv.Atoi(sqlresult4[0]["total"])
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
					//int tmpnumberperpage = numberperpage;
					numberperpage = fromcount + numberperpage
					fromcount = 0
				}
			}
			sqlcmd1 = fmt.Sprintf("select A.taskid,A.taskname,A.createtime,A.executor,A.creater,B.taskbegintime,B.taskendtime from scc_commontask A left join scc_taskadvanced B on A.taskid=B.taskid where A.taskstatus = '%v' and B.taskapprover='%v' order by A.taskid limit %v,%v", 0, json.Taskuser, fromcount, numberperpage)
			sqlresult = sccinfo.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"approvetaskinfo": tmpresult, "total": icount}})
}
func scccapprovetask(c *gin.Context) {
	var json Approvetask
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskstatus := 1
	if "no" == json.Approveornot {
		taskstatus = 4
	}
	sqlcmd1 := fmt.Sprintf("update scc_commontask set taskstatus = '%v' where taskid = '%v'", taskstatus, json.Taskid)
	sccinfo.tmpsql.Execsqlcmd(sqlcmd1, false)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func scccquerytaskbyid(c *gin.Context) {
	var json Querytask
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	/*sqlcmd1 := fmt.Sprintf("select A.taskname,A.createtime,A.executor,A.creater,A.textinfo,A.filepath,A.taskstatus,B.taskbegintime,B.taskendtime,B.tasktype.B.taskgrade,B.taskapprover,C.createtime,C.textnfo as exectextinfo,C.filepath as execfilepath from scc_commontask A inner join scc_taskadvanced B on A.taskid=B.taskid inner join scc_commontaskdoinfo C on A.taskid=C.taskid where A.taskid='%v'", json.Taskid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)*/

	//拼解sql语句 先查任务表
	sqlcmd1 := fmt.Sprintf("select taskname,createtime,executor,creater,textinfo,filepath,taskstatus from scc_commontask where taskid = '%v'", json.Taskid)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd1)

	//拼解sql语句 在查高级配置表
	sqlcmd1 = fmt.Sprintf("select taskbegintime,taskendtime,tasktype,taskgrade,taskapprover from scc_taskadvanced  where taskid = '%v'", json.Taskid)
	sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)

	//拼解sql语句 在查任务完成度表
	sqlcmd1 = fmt.Sprintf("select createtime as donetime,textinfo as doneinfo,filepath as donefilepath from scc_commontaskdoinfo where taskid = '%v'", json.Taskid)
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd1)

	//拼解sql语句 在查抄送表
	sqlcmd1 = fmt.Sprintf("select cc from scc_commontaskcc where taskid = '%v'", json.Taskid)
	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd1)

	//拼解sql语句 在查评论表
	sqlcmd1 = fmt.Sprintf("select id,commenttime,commenter,commentinfo,commentpath from scc_commontaskcomment  where taskworkid = '%v'", json.Taskid)
	sqlresult4 := sccinfo.tmpsql.SelectData(sqlcmd1)

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"taskinfo": sqlresult, "taskadvancedinfo": sqlresult1, "taskworkinfo": sqlresult2, "taskccinfo": sqlresult3, "taskcommentinfo": sqlresult4}})
}
func scccommenttask(c *gin.Context) {
	var json Commenttask
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd1 := fmt.Sprintf("insert into scc_commontaskcomment (taskworkid,commenttime,commenter,commentinfo,commentpath)values('%v','%v','%v','%v','%v');", json.Taskid, time.Now().Unix(), json.Commenter, json.Commentinfo, json.Commentpath)
	_ = sccinfo.tmpsql.Insertmql(sqlcmd1)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

//注释样式  接口名称  开发完成  自测完成
func setrouter(r *gin.Engine) {
	r.POST("/createcommtask", scccreatecommtask) //创建任务 √  √
	r.POST("/docommtask", scccdocommtask)        //任务提交 √ √
	r.POST("/approvetask", scccapprovetask)      //审批任务 √ √
	r.POST("/commenttask", scccommenttask)       //添加评论 √ √

	r.POST("/querytaskbyid", scccquerytaskbyid)                      //通过taskid查询任务详细√ √
	r.POST("/querycommtasktodo", sccquerycommtaskto)                 //待我执行的任务√ √
	r.POST("/querycommtaskcc", sccquerycommtaskcc)                   //抄送给我任务√  √
	r.POST("/querycommtaskdone", sccquerycommtaskdone)               //我已执行完成的任务  √
	r.POST("/querycommtaskcreated", sccquerycommtaskcreated)         //我发布的任务   √  √
	r.POST("/querycommtaskneedapprove", sccquerycommtaskneedapprove) //需要我审批后才发布的任务 √ √
}
