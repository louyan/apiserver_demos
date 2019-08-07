package user

import (
	"github.com/gin-gonic/gin"
	. "yannotes.cn/apiserver_demos/demo10/handler"
	"yannotes.cn/apiserver_demos/demo10/model"
	"yannotes.cn/apiserver_demos/demo10/pkg/auth"
	"yannotes.cn/apiserver_demos/demo10/pkg/errno"
	"yannotes.cn/apiserver_demos/demo10/pkg/token"
)

func Login(c *gin.Context) {
	var u model.UserModel
	if err := c.Bind(&u); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	// Get the user information by the login username.
	d, err := model.GetUser(u.Username)
	if err != nil {
		SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	// Compare the login password with the user password
	//校验密码是否匹配
	if err := auth.Compare(d.Password, u.Password); err != nil {
		SendResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	//Sign the json web token ,签发 token
	t, err := token.Sign(c, token.Context{ID: d.Id, Username: d.Username}, "")
	if err != nil {
		SendResponse(c, errno.ErrToken, nil)
		return
	}

	SendResponse(c, nil, model.Token{Token: t})
}
