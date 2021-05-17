package keys

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignPayload_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks2.NewMockStorage(ctrl)
	mockGetKeyUC := mocks.NewMockGetKeyUseCase(ctrl)
	ctx := log.Context(context.Background(), log.Default())
	address := "0xaddress"
	namespace := "namespace"
	data := base64.StdEncoding.EncodeToString([]byte("my data to sign"))

	mockGetKeyUC.EXPECT().WithStorage(mockStorage).Return(mockGetKeyUC).AnyTimes()

	usecase := NewSignUseCase(mockGetKeyUC).WithStorage(mockStorage)

	t.Run("should execute use case successfully: ECDSA", func(t *testing.T) {
		key := apputils.FakeKey()
		key.Curve = entities.Secp256k1
		key.Algorithm = entities.ECDSA
		key.PrivateKey = "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw="

		mockGetKeyUC.EXPECT().Execute(ctx, address, namespace).Return(key, nil)

		signature, err := usecase.Execute(ctx, address, namespace, data)

		assert.NoError(t, err)
		assert.Equal(t, "YzQeLIN0Sd43Nbb0QCsVSqChGNAuRaKzEfujnERAJd0523aZyz2KXK93KKh+d4ws3MxAhc8qNG43wYI97Fzi7Q==", signature)
	})

	t.Run("should execute use case successfully: EDDSA", func(t *testing.T) {
		key := apputils.FakeKey()
		key.Curve = entities.BN254
		key.Algorithm = entities.EDDSA

		mockGetKeyUC.EXPECT().Execute(ctx, address, namespace).Return(key, nil)

		signature, err := usecase.Execute(ctx, address, namespace, data)

		assert.NoError(t, err)
		assert.Equal(t, "tdpR9JkX7lKSugSvYJX2icf6/uQnCAmXG9v/FG26vS0AcBqg6eVakZQNYwfic/Ec3LWqzSbXg54TBteQq6grdw==", signature)
	})

	t.Run("should fail with same error if Get Account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetKeyUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, data)

		assert.Empty(t, signature)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail if creation of EDDSA private key fails", func(t *testing.T) {
		key := apputils.FakeKey()
		key.Curve = entities.BN254
		key.Algorithm = entities.EDDSA
		key.PrivateKey = "account.PrivateKey"

		mockGetKeyUC.EXPECT().Execute(ctx, address, namespace).Return(key, nil)

		signature, err := usecase.Execute(ctx, address, namespace, data)

		assert.Empty(t, signature)
		assert.Error(t, err)
	})

	t.Run("should fail if creation of ECDSA private key fails", func(t *testing.T) {
		key := apputils.FakeKey()
		key.Curve = entities.Secp256k1
		key.Algorithm = entities.ECDSA
		key.PrivateKey = "account.PrivateKey"

		mockGetKeyUC.EXPECT().Execute(ctx, address, namespace).Return(key, nil)

		signature, err := usecase.Execute(ctx, address, namespace, data)

		assert.Empty(t, signature)
		assert.Error(t, err)
	})

	t.Run("should fail with InvalidParameterError if data is not a base64 string", func(t *testing.T) {
		key := apputils.FakeKey()
		key.Curve = entities.Secp256k1
		key.Algorithm = entities.ECDSA

		signature, err := usecase.Execute(ctx, address, namespace, "invalid data")

		assert.Empty(t, signature)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
