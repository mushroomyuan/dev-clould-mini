package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/websocket"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/binding"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/secret"
)

func init() {
	ioc.Api().Registry(&SecretApiHandler{})
}

type SecretApiHandler struct {
	ioc.ObjectImpl
}

func (r *SecretApiHandler) Name() string {
	return secret.AppName
}

func (r *SecretApiHandler) Init() error {
	// 获取webservice
	ws := gorestful.ObjectRouter(r)
	tags := []string{"凭证管理"}
	ws.Route(ws.GET("").To(r.QuerySecret).Doc("凭证列表").
		Param(ws.QueryParameter("page_size", "分页大小").DataType("integer")).
		Param(ws.QueryParameter("page_number", "页码").DataType("integer")).
		Param(ws.QueryParameter("keywords", "关键字过滤").DataType("integer")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(SecretSet{}).
		Returns(200, "OK", SecretSet{}).
		Returns(404, "Not Found", exception.NewNotFound("")))

	ws.Route(ws.GET("/{id}").To(r.DescribeSecret).Doc("凭证详情").
		Param(ws.PathParameter("id", "凭证id").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(secret.Secret{}).
		Returns(200, "OK", secret.Secret{}).
		Returns(404, "Not Found", exception.NewNotFound("")))

	ws.Route(ws.GET("/{id}/sync").To(r.SyncResource).Doc("资源同步").
		Param(ws.PathParameter("id", "凭证Id").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(secret.ResourceResponse{}).
		Returns(200, "OK", secret.ResourceResponse{}).
		Returns(404, "Not Found", exception.NewNotFound("")))
	return nil
}

type SecretSet struct {
	Total int64            `json:"total"`
	Items []*secret.Secret `json:"items"`
}

func (r *SecretApiHandler) QuerySecret(in *restful.Request, out *restful.Response) {
	req := secret.NewQuerySecretRequest()
	// 获取绑定参数
	if err := binding.Query.Bind(in.Request, req); err != nil {
		response.Failed(out, exception.NewBadRequest(err.Error()))
		return
	}
	// 参数校验
	// 调用查询方法
	set, err := secret.GetService().QuerySecret(in.Request.Context(), req)
	if err != nil {
		response.Failed(out, err)
	}

	response.Success(out, set)
}

func (r *SecretApiHandler) DescribeSecret(in *restful.Request, out *restful.Response) {
	req := secret.NewDescribeSecretRequeset(in.PathParameter("id"))

	// 参数校验
	// 调用查询方法
	ins, err := secret.GetService().DescribeSecret(in.Request.Context(), req)
	if err != nil {
		response.Failed(out, err)
	}

	response.Success(out, ins)
}

var upgrader = websocket.Upgrader{} // use default options
// websocket api
// http 协议(短链接)--> websocket(长链接) conn.write()
// websocket: https://github.com/gorilla/websocket/tree/main/examples/chat
// Go Websocket Client: https://github.com/gorilla/websocket/blob/main/examples/echo/client.go
// Web Browser Websocket Client: https://github.com/gorilla/websocket/blob/main/examples/chat/home.html
func (r *SecretApiHandler) SyncResource(req *restful.Request, resp *restful.Response) {
	sr := secret.NewSyncResourceRequest(req.PathParameter("id"))

	// websocket upgrade
	conn, err := upgrader.Upgrade(resp, req.Request, nil)
	if err != nil {
		response.Failed(resp, err)
		return
	}

	// 业务逻辑
	err = secret.GetService().SyncResource(req.Request.Context(), sr, func(rr secret.ResourceResponse) {
		conn.WriteJSON(rr)
	})
	if err != nil {
		response.Failed(resp, err)
		return
	}
}
