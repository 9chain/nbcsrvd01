package web

import (
	"fmt"
	"github.com/9chain/nbcsrvd01/config"
	"github.com/9chain/nbcsrvd01/primitives"
	"github.com/9chain/nbcsrvd01/state"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"strings"
	"time"
)

type User = state.User

func init() {
	// 注册标准json2rpc处理函数
	handlers["register"] = handleV1Register
	handlers["login"] = handleV1Login
	handlers["reset-password"] = handleV1ResetPassword
	handlers["reset-apikey"] = handleV1ResetApiKey
	handlers["forget-password"] = handleV1ForgetPassword
	handlers["userinfo"] = handleV1UserInfo
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

	var user User
	if err := state.DB.Take(&user, "username=?", p.Username).Error; err != nil {
		return nil, primitives.NewCustomInternalError("invalid username/password")
	}

	if user.State == 0 {
		return nil, primitives.NewCustomInternalError("not confirmed yet")
	}

	if p.Password != user.Password {
		return nil, primitives.NewCustomInternalError("invalid username/password")
	}

	//　不管以前有没有登陆，重写
	session, _ := store.Get(ctx.Request, "session")
	session.Options = &sessions.Options{MaxAge: 60 * config.Cfg.Session.MaxAgeMin}
	session.Values["username"] = p.Username

	if err := session.Save(ctx.Request, ctx.Writer); err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	user.Password = "******"
	return gin.H{"user": user}, nil
}

// 获取用户信息
func handleV1UserInfo(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	username, err := getUsername(ctx)
	if err != nil {
		return nil, err
	}

	var user User
	db := state.DB
	if err := db.Take(&user, "username=?", username).Error; err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	user.Password = "******"
	return gin.H{"user": user}, nil
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
			fmt.Println("send email fail", err)
			return primitives.NewCustomInternalError("send email fail") // TODO
		}
		return nil
	}

	// 是否已经注册
	var user User
	db, stat := state.DB, -1
	if err := db.Take(&user, "username=?", p.Username).Error; err != nil {
		stat = -1
	} else {
		stat = user.State
	}

	switch stat {
	case -1: // not found
		// 写数据库
		newUser := User{
			Username:  p.Username,
			Password:  p.Password,
			Email:     p.Email,
			ApiKey:    genRandomApiKey(),
			EmailedAt: time.Now(),
		}

		if err := db.Create(&newUser).Error; err != nil {
			return nil, primitives.NewCustomInternalError(err.Error())
		}

		// 发邮件
		if err := sendRegisterEmail(p.Username); err != nil {
			return nil, err
		}

		break
	case 0: // not confirm
		now := time.Now()

		// 两次发送时间间隔
		if now.Before(user.EmailedAt.Add(2 * time.Minute)) {
			return nil, primitives.NewCustomInternalError("wait a minute")
		}

		if err := db.Model(&user).Updates(User{EmailedAt: time.Now()}).Error; err != nil {
			return nil, primitives.NewCustomInternalError(err.Error())
		}

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

	var user User
	db := state.DB
	if err := db.Take(&user, "username=?", p.Username).Error; err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	if user.Email != p.Email {
		return nil, primitives.NewCustomInternalError("invalid username/email")
	}

	now := time.Now()
	smtpCfg := config.Cfg.SMTP

	// 两次发送时间间隔
	if now.Before(user.EmailedAt.Add(2 * time.Minute)) {
		return nil, primitives.NewCustomInternalError("wait a minute")
	}

	if err := db.Model(&user).Updates(User{EmailedAt: time.Now()}).Error; err != nil {
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

	password := p.Password
	if !(len(password) >= 4 && len(password) <= 64) {
		return nil, primitives.NewInvalidParamsError()
	}

	var user User
	db := state.DB
	if err := db.Take(&user, "username=?", username).Error; err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	if password == user.Password {
		return gin.H{"message": "success"}, nil
	}

	if err := db.Model(&user).Updates(User{Password: password, UpdatedAt: time.Now()}).Error; err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}
	state.BackupUserConfig()

	user.Password = "******"
	return gin.H{"user": user}, nil
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

// 重设apiKey
func handleV1ResetApiKey(ctx *gin.Context, params interface{}) (interface{}, *JSONError) {
	username, err := getUsername(ctx)
	if err != nil {
		return nil, err
	}

	var user User
	db := state.DB
	if err := db.Take(&user, "username=?", username).Error; err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	if err := db.Model(&user).Updates(User{ApiKey: genRandomApiKey(), UpdatedAt: time.Now()}).Error; err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}
	state.BackupUserConfig()

	user.Password = "******"
	return gin.H{"user": user}, nil
}
