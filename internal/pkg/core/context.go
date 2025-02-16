package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/core/validator"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// MarshalOptions is a configurable JSON format marshaller.
var MarshalOptions = protojson.MarshalOptions{
	UseProtoNames:   true,
	EmitUnpopulated: true,
}

type Context struct {
	Context *gin.Context
}

func New(ctx *gin.Context) *Context {
	return &Context{ctx}
}

// Unauthorized 未认证
func (c *Context) Unauthorized(message string) error {
	c.Context.AbortWithStatusJSON(http.StatusUnauthorized, &Response{
		Code:    http.StatusUnauthorized,
		Message: message,
	})

	return nil
}

// Forbidden 未授权
func (c *Context) Forbidden(message string) error {
	c.Context.AbortWithStatusJSON(http.StatusForbidden, &Response{
		Code:    http.StatusForbidden,
		Message: message,
	})

	return nil
}

// InvalidParams 参数错误
func (c *Context) InvalidParams(message any) error {
	// 初始化响应，设置状态码为 400
	resp := &Response{Code: 400, Message: "invalid params"}

	// 根据 message 的类型进行处理
	switch msg := message.(type) {
	case error:
		// 如果 message 是 error 类型，使用 validator.Translate 翻译错误信息
		resp.Message = validator.Translate(msg)
	case string:
		// 如果 message 是字符串类型，直接赋值给 resp.Message
		resp.Message = msg
	default:
		// 如果 message 是其他类型，使用 fmt.Sprintf 转换为字符串
		resp.Message = fmt.Sprintf("%v", msg)
	}

	// 中止请求，返回 400 状态码和响应体
	c.Context.AbortWithStatusJSON(http.StatusBadRequest, resp)

	return nil
}

// Error 错误信息响应
func (c *Context) Error(err error) error {
	resp := &Response{Code: 400, Message: err.Error()}

	var e *errorx.Error
	if errors.As(err, &e) {
		resp.Code = e.Code
		resp.Message = e.Message
		c.Context.AbortWithStatusJSON(http.StatusBadRequest, resp)
	} else {
		resp.Code = 500
		resp.Message = err.Error()
		c.Context.AbortWithStatusJSON(http.StatusInternalServerError, resp)
	}

	return nil
}

// Success 成功响应(Json 数据)
func (c *Context) Success(data any, message ...string) error {
	resp := &Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}

	if len(message) > 0 {
		resp.Message = message[0]
	}

	// 检测是否是 proto 对象
	if value, ok := data.(proto.Message); ok {
		bt, _ := MarshalOptions.Marshal(value)

		var body map[string]any
		_ = json.Unmarshal(bt, &body)
		resp.Data = body
	}

	c.Context.AbortWithStatusJSON(http.StatusOK, resp)
	return nil
}

// Raw 成功响应(原始数据)
func (c *Context) Raw(value string) error {
	c.Context.Abort()
	c.Context.String(http.StatusOK, value)
	return nil
}

// UserId 返回登录用户的UID
func (c *Context) UserId() int {
	if session := c.JwtSession(); session != nil {
		return session.Uid
	}

	return 0
}

// JwtSession 返回登录用户的JSession
func (c *Context) JwtSession() *middleware.JSession {
	data, isOk := c.Context.Get(middleware.JWTSessionConst)
	if !isOk {
		return nil
	}

	return data.(*middleware.JSession)
}

// IsGuest 是否是游客(未登录状态)
func (c *Context) IsGuest() bool {
	return c.UserId() == 0
}

func (c *Context) Ctx() context.Context {
	return c.Context.Request.Context()
}
