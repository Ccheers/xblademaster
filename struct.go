package xblademaster

import (
	"fmt"
	"github.com/go-kratos/kratos/pkg/log"
	"net/http"
	"strconv"
	"time"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

var message = map[int]string{
	0:         "成功",
	100000001: "尚未设置错误码",
	400000000: "错误哦",
	400000001: "token失效,请检查后再试",

	500000000: "系统内部错误,请检查后再试",
}

type empty struct {
}

func JsonResponse(ctx *Context, data interface{}, err error) {
	if data == nil {
		data = empty{}
	}
	beginTime, _ := strconv.ParseInt(ctx.Writer.Header().Get("X-Begin-Time"), 10, 64)
	duration := time.Now().Sub(time.Unix(0, beginTime))
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	roundedStr := fmt.Sprintf("%.3fms", rounded)
	ctx.Writer.Header().Set("X-Response-time", roundedStr)
	ctx.JSON(data, err)
	log.Info("message: %s\t%s\t%s\t%s\t%s", "JSON Response", "errMsg", err.Error(), "took", roundedStr)
	ctx.Abort()
	return
}

func RenderResponse(ctx *Context, templateName string, errCode int, data H) {
	if data == nil {
		panic("常规信息应该设置")
	}
	msg := GetMsg(errCode)
	beginTime, _ := strconv.ParseInt(ctx.Writer.Header().Get("X-Begin-Time"), 10, 64)

	duration := time.Now().Sub(time.Unix(0, beginTime))
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	roundedStr := fmt.Sprintf("%.3fms", rounded)
	ctx.Writer.Header().Set("X-Response-time", roundedStr)
	requestUrl := ctx.Request.URL.String()
	requestMethod := ctx.Request.Method
	log.Info("message %s\t%s\t%s\t%s\t%s", "Index Response", "Request Url", requestUrl, "Request method", requestMethod, "code", errCode, "errMsg", msg, "took", roundedStr)
	if errCode == 500 {
		ctx.HTML(http.StatusInternalServerError, "5xx.tmpl", data)
	} else if errCode == 404 {
		ctx.HTML(http.StatusNotFound, "4xx.tmpl", data)
	} else if errCode == 0 {
		ctx.HTML(http.StatusOK, templateName, data)
	} else {
		ctx.HTML(http.StatusServiceUnavailable, "5xx.tmpl", nil)
	}
	ctx.Abort()
	return
}

func GetMsg(code int) string {
	if msg, ok := message[code]; ok {
		return msg
	}
	return message[100000001]
}
