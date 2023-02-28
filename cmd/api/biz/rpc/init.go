package rpc

import (
	"simple_tiktok/cmd/api/biz/rpc/base"
	"simple_tiktok/cmd/api/biz/rpc/interact"
)

func Init() {
	base.InitBase()
	interact.InitInteract()
}
