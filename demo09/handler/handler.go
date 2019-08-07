/*
* API 返回入口函数，供所有的服务模块返回时调用
 */
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"yannotes.cn/apiserver_demos/demo09/pkg/errno"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}
