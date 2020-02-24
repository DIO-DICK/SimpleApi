package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"simpleapi/app/conf"
	"time"
)

var Log_dir string
var Log_file string


func StatDir(dir1 string) (b bool, e error) {
	dir, path_err := os.Getwd()
	if path_err != nil {
		return false, path_err
	}
	log_dir := path.Join(dir, dir1)
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

func StatFile(dir1 string) (b bool, e error) {
	t := time.Now()
	t1 := t.Format("20060102")
	ok, err := StatDir(dir1)
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

var Log *logrus.Logger

func init() {
	var ls PathMap
	if b, _ := StatFile("/log"); !b {
		fmt.Println("logfile create failed")
	} else {
		ls = PathMap{
			logrus.InfoLevel: Log_file,
			logrus.ErrorLevel: Log_file,
		}
	}
	Log = logrus.New()
	file_log := NewFileHook(ls, &logrus.JSONFormatter{})
	Log.Hooks.Add(file_log)
}