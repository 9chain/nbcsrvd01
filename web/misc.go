package web

import (
	"github.com/9chain/nbcsrvd01/primitives"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
	"crypto/md5"
	"encoding/hex"
	"math/rand"

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


func MapToObject(source interface{}, dst interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

func genRandomApiKey() string {
	nonce := fmt.Sprintf("%x-%s", time.Now().UnixNano(), getRandomString(8))
	ctx := md5.New()
	ctx.Write([]byte(nonce))
	return hex.EncodeToString(ctx.Sum(nil))
}

func getRandomString(total int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < total; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
