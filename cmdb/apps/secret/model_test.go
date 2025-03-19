package secret

import (
	"encoding/base64"
	"github.com/infraboard/mcube/v2/crypto/cbc"
	"testing"
)

func TestMustGenRandomKey(t *testing.T) {
	t.Logf("%s", base64.StdEncoding.EncodeToString(cbc.MustGenRandomKey(cbc.AES_KEY_LEN_32)))
}
