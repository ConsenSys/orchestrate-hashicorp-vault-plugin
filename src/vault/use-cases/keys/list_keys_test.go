package keys

import (
	"context"
	"fmt"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListKeys_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())

	usecase := NewListKeysUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		key := apputils.FakeKey()
		expectedKeys := []string{key.ID}
		mockStorage.EXPECT().List(ctx, storage.ComputeKeysStorageKey("", key.Namespace)).Return(expectedKeys, nil)

		keys, err := usecase.Execute(ctx, key.Namespace)

		assert.NoError(t, err)
		assert.Equal(t, expectedKeys, keys)
	})

	t.Run("should fail with same error if List fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockStorage.EXPECT().List(ctx, gomock.Any()).Return(nil, expectedErr)

		keys, err := usecase.Execute(ctx, "namespace")

		assert.Nil(t, keys)
		assert.Equal(t, expectedErr, err)
	})
}
