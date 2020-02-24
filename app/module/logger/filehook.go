package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"reflect"
	"sync"
)

var defaultformatter = &logrus.TextFormatter{DisableColors:true}

type WriterMap map[logrus.Level]io.Writer
type PathMap map[logrus.Level]string

type FileHook struct {
	Paths PathMap
	Writers WriterMap
	Levells []logrus.Level
	Lock *sync.Mutex
	Formatter logrus.Formatter
	DefaultPath string
	DefaultWriter io.Writer
	HasDefaultPath bool
	HasDefaultWriter bool
}

func NewFileHook(output interface{}, formatt logrus.Formatter) *FileHook {
	hook := &FileHook{
		Lock:new(sync.Mutex),
	}

	hook.SetFormatter(formatt)

	switch output.(type) {
	case PathMap:
		hook.Paths = output.(PathMap)
		for level :=range output.(PathMap) {
			hook.Levells = append(hook.Levells, level)
		}
		break
	case WriterMap:
		hook.Writers = output.(WriterMap)
		for level := range output.(WriterMap) {
			hook.Levells = append(hook.Levells, level)
		}
		break
	case string:
		hook.SetDefaultPath(output.(string))
		break
	case io.Writer:
		hook.SetDefaultWriter(output.(io.Writer))
		break
	default:
		panic(fmt.Sprintf("Unsupport level map type：%v", reflect.TypeOf(output)))
	}
	return hook
}

func (hook *FileHook) SetDefaultPath(defaultpath string) {
	hook.Lock.Lock()
	defer hook.Lock.Unlock()
	hook.DefaultPath = defaultpath
	hook.HasDefaultPath = true
}

func (hook *FileHook) SetDefaultWriter(defaultwriter io.Writer) {
	hook.Lock.Lock()
	defer hook.Lock.Unlock()
	hook.DefaultWriter = defaultwriter
	hook.HasDefaultWriter = true
}

func (hook *FileHook) SetFormatter(formatt logrus.Formatter) {
	hook.Lock.Lock()
	defer hook.Lock.Unlock()
	if formatt == nil {
		formatt = defaultformatter
	} else {
		switch formatt.(type) {
		case *logrus.TextFormatter:
			textformatter := formatt.(*logrus.TextFormatter)
			textformatter.DisableColors = true
		}
	}
	hook.Formatter = formatt
}

func (hook *FileHook) ioWriter(entry *logrus.Entry) error {
	var (
		write io.Writer
		msg []byte
		ok bool
		err error
	)
	if write, ok = hook.Writers[entry.Level]; !ok {
		if hook.HasDefaultWriter {
			write = hook.DefaultWriter
		} else {
			return nil
		}
	}
	msg, err = hook.Formatter.Format(entry)
	if err != nil {
		fmt.Println("failed to cover format of entry")
		return err
	}
	_, err = write.Write(msg)
	return err
}

func (hook *FileHook) filewrite(entry *logrus.Entry) error {
	var (
		path string
		ok bool
		fd *os.File
		err error
		msg []byte
	)
	if path, ok = hook.Paths[entry.Level]; !ok {
		if hook.HasDefaultPath {
			path = hook.DefaultPath
		} else {
			return nil
		}
	}

	fd, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer fd.Close()

	if err != nil {
		fmt.Println("file open error")
		return err
	}
	msg, err = hook.Formatter.Format(entry)
	if err != nil {
		fmt.Println("failed to cover format of entry")
		return err
	}

	fd.Write(msg)
	return nil
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	hook.Lock.Lock()
	defer hook.Lock.Unlock()
	if hook.Writers != nil || hook.HasDefaultWriter {
		entry.Data["source"] = getFunc(0)
		return hook.ioWriter(entry)
	} else if hook.Paths != nil || hook.HasDefaultPath {
		entry.Data["source"] = getFunc(0)
		return hook.filewrite(entry)
	}
	return nil
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}