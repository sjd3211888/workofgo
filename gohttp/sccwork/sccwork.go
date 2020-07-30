package sccwork

import (
	"fmt"
	sccsql "golearn/gomysql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type coreinfo struct {
	tmpsql sccsql.Mysqlconnectpool
}

var sccinfo coreinfo

func init() {
	go func() {
		sccinfo.tmpsql.Initmysql("127.0.0.1", "root", "root", "sccwork", 3306)
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		setrouter(r)
		if err := r.Run(":9980"); err != nil {
			fmt.Println("startup service failed, err:%v\n", err)
		}
	}()

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
		Trade        int            `json:"trade" binding:"required"`
		Ptype        int            `json:"type" binding:"required"`
		Creater      int            `json:"creater" binding:"required"`
		Templateuer  int            `json:"templateuer" binding:"required"`
		Resrearcher  int            `json:"resrearcher" binding:"required"`
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
		for k, v := range json.Apporlist {
			sqlcmd := fmt.Sprintf("insert into scc_approver (workid,approver,approvertype)values(%v,%v,%v);", workid, v.Apporver, k+2)

			sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
			sqlcmd = fmt.Sprintf("insert into scc_approver (workid,approver,approvertype)values(%v,%v,%v);", workid, v.Apporver, k+2+100)

			sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
		}

	}
	countcclist := len(json.Cclist)
	if 0 == countcclist {
		fmt.Println("CC list is null")
	} else {
		for _, v := range json.Cclist {
			sqlcmd := fmt.Sprintf("insert into scc_cc(workid,cc)values(%v,%v);", workid, v.Cc)
			sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
		}

	}
	sqlcmd := fmt.Sprintf("insert into scc_worktempplate (workid,templatename,trade,p_type,createtime,approverlist,cclist,creater,templateuser,researcher)values(%v,'%v',%v,%v,%v,%v,%v,%v,%v,%v);", workid, json.Templatename, json.Trade, json.Ptype, workid, workid, workid, json.Creater, json.Templateuer, json.Resrearcher)
	fmt.Println(sqlcmd)
	sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
	fmt.Println(json.Apporlist, json.Cclist)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"templateid": workid}})
}

