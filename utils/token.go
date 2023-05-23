package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var tokenmanager = &tokenManager{}

type tokenManager struct {
	sync.RWMutex

	tokens map[string]string
	ready  bool
}

func (this *tokenManager) Constructor() {
	this.tokens = make(map[string]string)
	this.ready = true
}
func (this *tokenManager) Dispose() {
	this.tokens = nil
	this.ready = false
}

func CreateToken(id string) string {
	var err error
	if !tokenmanager.ready {
		tokenmanager.Constructor()
	}

	var rs = rand.New(rand.NewSource(time.Now().UnixNano()))
	var hash = md5.New()
	_, err = io.WriteString(hash, strconv.FormatInt(rs.Int63(), 10))
	if err == nil {
		var token = fmt.Sprintf("%x", hash.Sum(nil))
		tokenmanager.Lock()
		defer tokenmanager.Unlock()
		tokenmanager.tokens[id] = token
		return token
	}
	return ""
}
func UpdateToken(id string, token string) (string, error) {
	var err error
	if CheckToken(id, token) {
		var token = CreateToken(id)
		return token, nil
	} else {
		err = errors.New("Token 更新失败.")
	}
	return "", err
}
func DeleteToken(id string, token string) error {
	var err error
	if CheckToken(id, token) {
		if _, has := tokenmanager.tokens[id]; has {
			tokenmanager.Lock()
			defer tokenmanager.Unlock()
			delete(tokenmanager.tokens, id)
		}
	} else {
		err = errors.New("Token 删除失败.")
	}
	return err
}
func CheckToken(id string, token string) bool {
	if !tokenmanager.ready {
		tokenmanager.Constructor()
	}

	if tk, ok := tokenmanager.tokens[id]; ok {
		if tk == token {
			return true
		}
	}
	return false
}
