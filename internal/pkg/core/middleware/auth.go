package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go-chat/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const JWTSessionConst = "__JWT_SESSION__"

var (
	ErrNoAuthorize = errors.New("授权异常，请登录后操作! ")
)

type IStorage interface {
	// IsBlackList 判断是否是黑名单
	IsBlackList(ctx context.Context, token string) bool
}

type JSession struct {
	Uid       int    `json:"uid"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// Auth 授权中间件
func Auth(secret string, guard string, storage IStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 Authorization 字段, 并去除 Bearer 前缀，获取token
		token := AuthHeaderToken(c)

		// 验证token
		claims, err := verify(guard, secret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
			return
		}

		if storage.IsBlackList(c.Request.Context(), token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "请登录再试"})
			return
		}

		uid, err := strconv.Atoi(claims.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "解析 jwt 失败"})
			return
		}

		// 设置jwt session到gin框架的context中, 用于后续的请求
		c.Set(JWTSessionConst, &JSession{
			Uid:       uid,
			Token:     token,
			ExpiresAt: claims.ExpiresAt.Unix(),
		})

		c.Next()
	}
}

func AuthHeaderToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))

	// Headers 中没有授权信息则读取 url 中的 token
	if token == "" {
		token = c.DefaultQuery("token", "")
	}

	return token
}

// verify 验证token
// guard 守卫名称, 比如 api, admin, open
// secret 密钥, 比如 asdfasdfasdfadsf
// token 令牌
func verify(guard string, secret string, token string) (*jwt.AuthClaims, error) {

	if token == "" {
		return nil, ErrNoAuthorize
	}

	claims, err := jwt.ParseToken(token, secret)
	if err != nil {
		return nil, err
	}

	// 判断权限认证守卫是否一致
	if claims.Guard != guard || claims.Valid() != nil {
		return nil, ErrNoAuthorize
	}

	return claims, nil
}
