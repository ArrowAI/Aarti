package clifmt

import (
	"bytes"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type TimeFormat string

const (
	NoneTimeFormat     TimeFormat = "none"
	FullTimeFormat     TimeFormat = "full"
	RelativeTimeFormat TimeFormat = "relative"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	white  = 39
	gray   = 90
)

func New(f TimeFormat) logrus.Formatter {
	return &clifmt{start: time.Now(), format: f}
}

type clifmt struct {
	start  time.Time
	format TimeFormat
}

func (f *clifmt) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer
	var c *color.Color
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		c = color.New(gray)
	case logrus.WarnLevel:
		c = color.New(yellow)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		c = color.New(red)
	default:
		c = color.New(white)
	}
	msg := entry.Message
	if len(entry.Message) > 0 && entry.Level < logrus.DebugLevel {
		msg = strings.ToTitle(string(msg[0])) + msg[1:]
	}

	var err error
	switch f.format {
	case FullTimeFormat:
		_, err = c.Fprintf(&b, "[%s] %s\n", entry.Time.Format("2006-01-02 15:04:05"), entry.Message)
	case RelativeTimeFormat:
		_, err = c.Fprintf(&b, "[%5v] %s\n", entry.Time.Sub(f.start).Truncate(time.Second).String(), msg)
	case NoneTimeFormat:
		fallthrough
	default:
		_, err = c.Fprintln(&b, msg)
	}
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
