/*参考新浪开放平台 [Error code](https://link.juejin.im?target=http%3A%2F%2Fopen.weibo.com%2Fwiki%2FError_code) 的设计
* 服务级别错误：1 为系统级错误；2 为普通错误，通常是由用户非法操作引起的
*
20502
		2					05			02
服务级错误（1为系统级错误）	服务模块代码	具体错误代码
*/
package errno

var (
	//Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error."}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	ErrValidation = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase   = &Errno{Code: 20002, Message: "Database error."}
	ErrToken      = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token."}

	// user errors
	ErrEncrypt           = &Errno{Code: 20101, Message: "Error occurred while encrypting the user password."}
	ErrUserNotFound      = &Errno{Code: 20102, Message: "The user was not found."}
	ErrTokenInvalid      = &Errno{Code: 20103, Message: "The token was invalid."}
	ErrPasswordIncorrect = &Errno{Code: 20104, Message: "The password was incorrect."}
)
