package commconworkflow

import (
	"fmt"
	sccsql "golearn/gomysql"
	"net/http"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

type coreinfo struct {
	tmpsql sccsql.Mysqlconnectpool
}

var coworkflow coreinfo

func reverse(s []map[string]string) []map[string]string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
func init() {

	var conf map[string]map[string]string
	if _, err := toml.DecodeFile("./sccconfig.toml", &conf); err != nil {
		// handle error
	}
	Host := conf["sccCoworkflow"]["Host"]
	Username := conf["sccCoworkflow"]["Username"]
	Password := conf["sccCoworkflow"]["Password"]
	Dbname := conf["sccCoworkflow"]["Dbname"]
	Port := conf["sccCoworkflow"]["Port"]
	iport, _ := strconv.Atoi(Port)
	Serhost := conf["sccCoworkflow"]["Httpserverhost"]
	//fmt.Println("Hostxxxxxxxxxxxxxx", Host)
	go func(Host string, Username string, Password string, Dbname string, Serhost string, iport int) {
		coworkflow.tmpsql.Initmysql(Host, Username, Password, Dbname, iport)
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		setrouter(r)
		if err := r.Run(Serhost); err != nil {
			fmt.Println("startup service failed, err:\n", err)
		}
	}(Host, Username, Password, Dbname, Serhost, iport)

}
func scccreatetemplate(c *gin.Context) {
	var json CreateCoworkflow
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	iApprovelen := len(json.Approve)
	//fmt.Println("approve len  is iApprovelen")
	if iApprovelen <= 0 || iApprovelen > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no approver1"})
		return
	}
	var tmpapperver string
	for i := 0; i < iApprovelen; i++ {
		tmpapperver = tmpapperver + json.Approve[i].ApproverUsers
		if i != iApprovelen-1 {
			tmpapperver = tmpapperver + "&&"
		}

	}

	var strtmpcc string
	ilencc := len(json.Cc)
	for i := 0; i < ilencc; i++ {
		strtmpcc = strtmpcc + json.Cc[i].CcUsers

		if i != ilencc-1 {
			strtmpcc = strtmpcc + "&&"
		}

	}
	if "" == strtmpcc {
		strtmpcc = "0"
	}
	createtime := time.Now().Unix()
	sqlcmd1 := fmt.Sprintf("insert into scc_commonworkflow (subject,workflowtype,createtime,approverlist,cclist,creater,textinfo,filepath)values('%v','%v','%v','%v','%v','%v','%v','%v');", json.Subject, json.Workflowtype, createtime, tmpapperver, strtmpcc, json.Creater, json.Textinfo, json.Filepath)
	fmt.Println(sqlcmd1)
	workid := coworkflow.tmpsql.Insertmql(sqlcmd1)
	if 0 == workid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "create workflow failed"})
	}
	for i := 0; i < ilencc; i++ {
		sqlcmd1 = fmt.Sprintf("insert into scc_commonworkflowcc (workid,cctime,cc)values('%v','%v','%v');", workid, createtime, json.Cc[i].CcUsers)
		coworkflow.tmpsql.Execsqlcmd(sqlcmd1, false)
		fmt.Println(sqlcmd1)
	}
	for i := 0; i < iApprovelen; i++ {
		sqlcmd1 = fmt.Sprintf("insert into scc_commonworkflowapprove (workid,approver)values('%v','%v');", workid, json.Approve[i].ApproverUsers)
		coworkflow.tmpsql.Execsqlcmd(sqlcmd1, false)
		fmt.Println(sqlcmd1)
	}
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workid": workid}})
}
func sccqueryCoworkflowtodo(c *gin.Context) {
	var json Coworkflowtodo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sqlresult []map[string]string
	/*sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on B.workid=A.workid  where B.approver = '%v' and A.flowstatus =0", json.Sccid)
	sqlresult4 := coworkflow.tmpsql.SelectData(sqlcmd1)
	var pagenum = json.Pagenum
	var numberperpage = 30
	var fromcount = 0
	//var icount int
	if _, ok := sqlresult4[0]["total"]; ok {
		count, _ := strconv.Atoi(sqlresult4[0]["total"])
		//icount = count
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
			//查询和sccid有关并且当前sccid没审批的id
			sqlcmd1 = fmt.Sprintf("select B.id,B.workid,B.approvetype, B.approvetime,A.subject,A.workflowtype,A.createtime,A.approverlist,A.cclist,A.creater,A.textinfo,A.filepath from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on B.workid=A.workid  where  B.approver = '%v' and A.flowstatus =0 ", json.Sccid)
			sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
			//fmt.Println(sqlresult)
			for i := 0; i < len(sqlresult); {
				sqlcmd1 = fmt.Sprintf("select min(id) as id from scc_commonworkflowapprove where workid ='%v' and approvetype !=0 and approver = '%v'", sqlresult[i]["workid"], json.Sccid)
				sqlresult2 := coworkflow.tmpsql.SelectData(sqlcmd1)
				if "" == sqlresult2[0]["id"] && sqlresult2[0]["id"] != sqlresult[i]["id"] {
					fmt.Println(sqlresult2, "sss", sqlresult2[0]["id"], "zzzz", sqlresult[i]["id"])
					sqlresult = append(sqlresult[:i], sqlresult[i+1:]...)
				} else {
					i++
				}
			}
		}
	}*/

	sqlcmd1 := fmt.Sprintf("select B.id,B.workid,B.approvetype, B.approvetime,A.subject,A.workflowtype,A.createtime,A.approverlist,A.cclist,A.creater,A.textinfo,A.filepath from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on B.workid=A.workid  where  B.approver = '%v' and A.flowstatus =0 ", json.Sccid)
	sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
	//fmt.Println(sqlresult)
	for i := 0; i < len(sqlresult); {
		sqlcmd1 = fmt.Sprintf("select min(id) as id, approver from scc_commonworkflowapprove where workid ='%v' and approvetype =0 ", sqlresult[i]["workid"])
		sqlresult2 := coworkflow.tmpsql.SelectData(sqlcmd1)
		//fmt.Println(sqlresult2, "aaa", sqlresult2[0]["id"], "bbb", sqlresult[i]["id"], "ccc", sqlresult2[0]["approver"])
		if sqlresult2[0]["approver"] != json.Sccid {
			fmt.Println(sqlresult2, "sss", sqlresult2[0]["id"], "zzzz", sqlresult[i]["id"])
			sqlresult = append(sqlresult[:i], sqlresult[i+1:]...)
		} else {
			i++
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"todoinfo": tmpresult, "total": len(sqlresult)}})
}
func sccqueryCoworkflowcc(c *gin.Context) {
	var json Coworkflowcc
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commonworkflow A inner join scc_commonworkflowcc as B on B.workid=A.workid where B.cc= '%v' and A.flowstatus=1", json.Sccid)
	sqlresult4 := coworkflow.tmpsql.SelectData(sqlcmd1)
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
			sqlcmd1 = fmt.Sprintf("select A.subject,A.createtime,A.creater,A.textinfo,A.filepath,A.flowstatus,A.workid from scc_commonworkflow A inner join scc_commonworkflowcc as B on B.workid=A.workid where B.cc= '%v' and A.flowstatus=1 limit %v,%v", json.Sccid, fromcount, numberperpage)
			sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"ccinfo": tmpresult, "total": icount}})
}
func sccapproveCoworkflow(c *gin.Context) {
	var json Approve
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tmpapprove := 0
	if "yes" == json.Approve {
		tmpapprove = 2
	} else {
		tmpapprove = 1
	}
	//更新审批者表
	sqlcmd1 := fmt.Sprintf("update scc_commonworkflowapprove set approvetype =  '%v',approvetime='%v',approvetextinfo='%v',approvefilepath='%v' where approver = '%v' and workid = '%v'", tmpapprove, time.Now().Unix(), json.Textinfo, json.Fileptah, json.Sccid, json.Workid)
	coworkflow.tmpsql.Execsqlcmd(sqlcmd1, false)
	if 1 == tmpapprove {
		sqlcmd1 := fmt.Sprintf("update scc_commonworkflow set flowstatus =  '%v' where workid = '%v'", 2, json.Workid)
		coworkflow.tmpsql.Execsqlcmd(sqlcmd1, false)
	} else {
		//同时还要更新审批表
		sqlcmd1 = fmt.Sprintf("select approvetype from scc_commonworkflowapprove where workid = '%v' order by id desc  limit 0,1", json.Workid)
		mysqlret := coworkflow.tmpsql.SelectData(sqlcmd1)
		if "0" == mysqlret[0]["approvetype"] {
			//sqlcmd1 := fmt.Sprintf("update scc_commonworkflow set workflowtype =  '%v' where workid = '%v'", 0, json.Workid)
			//coworkflow.tmpsql.Execsqlcmd(sqlcmd1, false)
		} else if "1" == mysqlret[0]["approvetype"] {
			sqlcmd1 := fmt.Sprintf("update scc_commonworkflow set flowstatus =  '%v' where workid = '%v'", 0, json.Workid)
			coworkflow.tmpsql.Execsqlcmd(sqlcmd1, false)
		} else if "2" == mysqlret[0]["approvetype"] {
			sqlcmd1 := fmt.Sprintf("update scc_commonworkflow set flowstatus =  '%v' where workid = '%v'", 1, json.Workid) //最大的是2 证明已完结
			coworkflow.tmpsql.Execsqlcmd(sqlcmd1, false)
		}
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func queryCoworkflowbyworkid(c *gin.Context) {
	var json Coworkflowbyworkid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd1 := fmt.Sprintf("select approvetype,approvetime,approvetextinfo,approvefilepath,approver from scc_commonworkflowapprove where workid = '%v'", json.Workid)
	sqlret := coworkflow.tmpsql.SelectData(sqlcmd1)
	sqlcmd1 = fmt.Sprintf("select subject,workflowtype,createtime,approverlist,cclist,creater,textinfo,filepath,flowstatus from scc_commonworkflow where workid = '%v'", json.Workid)

	sqlret1 := coworkflow.tmpsql.SelectData(sqlcmd1)
	sqlcmd1 = fmt.Sprintf("select id,commentworkid,commentid,commenttime,commenter,commentinfo,commentpath from scc_commonworkflowcomment where commentworkid = '%v'", json.Workid)

	sqlret2 := coworkflow.tmpsql.SelectData(sqlcmd1)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"Coworkflowinfo": sqlret1, "Coworkflowapprove": sqlret, "CoworkflowComment": sqlret2}})
}
func sccqueryCoworkflowallbycreate(c *gin.Context) {
	var json Coworkflowdone
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commonworkflow where creater = '%v'", json.Sccid)
	sqlresult4 := coworkflow.tmpsql.SelectData(sqlcmd1)
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
			sqlcmd1 = fmt.Sprintf("select workid,subject,workflowtype,createtime,approverlist,cclist,creater,textinfo,filepath,flowstatus from scc_commonworkflow where creater = '%v' limit %v,%v", json.Sccid, fromcount, numberperpage)
			sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"flowinfo": tmpresult, "total": icount}})
}
func sccqueryCoworkflowallbycreatereject(c *gin.Context) {
	var json Coworkflowdone
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commonworkflow where creater = '%v' and flowstatus = 2", json.Sccid)
	sqlresult4 := coworkflow.tmpsql.SelectData(sqlcmd1)
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
			sqlcmd1 = fmt.Sprintf("select workid,subject,workflowtype,createtime,approverlist,cclist,creater,textinfo,filepath,flowstatus from scc_commonworkflow where creater = '%v' and flowstatus = 2 limit %v,%v", json.Sccid, fromcount, numberperpage)
			fmt.Println(sqlcmd1)
			sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"flowinfo": tmpresult, "total": icount}})
}
func sccqueryCoworkflowallbycreatedoing(c *gin.Context) {
	var json Coworkflowdone
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commonworkflow where creater = '%v' and flowstatus = 0", json.Sccid)
	sqlresult4 := coworkflow.tmpsql.SelectData(sqlcmd1)
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
			sqlcmd1 = fmt.Sprintf("select workid,subject,workflowtype,createtime,approverlist,cclist,creater,textinfo,filepath,flowstatus from scc_commonworkflow where creater = '%v' and flowstatus = 0 limit %v,%v", json.Sccid, fromcount, numberperpage)
			fmt.Println(sqlcmd1)
			sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"flowinfo": tmpresult, "total": icount}})
}
func sccqueryCoworkflowallbycreatedone(c *gin.Context) {
	var json Coworkflowdone
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commonworkflow where creater = '%v' and flowstatus = 1", json.Sccid)
	sqlresult4 := coworkflow.tmpsql.SelectData(sqlcmd1)
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
			sqlcmd1 = fmt.Sprintf("select workid,subject,workflowtype,createtime,approverlist,cclist,creater,textinfo,filepath,flowstatus from scc_commonworkflow where creater = '%v' and flowstatus = 1 limit %v,%v", json.Sccid, fromcount, numberperpage)
			fmt.Println(sqlcmd1)
			sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"flowinfo": tmpresult, "total": icount}})
}
func sccqueryCoworkflowallbyapprove(c *gin.Context) {
	var json Coworkflowdone
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on A.workid = B.workid where B.approver='%v'", json.Sccid) //mod by sjd 在多人审批中有问题 放开审批中限制
	sqlresult4 := coworkflow.tmpsql.SelectData(sqlcmd1)
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
			sqlcmd1 = fmt.Sprintf("select A.workid,A.subject,A.workflowtype,A.createtime,A.approverlist,A.cclist,A.creater,A.textinfo,A.filepath,A.flowstatus from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on A.workid = B.workid where B.approver='%v' limit %v,%v ", json.Sccid, fromcount, numberperpage)
			fmt.Println(sqlcmd1)
			sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"flowinfo": tmpresult, "total": icount}})
}
func sccqueryCoworkflowallbyapprovereject(c *gin.Context) {
	var json Coworkflowdone
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on A.workid = B.workid where B.approver='%v' and A.flowstatus =2", json.Sccid)
	sqlresult4 := coworkflow.tmpsql.SelectData(sqlcmd1)
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
			sqlcmd1 = fmt.Sprintf("select A.workid,A.subject,A.workflowtype,A.createtime,A.approverlist,A.cclist,A.creater,A.textinfo,A.filepath from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on A.workid = B.workid where B.approver='%v' and A.flowstatus =2 limit %v,%v", json.Sccid, fromcount, numberperpage)
			fmt.Println(sqlcmd1)
			sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"flowinfo": tmpresult, "total": icount}})
}
func sccqueryCoworkflowallbyapprovedone(c *gin.Context) {
	var json Coworkflowdone
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on A.workid = B.workid where B.approver='%v' and A.flowstatus =1", json.Sccid)
	sqlresult4 := coworkflow.tmpsql.SelectData(sqlcmd1)
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
			sqlcmd1 = fmt.Sprintf("select A.workid,A.subject,A.workflowtype,A.createtime,A.approverlist,A.cclist,A.creater,A.textinfo,A.filepath from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on A.workid = B.workid where B.approver='%v' and A.flowstatus =1 limit %v,%v", json.Sccid, fromcount, numberperpage)
			fmt.Println(sqlcmd1)
			sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"flowinfo": tmpresult, "total": icount}})
}

