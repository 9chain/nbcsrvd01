package api

import (
	"errors"
	"fmt"
	"github.com/9chain/nbcsrvd01/primitives"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type JSONError = primitives.JSONError
type JSON2Request = primitives.JSON2Request
type JSON2Response = primitives.JSON2Response

var (
	handlers    = make(map[string]func(ctx *gin.Context, params interface{}) (interface{}, *JSONError))
	notifyChans = userNotify{channels: make(map[string]chan interface{}), rwLock: sync.RWMutex{}}
)

type userNotify struct {
	channels map[string]chan interface{}
	rwLock   sync.RWMutex
}

func (n *userNotify) Notify(username string, msg interface{}) {
	n.rwLock.RLock()
	defer n.rwLock.RUnlock()

	ch, ok := n.channels[username]
	if !ok {
		fmt.Println("ignore to notify", username, msg)
		return
	}

	ch <- msg
}

func (n *userNotify) getChannel(username string) chan interface{} {
	lock, channels := n.rwLock, n.channels
	lock.Lock()
	defer lock.Unlock()

	lastChan, ok := channels[username]
	if ok {
		fmt.Println("last chan", lastChan)
		close(lastChan)
	}

	ch := make(chan interface{}, 10)
	channels[username] = ch
	return ch
}

func (n *userNotify) close(username string) {
	lock, channels := n.rwLock, n.channels
	lock.Lock()
	defer lock.Unlock()

	ch, ok := channels[username]
	if !ok {
		return
	}

	delete(channels, username)
	close(ch)
}

func InitApi(r *gin.RouterGroup) {
	initState()
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

	r.GET("v1/notify", handleNotify)

	go func() {
		for {
			time.Sleep(time.Second * 5)
			notifyChans.Notify("kitty", time.Now())
		}
	}()
}

func checkApiKey(ctx *gin.Context) bool {
	username, apiKey := ctx.GetHeader("X-Username"), ctx.GetHeader("X-Api-Key")
	if len(username) == 0 || len(apiKey) == 0 {
		return false
	}

	if !checkApiKeyInternal(username, apiKey) {
		return false
	}

	return true
}

func handleV1Request(ctx *gin.Context, j *JSON2Request) (*JSON2Response, *JSONError) {
	var resp interface{}
	var jsonError *JSONError

	if !checkApiKey(ctx) {
		return nil, primitives.NewCustomInternalError("invalid username/apikey")
	}

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

func handleNotify(ctx *gin.Context) {
	if !checkApiKey(ctx) {
		ctx.AbortWithError(400, errors.New("invalid username/apikey"))
		return
	}

	var upgrader = websocket.Upgrader{} // use default options

	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	username := ctx.GetHeader("X-Username")

	ch := notifyChans.getChannel(username)

	for {
		msg, ok := <-ch
		if msg != nil {
			fmt.Printf("recv: %+v\n", msg)
			err = c.WriteJSON(msg)
			if err != nil {
				fmt.Println("close chan . write error", err, ch)
				notifyChans.close(username)
				break
			}
		}

		if !ok {
			fmt.Println("already close ", username, ok, ch)
			break
		}
	}
}
