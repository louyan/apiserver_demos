package user

import (
	"github.com/gin-gonic/gin"
	. "yannotes.cn/apiserver_demos/demo10/handler"
	"yannotes.cn/apiserver_demos/demo10/pkg/errno"
	"yannotes.cn/apiserver_demos/demo10/service"
)

func List(c *gin.Context) {
	var r ListRequest
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	infos, count, err := service.ListUser(r.Username, r.Offset, r.Limit)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}

	SendResponse(c, nil, ListResponse{
		TotalCount: count,
		UserList:   infos,
	})
}
