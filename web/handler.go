package web

import (
	"encoding/json"
	"github.com/9chain/nbcsrvd01/primitives"
	"github.com/9chain/nbcsrvd01/state"
	"github.com/gin-gonic/gin"
)

func init() {
	handlers["register"] = handleV1Register
	handlers["login"] = handleV1Login
	handlers["reset-password"] = handleV1ResetPassword
	handlers["apikey"] = handleV1ApiKey
	handlers["reset-apikey"] = handleV1ResetApiKey
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
