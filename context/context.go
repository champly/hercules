package context

type Context struct {
	Request  *Request
	Response *Response
}

func NewContext() *Context {
	return &Context{}
}
