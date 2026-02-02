package completion

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"github.com/jessevdk/go-flags"
)

var EnableLogging atomic.Bool

func PrintCompletionList(items []flags.Completion) {
	if IsZshShellFunc() {
		for _, v := range items {
			fmt.Println(v.Item)
		}
		os.Exit(0)
	}

	// normal bash shell
	for _, v := range items {
		if strings.HasSuffix(v.Item, NoSpace) {
			v.Item = strings.TrimSuffix(v.Item, NoSpace)
			fmt.Println(v.Item)
		} else {
			fmt.Println(v.Item + " ") // add space to the end
		}
	}
	os.Exit(0)
}

func WriteToLog(format string, args ...any) {
	if !EnableLogging.Load() {
		return
	}

	if os.Getenv("GO_TEST") != "" {
		return
	}

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(file, format, args...)
	if err != nil {
		panic(err)
	}

	if err := file.Close(); err != nil {
		panic(err)
	}
}

func PrintArray(values []string) string {
	var buf bytes.Buffer
	for index, v := range values {
		if index > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%q", v))
	}
	return buf.String()
}
