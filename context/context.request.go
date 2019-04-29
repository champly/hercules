package context

import "github.com/gin-gonic/gin"

type Request struct {
	Http *gin.Context
}

func NewRequest(ctx *gin.Context) *Request {
	return &Request{Http: ctx}
}
