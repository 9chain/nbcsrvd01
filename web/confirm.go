package web

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

type ConfirmMessage struct {
	Action    string
	Username  string
	Timestamp time.Time
}

func encodeContent(p interface{}, salt []byte) (string, error) {
	src, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	bs := src[:]
	bs = append(bs, salt...)
	ctx := md5.New()
	ctx.Write([]byte(bs))
	sum := hex.EncodeToString(ctx.Sum(nil))
	content := append(src, sum...)
	return base64.URLEncoding.EncodeToString(content), nil
}

func decodeContent(src string, salt []byte, res interface{}) error {
	bs, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		return err
	}

	if len(bs) <= 32 {
		return errors.New("invalid length")
	}

	content := []byte(string(bs[:len(bs)-32]))
	all := append(content, salt...)
	ctx := md5.New()
	ctx.Write([]byte(all))
	sum := hex.EncodeToString(ctx.Sum(nil))

	targetSum := string(bs[len(bs)-32:])
	if sum == targetSum {
		if err := json.Unmarshal(content, &res); err != nil {
			return err
		}
		return nil
	}
	return errors.New("invalid sum")
}
