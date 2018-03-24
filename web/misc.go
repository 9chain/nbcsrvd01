package web

import (
	"github.com/9chain/nbcsrvd01/primitives"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

// 解析json2rpc参数
func parseJSON2Request(ctx *gin.Context) (*JSON2Request, error) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}

	j, err := primitives.ParseJSON2Request(body)
	if err != nil {
		return nil, err
	}

	return j, nil
}
