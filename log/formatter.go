package log

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/zimmski/logrus"
)

/*

	This is pretty much a copy of https://github.com/Sirupsen/logrus/blob/master/text_formatter.go

*/

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
)

// TextFormatter implements a text formatter for logrus
// Only the log level and the log message are printed out for every log entry.
type TextFormatter struct{}

// Format takes a logrus log entry and transforms it to its text format.
// The error return argument is not nil if the internal buffer is unwritable.
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	prefixFieldClashes(entry)

	if logrus.IsTerminal() {
		levelText := strings.ToUpper(entry.Data["level"].(string))[0:4]

		levelColor := blue

		if entry.Data["level"] == "info" {
			levelColor = green
		} else if entry.Data["level"] == "warning" {
			levelColor = yellow
		} else if entry.Data["level"] == "error" ||
			entry.Data["level"] == "fatal" ||
			entry.Data["level"] == "panic" {
			levelColor = red
		}

		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m %-44s ", levelColor, levelText, entry.Data["msg"])

		var keys []string
		for k := range entry.Data {
			if k != "level" && k != "msg" && k != "time" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := entry.Data[k]
			fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=%v", levelColor, k, v)
		}
	} else {
		appendKeyValue(b, "level", entry.Data["level"].(string))
		appendKeyValue(b, "msg", entry.Data["msg"].(string))

		for k, v := range entry.Data {
			if k != "level" && k != "msg" && k != "time" {
				appendKeyValue(b, k, v)
			}
		}
	}

	if err := b.WriteByte('\n'); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func appendKeyValue(b *bytes.Buffer, key, value interface{}) {
	if _, ok := value.(string); ok {
		fmt.Fprintf(b, "%v=%q ", key, value)
	} else {
		fmt.Fprintf(b, "%v=%v ", key, value)
	}
}

func prefixFieldClashes(entry *logrus.Entry) {
	_, ok := entry.Data["time"]
	if ok {
		entry.Data["fields.time"] = entry.Data["time"]
	}

	entry.Data["time"] = entry.Time.Format(time.RFC3339)

	_, ok = entry.Data["msg"]
	if ok {
		entry.Data["fields.msg"] = entry.Data["msg"]
	}

	entry.Data["msg"] = entry.Message

	_, ok = entry.Data["level"]
	if ok {
		entry.Data["fields.level"] = entry.Data["level"]
	}

	entry.Data["level"] = entry.Level.String()
}
