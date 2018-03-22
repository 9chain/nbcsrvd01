package api

import (
	"github.com/gin-gonic/gin"
	"github.com/9chain/nbcsrvd01/primitives"
	"io/ioutil"
)

func InitApi(r *gin.RouterGroup) {
	r.GET("v1", handleV1)
	r.POST("v1", handleV1)
}

func HandleV2Error(ctx *gin.Context, j *primitives.JSON2Request, err *primitives.JSONError) {
	resp := primitives.NewJSON2Response()
	if j != nil {
		resp.ID = j.ID
	} else {
		resp.ID = nil
	}
	resp.Error = err

	ctx.JSON(400, resp)	// TODO 400 or 200?
}

func handleV1(ctx *gin.Context) {
	// check username & apikey
	//if err := checkAuthHeader(ctx.Request); err != nil {
	//	remoteIP := ""
	//	remoteIP += strings.Split(ctx.Request.RemoteAddr, ":")[0]
	//	fmt.Printf("Unauthorized API client connection attempt from %s %s\n", remoteIP, err)
	//	ctx.ResponseWriter.Header().Add("WWW-Authenticate", `Basic realm="factomd RPC"`)
	//	http.Error(ctx.ResponseWriter, "401 Unauthorized.", http.StatusUnauthorized)
	//	return
	//}
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
		return
	}

	j, err := primitives.ParseJSON2Request(body)
	if err != nil {
		HandleV2Error(ctx, nil, primitives.NewInvalidRequestError())
		return
	}

	jsonResp, jsonError := HandleV1Request(j)

	if jsonError != nil {
		HandleV2Error(ctx, j, jsonError)
		return
	}

	ctx.JSON(200, []byte(jsonResp.String()))
}

func HandleV1Request(j *primitives.JSON2Request) (*primitives.JSON2Response, *primitives.JSONError) {
	var resp interface{}
	var jsonError *primitives.JSONError
	params := j.Params
	switch j.Method {
	case "create-chain":
		resp, jsonError = handleV1CreateChain(params)
		break
	case "create-entry":
		resp, jsonError = handleV1CreateEntry(params)
		break
	case "entry":
		resp, jsonError = handleV1Entry(params)
		break
	case "validate":
		resp, jsonError = handleV1Validate(params)
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

// check and send to sdksrvd
func handleV1CreateChain(params interface{}) (interface{}, *primitives.JSONError) {
	return gin.H{"message": "success"}, nil
}

// check and send to sdksrvd
func handleV1CreateEntry(params interface{}) (interface{}, *primitives.JSONError) {
	return gin.H{"message": "success"}, nil
}

// check and send to sdksrvd
func handleV1Entry(params interface{}) (interface{}, *primitives.JSONError) {
	return gin.H{"message": "success"}, nil
}

// check and send to sdksrvd
func handleV1Validate(params interface{}) (interface{}, *primitives.JSONError) {
	return gin.H{"message": "success"}, nil
}
