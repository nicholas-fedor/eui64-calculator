package mocks

import "github.com/gin-gonic/gin"

// RequestContext mocks the RequestContext interface for testing.
type RequestContext struct {
	ginContext *gin.Context
}

// NewRequestContext creates a new RequestContext with the given gin.Context.
func NewRequestContext(ginContext *gin.Context) *RequestContext {
	return &RequestContext{ginContext: ginContext}
}

func (m *RequestContext) FormValue(key string) string {
	return m.ginContext.PostForm(key)
}

func (m *RequestContext) GetContext() *gin.Context {
	return m.ginContext
}
