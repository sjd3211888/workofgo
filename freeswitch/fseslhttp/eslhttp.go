package fshttp

import (
	"fmt"
	"net/http"

	fstoesl "golearn/freeswitch/fstoesl"

	"github.com/gin-gonic/gin"
)

var sccfsinfo *fstoesl.Fseslinfo

func Setsccfsinfo(info *fstoesl.Fseslinfo) {
	sccfsinfo = info
}
func init() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		setrouter(r)
		if err := r.Run(":19980"); err != nil {
			fmt.Println("startup service failed, err:\n", err)
		}
	}()

}
func scchangupuser(c *gin.Context) {
	type hangupuser struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		UUID string `json:"uuid" binding:"required"`
	}
	var json hangupuser
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sccfsinfo.Hangupuser(json.UUID)
}
func sccServercall(c *gin.Context) {
	type hangupuser struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		CALLERL  string `json:"Callerl" binding:"required"`
		CALLERR  string `json:"callerr" binding:"required"`
		CALLEEID string `json:"calleeid" binding:"required"`
	}
	var json hangupuser
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sccfsinfo.Servercall(json.CALLERL, json.CALLERR, json.CALLEEID)
}
func sccMonitoruser(c *gin.Context) {
	type monitoruser struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		CALLERL  string `json:"Callerl" binding:"required"`
		CALLERR  string `json:"callerr" binding:"required"`
		CALLEEID string `json:"calleeid" binding:"required"`
	}
	var json monitoruser
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sccfsinfo.Monitoruser(json.CALLERL, json.CALLERR, json.CALLEEID)
}
func sccYellinguser(c *gin.Context) {
	type yellinguser struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		CALLERL  string `json:"Callerl" binding:"required"`
		CALLERR  string `json:"callerr" binding:"required"`
		CALLEEID string `json:"calleeid" binding:"required"`
	}
	var json yellinguser
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sccfsinfo.Yellinguser(json.CALLERL, json.CALLERR, json.CALLEEID)
}
func setrouter(r *gin.Engine) {
	r.POST("/hangupuser", scchangupuser)
	r.POST("/Servercall", sccServercall)
	r.POST("/Monitoruser", sccMonitoruser)
	r.POST("/Yellinguser", sccYellinguser)
}
