package main

import (
	"fmt"
	sccsql "golearn/gomysql"
	sccredis "golearn/goredis"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1.创建路由
	// 默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()
	var testsccmysql sccsql.Mysqlconnectpool
	testsccmysql.Initmysql("192.168.1.124", "root", "root", "freeswitch", 3306)
	test := testsccmysql.SelectData("select username,sex from person")
	for k, v := range test {
		fmt.Println("sssss is vvv is ", k, v)
	}

	var sjdtest sccredis.Redisconnectpool
	sjdtest.Redisip = "192.168.1.124:6379"
	sjdtest.Redispassword = "123456"
	sjdtest.ConnectRedis()
	rr, _ := sjdtest.SccredisGetmembernearby("118.32155", "31.123", 1000)
	sjdtest.SccredisGetmembergpsinf(rr)

	shoppingGroup := r.Group("/shopping")
	{
		shoppingGroup.GET("/index", shopIndexHandler)
		shoppingGroup.GET("/home", shopHomeHandler)
	}
	r.Run(":8000")
}

func shopIndexHandler(c *gin.Context) {
	time.Sleep(5 * time.Second)
}

func shopHomeHandler(c *gin.Context) {
	time.Sleep(3 * time.Second)
}
