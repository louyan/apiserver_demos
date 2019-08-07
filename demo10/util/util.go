package util

import (
	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
)

/*`
go test -coverprofile=cover.out`：在测试文件目录下运行测试并统计测试覆盖率

`go tool cover -func=cover.out`：分析覆盖率文件，可以看出哪些函数没有测试，哪些函数内部的分支没有测试完全，cover 工具会通过执行代码的行数与总行数的比例表示出覆盖率
*/

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestId, ok := v.(string); ok {
		return requestId
	}
	return ""
}
