package msgo

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	greenBg   = "\033[97;42m"
	whiteBg   = "\033[90;47m"
	yellowBg  = "\033[90;43m"
	redBg     = "\033[97;41m"
	blueBg    = "\033[97;44m"
	magentaBg = "\033[97;45m"
	cyanBg    = "\033[97;46m"
	green     = "\033[32m"
	white     = "\033[37m"
	yellow    = "\033[33m"
	red       = "\033[31m"
	blue      = "\033[34m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	reset     = "\033[0m"
)

var DefaultWriter io.Writer = os.Stdout

type LoggerConfig struct {
	Formatter LoggerFormatter
	out       io.Writer
	IsColor   bool
}

type LoggerFormatter func(params LogFormatterParams) string

type LogFormatterParams struct {
	Request    *http.Request
	TimeStamp  time.Time
	StatusCode int
	Latency    time.Duration
	ClientIP   net.IP
	Method     string
	Path       string
}

func Logging(next HandlerFunc) HandlerFunc {
	return LoggerWithConfig(LoggerConfig{}, next)
}

func LoggerWithConfig(conf LoggerConfig, next HandlerFunc) HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}
	out := conf.out
	if out == nil {
		out = DefaultWriter
		conf.IsColor = true
	}
	return func(ctx *Context) {
		param := LogFormatterParams{
			Request: ctx.Request,
		}
		// Start timer
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		//执行业务
		next(ctx)
		// stop timer
		stop := time.Now()
		latency := stop.Sub(start)
		ip, _, _ := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr))
		clientIP := net.ParseIP(ip)
		method := ctx.Request.Method
		statusCode := ctx.StatusCode

		if raw != "" {
			path = path + "?" + raw
		}

		param.ClientIP = clientIP
		param.TimeStamp = stop
		param.Latency = latency
		param.StatusCode = statusCode
		param.Method = method
		param.Path = path
		fmt.Fprint(out, formatter(param))
	}
}

var defaultLogFormatter = func(params LogFormatterParams) string {
	statusCodeColor := params.StatusCodeColor()
	resetColor := params.ResetColor()
	if params.Latency > time.Minute {
		params.Latency = params.Latency.Truncate(time.Second)
	}
	// 这里是通过 ANSI转义序列，实现彩色展示的
	// ANSI转义序列以ESC字符（\033或\x1b）开头，后跟一个[字符，然后是参数列表（由分号分隔的数字组成），最后以m字符结束。例如，\033[31m是一个将文本颜色设置为红色的ANSI转义序列。
	// resetColor 变量所代表的值，则代表重置样式
	// 组合过程 样式 + 字符串 + 重置样式 + 字符串
	return fmt.Sprintf("%s[msgo]%s |%s %v %s| %s %3d %s |%s %13v %s| %15s  |%s %-7s %s %s %#v %s",
		yellow, resetColor,
		blue, params.TimeStamp.Format("2006/01/02 - 15:04:05"), resetColor,
		statusCodeColor, params.StatusCode, resetColor,
		red, params.Latency, resetColor,
		params.ClientIP,
		magenta, params.Method, resetColor,
		cyan, params.Path, resetColor,
	)
}

func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode
	switch code {
	case http.StatusOK:
		return green
	default:
		return red
	}
}
func (p *LogFormatterParams) ResetColor() string {
	return reset
}
