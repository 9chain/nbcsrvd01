package web

import (
	"errors"
	"fmt"
	"github.com/9chain/nbcsrvd01/config"
	"github.com/9chain/nbcsrvd01/primitives"
	"github.com/9chain/nbcsrvd01/state"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"os"
	"time"
)

type JSONError = primitives.JSONError
type JSON2Request = primitives.JSON2Request
type JSON2Response = primitives.JSON2Response

var (
	store    *sessions.FilesystemStore
	handlers = make(map[string]func(ctx *gin.Context, params interface{}) (interface{}, *JSONError))
)

// 初始化session
func initSession() {
	sessionCfg := config.Cfg.Session
	dir, key := sessionCfg.SessionDir, sessionCfg.SessionKey
	if _, err := os.Stat(dir); err != nil {
		fmt.Println("mk sesison dir", dir)
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic("mkdir for session fail " + err.Error())
		}
	}
	store = sessions.NewFilesystemStore(dir, []byte(key))
}

func InitWeb(r *gin.RouterGroup) {
	// 初始化session
	initSession()

	// 标准json2rpc处理
	r.POST("v1", func(ctx *gin.Context) {
		j, err := parseJSON2Request(ctx)
		if err != nil {
			handleV1Error(ctx, nil, primitives.NewInvalidRequestError())
			return
		}

		jsonResp, jsonError := handleV1Request(ctx, j)
		if jsonError != nil {
			handleV1Error(ctx, j, jsonError)
			return
		}

		ctx.JSON(200, jsonResp)
	})

	// 点击　注册、忘记密码邮件确认连接
	r.GET("v1/confirm", handleV1Confirm)

	// 忘记密码：重置
	r.POST("v1/resetpassword", handleForgetResetPassword)
}

func handleV1Error(ctx *gin.Context, j *JSON2Request, err *JSONError) {
	resp := primitives.NewJSON2Response()
	if j != nil {
		resp.ID = j.ID
	} else {
		resp.ID = nil
	}
	resp.Error = err

	ctx.JSON(200, resp)
}

func handleForgetResetPassword(ctx *gin.Context) {
	j, err := parseJSON2Request(ctx)
	if err != nil {
		handleV1Error(ctx, nil, primitives.NewInvalidRequestError())
		return
	}

	message, err := decodeMessage(ctx)
	if err != nil {
		handleV1Error(ctx, j, primitives.NewCustomInternalError(err.Error()))
		return
	}

	if message.Action != "forget" {
		handleV1Error(ctx, j, primitives.NewCustomInternalError("invalid action"))
		return
	}

	type regParam struct {
		Password string `json:"password"`
	}

	p := new(regParam)
	if err := MapToObject(j.Params, p); err != nil {
		handleV1Error(ctx, j, primitives.NewCustomInternalError(err.Error()))
		return
	}

	if !(len(p.Password) >= 4 && len(p.Password) <= 32) {
		handleV1Error(ctx, j, primitives.NewCustomInternalError("invalid password"))
		return
	}

	var user User
	db := state.DB
	if err := db.Take(&user, "username=?", message.Username).Error; err != nil {
		handleV1Error(ctx, j, primitives.NewCustomInternalError(err.Error()))
		return
	}

	if err := db.Model(&user).Updates(User{Password:p.Password, UpdatedAt:time.Now()}).Error; err != nil {
		handleV1Error(ctx, j, primitives.NewCustomInternalError(err.Error()))
		return
	}

	state.BackupUserConfig()

	jsonResp := primitives.NewJSON2Response()
	jsonResp.ID = j.ID
	jsonResp.Result = gin.H{"message": "success"}
	ctx.JSON(200, jsonResp)
}

// 解析、验证content是否合法
func decodeMessage(ctx *gin.Context) (*ConfirmMessage, error) {
	content := ctx.Query("content")

	smtpCfg := config.Cfg.SMTP
	var message ConfirmMessage

	// 签名验证
	if err := decodeContent(content, []byte(smtpCfg.Salt), &message); err != nil {
		return nil, err
	}

	//　超时验证
	now := time.Now()
	if now.After(message.Timestamp.Add(time.Duration(smtpCfg.TimeoutMin) * time.Minute)) {
		return nil, errors.New("timeout")
	}

	return &message, nil
}

// 注册、忘记密码　邮件连接点击　(GET)
func handleV1Confirm(ctx *gin.Context) {
	message, err := decodeMessage(ctx)
	if err != nil {
		ctx.Redirect(302, "/") // 重定向到登陆页面 TODO
		return
	}

	redirectIndex := func() {
		ctx.Redirect(302, "/") // 重定向到登陆页面 TODO
	}

	switch message.Action {
	case "register":
		// 注册邮件确认
		var user User
		db := state.DB
		if err := db.Take(&user, "username=?", message.Username).Error; err != nil {
			redirectIndex()
			return
		}

		if user.State > 0 {
			redirectIndex()
			return
		}

		// 修改状态为已经确认(1)
		if err := db.Model(&user).Updates(User{State:1, UpdatedAt:time.Now()}).Error; err != nil {
			redirectIndex()
			return
		}

		state.BackupUserConfig()

		redirectIndex()
		return
	case "forget":
		// 忘记密码邮件确认： 重定向到修改密码页面
		url := fmt.Sprintf("%s?%s", config.Cfg.SMTP.PageForgetPassord, ctx.Request.URL.RawQuery)
		ctx.Redirect(302, url)
		return
	default:
		redirectIndex()
		return
	}
}

func handleV1Request(ctx *gin.Context, j *JSON2Request) (*JSON2Response, *JSONError) {
	var resp interface{}
	var jsonError *JSONError

	//　查找、调用　处理函数
	if f, ok := handlers[j.Method]; ok {
		resp, jsonError = f(ctx, j.Params)
	} else {
		jsonError = primitives.NewMethodNotFoundError()
	}

	if jsonError != nil {
		return nil, jsonError
	}

	jsonResp := primitives.NewJSON2Response()
	jsonResp.ID = j.ID
	jsonResp.Result = resp

	return jsonResp, nil
}
