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

	// user errors
	ErrUserNotFound = &Errno{Code: 20102, Message: "The user was not found."}
)