func sccqueryCoworkflowallbyapproveding(c *gin.Context) {
	var json Coworkflowdone
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sqlresult []map[string]string
	sqlcmd1 := fmt.Sprintf("select count(*) as total from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on A.workid = B.workid where B.approver='%v' and A.flowstatus =0", json.Sccid)
	sqlresult4 := coworkflow.tmpsql.SelectData(sqlcmd1)
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
			sqlcmd1 = fmt.Sprintf("select A.workid,A.subject,A.workflowtype,A.createtime,A.approverlist,A.cclist,A.creater,A.textinfo,A.filepath from scc_commonworkflowapprove as B inner join scc_commonworkflow as A on A.workid = B.workid where B.approver='%v' and A.flowstatus =0 limit %v,%v", json.Sccid, fromcount, numberperpage)
			fmt.Println(sqlcmd1)
			sqlresult = coworkflow.tmpsql.SelectData(sqlcmd1)
		}
	}
	tmpresult := reverse(sqlresult)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"flowinfo": tmpresult, "total": icount}})
}
func sccaddcomment(c *gin.Context) {
	var json Comment
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd1 := fmt.Sprintf("insert into scc_commonworkflowcomment (commentworkid,commentid,commenttime,commenter,commentinfo,commentpath)values('%v','%v','%v','%v','%v','%v');", json.Commentworkid, json.Commentid, time.Now().Unix(), json.Commenter, json.Textinfo, json.Fileptah)
	//fmt.Println(sqlcmd1)
	coworkflow.tmpsql.Execsqlcmd(sqlcmd1, true)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func setrouter(r *gin.Engine) {
	r.POST("/createCoworkflow", scccreatetemplate)                                   //创建通用审批
	r.POST("/queryCoworkflowtodo", sccqueryCoworkflowtodo)                           //查询待我审批的
	r.POST("/queryCoworkflowcc", sccqueryCoworkflowcc)                               //查询抄送给我的审批 并且已完成的
	r.POST("/approveCoworkflow", sccapproveCoworkflow)                               //审批操作
	r.POST("/queryCoworkflowbyworkid", queryCoworkflowbyworkid)                      //根据workid 查询审批的详细信息
	r.POST("/queryCoworkflowallbycreate", sccqueryCoworkflowallbycreate)             //我创建的所有审批
	r.POST("/queryCoworkflowallbycreatereject", sccqueryCoworkflowallbycreatereject) //我创建的所有审批 被拒绝的
	r.POST("/queryCoworkflowallbycreatedoing", sccqueryCoworkflowallbycreatedoing)   //我创建的所有审批 过程中
	r.POST("/queryCoworkflowallbycreatedone", sccqueryCoworkflowallbycreatedone)     //我创建的所有审批 已同意的

	r.POST("/queryCoworkflowallbyapprove", sccqueryCoworkflowallbyapprove)             //我审批的所有审批
	r.POST("/queryCoworkflowallbyapprovereject", sccqueryCoworkflowallbyapprovereject) //我拒绝的所有审批 被拒绝的
	r.POST("/queryCoworkflowallbyapprovedone", sccqueryCoworkflowallbyapprovedone)     //我审批的所有审批 已同意的
	r.POST("/queryCoworkflowallbyapproveding", sccqueryCoworkflowallbyapproveding)     //我审批的所有审批 处理中

	r.POST("/addcomment", sccaddcomment) //评论审批
}
