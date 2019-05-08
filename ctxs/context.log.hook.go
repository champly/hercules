package ctxs

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/champly/hercules/configs"
	"github.com/sirupsen/logrus"
)

type FileSplitHook struct {
	out  map[string]io.Writer
	lock sync.Mutex
	sid  string
}

func NewFileSplitHook() *FileSplitHook {
	return &FileSplitHook{
		out: map[string]io.Writer{},
		sid: GetGuid(),
	}
}

func (f *FileSplitHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *FileSplitHook) Fire(e *logrus.Entry) (err error) {
	e.Logger.Out, err = f.getStdOut(e.Time.Format("2006-01-02"))
	if err != nil {
		return err
	}
	e.Data["sid"] = f.sid
	return nil
}

func (f *FileSplitHook) getStdOut(t string) (io.Writer, error) {
	fp, ok := f.out[t]
	if ok {
		return fp, nil
	}
	f.lock.Lock()
	defer f.lock.Unlock()

	if fp, ok = f.out[t]; ok {
		return fp, nil
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

	iw := io.MultiWriter(mw...)
	f.out[t] = iw
	return iw, nil
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))

}

//生成Guid字串
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))

}
