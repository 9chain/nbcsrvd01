package web

import (
	"encoding/json"
	"fmt"
	"github.com/9chain/nbcsrvd01/config"
	"github.com/9chain/nbcsrvd01/primitives"
	"github.com/9chain/nbcsrvd01/state"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"strings"
	"time"
)

func init() {
	// 注册标准json2rpc处理函数
	handlers["register"] = handleV1Register
	handlers["login"] = handleV1Login
	handlers["reset-password"] = handleV1ResetPassword
	handlers["apikey"] = handleV1ApiKey
	handlers["reset-apikey"] = handleV1ResetApiKey
	handlers["forget-password"] = handleV1ForgetPassword
}

func MapToObject(source interface{}, dst interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

//　登陆
func handleV1Login(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	type loginParam struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	p := new(loginParam)
	if err := MapToObject(params, p); err != nil {
		return nil, primitives.NewInvalidParamsError()
	}

	if err := state.State.CheckLogin(p.Username, p.Password); err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	//　不管以前有没有登陆，重写
	session, _ := store.Get(ctx.Request, "session")
	session.Options = &sessions.Options{MaxAge: 60 * config.Cfg.Session.MaxAgeMin}
	session.Values["username"] = p.Username

	if err := session.Save(ctx.Request, ctx.Writer); err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	return gin.H{"message": "success"}, nil
}

// 注册，不需要 session
func handleV1Register(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	type regParam struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	p := new(regParam)
	if err := MapToObject(params, p); err != nil {
		return nil, primitives.NewInvalidParamsError()
	}

	// 简单检查 TODO
	if !(len(p.Username) > 4 && len(p.Username) <= 32 &&
		len(p.Password) > 4 && len(p.Password) <= 32 &&
		len(p.Email) > 4 && len(p.Email) < 64 && strings.Contains(p.Email, "@")) {
		return nil, primitives.NewInvalidParamsError()
	}

	smtpCfg := config.Cfg.SMTP

	sendRegisterEmail := func(username string) *JSONError {
		action := ConfirmMessage{
			Action:    "register",
			Username:  username,
			Timestamp: time.Now(),
		}

		// 添加签名
		b64, err := encodeContent(action, []byte(smtpCfg.Salt))
		if err != nil {
			return primitives.NewCustomInternalError(err.Error())
		}

		link := fmt.Sprintf("%s?content=%s", smtpCfg.ConfirmUrl, b64)
		content := fmt.Sprintf("点击连接激活: %s", link)
		email := NewEmail(p.Email, smtpCfg.ActiveTitle, content)
		if err := SendEmail(email); err != nil {
			return primitives.NewCustomInternalError("send email fail") // TODO
		}
		return nil
	}

	// 是否已经注册
	stat, _ := state.State.GetUserState(p.Username)
	switch stat {
	case -1: // not found
		// 写数据库
		if err := state.State.AddNewUser(p.Username, p.Password, p.Email); err != nil {
			return nil, primitives.NewCustomInternalError(err.Error())
		}

		// 发邮件
		if err := sendRegisterEmail(p.Username); err != nil {
			return nil, err
		}

		break
	case 0: // not confirm
		_, t, err := state.State.EmailInfo(p.Username)
		if err != nil {
			return nil, primitives.NewCustomInternalError(err.Error())
		}

		now := time.Now()
		// 两次发送时间间隔
		if now.Before(t.Add(time.Duration(smtpCfg.TimeoutMin) * time.Minute)) {
			return nil, primitives.NewCustomInternalError("wait a minute")
		}

		state.State.UpdateUserInfo(p.Username, p.Password, p.Email)
		if err := sendRegisterEmail(p.Username); err != nil {
			return nil, err
		}
		break

	default: // others
		return nil, primitives.NewCustomInternalError("already register")
	}

	return gin.H{"message": "success"}, nil
}

// 忘记密码，不需要 session
func handleV1ForgetPassword(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	type Param struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	p := new(Param)
	if err := MapToObject(params, p); err != nil {
		return nil, primitives.NewInvalidParamsError()
	}

	if !(len(p.Username) >= 4 && len(p.Username) <= 32 &&
		len(p.Email) >= 4 && len(p.Email) <= 64 && strings.Contains(p.Email, "@")) {
		return nil, primitives.NewCustomInternalError("invalid username/email")
	}

	emailAddr, t, err := state.State.EmailInfo(p.Username)
	if err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	if emailAddr != p.Email {
		return nil, primitives.NewCustomInternalError("invalid username/email")
	}

	now := time.Now()
	smtpCfg := config.Cfg.SMTP

	// 两次发送时间间隔
	if now.Before(t.Add(time.Duration(smtpCfg.TimeoutMin) * time.Minute)) {
		return nil, primitives.NewCustomInternalError("wait a minute")
	}

	if err := state.State.UpdateEmailTime(p.Username); err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	action := ConfirmMessage{
		Action:    "forget",
		Username:  p.Username,
		Timestamp: time.Now(),
	}

	// 添加签名
	b64, err := encodeContent(action, []byte(smtpCfg.Salt))
	if err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	link := fmt.Sprintf("%s?content=%s", smtpCfg.ConfirmUrl, b64)
	content := fmt.Sprintf("点击连接重设密码: %s", link)
	email := NewEmail(p.Email, smtpCfg.ResetPasswordTitle, content)
	if err := SendEmail(email); err != nil {
		return nil, primitives.NewCustomInternalError("send email fail")
	}

	return gin.H{"message": "success"}, nil
}

// 登陆后重设密码
func handleV1ResetPassword(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	username, err := getUsername(ctx)
	if err != nil {
		return nil, err
	}

	type Param struct {
		Password string `json:"password"`
	}

	p := new(Param)
	if err := MapToObject(params, p); err != nil {
		return nil, primitives.NewInvalidParamsError()
	}

	if err := state.State.ResetPassword(username, p.Password); err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	return gin.H{"message": "success"}, nil
}

// 从session获取用户名
func getUsername(ctx *gin.Context) (string, *primitives.JSONError) {
	session, err := store.Get(ctx.Request, "session")
	if err != nil || session.IsNew {
		return "", primitives.NewCustomInternalError("not authorized")
	}

	username, ok := session.Values["username"]
	if !ok {
		return "", primitives.NewCustomInternalError("not authorized")
	}

	return username.(string), nil
}

// 获取apiKey TODO
func handleV1ApiKey(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	username, err := getUsername(ctx)
	if err != nil {
		return nil, err
	}

	apikey, ok := state.State.GetUserApiKey(username)
	if !ok {
		return nil, primitives.NewCustomInternalError("not found apikey")
	}

	return gin.H{"apikey": apikey}, nil
}

// 重设apiKey
func handleV1ResetApiKey(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	username, err := getUsername(ctx)
	if err != nil {
		return nil, err
	}

	if apiKey, err := state.State.ResetUserApiKey(username); err == nil {
		return gin.H{"message": "success", "apikey": apiKey}, nil
	}

	return nil, primitives.NewCustomInternalError(err.Error())
}
