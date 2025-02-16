package v1

import (
	"errors"
	"strconv"
	"time"

	"go-chat/api/pb/admin/v1"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/encrypt/rsautil"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"

	"github.com/mojocn/base64Captcha"
	"gorm.io/gorm"
)

type Auth struct {
	Config          *config.Config
	AdminRepo       *repo.Admin
	JwtTokenStorage *cache.JwtTokenStorage //JWT token存储(用于管理token黑名单)
	ICaptcha        *base64Captcha.Captcha //验证码工具
	Rsa             rsautil.IRsa           //RSA工具
}

// Login 管理员登录接口
func (c *Auth) Login(ctx *core.Context) error {

	// 解析json请求体到AuthLoginRequest结构体
	var in admin.AuthLoginRequest
	if err := ctx.Context.ShouldBindJSON(&in); err != nil {
		return ctx.InvalidParams(err)
	}
	// TODO 这里是如何验证图形验证码的？理论是有一个session？是存在哪里的？
	// 验证图形验证码
	if !c.ICaptcha.Verify(in.CaptchaVoucher, in.Captcha, true) {
		return ctx.InvalidParams("验证码填写不正确")
	}
	// 根据用户名或邮箱查找管理员
	adminInfo, err := c.AdminRepo.FindByWhere(ctx.Ctx(), "username = ? or email = ?", in.Username, in.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.InvalidParams("账号不存在或密码填写错误!")
		}

		return ctx.Error(err)
	}

	// 解密密码
	// TODO 这里其实可以改一下，不需要解密，直接验证密码。存储里就应该存的是加密后的数据，这样才安全
	password, err := c.Rsa.Decrypt(in.Password)
	if err != nil {
		return ctx.Error(err)
	}
	// 验证密码
	if !encrypt.VerifyPassword(adminInfo.Password, string(password)) {
		return ctx.InvalidParams("账号不存在或密码填写错误!")
	}
	// TODO 这里只允许管理员登录吗？
	// 验证管理员状态
	if adminInfo.Status != model.AdminStatusNormal {
		return ctx.Error(entity.ErrAccountDisabled)
	}
	// TODO 后面可以加上refresh token等等登录机制，将生产环境的逻辑在这里写一遍
	// 设置token过期时间
	expiresAt := time.Now().Add(12 * time.Hour)

	// 生成登录凭证
	token := jwt.GenerateToken("admin", c.Config.Jwt.Secret, &jwt.Options{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        strconv.Itoa(adminInfo.Id),
		Issuer:    "im.admin",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return ctx.Success(&admin.AuthLoginResponse{
		Auth: &admin.AccessToken{
			Type:        "Bearer",
			AccessToken: token,
			ExpiresIn:   int32(expiresAt.Unix() - time.Now().Unix()),
		},
	})
}

// Captcha 图形验证码
func (c *Auth) Captcha(ctx *core.Context) error {
	// 生成验证码和凭证
	voucher, captcha, _, err := c.ICaptcha.Generate()
	if err != nil {
		return ctx.Error(err)
	}
	// 返回验证码信息
	return ctx.Success(&admin.AuthCaptchaResponse{
		Voucher: voucher,
		Captcha: captcha,
	})
}

// Logout 退出登录接口
func (c *Auth) Logout(ctx *core.Context) error {

	// 获得当前登录用户
	session := ctx.JwtSession()
	if session != nil {
		// 如果session存在，并且没有过期，则将token加入黑名单
		if ex := session.ExpiresAt - time.Now().Unix(); ex > 0 {
			_ = c.JwtTokenStorage.SetBlackList(ctx.Ctx(), session.Token, time.Duration(ex)*time.Second)
		}
	}
	return ctx.Success(nil)
}

// Refresh Token 刷新接口
func (c *Auth) Refresh(ctx *core.Context) error {

	// TODO 业务逻辑 ...

	return ctx.Success(nil)
}
