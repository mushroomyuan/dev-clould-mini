package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/binding"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
)

func init() {
	ioc.Api().Registry(&ResourceApiHandler{})
}

type ResourceApiHandler struct {
	ioc.ObjectImpl
}

func (r *ResourceApiHandler) Name() string {
	return resource.AppName
}

func (r *ResourceApiHandler) Init() error {
	// 获取webservice
	ws := gorestful.ObjectRouter(r)
	tags := []string{"资源管理"}
	ws.Route(ws.GET("").To(r.Search).Doc("资源检索").
		Param(ws.PathParameter("page_size", "分页大小").DataType("integer")).
		Param(ws.PathParameter("page_number", "页码").DataType("integer")).
		Param(ws.PathParameter("keywords", "关键字过滤").DataType("integer")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(resource.ResourceSet{}).
		Returns(200, "OK", resource.ResourceSet{}).
		Returns(404, "Not Found", exception.ApiException{}))
	return nil
}

func (r *ResourceApiHandler) Search(in *restful.Request, out *restful.Response) {
	req := resource.NewSearchRequest()
	// 获取绑定参数
	if err := binding.Query.Bind(in.Request, req); err != nil {
		response.Failed(out, exception.NewBadRequest(err.Error()))
		return
	}
	// 参数校验
	// 调用查询方法
	set, err := resource.GetService().Search(in.Request.Context(), req)
	if err != nil {
		response.Failed(out, err)
	}

	response.Success(out, set)
}
