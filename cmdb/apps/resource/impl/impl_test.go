package impl_test

import (
	"context"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/test"
)

var (
	ctx = context.Background()
	svc = resource.GetService()
)

func init() {
	test.Setup()
}
