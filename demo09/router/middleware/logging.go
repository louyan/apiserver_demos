package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/willf/pad"

	"yannotes.cn/apiserver_demos/demo09/pkg/errno"

	"github.com/lexkong/log"

	"yannotes.cn/apiserver_demos/demo09/handler"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logging is a middleware function that logs the each request.
func Logging() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path

		reg := regexp.MustCompile("(/v1/user|/login)")
		if !reg.MatchString(path) {
			return
		}

		// Skip for the health check requests
		if path == "/sd/health" || path == "/sd/ram" || path == "/sd/cpu" || path == "/sd/disk" {
			return
		}
		//读取Body内容
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		//HTTP 的请求 Body，在读取过后会被置空
		//读取完后会重新赋值
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		method := c.Request.Method
		ip := c.ClientIP()

		//log.Debugf("New request come in, path: %s, Method: %s, body `%s`", path, method, string(bodyBytes))
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		//计算时间
		end := time.Now().UTC()
		latency := end.Sub(start)

		code, message := -1, ""
		// get code and message
		var response handler.Response
		if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
			log.Errorf(err, "response body can not unmarshal to model.Response struct, body: `%s`", blw.body.Bytes())
			code = errno.InternalServerError.Code
			message = err.Error()
		} else {
			code = response.Code
			message = response.Message
		}

		log.Infof("%-13s | %-12s | %s %s | {code: %d, message: %s}", latency, ip, pad.Right(method, 5, ""), path, code, message)
	}
}
