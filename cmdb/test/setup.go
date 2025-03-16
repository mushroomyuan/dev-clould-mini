package test

import (
	"fmt"
	"github.com/infraboard/mcube/v2/ioc"
	_ "github.com/mushroomyuan/dev-clould-mini/cmdb/apps"
	"os"
)

func Setup() {
	fmt.Println(os.Getwd())
	ioc.DevelopmentSetupWithPath("/home/yfz/dev-clould-mini/cmdb/etc/application.toml")
}
