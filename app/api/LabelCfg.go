package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simpleapi/app/models"
)

func GetLabel(c *gin.Context) {
	var la models.T02DabsUpLabelCfg
	result, err := la.GetAllLabelCfg()

	if err != nil {
		fmt.Println("获取数据失败")
		c.JSON(
			http.StatusOK,
			gin.H{
				"status" : -1,
				"data" : "读取数据失败",
			})
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status" : 1,
			"data" : result,
		})
}

func GetLabelByField(c *gin.Context) {
	id := c.Param("field")
	var la models.T02DabsUpLabelCfg
	result, err := la.GetLabelCfgByField(id)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"status" : -1,
				"data" : "获取数据失败",
			})
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status" : 1,
			"data" : result,
		})
}