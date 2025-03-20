package impl

import (
	"context"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *ResourceServiceImpl) Search(ctx context.Context, in *resource.SearchRequest) (*resource.ResourceSet, error) {
	set := resource.NewResourceSet()
	filter := bson.M{}
	if in.Keywords != "" {
		filter["name"] = bson.M{"$regex": in.Keywords, "$options": "im"}
	}
	if in.Type != nil {
		filter["type"] = in.Type
	}
	for k, v := range in.Tag {
		filter[k] = v
	}

	result, err := s.col.Find(ctx, filter, options.Find().SetLimit(in.PageSize).SetSkip(in.Skip()))
	if err != nil {
		return nil, err
	}

	for result.Next(ctx) {
		res := &resource.Resource{}
		if err := result.Decode(res); err != nil {
			return nil, err
		}
		set.Items = append(set.Items, res)
	}
	return set, nil
}
func (s *ResourceServiceImpl) Save(ctx context.Context, in *resource.Resource) (*resource.Resource, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	// 保持数据需要从ioc里获取一个mongodb实例
	_, err := s.col.UpdateOne(ctx, bson.M{"_id": in.Meta.Id}, bson.M{"$set": in}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return in, nil
}

func (s *ResourceServiceImpl) DeleteResource(ctx context.Context, in *resource.DeleteResourceRequest) error {
	_, err := s.col.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": in.ResourceIds}})
	if err != nil {
		return err
	}
	return nil
}
