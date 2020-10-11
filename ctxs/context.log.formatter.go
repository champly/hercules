package ctxs

import (
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	defaultLogFormat       = "[%time%] [%level%] [%sid%] %msg%"
	defaultTimestampFormat = "2006-01-02 15:04:05.000"
)

// const (
// colorRed = uint8(iota + 91)
// colorGreen
// colorYellow
// colorBlue
// colorMagenta //洋红
// colorCyan
// )

// var colorList map[string]uint8

// func init() {
// colorList = make(map[string]uint8)

// colorList[logrus.TraceLevel.String()] = colorBlue
// colorList[logrus.DebugLevel.String()] = colorCyan
// colorList[logrus.InfoLevel.String()] = colorGreen
// colorList[logrus.WarnLevel.String()] = colorYellow
// colorList[logrus.ErrorLevel.String()] = colorRed
// colorList[logrus.FatalLevel.String()] = colorMagenta
// colorList[logrus.PanicLevel.String()] = colorMagenta
// }

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

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)
	output = strings.Replace(output, "%msg%", entry.Message, 1)
	output = strings.Replace(output, "%level%", strings.ToUpper(entry.Level.String()[:4]), 1)

	for k, v := range entry.Data {
		if s, ok := v.(string); ok {
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}

	// output = fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorList[entry.Level.String()], output)
	output += "\n"

	return []byte(output), nil
}
