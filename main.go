package main

import (
	"container/list"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func isValid(s string) bool {
	stack := list.New()
	//b := []byte(s)
	for _, v := range s {

		if '{' == v {
			stack.PushBack('}')
		} else if '[' == v {
			stack.PushBack(']')
		} else if '(' == v {
			stack.PushBack(')')
		} else if ')' == v || ']' == v || '}' == v {
			if 0 == stack.Len() {
				return false
			}
			element := stack.Back()
			if v == element.Value {
				stack.Remove(element)
			} else {
				return false
			}
		}

	}
	if stack.Len() != 0 {
		return false
	}
	return true
}
func twoSum(nums []int, target int) []int {
	xxx := []int{0, 0}
	for s, v := range nums {
		xxx[0] = v
		for index, value := range nums {
			if index == s {
				continue
			} else {
				if target == value+v {
					xxx[1] = value
					return xxx
				}
			}
		}
	}
	return xxx
}

func longestPalindrome(s string) string {
	stirnglen := len(s)
	bstatt := false
	var mid int
	mid = 0
	sjdlen := 0
	maxmid := 0
	maxsjdle := 0
	for it := range s {
		if it == 0 || it == stirnglen-1 {
			continue
		} else {
			mid = it
			q := 0
			for it = mid; it-q >= 0 && it+q < stirnglen; q++ {
				if s[mid-q] == s[mid+q] {
					sjdlen = q
				} else if s[mid-1] == s[mid] {
					if maxsjdle < 1 {
						bstatt = true
						maxmid = mid
					}
				} else {
					if maxsjdle < sjdlen {
						maxsjdle = sjdlen
						maxmid = mid - 1
					}
					break
				}
			}
		}
	}
	var sjd string
	fmt.Println("mid is len is", maxmid, maxsjdle)
	if bstatt {
		sjd = s[maxmid-maxsjdle-1 : maxmid+maxsjdle+1]
	} else {
		sjd = s[maxmid-maxsjdle : maxmid+maxsjdle+1]
	}

	return sjd
}

type bar struct {
	thingOne string
	thingTwo int
}

func init() {
	fmt.Println("THIs is first")
}
func xxxsjd(xxxjsd *bar) int {
	xxxjsd.thingTwo = 3
	return xxxjsd.thingTwo
}

func xxxsjd1(xxxjsd []bar) int {
	xxxjsd[7].thingTwo = 9
	return 4
}

type People interface {
	Speak(string) string
}

type Student struct{}

func (stu *Student) Speak(think string) (talk string) {
	if think == "sb" {
		talk = "你是个大帅比"
	} else {
		talk = "您好"
	}
	return
}

// 定义中间
func myTime(c *gin.Context) {
	start := time.Now()
	c.Abort()
	// 统计时间
	since := time.Since(start)
	fmt.Println("程序用时：", since)
}

func main() {
	// 1.创建路由
	// 默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()
	// 注册中间件
	r.Use(myTime)
	// {}为了代码规范
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

/*func main() {
	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", FirstMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})
	// 匹配 /user/geektutu
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	r.GET("/users", func(c *gin.Context) {
		//name := c.Query("name")
		role := c.DefaultQuery("role", "teacher")
		if "teacher" == role {
			time.Sleep(time.Second * 10)
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"status_code": http.StatusOK,
				"status":      "ok",
				"caonima": gin.H{
					"hhh": 1,
				},
			},
			"result": 1,
		})
	})

	r.POST("/upload1", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"username": "ss",
			"password": "zzzzz",
		})
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")

}*/

func write(ch chan string) {
	for {
		select {
		// 写数据
		case ch <- "hello":
			fmt.Println("write hello")
		default:
			fmt.Println("channel full")
		}

		time.Sleep(time.Millisecond * 500)
	}

}
func FirstMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("first middleware before next()")
		isAbort := c.Query("isAbort")
		bAbort, err := strconv.ParseBool(isAbort)
		if err != nil {
			fmt.Printf("is abort value err, value %s\n", isAbort)
			c.Abort() // (2)
		}
		if bAbort {
			fmt.Println("first middleware abort") //(3)
			c.Abort()
			//c.AbortWithStatusJSON(http.StatusOK, "abort is true")
			return
		} else {
			fmt.Println("first middleware doesnot abort") //(4)
			return
		}

		fmt.Println("first middleware aftessr next()")
	}
}
