package log

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
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

type LogFormatter struct{}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

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
		f.AppendKeyValue(b, "level", entry.Data["level"].(string))
		f.AppendKeyValue(b, "msg", entry.Data["msg"].(string))

		for k, v := range entry.Data {
			if k != "level" && k != "msg" && k != "time" {
				f.AppendKeyValue(b, k, v)
			}
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *LogFormatter) AppendKeyValue(b *bytes.Buffer, key, value interface{}) {
	if _, ok := value.(string); ok {
		fmt.Fprintf(b, "%v=%q ", key, value)
	} else {
		fmt.Fprintf(b, "%v=%v ", key, value)
	}
}
