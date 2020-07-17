package coreservice

import (
	"fmt"
	sccsql "golearn/gomysql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	var tmpsql sccsql.Mysqlconnectpool
	tmpsql.Initmysql("49.235.86.39", "root", "root", "SCC", 3306)

	r := gin.Default()
	setrouter(r)
	if err := r.Run(":9888"); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello www.topgoer.com!",
	})
}
func querysccdepartment(c *gin.Context) {

}
func querydepartmentuser(c *gin.Context) {

}
func querygroup(c *gin.Context) {

}
func queryuser(c *gin.Context) {

}
func querygroupuser(c *gin.Context) {

}
func queryofflinemsg(c *gin.Context) {

}
func queryRecnetSession(c *gin.Context) {

}
func reportgps(c *gin.Context) {

}
func querygps(c *gin.Context) {

}
func uerypersonhistoryim(c *gin.Context) {

}
func querygrouphistoryim(c *gin.Context) {

}
func moduserdetail(c *gin.Context) {

}
func querynearbyscc(c *gin.Context) {

}
func setrouter(r *gin.Engine) {
	r.GET("/topgoer", helloHandler)
	r.GET("/querydepartment", querysccdepartment)
	r.GET("/querydepartmentuser", querydepartmentuser)
	r.GET("/querygroup", querygroup)
	r.GET("/queryuser", queryuser)
	r.GET("/querygroupuser", querygroupuser)
	r.GET("/queryofflinemsg", queryofflinemsg)
	r.GET("/queryRecnetSession", queryRecnetSession)
	r.GET("/reportgps", reportgps)
	r.GET("/querygps", querygps)
	r.GET("/uerypersonhistoryim", uerypersonhistoryim)
	r.GET("/querygrouphistoryim", querygrouphistoryim)
	r.GET("/moduserdetail", moduserdetail)
	r.GET("/querynearbyscc", querynearbyscc)
}
