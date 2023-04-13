package pkg

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogFormatter struct{}

var levelList = []string{
	"PANIC",
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

func (lf *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	level := levelList[int(entry.Level)]
	strList := strings.Split(entry.Caller.File, "/")
	fileName := strList[len(strList)-1]

	traceId := "(N/A)"
	if entry.Context != nil && entry.Context.Value("traceid") != nil {
		traceId = fmt.Sprintf("(%s)", entry.Context.Value("traceid").(string))
	}

	b.WriteString(fmt.Sprintf("%s [%s:%d]-[%v] - %s %s\n",
		entry.Time.Format("2006-01-02 15:04:05.000"), fileName,
		entry.Caller.Line, level, traceId, entry.Message))
	return b.Bytes(), nil
}

func NewLogger() error {
	var path string

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&LogFormatter{})

	// Please kindly change the path to the preferred path that exists on your device
	switch runtime.GOOS {
	case "windows":
		path = os.TempDir()
	default:
		path = os.Getenv("LOG_LOCATION")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			logrus.Error("Could not make dir: ", err.Error())
			return err
		}
	}

	file, err := os.OpenFile(path+"/log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	logrus.SetOutput(file)
	logrus.SetOutput(os.Stdout)

	return nil
}
