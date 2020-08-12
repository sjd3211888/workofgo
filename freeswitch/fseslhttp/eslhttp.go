package fshttp

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type sadas struct {
}

var xcxx sadas

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
func (ss sadas) scchangupuser(c *gin.Context) {

}
func setrouter(r *gin.Engine) {
	r.POST("/hangupuser", xcxx.scchangupuser)
}
