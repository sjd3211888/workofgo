package zlmediahook

import (
	"encoding/json"
	"fmt"
	sccsql "golearn/gomysql"
	sccredis "golearn/goredis"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

type Reportroom func(msg map[string]interface{})

type coreinfo struct {
	tmpsql   sccsql.Mysqlconnectpool
	tmpredis sccredis.Redisconnectpool
	reportup Reportroom
}

var sccinfo coreinfo
var Httprequesturl string
var Secret string

func Setcallback(callback Reportroom) {
	sccinfo.reportup = callback
}
func init() {

	var conf map[string]map[string]string
	if _, err := toml.DecodeFile("./sccconfig.toml", &conf); err != nil {
		// handle error
	}
	Host := conf["liveroom"]["Host"]
	Username := conf["liveroom"]["Username"]
	Password := conf["liveroom"]["Password"]
	Dbname := conf["liveroom"]["Dbname"]
	Port := conf["liveroom"]["Port"]
	iport, _ := strconv.Atoi(Port)
	Serhost := conf["liveroom"]["Httpserverhost"]
	Redisip := conf["liveroom"]["Redisip"]
	Httprequesturl = conf["liveroom"]["Httprequesturl"]
	Secret = conf["liveroom"]["Secret"]
	fmt.Println(Serhost)
	go func(Host string, Username string, Password string, Dbname string, Serhost string, Redisip string, iport int) {
		sccinfo.tmpsql.Initmysql(Host, Username, Password, Dbname, iport)
		sccinfo.tmpredis.Redisip = (Redisip)
		sccinfo.tmpredis.ConnectRedis()
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		setrouter(r)
		if err := r.Run(Serhost); err != nil {
			fmt.Println("startup service failed, err:\n", err)
		}
	}(Host, Username, Password, Dbname, Serhost, Redisip, iport)

}
func onpublish(c *gin.Context) {

	type onpublish struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		App    string `json:"app" `
		Id     string `json:"id" `
		Ip     string `json:"ip" `
		Params string `json:"params"`
		Port   int    `json:"port" `
		Vchema string `json:"schema" `
		Vtream string `json:"stream" `
		Vhost  string `json:"vhost"`
	}
	var json onpublish
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(json)
	c.JSON(http.StatusOK, gin.H{"msg": "success", "code": 0, "enableHls": false, "enableMP4": true, "enableRtxp": true})

}
func onstreamnonereader(c *gin.Context) {
	type noreader struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		App    string `json:"app" `
		Schema string `json:"schema" `
		Stream string `json:"stream" `
		Vhost  string `json:"vhost"`
	}
	var json noreader
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"close": false, "code": 0})
}
func onplay(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "success", "code": 0})
}
func onflowreport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "success", "code": 0})
}
func onstreamchanged(c *gin.Context) {
	/*type streamchang struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		App    string `json:"app" `
		Regist bool   `json:"regist" `
		Schema string `json:"schema" `
		Stream string `json:"stream"`
		Vhost  string `json:"vhost" `
	}
	var json streamchang
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}*/
	//c.Request.Body
	body, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("---body/--- \r\n " + string(body))

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &dat); err == nil {
		fmt.Println("3333333333333333333")
		sccinfo.reportup(dat)
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success", "code": 0})

}
func serverstarted(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "success", "code": 0})
}
func onrecordmp4(c *gin.Context) {
	type recordinfo struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		App       string `json:"app" `
		Filename  string `json:"file_name" `
		Filepath  string `json:"file_path" `
		Filesize  int    `json:"file_size"`
		Folder    string `json:"folder" `
		Starttime int    `json:"start_time" `
		Stream    string `json:"stream" `
		Timelen   int    `json:"time_len"`
		URL       string `json:"url"`
		Vost      string `json:"vhost"`
	}
	var json recordinfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success", "code": 0})
}
func onstreamnotfound(c *gin.Context) {
	type streamnotfound struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		App    string `json:"app" `
		Id     string `json:"id"`
		Ip     string `json:"ip"`
		Params string `json:"params"`
		Port   int    `json:"port" `
		Schema string `json:"schema" `
		Stream string `json:"stream" `
		Vhost  string `json:"vhost"`
	}
	var json streamnotfound
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success", "code": 0})
}
func getMediaList(c *gin.Context) {
	type listinfo struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		App    string `json:"app"`
		Schema string `json:"schema"  binding:"required"`
	}
	var json1 listinfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json1); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	host := Httprequesturl + "index/api/getMediaList?"
	client := http.Client{}
	q := url.Values{}
	q.Set("schema", json1.Schema)
	if "" != json1.App {
		q.Set("app", json1.App)
	}
	q.Set("secret", Secret)
	req, _ := http.NewRequest("POST", host+q.Encode(), nil)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &dat); err == nil {
		//sccinfo.reportup(dat)
	}

	c.JSON(http.StatusOK, gin.H{"msg": "success", "data": gin.H{"listinfo": dat}})
}
func closestreams(c *gin.Context) {
	type closeinfo struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		App    string `json:"app"`
		Schema string `json:"schema"  binding:"required"`
		Stream string `json:"stream"`
		Force  string `json:"force"`
	}
	var json1 closeinfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json1); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	host := Httprequesturl + "index/api/close_streams?"
	client := http.Client{}
	q := url.Values{}
	q.Set("schema", json1.Schema)
	if "" != json1.App {
		q.Set("app", json1.App)
	}
	if "" != json1.Stream {
		q.Set("stream", json1.Stream)
	}

	q.Set("force", json1.Force)
	q.Set("secret", Secret)
	req, _ := http.NewRequest("POST", host+q.Encode(), nil)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &dat); err == nil {
	}

	c.JSON(http.StatusOK, gin.H{"msg": "success", "data": gin.H{"listinfo": dat}})
}
func addStreamProxy(c *gin.Context) {
	type closeinfo struct {
		// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
		App       string `json:"app"`
		Schema    string `json:"schema"  binding:"required"`
		Stream    string `json:"stream" binding:"required"`
		URL       string `json:"url" binding:"required"`
		Enablemp4 int64  `json:"enable_mp4"`
		Rtptype   int64  `json:"rtp_type"`
	}
	var json1 closeinfo
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&json1); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	host := Httprequesturl + "index/api/addStreamProxy?"
	client := http.Client{}
	q := url.Values{}
	q.Set("schema", json1.Schema)
	if "" != json1.App {
		q.Set("app", json1.App)
	}
	if "" != json1.Stream {
		q.Set("stream", json1.Stream)
	}
	q.Set("url", json1.URL)
	enablemp4 := strconv.FormatInt(json1.Enablemp4, 10)
	rtptype := strconv.FormatInt(json1.Rtptype, 10)
	q.Set("enable_mp4", enablemp4)
	q.Set("rtp_type", rtptype)
	q.Set("secret", Secret)
	q.Set("vhost", "__defaultVhost__")
	q.Set("enable_rtsp", "1")
	q.Set("enable_rtmp", "1")
	req, _ := http.NewRequest("POST", host+q.Encode(), nil)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &dat); err == nil {
	}

	c.JSON(http.StatusOK, gin.H{"msg": "success", "data": gin.H{"listinfo": dat}})
}
func setrouter(r *gin.Engine) {
	r.POST("/index/hook/on_server_started", serverstarted)
	r.POST("/index/hook/on_publish", onpublish)
	r.POST("/index/hook/on_stream_none_reader", onstreamnonereader)
	r.POST("/index/hook/on_play", onplay)
	r.POST("/index/hook/on_flow_report", onflowreport)
	r.POST("/index/hook/on_stream_changed", onstreamchanged)
	r.POST("/index/hook/on_record_mp4", onrecordmp4)
	r.POST("/index/hook/on_stream_not_found", onstreamnotfound)
	r.POST("/index/api/getMediaList", getMediaList)
	r.POST("/index/api/closestreams", closestreams)
	r.POST("/index/api/addStreamProxy", addStreamProxy)
}
