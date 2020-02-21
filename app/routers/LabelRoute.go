package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"simpleapi/app/api"
	"simpleapi/app/module/logger"
	"time"
)

// 设置hook当请求结束输出请求的相关信息
func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		start_time := time.Now()
		c.Next()
		end_time := time.Now()
		delay := end_time.Sub(start_time)
		logger.Log.WithFields(logrus.Fields{"name":"zheng"}).Info(fmt.Sprintf("请求页面：%s, 请求IP地址为：%s, 请求方式为：%s, 执行时长为：%s\n",c.Request.URL, c.ClientIP(), c.Request.Method, delay))
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
