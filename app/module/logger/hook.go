package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

// 函数行号钩子，用于错误时输出
type lineHook struct {
	Field string
	Skip int
	level []logrus.Level
}


func (hook *lineHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (hook *lineHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = getFunc(hook.Skip)
	return nil
}

func NewHook() logrus.Hook {
	hook := lineHook{
		"source",
		0,
		logrus.AllLevels,
	}
	return &hook
}

func getFunc(skip int) string {
	file := ""
	line := 0
	fnName := ""
	var pc uintptr

	for i := 0; i<9; i++ {
		file, line, pc = getCaller(skip + i)
	}
	fullnme := runtime.FuncForPC(pc)
	if fullnme != nil {
		fnNameStr := fullnme.Name()
		parts := strings.Split(fnNameStr,".")
		fnName = parts[len(parts) - 1]
		switch fnName {
		case "0":
			fnName = parts[len(parts) - 2]
		default:
			fnName = fnName
		}
	}

	return fmt.Sprintf("[%s:%d:%s()]", file, line, fnName)
}

func getCaller(skip int) (string, int, uintptr) {
	pc, package_name, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, pc
	}

	n:= 0

	for i :=  len(package_name) - 1; i>0; i-- {
		if package_name[i] == '/' {
			n++
			if n >= 2 {
				package_name = package_name[i + 1:]
				break
			}
		}
	}
	return package_name, line, pc
}
