package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JsonErrStruct struct {
	Code int         `JSON:code`
	Msg  interface{} `JSON:msg`
}

type JsonStruct struct {
	JsonErrStruct
	Data  interface{} `JSON:data`
	Count int64       `JSON:count`
}

func ReturnSuccess(
	c *gin.Context,
	code int,
	msg interface{},
	data interface{},

	count int64) {
	json := JsonStruct{
		JsonErrStruct: JsonErrStruct{
			Code: code,
			Msg:  msg,
		},
		Data:  data,
		Count: count,
	}
	c.JSON(http.StatusOK, json)
}
func ReturnError(
	c *gin.Context,
	code int,
	msg interface{},
) {
	json := JsonErrStruct{
		Code: code,
		Msg:  msg,
	}

	c.JSON(http.StatusOK, json)
}
func EncryMd5(s string) string {
	ctx := md5.New()
	ctx.Write([]byte(s))
	return hex.EncodeToString(ctx.Sum(nil))
}
