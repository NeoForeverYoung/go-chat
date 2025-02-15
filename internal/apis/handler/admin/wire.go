package admin

import (
	v12 "go-chat/internal/apis/handler/admin/v1"

	"github.com/google/wire"
)

// TODO 这个wire完全不知道是干啥的，也没见过
var ProviderSet = wire.NewSet(
	v12.NewIndex,
	wire.Struct(new(v12.Auth), "*"),

	wire.Struct(new(V1), "*"),
	wire.Struct(new(V2), "*"),
)
