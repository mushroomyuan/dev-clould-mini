package main

import (
	// 引入模块 iam
	_ "github.com/infraboard/modules/iam/init"

	// 启动Ioc
	"github.com/infraboard/mcube/v2/ioc/server/cmd"
	// 非功能组件
	_ "github.com/infraboard/mcube/v2/ioc/apps/apidoc/restful"
)

func main() {
	cmd.Start()
}
