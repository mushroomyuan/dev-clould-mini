package impl_test

import (
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/secret"
	"testing"
)

func TestCreateSecret(t *testing.T) {
	req := secret.NewCreateSecretRequest()
	req.Name = "阿里云只读账号"
	req.Vendor = resource.VENDOR_ALIYUN
	req.ApiKey = "xxx"
	req.ApiSecret = ""
	req.Regions = []string{"ap-shanghai", "ap-guangzhou"}
	ins, err := svc.CreateSecret(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestQuerySecret(t *testing.T) {
	req := secret.NewQuerySecretRequest()
	set, err := svc.QuerySecret(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestDescribeSecret(t *testing.T) {
	req := secret.NewDescribeSecretRequeset("ebbde3e6-1525-3a3b-ae2f-16c101f15cd5")
	ins, err := svc.DescribeSecret(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
