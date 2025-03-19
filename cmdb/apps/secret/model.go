package secret

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/infraboard/mcube/v2/crypto/cbc"
	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/mcube/v2/types"
	"github.com/mushroomyuan/dev-clould-mini/cmdb/apps/resource"
)

func NewSecretSet() *types.Set[*Secret] {
	return types.New[*Secret]()
}

func NewSecret(in *CreateSecretRequest) *Secret {
	//  hash版本的UUID
	// 	Vendor Address ApiKey
	uid := uuid.NewMD5(uuid.Nil, fmt.Appendf(nil, "%d.%s.%s", in.Vendor, in.Address, in.ApiKey)).String()
	return &Secret{
		Id:                  uid,
		UpdateAt:            time.Now().Unix(),
		CreateSecretRequest: *in,
	}
}

type Secret struct {
	Id                  string `json:"id" bson:"_id"`
	UpdateAt            int64  `json:"update_at" bson:"update_at"`
	CreateSecretRequest `bson:"inline"`
}

func (s *Secret) SetDefault() *Secret {
	if s.SyncLimit == 0 {
		s.SyncLimit = 10
	}
	return s
}

func (s *Secret) String() string {
	return pretty.ToJSON(s)
}

func NewCreateSecretRequest() *CreateSecretRequest {
	return &CreateSecretRequest{
		Regions:   []string{},
		SyncLimit: 10,
	}
}

type CreateSecretRequest struct {
	// 名称
	Name string `json:"name"`
	//
	Vendor resource.VENDOR `json:"vendor"`
	// Vmware
	Address string `json:"address"`
	// 需要被脱敏
	// Musk
	ApiKey string `json:"api_key"`
	//
	ApiSecret string `json:"api_secret" mask:",5,4"`
	//
	isEncrypted bool

	// 资源所在区域
	Regions []string `json:"regions"`
	// 通过分页大小
	SyncLimit int64 `json:"sync_limit"`
}

func (r *CreateSecretRequest) SetIsEncrypted(v bool) {
	r.isEncrypted = v
}

func (r *CreateSecretRequest) GetSyncLimit() int64 {
	if r.SyncLimit == 0 {
		return 10
	}
	return r.SyncLimit
}

func (r *CreateSecretRequest) EncryptedApiSecret() error {
	if r.isEncrypted {
		return nil
	}
	// Hash, 对称，非对称
	// 对称加密 AES(cbc)
	// @v1,xxxx@xxxxx

	key, err := base64.StdEncoding.DecodeString(SECRET_KEY)
	if err != nil {
		return err
	}

	cipherText, err := cbc.MustNewAESCBCCihper(key).Encrypt([]byte(r.ApiSecret))
	if err != nil {
		return err
	}
	r.ApiSecret = base64.StdEncoding.EncodeToString(cipherText)
	r.SetIsEncrypted(true)
	return nil

}

func (r *CreateSecretRequest) DecryptedApiSecret() error {
	if r.isEncrypted {
		cipherdText, err := base64.StdEncoding.DecodeString(r.ApiSecret)
		if err != nil {
			return err
		}

		key, err := base64.StdEncoding.DecodeString(SECRET_KEY)
		if err != nil {
			return err
		}

		plainText, err := cbc.MustNewAESCBCCihper(key).Decrypt([]byte(cipherdText))
		if err != nil {
			return err
		}
		r.ApiSecret = string(plainText)
		r.SetIsEncrypted(false)
	}
	return nil
}
