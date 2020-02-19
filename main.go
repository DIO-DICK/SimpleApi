package main

import (
	"fmt"
	orm "simpleapi/app/databases"
	"simpleapi/app/module/logger"
	"simpleapi/app/routers"
)

func main() {
	defer orm.DbMysql.Close()
	fmt.Println(logger.Log_file)
	router := routers.InitRouter()
	router.Run("127.0.0.1:9090")
}