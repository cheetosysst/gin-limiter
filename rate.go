package limiter

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

// global ip rate
type GlobalRate struct {
	Command        string
	Period         time.Duration
	Limit          int
	deadLine       int64
	deadLineFormat string
}

func (gr *GlobalRate) UpdateDeadLine() {
	gr.deadLine = time.Now().Add(gr.Period).Unix()
	gr.deadLineFormat = time.Unix(gr.deadLine, 0).Format(TimeFormat)
}

func (gr *GlobalRate) GetDeadLine() int64 {
	return gr.deadLine
}

// local ip rate
type singleRate struct {
	Path           string
	Method         string
	Command        string
	Period         time.Duration
	Limit          int
	deadLine       int64
	deadLineFormat string
}

func (sr *singleRate) UpdateDeadLine() {
	sr.deadLine = time.Now().Add(sr.Period).Unix()
	sr.deadLineFormat = time.Unix(sr.deadLine, 0).Format(TimeFormat)
}

func (sr *singleRate) GetDeadLine() int64 {
	return sr.deadLine
}

type Rates []*singleRate

func (rs Rates) Append(sr *singleRate) {
	rs = append(rs, sr)
}

func (rs Rates) getLimit(path, method string) int {
	for _, rate := range rs {
		if (rate.Path == path) && (rate.Method == method) {
			return rate.Limit
		}
	}
	return -1
}

func (rs Rates) UpdateDeadLine() {
	for _, rate := range rs {
		rate.UpdateDeadLine()
	}
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
func newGlobalRate(command string, limit int) (*GlobalRate, error) {
	var gRate GlobalRate
	var period time.Duration

	values := strings.Split(command, "-")
	if len(values) != 2 {
		log.Println("Some error with your input command!, the len of your command is ", len(values))
		return nil, FormatError
	}

	unit, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, FormatError
	}
	if unit <= 0 {
		return nil, CommandError
	}

	// limit should > 0
	if limit <= 0 {
		return nil, LimitError
	}

	if t, ok := timeDict[strings.ToUpper(values[1])]; ok {
		period = time.Duration(unit) * t
	} else {
		return nil, FormatError
	}

	gRate.Command = command
	gRate.Period = period
	gRate.Limit = limit
	return &gRate, nil
}

func newSingleRate(path, command, method string, limit int) (*singleRate, error) {
	var sRate singleRate
	var period time.Duration

	values := strings.Split(command, "-")
	if len(values) != 2 {
		log.Println("Some error with your input command!, the len of your command is ", len(values))
		return nil, FormatError
	}

	unit, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, FormatError
	}
	if unit <= 0 {
		return nil, CommandError
	}

	// limit should > 0
	if limit <= 0 {
		return nil, LimitError
	}

	if t, ok := timeDict[strings.ToUpper(values[1])]; ok {
		period = time.Duration(unit) * t
	} else {
		return nil, FormatError
	}

	if _, ok := methodDict[strings.ToUpper(method)]; !ok {
		return nil, MethodError
	}

	sRate.Path = path
	sRate.Method = method
	sRate.Command = command
	sRate.Period = period
	sRate.Limit = limit
	return &sRate, nil
}
