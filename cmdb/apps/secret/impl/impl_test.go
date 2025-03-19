package impl_test

import (
	"context"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/secret"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/test"
)

var (
	ctx = context.Background()
	svc = secret.GetService()
)

func init() {
	test.Setup()
}
