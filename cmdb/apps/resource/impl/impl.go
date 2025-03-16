package impl

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/grpc"
	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	ioc.Controller().Registry(&ResourceServiceImpl{})
}

type ResourceServiceImpl struct {
	resource.UnimplementedRpcServer
	ioc.ObjectImpl
	col *mongo.Collection
}

func (s *ResourceServiceImpl) Name() string {
	return resource.AppName
}

func (s *ResourceServiceImpl) Init() error {
	// 注册给grpc
	resource.RegisterRpcServer(grpc.Get().Server(), s)
	// 定义使用的集合
	s.col = ioc_mongo.DB().Collection("resources")
	return nil
}
