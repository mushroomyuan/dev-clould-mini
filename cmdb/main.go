package main

import (
	"github.com/infraboard/mcube/v2/ioc/server/cmd"
	_ "github.com/mushroomyuan/dev-clould-mini/cmdb/apps"
)

func main() {
	cmd.Start()
}
