package filestate

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

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
