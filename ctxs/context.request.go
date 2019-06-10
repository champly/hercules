package ctxs

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Request struct {
	*gin.Context
}

func NewRequest(ctx *gin.Context) *Request {
	return &Request{ctx}
}

func (r *Request) GetBody() string {
	b, err := ioutil.ReadAll(r.Request.Body)
	if err != nil {
		return ""
	}
	return string(b)
}
