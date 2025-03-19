package secret

import (
	"context"
	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/mcube/v2/types"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
)

const (
	AppName    = "secret"
	SECRET_KEY = "KxU1RwlT8JG5JnOcgLEYks/jH8v5Z4Tp9Ha28fA7474="
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 用于Secret的管理(后台管理员配置)
	// 创建secret
	CreateSecret(context.Context, *CreateSecretRequest) (*Secret, error)
	// 查询secret
	QuerySecret(context.Context, *QuerySecretRequest) (*types.Set[*Secret], error)
	// 查询详情, 已解密, API层需要脱敏
	DescribeSecret(context.Context, *DescribeSecretRequeset) (*Secret, error)

	// 基于云商凭证来同步资源
	// 怎么API怎么设计
	// 同步阿里云所有资源, 10分钟，30分钟 ...
	// 这个接口调用持续30分钟...
	// Req ---> <---- Resp:   能快速响应的同步调用
	// SyncResource(Req) Resp

	// Stream API, websocket --> UI 当前资源同步的进度
	//SyncResource(context.Context, *SyncResourceRequest, SyncResourceHandleFunc) error
}

type SyncResourceHandleFunc func(ResourceResponse)

type ResourceResponse struct {
	Success    bool
	InstanceId string             `json:"instance_id"`
	Resource   *resource.Resource `json:"resource"`
	Message    string             `json:"message"`
}

func (t ResourceResponse) String() string {
	return pretty.ToJSON(t)
}

func NewQuerySecretRequest() *QuerySecretRequest {
	return &QuerySecretRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QuerySecretRequest struct {
	// 分页请求
	*request.PageRequest
}

func NewDescribeSecretRequeset(id string) *DescribeSecretRequeset {
	return &DescribeSecretRequeset{
		Id: id,
	}
}

type DescribeSecretRequeset struct {
	Id string `json:"id"`
}

func NewSyncResourceRequest(id string) *SyncResourceRequest {
	return &SyncResourceRequest{
		Id: id,
	}
}

type SyncResourceRequest struct {
	Id string `json:"id"`
}
