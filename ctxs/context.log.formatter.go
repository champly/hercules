package ctxs

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	timeVar                = "time"
	logLevelVar            = "level"
	sessionIDVar           = "sid"
	contentVar             = "msg"
	warpStr                = "%"
	defaultLogFormat       = fmt.Sprintf("[%s] [%s] [%s] %s", warp(timeVar), warp(logLevelVar), warp(sessionIDVar), warp(contentVar))
	defaultTimestampFormat = "2006-01-02 15:04:05.000"
)

type Formatter struct {
	TimestampFormat string
	LogFormat       string
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {

	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, warp(timeVar), entry.Time.Format(timestampFormat), 1)
	output = strings.Replace(output, warp(contentVar), entry.Message, 1)
	output = strings.Replace(output, warp(logLevelVar), strings.ToUpper(entry.Level.String()[:4]), 1)

	for k, v := range entry.Data {
		if s, ok := v.(string); ok {
			output = strings.Replace(output, warp(k), s, 1)
		}
	}

	// output = fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorList[entry.Level.String()], output)
	output += "\n"

	return []byte(output), nil
}

func warp(v string) string {
	return warpStr + v + warpStr
}
