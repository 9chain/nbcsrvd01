package web

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"github.com/9chain/nbcsrvd01/primitives"
	"github.com/gorilla/sessions"
	"fmt"
)

//var store = sessions.NewCookieStore([]byte("something-very-secret"))
var handlers = make(map[string]func(params interface{}) (interface{}, *primitives.JSONError))
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

	if f, ok := handlers[j.Method]; ok {
		resp, jsonError = f(params)
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

