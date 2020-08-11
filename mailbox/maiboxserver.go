package main

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

func scclicense(c *gin.Context) {
	now := time.Now()        //获取当前时间
	timestamp1 := now.Unix() //时间戳
	sqlcmd := fmt.Sprintf("Select expire,count,updatedate,applicant,cause,licensor  from license_records where expire<2999")
	sqlresult := sccinfo.tmpsql.SelectData(sqlcmd)
	workflworkresult := make([]map[string]string, 0)
	for k := range sqlresult {
		formatTimeStr := sqlresult[k]["updatedate"]
		formatTime, err := time.Parse("2006-01-02 15:04:05", formatTimeStr)
		if err == nil {
			//fmt.Println(sqlresult[k]["expire"])

			tmpcontent := sqlresult[k]["expire"][0 : len(sqlresult[k]["expire"])-1] //由于它这个字符多带个了\n,所以去掉一位
			expire, _ := strconv.ParseInt(tmpcontent, 10, 64)

			//utc 我们是东八区 领先UTC 8个消息
			licenseexpiretime := formatTime.Unix() + (expire * 3600 * 24) - (3600 * 8)
			timeremain := licenseexpiretime - timestamp1
			if -15*3600*24 <= timeremain && timeremain <= 8*3600*24 {
				leavetime := strconv.FormatInt(timeremain/3600/24, 10)
				fmt.Println(timeremain, licenseexpiretime, timestamp1, expire)
				sqlresult[k]["leavetime"] = leavetime
				workflworkresult = append(workflworkresult, sqlresult[k])
			}
		}
	}

	var msgneedtosend string
	for index := range workflworkresult {
		tmpstring := fmt.Sprintf("licensor-->%v,updatedate-->%v,expire-->%v,count-->%v,applicant-->%v,cause-->%v,timeleave-->%v\r\n", workflworkresult[index]["licensor"], workflworkresult[index]["updatedate"], workflworkresult[index]["expire"], workflworkresult[index]["count"], workflworkresult[index]["applicant"], workflworkresult[index]["cause"], workflworkresult[index]["leavetime"])
		msgneedtosend = msgneedtosend + tmpstring
	}
	c.JSON(http.StatusOK, gin.H{"result": msgneedtosend})
}

func setrouter(r *gin.Engine) {
	r.GET("/queryscclicense", scclicense)

}
func main() {
	sccinfo.tmpsql.Initmysql("127.0.0.1", "root", "root", "license", 3306)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	setrouter(r)
	if err := r.Run(":10010"); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
