package completion

import (
	"bytes"
	"fmt"
	"os"
)

func WriteToLog(format string, args ...any) {
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
