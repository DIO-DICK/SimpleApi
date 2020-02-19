package databases

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"simpleapi/app/conf"
	_ "github.com/go-sql-driver/mysql"
)

var DbMysql *gorm.DB

func init() {
	Dburl:=fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.Get().Mysql.DB_user,
		conf.Get().Mysql.DB_password,
		conf.Get().Mysql.DB_ip,
		conf.Get().Mysql.DB_port,
		conf.Get().Mysql.DB,
		conf.Get().Mysql.DB_charset)

	var err error
	DbMysql,err = gorm.Open("mysql", Dburl)

	if err != nil {
		fmt.Println(err)
	}

	if DbMysql.Error != nil {
		fmt.Printf("database error %v", DbMysql.Error)
	}
}
