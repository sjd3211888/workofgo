package sccwork

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

var sccinfo coreinfo

func init() {

	var conf map[string]map[string]string
	if _, err := toml.DecodeFile("./sccconfig.toml", &conf); err != nil {
		// handle error
	}
	Host := conf["sccwork"]["Host"]
	Username := conf["sccwork"]["Username"]
	Password := conf["sccwork"]["Password"]
	Dbname := conf["sccwork"]["Dbname"]
	Port := conf["sccwork"]["Port"]
	iport, _ := strconv.Atoi(Port)
	Serhost := conf["sccwork"]["Httpserverhost"]
	//fmt.Println("Hostxxxxxxxxxxxxxx", Host)
	go func(Host string, Username string, Password string, Dbname string, Serhost string, iport int) {
		sccinfo.tmpsql.Initmysql(Host, Username, Password, Dbname, iport)
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		setrouter(r)
		if err := r.Run(Serhost); err != nil {
			fmt.Println("startup service failed, err:\n", err)
		}
	}(Host, Username, Password, Dbname, Serhost, iport)

}

func sccworklogin(c *gin.Context) {
	type scclogin struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var json scclogin
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd := fmt.Sprintf("select  userid from scc_user where userid='%v'", json.Username)
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
	c.SetCookie("abc", "123", 60, "/",
		"localhost", false, true)
	// 返回信息
	if 0 == len(sqlresult) {
		c.JSON(http.StatusBadRequest, gin.H{"result": "falied"})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "success"})
	}
}
func scccreatetemplate(c *gin.Context) {
	type apporverlist struct {
		Apporver string `json:"apporver" binding:"required"`
	}
	type sccccist struct {
		Cc string `json:"ccuser" binding:"required"`
	}
	type templatework struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Templatename string         `json:"name" binding:"required"`
		Trade        string         `json:"trade" binding:"required"`
		Ptype        string         `json:"type" binding:"required"`
		Creater      string         `json:"creater" binding:"required"`
		Templateuer  string         `json:"templateuer" binding:"required"`
		Resrearcher  string         `json:"resrearcher" binding:"required"`
		Apporlist    []apporverlist `json:"apporver" binding:"required"`
		Cclist       []sccccist     `json:"cc"`
	}
	var json templatework
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	countapprolist := len(json.Apporlist)

	workid := time.Now().Unix()
	if 0 == countapprolist {
		c.JSON(http.StatusOK, gin.H{"error": "Apporlist is null"})
	} else {
		sqlcmd2 := fmt.Sprintf("insert into scc_approver (workid,approver,approvertype)values(%v,%v,%v);", workid, json.Templateuer, 0)

		sccinfo.tmpsql.Execsqlcmd(sqlcmd2, false)
		sqlcmd1 := fmt.Sprintf("insert into scc_approver (workid,approver,approvertype)values(%v,%v,%v);", workid, json.Resrearcher, 1)

		sccinfo.tmpsql.Execsqlcmd(sqlcmd1, false)
		sqlcmd1 = fmt.Sprintf("insert into scc_approver (workid,approver,approvertype)values(%v,%v,%v);", workid, json.Resrearcher, 101)

		sccinfo.tmpsql.Execsqlcmd(sqlcmd1, false)
		for k, v := range json.Apporlist {
			sqlcmd := fmt.Sprintf("insert into scc_approver (workid,approver,approvertype)values(%v,%v,%v);", workid, v.Apporver, k+2)

			sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
			sqlcmd = fmt.Sprintf("insert into scc_approver (workid,approver,approvertype)values(%v,%v,%v);", workid, v.Apporver, k+2+100)

			sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
		}

	}
	countcclist := len(json.Cclist)
	if 0 == countcclist {
		//fmt.Println("CC list is null")
	} else {
		for _, v := range json.Cclist {
			sqlcmd := fmt.Sprintf("insert into scc_cc(workid,cc)values(%v,%v);", workid, v.Cc)
			sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
		}

	}
	sqlcmd := fmt.Sprintf("insert into scc_worktempplate (workid,templatename,trade,p_type,createtime,approverlist,cclist,creater,templateuser,researcher)values(%v,'%v',%v,%v,%v,%v,%v,%v,%v,%v);", workid, json.Templatename, json.Trade, json.Ptype, workid, workid, workid, json.Creater, json.Templateuer, json.Resrearcher)
	//fmt.Println(sqlcmd)
	sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
	//fmt.Println(json.Apporlist, json.Cclist)
	strworkid := strconv.FormatInt(workid, 10)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"templateid": strworkid}})
}

