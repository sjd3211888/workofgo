package main

import (
	"fmt"
	sccsql "golearn/gomysql"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1.创建路由
	// 默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()
	Db := sccsql.ConnectMysql("root", "root", "192.168.1.124", 3306, "freeswitch", "utf8mb4")
	defer Db.Close()
	xxx := sccsql.SelectData(Db, "select username,sex from person")
	for k, v := range xxx {
		fmt.Println(k, v)
	}

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
