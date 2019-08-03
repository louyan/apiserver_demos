package user

import (
	"github.com/gin-gonic/gin"
	. "yannotes.cn/apiserver_demos/demo08/handler"
	"yannotes.cn/apiserver_demos/demo08/model"
	"yannotes.cn/apiserver_demos/demo08/pkg/errno"
)

func Get(c *gin.Context) {
	username := c.Param("username")
	user, err := model.GetUser(username)
	if err != nil {
		SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	SendResponse(c, nil, user)
}
