package log

import (
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	checkedEnv = false
	showDebug  = false
)

const (
	y = "YES"
	t = "TRUE"
)

// Debug logs debug level messages
func Debug(v interface{}) {
	if showDebug {
		log.Println(v)
	}
}

// Debugf logs using a Printf method
func Debugf(format string, v interface{}) {
	if showDebug {
		log.Printf(format, v)
	}
}

// Info is the standard level for logs you want all the time
func Info(v interface{}) {
	log.Println(v)
}

// Infof is the standard logger, but using the Printf method
func Infof(format string, v interface{}) {
	log.Printf(format, v)
}

func init() {
	if checkedEnv {
		return
	}

	i, err := strconv.ParseInt(os.Getenv("DEBUG"), 10, 0)
	if err != nil {
		str := strings.ToUpper(os.Getenv("DEBUG"))
		if str == y || str == t {
			i = 1
		}
	}

	if i == 1 {
		showDebug = true
	}
	checkedEnv = true
}
