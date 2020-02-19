package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"simpleapi/app/api"
	"time"
)

// 设置hook当请求结束输出请求的相关信息
func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		start_time := time.Now()
		c.Next()
		end_time := time.Now()
		delay := end_time.Sub(start_time)
		fmt.Printf("请求页面：%s, 请求IP地址为：%s, 请求方式为：%s, 执行时长为：%s\n",c.Request.URL, c.ClientIP(), c.Request.Method, delay)
	}
}

// 初始化路由
func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(MiddleWare())
	v1 := router.Group("/label")
	{
		v1.GET("/LabelAll", api.GetLabel)
		v1.GET("/Label_Field/:field", api.GetLabelByField)
	}
	return router
}
