package impl

import (
	"context"
	"github.com/infraboard/mcube/v2/ioc/config/cache"
	"github.com/infraboard/mcube/v2/types"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/secret"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// CreateSecret implements secret.Service.
func (s *SecretServiceImpl) CreateSecret(ctx context.Context, in *secret.CreateSecretRequest) (*secret.Secret, error) {
	ins := secret.NewSecret(in)

	// 需要加密
	if err := ins.EncryptedApiSecret(); err != nil {
		return nil, err
	}

	// upsert, gorm save
	_, err := s.col.UpdateOne(ctx, bson.M{"_id": ins.Id}, bson.M{"$set": ins}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	return ins, nil
}

// DescribeSecret implements secret.Service.
// 缓存怎么做
// 1. 从缓存中去(内存， 公共的内存服务 Redis)
// 2. 能获取到，直接返回
// 3. 不能获取, 选好从本地获取，返回，再把他设置到缓存中去
// 4. 怎么实现: redis redis get(key)/set(key), obj -> JSON
// 5. https://github.com/redis/go-redis  get, set
// CacheGetter --> go-redis --> ObjectFinder
func (s *SecretServiceImpl) DescribeSecret(ctx context.Context, in *secret.DescribeSecretRequeset) (*secret.Secret, error) {
	// 封装过后的改良版
	ins := secret.NewSecret(secret.NewCreateSecretRequest())
	err := cache.NewGetter(ctx, func(ctx context.Context, objectId string) (any, error) {
		return s.describeSecret(ctx, in)
	}).Get(in.Id, ins)
	if err != nil {
		return nil, err
	}
	return ins.SetDefault(), nil

}

// @cached(ttl=30s)
// h = cached(ttl=30s) -> h
func (s *SecretServiceImpl) describeSecret(ctx context.Context, in *secret.DescribeSecretRequeset) (*secret.Secret, error) {
	// 取出后，需要解密
	e := secret.NewSecret(&secret.CreateSecretRequest{})
	// gorm take
	if err := s.col.FindOne(ctx, bson.M{"_id": in.Id}).Decode(e); err != nil {
		return nil, err
	}

	e.SetIsEncrypted(true)
	if err := e.DecryptedApiSecret(); err != nil {
		return nil, err
	}

	// 解密过后的数据
	return e, nil
}

// QuerySecret implements secret.Service.
func (s *SecretServiceImpl) QuerySecret(ctx context.Context, in *secret.QuerySecretRequest) (*types.Set[*secret.Secret], error) {
	set := secret.NewSecretSet()

	filter := bson.M{}
	cursor, err := s.col.Find(ctx, filter, options.Find().SetLimit(int64(in.PageSize)).SetSkip(in.ComputeOffset()))
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		e := secret.NewSecret(&secret.CreateSecretRequest{})
		if err := cursor.Decode(e); err != nil {
			return nil, err
		}
		e.SetDefault()
		set.Add(e)
	}

	return set, nil

}

// SyncResource implements secret.Service.
func (s *SecretServiceImpl) SyncResource(ctx context.Context, in *secret.SyncResourceRequest, cb secret.SyncResourceHandleFunc) error {
	ins, err := s.DescribeSecret(ctx, secret.NewDescribeSecretRequeset(in.Id))
	if err != nil {
		return err
	}

	return ins.Sync(func(rr secret.ResourceResponse) {
		// 进行必要数据的填充
		rr.Resource.Meta.SyncAt = time.Now().Unix()
		// 资源归属
		rr.Resource.Meta.Domain = "default"
		rr.Resource.Meta.Namespace = "default"

		// 调用resource模块来进行 资源的保存
		res, err := resource.GetService().Save(ctx, rr.Resource)
		if err != nil {
			rr.Success = false
			rr.Message = err.Error()
		} else {
			rr.Success = true
			rr.Resource = res
		}
		cb(rr)
	})
}
