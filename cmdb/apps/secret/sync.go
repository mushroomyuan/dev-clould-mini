package secret

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v6/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
)

func (s *Secret) Sync(cb SyncResourceHandleFunc) error {
	switch s.Vendor {
	case resource.VENDOR_TENCENT:
		// 省略腾讯云代码
	case resource.VENDOR_ALIYUN:
		client, err := CreateClient(s)
		if err != nil {
			return err
		}

		for _, region := range s.Regions {
			request := &ecs20140526.DescribeInstancesRequest{
				RegionId:   tea.String(region), // 遍历 region
				PageSize:   tea.Int32(int32(s.SyncLimit)),
				PageNumber: tea.Int32(1), // 初始页码
			}

			runtime := &util.RuntimeOptions{}

			hasNext := true
			for hasNext {
				response, err := client.DescribeInstancesWithOptions(request, runtime)
				if err != nil {
					return err
				}

				// 处理数据
				if response.Body.Instances != nil {
					for _, instance := range response.Body.Instances.Instance {
						fmt.Println(*instance)
						cb(ResourceResponse{
							Resource: TransferInstanceToResource(instance),
						})
					}
				}

				// 判断是否还有下一页数据
				if int32(len(response.Body.Instances.Instance)) < *request.PageSize {
					hasNext = false
				} else {
					*request.PageNumber += 1 // 增加页码
				}
			}
		}

		return nil
	}

	return nil
}

func CreateClient(s *Secret) (_result *ecs20140526.Client, _err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(s.ApiKey),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(s.ApiSecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Ecs
	config.Endpoint = tea.String("ecs.cn-beijing.aliyuncs.com")
	_result = &ecs20140526.Client{}
	_result, _err = ecs20140526.NewClient(config)
	return _result, _err
}

// 云商数据结构lighthouse.Instance --> Resource
func TransferInstanceToResource(ins *ecs20140526.DescribeInstancesResponseBodyInstancesInstance) *resource.Resource {
	res := resource.NewResource()
	// 具体的转化逻辑
	res.Meta.Id = GetValue(ins.InstanceId)
	res.Spec.Name = GetValue(ins.InstanceName)
	res.Spec.Cpu = GetValue(ins.Cpu)
	res.Spec.Memory = int64(GetValue(ins.Memory))
	res.Spec.Storage = GetValue(ins.LocalStorageCapacity)
	res.Status.PrivateAddress = tea.StringSliceValue(ins.InnerIpAddress.IpAddress)
	return res
}

func GetValue[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}

	return *ptr
}
