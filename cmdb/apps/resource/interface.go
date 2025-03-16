package resource

import (
	"context"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/validator"
)

const (
	AppName = "resource"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 需要对外暴露为rpc的
	RpcServer
	// 给内部程序调用的
	DeleteResource(context.Context, *DeleteResourceRequest) error
}

func NewDeleteResourceRequest() *DeleteResourceRequest {
	return &DeleteResourceRequest{}
}

type DeleteResourceRequest struct {
	ResourceIds []string `json:"resourceIds"`
}

func (r *Resource) Validate() error {
	err := validator.Validate(r)
	if err != nil {
		return exception.NewBadRequest("参数校验失败,%s", err)
	}
	return nil
}

func (s *SearchRequest) Skip() int64 {
	return (s.PageNumber - 1) * s.PageSize
}

func NewResourceSet() *ResourceSet {
	return &ResourceSet{
		Items: make([]*Resource, 0),
	}
}

func NewSearchRequest() *SearchRequest {
	return &SearchRequest{
		PageNumber: 1,
		PageSize:   10,
		Tag:        map[string]string{},
	}
}