func sccquerytemplate(c *gin.Context) {
	type querytemplatework struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Usertype string `json:"usertype" binding:"required"` // 1 模板创建者 2模板使用者
		Username string `json:"useid" binding:"required"`
	}
	var json querytemplatework
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if json.Usertype == "1" {
		sqlcmd1 := fmt.Sprintf("Select workid,templatename,trade,p_type,createtime,approverlist,cclist,creater,templateuser,active from scc_worktempplate where creater= '%v'", json.Username)
		sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
		c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"templateinfo": sqlresult1}})
	} else if json.Usertype == "2" {
		sqlcmd1 := fmt.Sprintf("Select workid,templatename,trade,p_type,createtime,approverlist,cclist,creater,templateuser,active from scc_worktempplate where templateuser= '%v'and active=1", json.Username)
		sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
		c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"templateinfo": sqlresult1}})
	}

}
func scccreateapply(c *gin.Context) {
	type querytemplatework struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Templateid   string `json:"templateid" binding:"required"`
		Username     string `json:"username" binding:"required"`
		Textinfo     string `json:"textinfo"`
		Filepath     string `json:"filepath"`
		Telephonenum string `json:"telephonenum"`
	}

	var json querytemplatework
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//先创建apply表
	sqlcmd := fmt.Sprintf("insert into scc_apply(templateid,createtime,textinfo,filepath,telephone,creater)values(%v,%v,'%v','%v','%v',%v);", json.Templateid, time.Now().Unix(), json.Textinfo, json.Filepath, json.Telephonenum, json.Username)
	appid := sccinfo.tmpsql.Insertmql(sqlcmd)

	//插入第一个生成表 插入研究员
	sqlcmd1 := fmt.Sprintf("insert into scc_workflow(appid,createtime,advise,templateid,appcurentnode,appnextnode)values(%v,%v,'%v',%v,%v,%v);", appid, time.Now().Unix(), "问题创建", json.Templateid, 0, 1)
	sccinfo.tmpsql.Execsqlcmd(sqlcmd1, false)

	//插入一条到满意统计表
	/*sqlcmd3 := fmt.Sprintf("insert into scc_statistica (appid,createtime,comment,commentuser)values(%v,%v,'%v','%v');", appid, time.Now().Unix(), json.Textinfo, json.Username)

	sccinfo.tmpsql.Execsqlcmd(sqlcmd3, false)*/

	c.JSON(http.StatusOK, gin.H{"result": "success", "appid": appid})
}
func sccqueryworktempplatebyworkid(workid string) []map[string]string {
	sqlcmd2 := fmt.Sprintf("Select workid,templatename,trade,p_type,creater from scc_worktempplate where workid= '%v'", workid)
	return sccinfo.tmpsql.SelectData(sqlcmd2)

}
func sccqueryapproveryworkid(workid string) []map[string]string {
	sqlcmd3 := fmt.Sprintf("Select approver,approvertype from scc_approver where workid= '%v'", workid)
	return sccinfo.tmpsql.SelectData(sqlcmd3)
}
func sccqueryccworkid(workid string) []map[string]string {
	sqlcmd4 := fmt.Sprintf("Select cc from scc_cc where workid= '%v'", workid)
	return sccinfo.tmpsql.SelectData(sqlcmd4)
}
func sccqueryapply(c *gin.Context) {
	type queryapply struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Applytype string `json:"type" binding:"required"` //4 我创建的  1  带我处理  2 已处理 3 抄送给我的 5 和我相关的 6 带我研究的 7 待执行
		Username  string `json:"username"`
	}

	var json queryapply
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	switch json.Applytype {
	case "4":
		{
			sqlcmd := fmt.Sprintf("Select t1.appid,t1.templateid,t1.createtime,t1.textinfo,t1.filepath,t1.telephone,t1.creater,t1.grade,t2.templatename from scc_apply t1 inner join  scc_worktempplate t2 on  t2.workid=t1.templateid  where t1.creater= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": sqlresult}})
				break
			}

			/*	workflworkresult := make([]map[string]string, 0)
				for k := range sqlresult {
					sqlcmd1 := fmt.Sprintf("Select id,appid,createtime,apperover,apperovertypei,advise,templateidm,appcurentnode,appnextnode from scc_workflow where appid = '%v'", sqlresult[k]["appid"])
					sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
					workflworkresult = append(workflworkresult, sqlresult1...)
				}*/
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": sqlresult}})
			//需要加开一个接口 通过appid查全部
			break
		}
	case "1":
		{
			sqlcmd := fmt.Sprintf("Select workid,approver,approvertype from scc_approver where approver= '%v'", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": sqlresult}})
				break
			}
			workflworkresult := make([]map[string]string, 0)
			for k := range sqlresult {
				sqlcmd1 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = '%v' and appnextnode = '%v' ", sqlresult[k]["workid"], sqlresult[k]["approvertype"])
				//(sqlcmd1)
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				for j := range sqlresult1 {
					sqlcmd2 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = '%v' and appcurentnode = '%v' and appid= '%v' order by id desc limit 1 ", sqlresult[k]["workid"], sqlresult[k]["approvertype"], sqlresult1[j]["appid"])
					sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
					if 0 == len(sqlresult2) {
						sqlcmd6 := fmt.Sprintf("Select t1.appid,t1.templateid,t1.createtime,t1.textinfo,t1.filepath,t1.telephone,t1.creater,t1.grade,t2.templatename from scc_apply t1 inner join scc_worktempplate t2 on t2.workid=t1.templateid where appid='%v'", sqlresult1[j]["appid"])
						sqlresult6 := sccinfo.tmpsql.SelectData(sqlcmd6)
						workflworkresult = append(workflworkresult, sqlresult6...)
					} else {
						tmpsql1id, err := strconv.Atoi(sqlresult1[j]["id"])
						if err != nil {
							tmpsql1id = 0
						}
						tmpsql2id, err := strconv.Atoi(sqlresult2[0]["id"])
						if err != nil {
							tmpsql2id = 0
						}
						if tmpsql2id > tmpsql1id { //最新审批过的id大于审批的id 证明被审批过了
							//证明呗审批过了
						} else {
							sqlcmd6 := fmt.Sprintf("Select t1.appid,t1.templateid,t1.createtime,t1.textinfo,t1.filepath,t1.telephone,t1.creater,t1.grade,t2.templatename from scc_apply t1 inner join scc_worktempplate t2 on t2.workid=t1.templateid where appid= '%v'", sqlresult2[0]["appid"])
							sqlresult6 := sccinfo.tmpsql.SelectData(sqlcmd6)
							workflworkresult = append(workflworkresult, sqlresult6...)
						}
					}

				}
			}
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": workflworkresult}})
			break
		}
	case "2":
		{
			sqlcmd := fmt.Sprintf("Select workid,approver,approvertype from scc_approver where approver= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": sqlresult}})
				break
			}
			workflworkresult := make([]map[string]string, 0)
			for k := range sqlresult {
				sqlcmd1 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = '%v' and appcurentnode = '%v'  limit 1", sqlresult[k]["workid"], sqlresult[k]["approvertype"])
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				if 0 != len(sqlresult1) {
					sqlcmd6 := fmt.Sprintf("Select t1.appid,t1.templateid,t1.createtime,t1.textinfo,t1.filepath,t1.telephone,t1.creater,t1.grade,t2.templatename from scc_apply t1 inner join scc_worktempplate t2 on t2.workid=t1.templateid where appid= '%v' ", sqlresult[k]["appid"])
					sqlresult6 := sccinfo.tmpsql.SelectData(sqlcmd6)
					workflworkresult = append(workflworkresult, sqlresult6...)
				}

			}
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": workflworkresult}})
			break
		}
	case "3":
		{
			sqlcmd := fmt.Sprintf("Select workid from scc_cc where cc= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			cccount := len(sqlresult)
			workflworkresult := make([]map[string]string, 0)
			for k := 0; k < cccount; k++ {
				tmpworkid := sqlresult[k]["workid"]
				inttmpworkid, _ := strconv.Atoi(tmpworkid)
				sqlcmd1 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = %v group by appid", inttmpworkid)
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				kccount := len(sqlresult1)
				for tmplen := 0; tmplen < kccount; tmplen++ {
					sqlcmd6 := fmt.Sprintf("Select t1.appid,t1.templateid,t1.createtime,t1.textinfo,t1.filepath,t1.telephone,t1.creater,t1.grade,t2.templatename,sw.id,sw.appcurentnode,sw.appnextnode from scc_apply t1 inner join scc_worktempplate t2 on t2.workid=t1.templateid inner join scc_workflow sw on sw.appid = t1.appid where t1.appid= %v order by id desc limit 1;", sqlresult1[tmplen]["appid"])
					sqlresult6 := sccinfo.tmpsql.SelectData(sqlcmd6)
					workflworkresult = append(workflworkresult, sqlresult6...)
				}
				if 0 != len(sqlresult1) {

				}
				//(workflworkresult)
			}
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": workflworkresult}})
			break
		}
	case "5":
		{

			//to do
			sqlcmd := fmt.Sprintf("Select DISTINCT workid from scc_approver where approver= '%v' and approvertype<100", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": sqlresult}})
				break
			}
			workflworkresult := make([]map[string]string, 0)
			for k := range sqlresult {
				tmpappid, _ := strconv.Atoi(sqlresult[k]["workid"])
				//sqlcmd1 := fmt.Sprintf("Select appid,templateid,createtime,textinfo,filepath,telephone,creater  from scc_apply where templateid = %v ", tmpappid)
				sqlcmd1 := fmt.Sprintf("select sc.appid,sc.templateid,sc.createtime,sc.textinfo,sc.filepath,sc.telephone,sc.creater,sc.grade,t2.advise,t2.id,t2.createtime as workflowtime,t2.appcurentnode,t2.appnextnode,t3.templatename from scc_apply sc inner join(select sw.appid,sw.advise,sw.id,sw.createtime,max(sw.appcurentnode) as appcurentnode,max(sw.appnextnode) as appnextnode from scc_workflow sw  group by sw.appid) t2 on sc.appid=t2.appid inner join scc_worktempplate t3 on t3.workid=sc.templateid where sc.templateid= %v group by sc.appid order by t2.id desc;", tmpappid)
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				//("SSSSSSSSSS", sqlresult1)
				if 0 != len(sqlresult1) {
					workflworkresult = append(workflworkresult, sqlresult1...)
				}
			}
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": workflworkresult}})
			break
		}
	case "6":
		{
			sqlcmd := fmt.Sprintf("Select workid from scc_worktempplate where researcher= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": sqlresult}})
				break
			}
			workflworkresult := make([]map[string]string, 0)
			for k := range sqlresult {
				sqlcmd1 := fmt.Sprintf("Select id,appid,templateid from scc_workflow where templateid = '%v' and appnextnode = 1 ", sqlresult[k]["workid"])
				//("wocao ", sqlcmd1)
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				for j := range sqlresult1 {
					sqlcmd2 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = '%v' and appnextnode = 2 and appid= '%v' order by id desc limit 1", sqlresult[k]["workid"], sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
					sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
					//("wocao1 ", sqlcmd2)
					if 0 == len(sqlresult2) {
						//fmt.Println("SS", sqlcmd2)
						sqlcmd5 := fmt.Sprintf("Select t1.appid,t1.templateid,t1.createtime,t1.textinfo,t1.filepath,t1.telephone,t1.creater,t1.grade,t2.templatename from scc_apply t1 inner join scc_worktempplate t2 on t2.workid=t1.templateid where appid= '%v' ", sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
						sqlresult5 := sccinfo.tmpsql.SelectData(sqlcmd5)
						workflworkresult = append(workflworkresult, sqlresult5...)
					} else {
						tmpproid, err := strconv.Atoi(sqlresult1[j]["id"])
						if err != nil {
							tmpproid = 0
						}
						tmpnext, err := strconv.Atoi(sqlresult2[0]["id"])
						if err != nil {
							tmpnext = 0
						}
						if tmpnext > tmpproid {
							//fmt.Println(tmpnext, tmpproid)
							//不加了
						} else {
							//fmt.Println(tmpnext, tmpproid)
							sqlcmd5 := fmt.Sprintf("Select t1.appid,t1.templateid,t1.createtime,t1.textinfo,t1.filepath,t1.telephone,t1.creater,t1.grade,t2.templatename from scc_apply t1 inner join scc_worktempplate t2 on t2.workid=t1.templateid where appid= '%v' ", sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
							sqlresult5 := sccinfo.tmpsql.SelectData(sqlcmd5)
							workflworkresult = append(workflworkresult, sqlresult5...)
						}
					}
					/*	}
						sqlcmd3 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = %v and appnextnode = 2 and appid=%v order by id desc limit 1", sqlresult[k]["workid"], sqlresult2[0]["appid"]) //查询当前appid中已经执行过研究的最大ID，如果执行研究的最大ID大于此ID，表明有呗打回，否则上面的ID应该是未研究的最大ID
						sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd3)
						fmt.Println("sss", sqlresult3)
						tmpproid, err := strconv.Atoi(sqlresult2[0]["id"])
						if err != nil {
							tmpproid = 0
						}
						if 0 == len(sqlresult3) {
							fmt.Println("cccc", sqlresult3)
							workflworkresult = append(workflworkresult, sqlresult1[j])
						} else {
							tmpnext, err := strconv.Atoi(sqlresult3[0]["id"])
							if err != nil {
								tmpnext = 0
							}
							if tmpnext < tmpproid {
								fmt.Println("bbb", sqlresult3)
								workflworkresult = append(workflworkresult, sqlresult1[j])
							}
						}

					}*/
				}
			}
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": workflworkresult}})
			break
		}
	case "7":
		{
			sqlcmd := fmt.Sprintf("Select workid from scc_worktempplate where researcher= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": sqlresult}})
				break
			}
			workflworkresult := make([]map[string]string, 0)
			for k := range sqlresult {
				sqlcmd1 := fmt.Sprintf("Select id,appid,templateid from scc_workflow where templateid = '%v' and appnextnode = 101 ", sqlresult[k]["workid"])
				//("wocao ", sqlcmd1)
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				for j := range sqlresult1 {
					sqlcmd2 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = '%v' and appnextnode = 102 and appid= '%v' order by id desc limit 1", sqlresult[k]["workid"], sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
					sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
					//fmt.Println("wocao1 ", sqlcmd2)
					if 0 == len(sqlresult2) {
						//fmt.Println("SS", sqlcmd2)
						sqlcmd5 := fmt.Sprintf("Select t1.appid,t1.templateid,t1.createtime,t1.textinfo,t1.filepath,t1.telephone,t1.creater,t1.grade,t2.templatename from scc_apply t1 inner join scc_worktempplate t2 on t2.workid=t1.templateid where appid= '%v'", sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
						sqlresult5 := sccinfo.tmpsql.SelectData(sqlcmd5)
						workflworkresult = append(workflworkresult, sqlresult5...)
					} else {
						tmpproid, err := strconv.Atoi(sqlresult1[j]["id"])
						if err != nil {
							tmpproid = 0
						}
						tmpnext, err := strconv.Atoi(sqlresult2[0]["id"])
						if err != nil {
							tmpnext = 0
						}
						if tmpnext > tmpproid {
							//fmt.Println(tmpnext, tmpproid)
							//不加了
						} else {
							fmt.Println(tmpnext, tmpproid)
							sqlcmd5 := fmt.Sprintf("Select t1.appid,t1.templateid,t1.createtime,t1.textinfo,t1.filepath,t1.telephone,t1.creater,t2.templatename from scc_apply t1 inner join scc_worktempplate t2 on t2.workid=t1.templateid where appid= '%v' ", sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
							sqlresult5 := sccinfo.tmpsql.SelectData(sqlcmd5)
							workflworkresult = append(workflworkresult, sqlresult5...)
						}
					}
					/*	}
						sqlcmd3 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = %v and appnextnode = 2 and appid=%v order by id desc limit 1", sqlresult[k]["workid"], sqlresult2[0]["appid"]) //查询当前appid中已经执行过研究的最大ID，如果执行研究的最大ID大于此ID，表明有呗打回，否则上面的ID应该是未研究的最大ID
						sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd3)
						fmt.Println("sss", sqlresult3)
						tmpproid, err := strconv.Atoi(sqlresult2[0]["id"])
						if err != nil {
							tmpproid = 0
						}
						if 0 == len(sqlresult3) {
							fmt.Println("cccc", sqlresult3)
							workflworkresult = append(workflworkresult, sqlresult1[j])
						} else {
							tmpnext, err := strconv.Atoi(sqlresult3[0]["id"])
							if err != nil {
								tmpnext = 0
							}
							if tmpnext < tmpproid {
								fmt.Println("bbb", sqlresult3)
								workflworkresult = append(workflworkresult, sqlresult1[j])
							}
						}

					}*/
				}
			}
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": workflworkresult}})
			break
		}
	default:
		{
			c.JSON(http.StatusOK, gin.H{"result": "success"})
			break
		}
	}
}
func sccuerytmplatebyworkid(c *gin.Context) {
	type workid struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Workid string `json:"workid" binding:"required"` //0 我创建的  1 已审批的  2 待我审批的  3 抄送给我的 5 和我相关的
	}

	var json workid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	templateinfo := sccqueryworktempplatebyworkid(json.Workid)
	pperorinfo := sccqueryapproveryworkid(json.Workid)
	ccinfo := sccqueryccworkid(json.Workid)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"templateinfo": templateinfo, "apperorinfo": pperorinfo, "ccinfo": ccinfo}})
}
func sccquerytemplateappandcc(c *gin.Context) {
	type workid struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Workid string `json:"workid" binding:"required"` //0 我创建的  1 已审批的  2 待我审批的  3 抄送给我的
	}

	var json workid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd2 := fmt.Sprintf("Select approver,approvertype from scc_approver where workid= '%v'", json.Workid)
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	sqlcmd3 := fmt.Sprintf("Select cc from scc_cc where workid= '%v'", json.Workid)
	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd3)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"approinfo": sqlresult2, "ccinfo": sqlresult3}})
}
func sccapprove(c *gin.Context) {
	type workid struct {
		Templateid string `json:"templateid" binding:"required"`
		Appid      string `json:"appid" binding:"required"`
		Status     string `json:"status" binding:"required"`
		Advise     string `json:"advise" binding:"required"`
		Username   string `json:"username" binding:"required"`
		Filepath   string `json:"filepath"`
	}
	//打回直接打回给研究人
	var json workid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if "2" == json.Status { //2 是打回
		//先查询下审批表
		sqlcmd2 := fmt.Sprintf("Select templateid,appnextnode,appid,id from scc_workflow where templateid= '%v' and  appid = %v order by id desc limit 1", json.Templateid, json.Appid)
		sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
		lenresult := len(sqlresult2)
		if 0 == lenresult {
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": "no user current1"})
			return
		}
		tmpnextnode, err1 := strconv.Atoi(sqlresult2[0]["appnextnode"])
		if nil != err1 {
			tmpnextnode = 0
		}
		sqlcmd3 := fmt.Sprintf("Select workid from scc_approver where workid= '%v' and  approver = %v and approvertype = %v", sqlresult2[0]["templateid"], json.Username, tmpnextnode)
		//fmt.Println(sqlcmd3)
		sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd3)
		if 0 == len(sqlresult3) {
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": "no user current2"})
			return
		}

		if tmpnextnode >= 100 {

			sqlcmd4 := fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode,filepath)values(%v,%v,'%v',%v,%v,%v,'%v');", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, 101, json.Filepath)
			sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
		} else {

			sqlcmd4 := fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode,filepath)values(%v,%v,'%v',%v,%v,%v,'%v');", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, 1, json.Filepath)
			sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
		}

	} else if "1" == json.Status {
		sqlcmd2 := fmt.Sprintf("Select templateid,appnextnode,appid,id from scc_workflow where templateid= '%v' and  appid = %v order by id desc limit 1", json.Templateid, json.Appid)
		sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)

		lenresult := len(sqlresult2)
		if 0 == lenresult {
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": "no user current1"})
			return
		}
		tmpapptype, err := strconv.Atoi(sqlresult2[0]["appnextnode"])

		if nil != err {
			tmpapptype = 0
		}
		sqlcmd3 := fmt.Sprintf("Select workid from scc_approver where workid= %v and  approver = %v and approvertype = %v ", sqlresult2[0]["templateid"], json.Username, tmpapptype)
		//fmt.Println(sqlcmd3)
		sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd3)
		if 0 == len(sqlresult3) {
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": "no user current2"})
			return
		}
		tmpnextnode, err1 := strconv.Atoi(sqlresult2[0]["appnextnode"])
		if nil != err1 {
			tmpnextnode = 0
		}

		if tmpnextnode >= 100 {
			sqlcmd4 := fmt.Sprintf("Select workid from scc_approver where workid= '%v'  and approvertype = %v", sqlresult2[0]["templateid"], tmpnextnode+1)
			sqlresult4 := sccinfo.tmpsql.SelectData(sqlcmd4)
			//fmt.Println(sqlcmd4)
			if 0 == len(sqlresult4) {
				//证明没人审批了 当前此条审批就是最后一条审批
				sqlcmd4 = fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode,filepath)values(%v,%v,'%v',%v,%v,%v,'%v');", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, 200, json.Filepath)
				sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
			} else {
				sqlcmd4 = fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode,filepath)values(%v,%v,'%v',%v,%v,%v,'%v');", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, tmpnextnode+1, json.Filepath)
				sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)

				sqlcmd5 := fmt.Sprintf("update scc_apply set status=1 where appid = %v;", json.Appid)
				sccinfo.tmpsql.Execsqlcmd(sqlcmd5, false)
			}
		} else {
			sqlcmd4 := fmt.Sprintf("Select workid from scc_approver where workid= '%v' and approvertype = %v", sqlresult2[0]["templateid"], tmpnextnode+1)
			sqlresult4 := sccinfo.tmpsql.SelectData(sqlcmd4)
			if 0 == len(sqlresult4) {
				//证明没人审批了 当前此条审批就是最后一条审批
				sqlcmd4 = fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode,filepath)values(%v,%v,'%v',%v,%v,%v,'%v');", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, 101, json.Filepath)
				sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
			} else {
				sqlcmd4 = fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode,filepath)values(%v,%v,'%v',%v,%v,%v,'%v');", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, tmpnextnode+1, json.Filepath)
				sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
			}
		}

	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccqueryworkflow(c *gin.Context) {
	type workflowid struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Appid string `json:"appid" binding:"required"`
	}
	var json workflowid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd2 := fmt.Sprintf("select appid,createtime,advise,templateid,appcurentnode,appnextnode,id,filepath from scc_workflow where appid= '%v'", json.Appid)
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)

	sqlcmd3 := fmt.Sprintf("select appid,templateid,createtime,textinfo,filepath,telephone,creater,grade from scc_apply where appid= '%v'", json.Appid)
	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd3)

	sqlcmd4 := fmt.Sprintf("select id,workflowid,createtime,comment,commentuser,filepath from scc_comment where appid= '%v'", json.Appid)
	sqlresult4 := sccinfo.tmpsql.SelectData(sqlcmd4)

	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"applinfo": sqlresult3, "workflowinfo": sqlresult2, "commentinfo": sqlresult4}})
}
func sccresearchworkflow(c *gin.Context) {
	type workid struct {
		Templateid string `json:"templateid" binding:"required"`
		Appid      string `json:"appid" binding:"required"`
		Advise     string `json:"advise" binding:"required"`
		Username   string `json:"username" binding:"required"`
		Grade      string `json:"grade"`
		Filepath   string `json:"filepath"`
	}
	//打回直接打回给研究人
	var json workid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd2 := fmt.Sprintf("Select approver,approvertype from scc_approver where workid= '%v'", json.Templateid)
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	lenresult := len(sqlresult2)
	if 0 == lenresult {
		c.JSON(http.StatusOK, gin.H{"result": "success"})
		return
	}
	sqlcmd := fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode,filepath)values(%v,%v,'%v',%v,%v,%v,'%v');", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, 1, 2, json.Filepath)

	sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
	if "" != json.Grade {
		sqlcmd1 := fmt.Sprintf("update scc_apply set grade=%v where appid=%v;", json.Grade, json.Appid)

		sccinfo.tmpsql.Execsqlcmd(sqlcmd1, false)
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func scchandleworkflow(c *gin.Context) {
	type workid struct {
		Templateid string `json:"templateid" binding:"required"`
		Appid      string `json:"appid" binding:"required"`
		Advise     string `json:"advise" binding:"required"`
		Username   string `json:"username" binding:"required"`
		Filepath   string `json:"filepath"`
	}
	//打回直接打回给研究人
	var json workid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd2 := fmt.Sprintf("Select approver,approvertype from scc_approver where workid= '%v'", json.Templateid)
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	lenresult := len(sqlresult2)
	if 0 == lenresult {
		c.JSON(http.StatusOK, gin.H{"result": "success"})
		return
	}
	sqlcmd := fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode，filepath)values(%v,%v,'%v',%v,%v,%v,'%v');", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, 101, 102, json.Filepath)

	sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccqueryapprolist(c *gin.Context) {
	type approlist struct {
		Approlistid string `json:"approlistid" binding:"required"`
	}
	var json approlist
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd2 := fmt.Sprintf("Select id,approver,approvertype from scc_approver where workid= '%v'", json.Approlistid)
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult2})
}
func sccquerycclist(c *gin.Context) {
	type cclist struct {
		Cclistid string `json:"cclistid" binding:"required"`
	}
	var json cclist
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd2 := fmt.Sprintf("Select id,cc  from scc_cc where workid= '%v'", json.Cclistid)
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult2})
}
func sccqueryuserinfo(c *gin.Context) {
	type useridinfo struct {
		Userid string `json:"userid" binding:"required"`
	}
	var json useridinfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd2 := fmt.Sprintf("Select id,usertype,phone,telephone,username,dsec from scc_approver where workid= '%v'", json.Userid)
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": sqlresult2})
}
func sccaddaddcomment(c *gin.Context) {
	type comment struct {
		Workflowid  string `json:"workflowid" binding:"required"`
		Comment     string `json:"comment" binding:"required"`
		Appid       string `json:"appid" binding:"required"`
		Commentuser string `json:"commentuser" binding:"required"`
		Filepath    string `json:"filepath"`
	}
	var json comment
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tempfilepath := ""
	if "" == json.Filepath {
		tempfilepath = " "
	} else {
		tempfilepath = json.Filepath
	}
	sqlcmd := fmt.Sprintf("insert into scc_comment (workflowid,createtime,comment,appid,commentuser,filepath)values(%v,%v,'%v','%v','%v','%v');", json.Workflowid, time.Now().Unix(), sccsql.MysqlRealEscapeString(json.Comment), json.Appid, json.Commentuser, tempfilepath)

	sccinfo.tmpsql.Execsqlcmd(sqlcmd, true)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccaddsatisfaction(c *gin.Context) {
	type comment struct {
		Satisfaction string `json:"satisfaction" binding:"required"`
		Comment      string `json:"comment"`
		Commentuser  string `json:"commentuser"`
	}
	var json comment
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if "yes" == json.Satisfaction {
		sqlcmd := fmt.Sprintf("insert into scc_statistica (appid,createtime,comment,commentuser)values(%v,%v,'%v','%v');", 0, time.Now().Unix(), json.Comment, json.Commentuser)

		sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
		c.JSON(http.StatusOK, gin.H{"result": "success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不满意请走工单提交流程"})
		return
	}
}

func sccscccreatetemplate(c *gin.Context) {
	type modtemplatework struct {
		Templateid string `json:"templateid" binding:"required"`
		Active     string `json:"active" binding:"required"`
	}
	var json modtemplatework
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Tmactive := 0
	if "yes" == json.Active {
		Tmactive = 1
	}
	sqlcmd1 := fmt.Sprintf("update scc_worktempplate set active = '%v' where workid = %v", Tmactive, json.Templateid)
	sccinfo.tmpsql.Execsqlcmd(sqlcmd1, false)
	c.JSON(http.StatusOK, gin.H{"result": "success"})

}
func setrouter(r *gin.Engine) {
	r.POST("/login", sccworklogin)
	r.POST("/createtemplate", scccreatetemplate)
	r.POST("/querytemplate", sccquerytemplate)
	r.POST("/querytemplateappandcc", sccquerytemplateappandcc) //这个接口估计是废物
	r.POST("/createapply", scccreateapply)
	r.POST("/queryapply", sccqueryapply)
	r.POST("/researchworkflow", sccresearchworkflow)
	r.POST("/querytmplatebyworkid", sccuerytmplatebyworkid)
	r.POST("/approve", sccapprove)
	r.POST("/queryworkflow", sccqueryworkflow)
	r.POST("/handleworkflow", scchandleworkflow)
	r.POST("/queryapprolist", sccqueryapprolist)
	r.POST("/querycclist", sccquerycclist)
	r.POST("/queryuserinfo", sccqueryuserinfo)
	r.POST("/addcomment", sccaddaddcomment)
	r.POST("/addsatisfaction", sccaddsatisfaction)
	r.POST("/modtemplate", sccscccreatetemplate)
}
