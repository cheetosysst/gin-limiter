package limiter

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

// global ip rate
type GlobalRate struct {
	Command  string
	Period   time.Duration
	Limit    int
	deadLine int64
}

func (gr *GlobalRate) SetDeadLine(deadline int64) {
	if deadline > 0 {
		gr.deadLine = deadline
	}
}

func (gr GlobalRate) GetDeadLine() int64 {
	return gr.deadLine
}

// local ip rate
type singleRate struct {
	Path    string
	Method  string
	Command string
	Period  time.Duration
	Limit   int
}

type Rates []singleRate

func (rs Rates) Append(sr singleRate) {
	rs = append(rs, sr)
	fmt.Println(rs)
}

func (rs Rates) getLimit(path, method string) int {
	for _, rate := range rs {
		if (rate.Path == path) && (rate.Method == method) {
			return rate.Limit
		}
	}
	return -1
}

var methodDict = map[string]bool{
	"GET":     true,
	"PUT":     true,
	"POST":    true,
	"HEAD":    true,
	"TRACE":   true,
	"PATCH":   true,
	"DELETE":  true,
	"CONNECT": true,
	"OPTIONS": true,
}

var timeDict = map[string]time.Duration{
	"S": time.Second,
	"M": time.Minute,
	"H": time.Hour,
	"D": time.Hour * 24,
}

var MethodError = errors.New("Please check the method is one of http method.")
var CommandError = errors.New("The command of first number should > 0.")
var FormatError = errors.New("Please check the format with your input.")
var LimitError = errors.New("Limit should > 0.")

// NewGlobalRate("10-m", 200)
// Each 10 minutes single ip address can request 200 times.
func newGlobalRate(command string, limit int) (GlobalRate, error) {
	var gRate GlobalRate
	var period time.Duration

	values := strings.Split(command, "-")
	if len(values) != 2 {
		log.Println("Some error with your input command!, the len of your command is ", len(values))
		return gRate, FormatError
	}

	unit, err := strconv.Atoi(values[0])
	if err != nil {
		return gRate, FormatError
	}
	if unit <= 0 {
		return gRate, CommandError
	}

	// limit should > 0
	if limit <= 0 {
		return gRate, LimitError
	}

	if t, ok := timeDict[strings.ToUpper(values[1])]; ok {
		period = time.Duration(unit) * t
	} else {
		return gRate, FormatError
	}

	gRate.Command = command
	gRate.Period = period
	gRate.Limit = limit
	return gRate, nil
}

func newSingleRate(path, command, method string, limit int) (singleRate, error) {
	var sRate singleRate
	var period time.Duration

	values := strings.Split(command, "-")
	if len(values) != 2 {
		log.Println("Some error with your input command!, the len of your command is ", len(values))
		return sRate, FormatError
	}

	unit, err := strconv.Atoi(values[0])
	if err != nil {
		return sRate, FormatError
	}
	if unit <= 0 {
		return sRate, CommandError
	}

	// limit should > 0
	if limit <= 0 {
		return sRate, LimitError
	}

	if t, ok := timeDict[strings.ToUpper(values[1])]; ok {
		period = time.Duration(unit) * t
	} else {
		return sRate, FormatError
	}

	if _, ok := methodDict[strings.ToUpper(method)]; !ok {
		return sRate, MethodError
	}

	sRate.Path = path
	sRate.Method = method
	sRate.Command = command
	sRate.Period = period
	sRate.Limit = limit
	return sRate, nil
}
