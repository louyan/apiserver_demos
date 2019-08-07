package user

import (
	"strconv"

	"yannotes.cn/apiserver_demos/demo10/pkg/errno"

	. "yannotes.cn/apiserver_demos/demo10/handler"

	"yannotes.cn/apiserver_demos/demo10/model"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	if err := model.DeleteUser(uint64(uid)); err != nil {
		SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	SendResponse(c, nil, nil)
}
