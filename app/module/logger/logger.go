package logger

import (
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
	Log = logrus.New()
	ll := &lineHook{"source",0,logrus.AllLevels}
	Log.AddHook(ll)
	ok, _ := StatFile("/log")
	if !ok {
		Log.WithFields(logrus.Fields{"name":"zheng"}).Error("目录与文件创建失败")
	}
	Log.WithFields(logrus.Fields{"name":"zheng"}).Info("目录与文件创建成功")
}