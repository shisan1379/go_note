package msgo

import (
	"errors"
	"fmt"
	"github.com/shisan1379/msgo/mserror"
	"net/http"
	"runtime"
	"strings"
)

func Recovery(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {

				//msError, ok := err.(mserror.MsError)
				//if ok {
				//	ctx.Logger.Error(detailMsg(msError))
				//	msError.ExecResult()
				//	return
				//}

				err2 := err.(error)
				if err2 != nil {
					var msError *mserror.MsError
					if errors.As(err2, &msError) {
						ctx.Logger.Error(detailMsg(err))

						msError.ExecResult()

						return
					}
				}
				ctx.Logger.Error(detailMsg(&err))
				ctx.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		next(ctx)
	}
}
func detailMsg(err any) string {
	var pcs [32]uintptr
	n := runtime.Callers(0, pcs[:])
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v\n", err))
	for _, pc := range pcs[0:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		sb.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return sb.String()
}

// printStackTrace 打印当前堆栈跟踪，包括文件名和行号
func printStackTrace() {
	var pcs [32]uintptr
	n := runtime.Callers(2, pcs[:]) // 跳过printStackTrace本身和它的调用者
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		fmt.Printf("%+v\n", frame) // 这将打印出文件名、函数名、行号等信息
		if !more {
			break
		}
	}
}