func sccquerytemplate(c *gin.Context) {
	type querytemplatework struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Usertype int `json:"usertype" binding:"required"` // 1 模板创建者 2模板使用者
		Username int `json:"useid" binding:"required"`
	}
	var json querytemplatework
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if json.Usertype == 1 {
		sqlcmd1 := fmt.Sprintf("Select workid,templatename,trade,p_type,createtime,approverlist,cclist,creater,templateuser from scc_worktempplate where creater= '%v'", json.Username)
		sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
		c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"templateinfo": sqlresult1}})
	} else if json.Usertype == 2 {
		sqlcmd1 := fmt.Sprintf("Select workid,templatename,trade,p_type,createtime,approverlist,cclist,creater,templateuser from scc_worktempplate where templateuser= '%v'", json.Username)
		sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
		c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"templateinfo": sqlresult1}})
	}

}
func scccreateapply(c *gin.Context) {
	type querytemplatework struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Templateid   int    `json:"templateid" binding:"required"`
		Username     int    `json:"username" binding:"required"`
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

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccqueryworktempplatebyworkid(workid int) []map[string]string {
	sqlcmd2 := fmt.Sprintf("Select workid,templatename,trade,p_type,creater from scc_worktempplate where workid= '%v'", workid)
	return sccinfo.tmpsql.SelectData(sqlcmd2)

}
func sccqueryapproveryworkid(workid int) []map[string]string {
	sqlcmd3 := fmt.Sprintf("Select approver,approvertype from scc_approver where workid= '%v'", workid)
	return sccinfo.tmpsql.SelectData(sqlcmd3)
}
func sccqueryccworkid(workid int) []map[string]string {
	sqlcmd4 := fmt.Sprintf("Select cc from scc_cc where workid= '%v'", workid)
	return sccinfo.tmpsql.SelectData(sqlcmd4)
}
func sccqueryapply(c *gin.Context) {
	type queryapply struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Applytype int `json:"type" binding:"required"` //4 我创建的  1  带我处理  2 已处理 3 抄送给我的 5 和我相关的 6 带我研究的 7 待执行
		Username  int `json:"username"`
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
	case 4:
		{
			sqlcmd := fmt.Sprintf("Select appid,templateid,createtime,textinfo,filepath,telephone,creater from scc_apply where creater= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "appidinfo": ""})
				break
			}

			/*	workflworkresult := make([]map[string]string, 0)
				for k := range sqlresult {
					sqlcmd1 := fmt.Sprintf("Select id,appid,createtime,apperover,apperovertypei,advise,templateidm,appcurentnode,appnextnode from scc_workflow where appid = '%v'", sqlresult[k]["appid"])
					sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
					workflworkresult = append(workflworkresult, sqlresult1...)
				}*/
			c.JSON(http.StatusOK, gin.H{"result": "success", "appidinfo": sqlresult})
			//需要加开一个接口 通过appid查全部
			break
		}
	case 1:
		{
			sqlcmd := fmt.Sprintf("Select workid,approver,approvertype from scc_approver where approver= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": ""})
				break
			}
			workflworkresult := make([]map[string]string, 0)
			for k := range sqlresult {
				sqlcmd1 := fmt.Sprintf("Select appid from scc_workflow where templateid = %v and appnextnode = %v ", sqlresult[k]["workid"], sqlresult[k]["approvertype"])
				fmt.Println(sqlcmd1)
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				for j := range sqlresult1 {
					sqlcmd2 := fmt.Sprintf("Select appid from scc_workflow where templateid = %v and appcurentnode = %v and appid=%v ", sqlresult[k]["workid"], sqlresult[k]["approvertype"], sqlresult1[j]["appid"])
					sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
					//fmt.Println(sqlcmd2, k)
					if 0 == len(sqlresult2) {
						sqlcmd6 := fmt.Sprintf("Select appid,templateid,createtime,textinfo,filepath,telephone,creater from scc_apply where appid= %v", sqlresult1[j]["appid"])
						sqlresult6 := sccinfo.tmpsql.SelectData(sqlcmd6)
						workflworkresult = append(workflworkresult, sqlresult6...)
					}
				}
			}
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": workflworkresult}})
			break
		}
	case 2:
		{
			sqlcmd := fmt.Sprintf("Select workid,approver,approvertype from scc_approver where approver= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": ""})
				break
			}
			workflworkresult := make([]map[string]string, 0)
			for k := range sqlresult {
				sqlcmd1 := fmt.Sprintf("Select appid from scc_workflow where templateid = %v and appcurentnode = %v ", sqlresult[k]["workid"], sqlresult[k]["approvertype"])
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				if 0 != len(sqlresult1) {
					sqlcmd6 := fmt.Sprintf("Select appid,templateid,createtime,textinfo,filepath,telephone,creater from scc_apply where appid= %v", sqlresult1[k]["appid"])
					sqlresult6 := sccinfo.tmpsql.SelectData(sqlcmd6)
					workflworkresult = append(workflworkresult, sqlresult6...)
				}

			}
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": workflworkresult}})
			break
		}
	case 3:
		{
			sqlcmd := fmt.Sprintf("Select workid from scc_cc where cc=%v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			cccount := len(sqlresult)
			workflworkresult := make([]map[string]string, 0)
			for k := 0; k < cccount; k++ {
				tmpworkid := sqlresult[k]["workid"]
				inttmpworkid, _ := strconv.Atoi(tmpworkid)
				sqlcmd1 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = %v", inttmpworkid)
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				if 0 != len(sqlresult1) {
					sqlcmd6 := fmt.Sprintf("Select appid,templateid,createtime,textinfo,filepath,telephone,creater from scc_apply where appid= %v", sqlresult1[k]["appid"])
					sqlresult6 := sccinfo.tmpsql.SelectData(sqlcmd6)
					workflworkresult = append(workflworkresult, sqlresult6...)
				}
				fmt.Println(workflworkresult)
			}
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": workflworkresult}})
			break
		}
	case 5:
		{

			//to do
			sqlcmd := fmt.Sprintf("Select workid from scc_approver where approver= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": ""})
				break
			}
			workflworkresult := make([]map[string]string, 0)
			for k := range sqlresult {
				tmpappid, _ := strconv.Atoi(sqlresult[k]["workid"])
				sqlcmd1 := fmt.Sprintf("Select appid,templateid,createtime,textinfo,filepath,telephone,creater  from scc_apply where templateid = %v ", tmpappid)
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				if 0 != len(sqlresult1) {
					workflworkresult = append(workflworkresult, sqlresult1...)
				}
			}
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflow": workflworkresult}})
			break
		}
	case 6:
		{
			sqlcmd := fmt.Sprintf("Select workid from scc_worktempplate where researcher= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflow": sqlresult}})
				break
			}
			workflworkresult := make([]map[string]string, 0)
			for k := range sqlresult {
				sqlcmd1 := fmt.Sprintf("Select id,appid,templateid from scc_workflow where templateid = %v and appnextnode = 1 ", sqlresult[k]["workid"])
				fmt.Println("wocao ", sqlcmd1)
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				for j := range sqlresult1 {
					sqlcmd2 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = %v and appnextnode = 2 and appid=%v order by id desc limit 1", sqlresult[k]["workid"], sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
					sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
					fmt.Println("wocao1 ", sqlcmd2)
					if 0 == len(sqlresult2) {
						fmt.Println("SS", sqlcmd2)
						sqlcmd5 := fmt.Sprintf("Select appid,templateid,createtime,textinfo,filepath,telephone,creater from scc_apply where appid=%v", sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
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
							fmt.Println(tmpnext, tmpproid)
							//不加了
						} else {
							fmt.Println(tmpnext, tmpproid)
							sqlcmd5 := fmt.Sprintf("Select appid,templateid,createtime,textinfo,filepath,telephone,creater from scc_apply where appid=%v", sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
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
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workapplyinfo": workflworkresult}})
			break
		}
	case 7:
		{
			sqlcmd := fmt.Sprintf("Select workid from scc_worktempplate where researcher= %v", json.Username)
			sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
			if 0 == len(sqlresult) {
				c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"workflowinfo": sqlresult}})
				break
			}
			workflworkresult := make([]map[string]string, 0)
			for k := range sqlresult {
				sqlcmd1 := fmt.Sprintf("Select id,appid,templateid from scc_workflow where templateid = %v and appnextnode = 101 ", sqlresult[k]["workid"])
				fmt.Println("wocao ", sqlcmd1)
				sqlresult1 := sccinfo.tmpsql.SelectData(sqlcmd1)
				for j := range sqlresult1 {
					sqlcmd2 := fmt.Sprintf("Select id,appid from scc_workflow where templateid = %v and appnextnode = 102 and appid=%v order by id desc limit 1", sqlresult[k]["workid"], sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
					sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)
					fmt.Println("wocao1 ", sqlcmd2)
					if 0 == len(sqlresult2) {
						fmt.Println("SS", sqlcmd2)
						sqlcmd5 := fmt.Sprintf("Select appid,templateid,createtime,textinfo,filepath,telephone,creater from scc_apply where appid=%v", sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
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
							fmt.Println(tmpnext, tmpproid)
							//不加了
						} else {
							fmt.Println(tmpnext, tmpproid)
							sqlcmd5 := fmt.Sprintf("Select appid,templateid,createtime,textinfo,filepath,telephone,creater from scc_apply where appid=%v", sqlresult1[j]["appid"]) //查询这个appid 中最大的当前执行研究的最大ID
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
		Workid int `json:"workid" binding:"required"` //0 我创建的  1 已审批的  2 待我审批的  3 抄送给我的 5 和我相关的
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
		Workid int `json:"workid" binding:"required"` //0 我创建的  1 已审批的  2 待我审批的  3 抄送给我的
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
		Templateid int    `json:"templateid" binding:"required"`
		Appid      int    `json:"appid" binding:"required"`
		Status     int    `json:"status" binding:"required"`
		Advise     string `json:"advise" binding:"required"`
		Username   int    `json:"username" binding:"required"`
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
	if 2 == json.Status { //2 是打回
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
		fmt.Println(sqlcmd3)
		sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd3)
		if 0 == len(sqlresult3) {
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": "no user current2"})
			return
		}

		if tmpnextnode >= 100 {

			sqlcmd4 := fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode)values(%v,%v,'%v',%v,%v,%v);", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, 101)
			sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
		} else {

			sqlcmd4 := fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode)values(%v,%v,'%v',%v,%v,%v);", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, 1)
			sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
		}

	} else if 1 == json.Status {
		sqlcmd2 := fmt.Sprintf("Select templateid,appnextnode,appid,id from scc_workflow where templateid= '%v' and  appid = %v order by id desc limit 1", json.Templateid, json.Appid)
		sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)

		lenresult := len(sqlresult2)
		if 0 == lenresult {
			c.JSON(http.StatusOK, gin.H{"result": "success", "data": "no user current1"})
		}
		tmpapptype, err := strconv.Atoi(sqlresult2[0]["appnextnode"])

		if nil != err {
			tmpapptype = 0
		}
		sqlcmd3 := fmt.Sprintf("Select workid from scc_approver where workid= %v and  approver = %v and approvertype = %v ", sqlresult2[0]["templateid"], json.Username, tmpapptype)
		fmt.Println(sqlcmd3)
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
			fmt.Println(sqlcmd4)
			if 0 == len(sqlresult4) {
				//证明没人审批了 当前此条审批就是最后一条审批
				sqlcmd4 = fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode)values(%v,%v,'%v',%v,%v,%v);", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, 200)
				sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
			} else {
				sqlcmd4 = fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode)values(%v,%v,'%v',%v,%v,%v);", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, tmpnextnode+1)
				sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
			}
		} else {
			sqlcmd4 := fmt.Sprintf("Select workid from scc_approver where workid= '%v' and approvertype = %v", sqlresult2[0]["templateid"], tmpnextnode+1)
			sqlresult4 := sccinfo.tmpsql.SelectData(sqlcmd4)
			if 0 == len(sqlresult4) {
				//证明没人审批了 当前此条审批就是最后一条审批
				sqlcmd4 = fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode)values(%v,%v,'%v',%v,%v,%v);", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, 101)
				sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
			} else {
				sqlcmd4 = fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode)values(%v,%v,'%v',%v,%v,%v);", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, tmpnextnode, tmpnextnode+1)
				sccinfo.tmpsql.Execsqlcmd(sqlcmd4, false)
			}
		}

	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccqueryworkflow(c *gin.Context) {
	type workflowid struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Appid int `json:"appid" binding:"required"`
	}
	var json workflowid
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlcmd2 := fmt.Sprintf("select appid,createtime,advise,templateid,appcurentnode,appnextnode,id from scc_workflow where appid= '%v'", json.Appid)
	sqlresult2 := sccinfo.tmpsql.SelectData(sqlcmd2)

	sqlcmd3 := fmt.Sprintf("select appid,templateid,createtime,textinfo,filepath,telephone,creater from scc_apply where appid= '%v'", json.Appid)
	sqlresult3 := sccinfo.tmpsql.SelectData(sqlcmd3)
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"applinfo": sqlresult3, "workflowinfo": sqlresult2}})
}
func sccresearchworkflow(c *gin.Context) {
	type workid struct {
		Templateid int    `json:"templateid" binding:"required"`
		Appid      int    `json:"appid" binding:"required"`
		Advise     string `json:"advise" binding:"required"`
		Username   int    `json:"username" binding:"required"`
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
	sqlcmd := fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode)values(%v,%v,'%v',%v,%v,%v);", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, 1, 2)

	sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func scchandleworkflow(c *gin.Context) {
	type workid struct {
		Templateid int    `json:"templateid" binding:"required"`
		Appid      int    `json:"appid" binding:"required"`
		Advise     string `json:"advise" binding:"required"`
		Username   int    `json:"username" binding:"required"`
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
	sqlcmd := fmt.Sprintf("insert into scc_workflow (appid,createtime,advise,templateid,appcurentnode,appnextnode)values(%v,%v,'%v',%v,%v,%v);", json.Appid, time.Now().Unix(), json.Advise, json.Templateid, 101, 102)

	sccinfo.tmpsql.Execsqlcmd(sqlcmd, false)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccqueryapprolist(c *gin.Context) {
	type approlist struct {
		Approlistid int `json:"approlistid" binding:"required"`
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
		Cclistid int `json:"cclistid" binding:"required"`
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
		Userid int `json:"userid" binding:"required"`
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
}
