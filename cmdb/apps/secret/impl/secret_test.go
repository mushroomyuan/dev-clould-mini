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
	req.ApiKey = ""
	req.ApiSecret = ""
	req.Regions = []string{"cn-beijing"}
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
	req := secret.NewDescribeSecretRequeset("12770802-fd7d-3378-91a9-16e12caa242e")
	ins, err := svc.DescribeSecret(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestSyncResource(t *testing.T) {
	req := secret.NewSyncResourceRequest("12770802-fd7d-3378-91a9-16e12caa242e")
	err := svc.SyncResource(ctx, req, func(rr secret.ResourceResponse) {
		t.Log(rr)
	})
	if err != nil {
		t.Fatal(err)
	}
}
