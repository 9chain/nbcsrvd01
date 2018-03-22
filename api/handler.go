package api

import (
	"github.com/9chain/nbcsrvd01/primitives"
	"github.com/gin-gonic/gin"
	"fmt"
)

func init() {
	handlers["create-chain"] = handleV1CreateChain
	handlers["create-entry"] = handleV1CreateEntry
	handlers["entry"] = handleV1Entry
	handlers["validate"] = handleV1Validate
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
	fmt.Println(1113, params)
	return gin.H{"message": "success"}, nil
}

// check and send to sdksrvd
func handleV1Validate(params interface{}) (interface{}, *primitives.JSONError) {
	fmt.Println(111, params)
	return gin.H{"message": "success"}, nil
}
