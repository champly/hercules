package ctxs

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type herContext struct {
	*gin.Context
}

func NewHerContext(c *gin.Context) *herContext {
	return &herContext{c}
}

func (h *herContext) GetBody() string {
	b, err := ioutil.ReadAll(h.Request.Body)
	if err != nil {
		return ""
	}
	return string(b)
}
