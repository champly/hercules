package ctxs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/champly/hercules/configs"
	"github.com/champly/lib4go/tool"
	"github.com/sirupsen/logrus"
)

var (
	lock sync.Mutex
	out  map[string]io.Writer
)

func init() {
	out = map[string]io.Writer{}
}

type FileSplitHook struct {
	sid string
}

func NewFileSplitHook() *FileSplitHook {
	return &FileSplitHook{
		sid: tool.GetGUID()[:8],
	}
}

func (f *FileSplitHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *FileSplitHook) Fire(e *logrus.Entry) (err error) {
	e.Logger.Out, err = getStdOut(e.Time.Format("2006-01-02"))
	if err != nil {
		return err
	}
	e.Data["sid"] = f.sid
	return nil
}

func getStdOut(t string) (io.Writer, error) {
	fp, ok := out[t]
	if ok {
		return fp, nil
	}
	lock.Lock()
	defer lock.Unlock()

	if fp, ok = out[t]; ok {
		return fp, nil
	}

	if strings.Trim(configs.LoggerInfo.Out, " ") == "" {
		out[t] = os.Stdout
		return os.Stdout, nil
	}

	mw := []io.Writer{}
	for _, ty := range strings.Split(configs.LoggerInfo.Out, "|") {
		switch strings.ToLower(strings.Trim(ty, " ")) {
		case "stdout":
			mw = append(mw, os.Stdout)
		case "file":
			if len(t) < 10 {
				break
			}
			fp, err := os.OpenFile(fmt.Sprintf("%s.log", t), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return nil, errors.New("logger out(file) create log file err:" + err.Error())
			}
			mw = append(mw, fp)
		default:
			panic("not support logger write function:" + ty)
		}
	}
	if len(mw) == 0 {
		panic("logger is not out")
	}

	iw := io.MultiWriter(mw...)
	out[t] = iw
	return iw, nil
}
