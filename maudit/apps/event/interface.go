package event

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
)

var (
	AppName = "event"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 存储
	SaveEvent(context.Context, *EventSet) error
	// 查询
	QueryEvent(context.Context, *QueryEventRequest) (*EventSet, error)
}

func NewQueryEventRequest() *QueryEventRequest {
	return &QueryEventRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryEventRequest struct {
	// 分页请求参数
	*request.PageRequest
}
