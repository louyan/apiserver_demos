package user

import (
	"yannotes.cn/apiserver_demos/demo07/model"

	"yannotes.cn/apiserver_demos/demo07/util"

	"github.com/lexkong/log/lager"

	. "yannotes.cn/apiserver_demos/demo07/handler"

	"github.com/lexkong/log"

	"yannotes.cn/apiserver_demos/demo07/pkg/errno"

	"github.com/gin-gonic/gin"
)

// Create creates a new user account
func Create(c *gin.Context) {
	log.Info("User Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest

	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	u := model.UserModel{
		Username: r.Username,
		Password: r.Password,
	}

	// Validate the data.
	if err := u.Validate(); err != nil {
		SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// Encrypt the user password
	if err := u.Encrypt(); err != nil {
		SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	//Insert the user to the database
	if err := u.Create(); err != nil {
		SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	rsp := CreateResponse{
		Username: r.Username,
	}
	//code, message := errno.DecodeErr(err)
	//c.JSON(http.StatusOK, gin.H{"code": code, "message": message})
	SendResponse(c, nil, rsp)
}

func (r *CreateRequest) checkParam() error {
	if r.Username == "" {
		return errno.New(errno.ErrValidation, nil).Add("username is empty.")
	}
	if r.Password == "" {
		return errno.New(errno.ErrValidation, nil).Add("password is empty.")
	}
	return nil
}
