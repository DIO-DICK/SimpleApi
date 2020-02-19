package logger

import (
	"fmt"
	"os"
	"path"
	"simpleapi/app/conf"
	"time"
)

var Log_dir string
var Log_file string


func StatDir() (b bool, e error) {
	dir, path_err := os.Getwd()
	if path_err != nil {
		return false, path_err
	}
	log_dir := path.Join(dir, "/log")
	_, err := os.Stat(log_dir)
	if err == nil {
		Log_dir = log_dir
		return true, nil
	} else {
		err := os.Mkdir(log_dir, os.ModePerm)
		if err != nil {
			return false, err
		}
		Log_dir = log_dir
		return true, nil
	}
}

func StatFile() (b bool, e error) {
	t := time.Now()
	t1 := t.Format("20060102")
	ok, err := StatDir()
	if ok {
		Log_file = path.Join(Log_dir, (conf.Get().Log_file+t1))
		_, err := os.Lstat(Log_file)
		if err != nil {
			f, err := os.Create(Log_file)
			if err != nil {
				return false, err
			} else {
				f.Close()
				return true, nil
			}
		}
		return true, nil
	}
	return false, err
}

func init() {
	ok, _ := StatFile()
	if ok {
		fmt.Println("日志文件初始化成功")
	} else {
		fmt.Println("日志文件初始化失败")
	}

}