package keys

import (
	"fmt"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

func (s *keysCtrlTestSuite) TestEthereumController_Import() {
	path := s.controller.Paths()[1]
	importOperation := path.Operations[logical.CreateOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, "keys/import", path.Pattern)
		assert.NotEmpty(t, importOperation)
	})

	s.T().Run("should define correct properties", func(t *testing.T) {
		properties := importOperation.Properties()

		assert.NotEmpty(t, properties.Description)
		assert.NotEmpty(t, properties.Summary)
		assert.NotEmpty(t, properties.Examples[0].Description)
		assert.NotEmpty(t, properties.Examples[0].Data)
		assert.NotEmpty(t, properties.Examples[0].Response)
		assert.NotEmpty(t, properties.Responses[200])
		assert.NotEmpty(t, properties.Responses[400])
		assert.NotEmpty(t, properties.Responses[500])
	})

	s.T().Run("handler should execute the correct use case", func(t *testing.T) {
		key := apputils.FakeKey()
		privKey := "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				formatters.NamespaceHeader: {key.Namespace},
			},
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.PrivateKeyLabel: privKey,
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.PrivateKeyLabel: {
					Type:        framework.TypeString,
					Description: "Private key in hexadecimal format",
					Required:    true,
				},
			},
		}

		s.createKeyUC.EXPECT().Execute(gomock.Any(), key.Namespace, key.ID, key.Algorithm, key.Curve, privKey, key.Tags).Return(key, nil)

		response, err := importOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, key.ID, response.Data["id"])
		assert.Equal(t, key.PublicKey, response.Data["publicKey"])
		assert.Equal(t, key.Namespace, response.Data["namespace"])
		assert.Equal(t, key.Curve, response.Data["curve"])
		assert.Equal(t, key.Algorithm, response.Data["algorithm"])
		assert.Equal(t, key.Tags, response.Data["tags"])
	})

	s.T().Run("should return same error if use case fails", func(t *testing.T) {
		privKey := "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.PrivateKeyLabel: privKey,
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.PrivateKeyLabel: {
					Type:        framework.TypeString,
					Description: "Private key in hexadecimal format",
					Required:    true,
				},
			},
		}
		expectedErr := fmt.Errorf("error")

		s.createKeyUC.EXPECT().Execute(gomock.Any(), "", "id", "algo", "curve", privKey, map[string]string{}).Return(nil, expectedErr)

		response, err := importOperation.Handler()(s.ctx, request, data)

		assert.Empty(t, response)
		assert.Equal(t, expectedErr, err)
	})
}
