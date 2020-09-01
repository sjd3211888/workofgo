package fshttp

import (
	"fmt"
	"net/http"

	. "golearn/freeswitch/fseslsql"
	fstoesl "golearn/freeswitch/fstoesl"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

var sccfsinfo *fstoesl.Fseslinfo

func Setsccfsinfo(info *fstoesl.Fseslinfo) {
	sccfsinfo = info
}
func init() {

	var conf map[string]map[string]string
	if _, err := toml.DecodeFile("./sccconfig.toml", &conf); err != nil {
		// handle error
	}
	Serhost := conf["sccfs"]["Httpserverhost"]
	go func(Serhost string) {
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		setrouter(r)
		if err := r.Run(Serhost); err != nil {
			fmt.Println("startup service failed, err:\n", err)
		}
	}(Serhost)

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
	c.JSON(http.StatusOK, gin.H{"result": "success"})
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
	c.JSON(http.StatusOK, gin.H{"result": "success"})
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
	c.JSON(http.StatusOK, gin.H{"result": "success"})
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
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccGetuseronline(c *gin.Context) {
	ret := Getuseronline()
	c.JSON(http.StatusOK, gin.H{"result": "success", "data": ret})
}
func sccBrocastuser(c *gin.Context) {
	//ret := Getuseronline()
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

/*func sccSendIMmessage(c *gin.Context) {
	//ret := Getuseronline()
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}*/
func sccMonitorcall(c *gin.Context) {
	//ret := Getuseronline()
	type monitorsipcall struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Monitoruid string `json:"monitoruid" binding:"required"`
		Calluuid   string `json:"calluuid" binding:"required"`
		Audiocall  string `json:"audiocall" binding:"required"`
	}
	var json monitorsipcall
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sccfsinfo.Eslmonitorcall(json.Monitoruid, json.Calluuid)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccInsertcall(c *gin.Context) {
	type insertsipcall struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Insertsipcalluid string `json:"insertsipcalluid" binding:"required"`
		Calluuid         string `json:"calluuid" binding:"required"`
		Audiocall        string `json:"audiocall" binding:"required"`
	}
	var json insertsipcall
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sccfsinfo.Eslinsertcall(json.Insertsipcalluid, json.Calluuid)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccDisassembecall(c *gin.Context) {
	type sipcalldisassembely struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Replaceuid string `json:"replaceuid" binding:"required"`
		Calluuid   string `json:"calluuid" binding:"required"`
		Aleg       string `json:"aleg" binding:"required"`
	}
	var json sipcalldisassembely
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ialeg := 0
	if "yes" == json.Aleg {
		ialeg = 1
	}
	sccfsinfo.Sccdisassembecall(json.Replaceuid, ialeg, json.Calluuid)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccTranfercall(c *gin.Context) {

	type tranfersipcall struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Tranferedid string `json:"tranferedid" binding:"required"`
		Calluuid    string `json:"calluuid" binding:"required"`
	}
	var json tranfersipcall
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sccfsinfo.Esltranfercall(json.Tranferedid, json.Calluuid)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func sccRecordcall(c *gin.Context) {
	type recordCall struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		Filepath string `json:"filepath" binding:"required"`
		Calluuid string `json:"calluuid" binding:"required"`
		Bstart   string `json:"bstart" binding:"required"`
	}
	var json recordCall
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	istart := 0
	if "yes" == json.Bstart {
		istart = 1
	}
	sccfsinfo.Eslrecordcall(json.Calluuid, json.Filepath, istart)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
func setrouter(r *gin.Engine) {
	r.POST("/Hangupuser", scchangupuser)
	r.POST("/Servercall", sccServercall)
	r.POST("/Monitoruser", sccMonitoruser)
	r.POST("/Yellinguser", sccYellinguser)
	r.POST("/Getuseronline", sccGetuseronline)
	r.POST("/Brocastuser", sccBrocastuser)
	//	r.POST("/SendIMmessage", sccSendIMmessage)
	r.POST("/Monitorcall", sccMonitorcall)
	r.POST("/Insertcall", sccInsertcall)
	r.POST("/Disassembecall", sccDisassembecall)
	r.POST("/Tranfercall", sccTranfercall)
	r.POST("/Recordcall", sccRecordcall)

}
