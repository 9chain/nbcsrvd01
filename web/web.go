package web

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"github.com/9chain/nbcsrvd01/primitives"
	"github.com/9chain/nbcsrvd01/state"
	"github.com/gorilla/sessions"
	"encoding/json"
	"net/smtp"
	"strings"
	"fmt"
)

//var store = sessions.NewCookieStore([]byte("something-very-secret"))

var store *sessions.FilesystemStore
func InitWeb(r *gin.RouterGroup) {
	store = sessions.NewFilesystemStore("/tmp/nbcsrvd_session", []byte("secret"))

	r.GET("v1", handleV1)
	r.POST("v1", handleV1)
	r.GET("v1/confirm", handleV1Confirm)
}

func handleV1Error(ctx *gin.Context, j *primitives.JSON2Request, err *primitives.JSONError) {
	resp := primitives.NewJSON2Response()
	if j != nil {
		resp.ID = j.ID
	} else {
		resp.ID = nil
	}
	resp.Error = err

	ctx.JSON(200, resp)
}

func handleV1(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
		return
	}

	j, err := primitives.ParseJSON2Request(body)
	if err != nil {
		handleV1Error(ctx, nil, primitives.NewInvalidRequestError())
		return
	}

	jsonResp, jsonError := handleV1Request(j)
	if jsonError != nil {
		handleV1Error(ctx, j, jsonError)
		return
	}

	ctx.JSON(200, jsonResp)
}

func handleV1Confirm(ctx *gin.Context) {
	//body, err := ioutil.ReadAll(ctx.Request.Body)
	//if err != nil {
	//	panic(err)
	//	return
	//}

	//_ = body
	fmt.Printf("=== %+v\n", ctx.Params)
	action := ctx.Query("action")		// register, forget, reset
	content := ctx.Query("content")
	fmt.Println(handleV1Confirm, action, content)
	_ = content
	// TODO validate

	// TODO if is forget-password, redirect to reset page

	ctx.JSON(200, gin.H{"message":"success"})
}

func handleV1Request(j *primitives.JSON2Request) (*primitives.JSON2Response, *primitives.JSONError) {
	var resp interface{}
	var jsonError *primitives.JSONError
	params := j.Params
	switch j.Method {
	case "register":
		resp, jsonError = handleV1Register(params)
		break
	case "login":
		resp, jsonError = handleV1Login(params)
		break
	case "forget-password":
		resp, jsonError = handleV1ForgetPassword(params)
		break
	case "reset-password":
		resp, jsonError = handleV1ResetPassword(params)
		break
	case "apikey":
		resp, jsonError = handleV1ApiKey(params)
		break
	case "reset-apikey":
		resp, jsonError = handleV1ResetApiKey(params)
		break
	default:
		break
	}
	if jsonError != nil {
		return nil, jsonError
	}

	jsonResp := primitives.NewJSON2Response()
	jsonResp.ID = j.ID
	jsonResp.Result = resp

	return jsonResp, nil
}

func MapToObject(source interface{}, dst interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}


func handleV1Login(params interface{}) (interface{}, *primitives.JSONError) {
	type loginParam struct {
		Username string 	`json:"username"`
		Password string 	`json:"password"`
	}

	p := new(loginParam)
	if err := MapToObject(params, p); err != nil {
		return nil, primitives.NewInvalidParamsError()
	}

	if err := state.State.CheckLogin(p.Username, p.Password); err != nil {
		return nil, primitives.NewCustomInternalError(err.Error())
	}

	return gin.H{"message": "success"}, nil
}


func handleV1Register(params interface{}) (interface{}, *primitives.JSONError) {
	//type loginParam struct {
	//	Username string		`json:"username"`
	//	Password string 	`json:"password"`
	//	Email string 		`json:"email"`
	//}
	//
	//p := new(loginParam)
	//if err := MapToObject(params, p); err != nil {
	//	return nil, primitives.NewInvalidParamsError()
	//}
	//
	//if err := state.State.CheckLogin(p.Username, p.Password); err != nil {
	//	return nil, primitives.NewCustomInternalError(err.Error())
	//}
	SendEmail(NewEmail("329365307@qq.com", "test email", "test content: http://www.ninechain.net/panel"))
	return gin.H{"message": "success"}, nil
}

func handleV1ForgetPassword(params interface{}) (interface{}, *primitives.JSONError) {
	//type loginParam struct {
	//	Username string		`json:"username"`
	//	Password string 	`json:"password"`
	//	Email string 		`json:"email"`
	//}
	//
	//p := new(loginParam)
	//if err := MapToObject(params, p); err != nil {
	//	return nil, primitives.NewInvalidParamsError()
	//}
	//
	//if err := state.State.CheckLogin(p.Username, p.Password); err != nil {
	//	return nil, primitives.NewCustomInternalError(err.Error())
	//}
	SendEmail(NewEmail("329365307@qq.com", "test email", "test content: http://www.ninechain.net/panel"))
	return gin.H{"message": "success"}, nil
}

func handleV1ResetPassword(params interface{}) (interface{}, *primitives.JSONError) {
	return gin.H{"message": "success"}, nil
}

func handleV1ApiKey(params interface{}) (interface{}, *primitives.JSONError) {
	return gin.H{"message": "success"}, nil
}

func handleV1ResetApiKey(params interface{}) (interface{}, *primitives.JSONError) {
	return gin.H{"message": "success"}, nil
}

const (
	HOST        = "smtp.163.com"
	SERVER_ADDR = "smtp.163.com:25"
	USER        = "aquariusye@163.com" //发送邮件的邮箱
	PASSWORD    = "helloshiki"         //发送邮件邮箱的密码
)

type Email struct {
	to      string "to"
	subject string "subject"
	msg     string "msg"
}

func NewEmail(to, subject, msg string) *Email {
	return &Email{to: to, subject: subject, msg: msg}
}

func SendEmail(email *Email) error {
	auth := smtp.PlainAuth("", USER, PASSWORD, HOST)
	sendTo := strings.Split(email.to, ";")

	go func() {
		for _, v := range sendTo {
			str := strings.Replace("From: "+USER+"~To: "+v+"~Subject: "+email.subject+"~~", "~", "\r\n", -1) + email.msg
			fmt.Println(str)
			err := smtp.SendMail(
				SERVER_ADDR,
				auth,
				USER,
				[]string{v},
				[]byte(str),
			)
			fmt.Println(11111, err)
		}
	}()

	return nil
}