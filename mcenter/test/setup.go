package test

import (
	"fmt"
	"os"

	"github.com/infraboard/mcube/v2/ioc"
)

func SetUp() {
	fmt.Println(os.Getenv("McenterConfigPath"))
	ioc.DevelopmentSetupWithPath(os.Getenv("McenterConfigPath"))
}
