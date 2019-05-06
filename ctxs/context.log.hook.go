package ctxs

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"

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

func (f *FileSplitHook) Fire(e *logrus.Entry) error {
	e.Logger.Out = f.getStdOut(e.Time.Format("2006-01-02"))
	e.Data["sid"] = f.sid
	return nil
}

func (f *FileSplitHook) getStdOut(t string) io.Writer {
	if len(t) < 10 {
		return os.Stdout
	}
	fp, ok := f.out[t]
	if ok {
		return fp
	}
	f.lock.Lock()
	defer f.lock.Unlock()

	if fp, ok = f.out[t]; ok {
		return fp
	}

	fp, _ = os.OpenFile(fmt.Sprintf("%s.log", t), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	iw := io.MultiWriter(os.Stdout, fp)
	f.out[t] = iw
	return iw
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
