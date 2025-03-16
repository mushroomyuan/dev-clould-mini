package impl_test

import (
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	resp, err := svc.Save(ctx, &resource.Resource{
		Spec: &resource.Spec{
			Name: "test",
		},
		Status: &resource.Status{},
		Meta: &resource.Meta{
			Id:        "ins-001",
			Domain:    "test",
			Namespace: "default",
			SyncAt:    time.Now().Unix(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
