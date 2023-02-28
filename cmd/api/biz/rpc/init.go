package rpc

import (
	"simple_tiktok/cmd/api/biz/rpc/base"
	"simple_tiktok/cmd/api/biz/rpc/interact"
	"simple_tiktok/cmd/api/biz/rpc/social"
)

func Init() {
	base.InitBase()
	interact.InitInteract()
	social.InitSocial()
}
