package v1

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
	"go-chat/internal/service/organize"
)

type User struct {
	service      *service.UserService
	smsService   *service.SmsService
	organizeServ *organize.OrganizeService
}

func NewUser(service *service.UserService, smsService *service.SmsService, organizeServ *organize.OrganizeService) *User {
	return &User{service: service, smsService: smsService, organizeServ: organizeServ}
}

// Detail 个人用户信息
func (u *User) Detail(ctx *ichat.Context) error {

	user, err := u.service.Dao().FindById(ctx.UserId())
	if err != nil {
		return ctx.Error(err.Error())
	}

	return ctx.Success(&web.UserDetailResponse{
		Id:       int32(user.Id),
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
		Email:    user.Email,
		Birthday: user.Birthday,
	})
}

// Setting 用户设置
func (u *User) Setting(ctx *ichat.Context) error {

	uid := ctx.UserId()

	user, _ := u.service.Dao().FindById(uid)

	isOk, _ := u.organizeServ.Dao().IsQiyeMember(uid)

	return ctx.Success(&web.UserSettingResponse{
		UserInfo: &web.UserSettingResponse_UserInfo{
			Uid:      int32(user.Id),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Motto:    user.Motto,
			Gender:   int32(user.Gender),
			IsQiye:   isOk,
			Mobile:   user.Mobile,
			Email:    user.Email,
		},
		Setting: &web.UserSettingResponse_ConfigInfo{},
	})
}

// ChangeDetail 修改个人用户信息
func (u *User) ChangeDetail(ctx *ichat.Context) error {

	params := &web.UserDetailUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if params.Birthday != "" {
		if !timeutil.IsDateFormat(params.Birthday) {
			return ctx.InvalidParams("birthday 格式错误")
		}
	}

	_, _ = u.service.Dao().BaseUpdate(&model.Users{}, entity.MapStrAny{
		"id": ctx.UserId(),
	}, entity.MapStrAny{
		"nickname": params.Nickname,
		"avatar":   params.Avatar,
		"gender":   params.Gender,
		"motto":    params.Motto,
		"birthday": params.Birthday,
	})

	return ctx.Success(nil, "个人信息修改成功！")
}

// ChangePassword 修改密码接口
func (u *User) ChangePassword(ctx *ichat.Context) error {

	params := &web.UserPasswordUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	if uid == 2054 || uid == 2055 {
		return ctx.BusinessError("预览账号不支持修改密码！")
	}

	if err := u.service.UpdatePassword(ctx.UserId(), params.OldPassword, params.NewPassword); err != nil {
		return ctx.BusinessError("密码修改失败！")
	}

	return ctx.Success(nil, "密码修改成功！")
}

// ChangeMobile 修改手机号接口
func (u *User) ChangeMobile(ctx *ichat.Context) error {

	params := &web.UserMobileUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	if uid == 2054 || uid == 2055 {
		return ctx.BusinessError("预览账号不支持修改手机号！")
	}

	if !u.smsService.CheckSmsCode(ctx.Ctx(), entity.SmsChangeAccountChannel, params.Mobile, params.SmsCode) {
		return ctx.BusinessError("短信验证码填写错误！")
	}

	user, _ := u.service.Dao().FindById(uid)

	if user.Mobile != params.Mobile {
		return ctx.BusinessError("手机号与原手机号一致无需修改！")
	}

	if !encrypt.VerifyPassword(user.Password, params.Password) {
		return ctx.BusinessError("账号密码填写错误！")
	}

	_, err := u.service.Dao().BaseUpdate(&model.Users{}, entity.MapStrAny{"id": user.Id}, entity.MapStrAny{"mobile": params.Mobile})
	if err != nil {
		return ctx.BusinessError("手机号修改失败！")
	}

	return ctx.Success(nil, "手机号修改成功！")
}

// ChangeEmail 修改邮箱接口
func (u *User) ChangeEmail(ctx *ichat.Context) error {

	params := &web.UserEmailUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	// todo 1.验证邮件激活码是否正确

	return nil
}
