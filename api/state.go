package api

import (
	"fmt"
	"github.com/9chain/nbcsrvd01/state"
	"sync"
	"time"
)

type apiKeyInfo struct {
	apiKey string
	valid  bool
	active time.Time
}

var (
	apiKeyMap = make(map[string]apiKeyInfo)
	lock      sync.Mutex
)

func initState() {
	go clearTimeoutApiKey()
}

func checkApiKeyInternal(username, apikey string) bool {
	lock.Lock()
	defer lock.Unlock()

	info, ok := apiKeyMap[username]
	if ok {
		if !(info.valid && info.apiKey == apikey) {
			//fmt.Println(333, info.valid, info.apiKey, apikey)
			return false
		}
		return true
	}

	info = apiKeyInfo{apiKey: "", valid: false, active: time.Now()}
	apiKeyMap[username] = info

	var user state.User
	if err := state.DB.Take(&user, "username=?", username).Error; err != nil {
		//fmt.Println(111)
		return false
	}

	if user.State == 0 || apikey != user.ApiKey {
		//fmt.Println(222)
		return false
	}

	info.valid, info.apiKey = true, user.ApiKey
	apiKeyMap[username] = info
	return true
}

func clearTimeoutApiKey() {
	tick := time.Tick(5 * time.Second)

	clear := func() {
		now := time.Now()
		var deleteKeys []string

		lock.Lock()
		defer lock.Unlock()
		for k, info := range apiKeyMap {
			if now.After(info.active.Add(time.Second * 30)) {
				deleteKeys = append(deleteKeys, k)
			}
		}

		for _, k := range deleteKeys {
			fmt.Println("delete", k)
			delete(apiKeyMap, k)
		}
	}
	for {
		select {
		case <-tick:
			if len(apiKeyMap) > 0 {
				clear()
			}
			break
		}
	}
}
