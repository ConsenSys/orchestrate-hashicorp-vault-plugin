package keys

import (
	"context"
	"encoding/base64"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/consensys/gnark-crypto/crypto/hash"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/hashicorp/go-hclog"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/consensys/quorum/crypto"
	"github.com/hashicorp/vault/sdk/logical"
)

type signPayloadUseCase struct {
	getKeyUC usecases.GetKeyUseCase
}

func NewSignUseCase(getKeyUC usecases.GetKeyUseCase) usecases.KeysSignUseCase {
	return &signPayloadUseCase{
		getKeyUC: getKeyUC,
	}
}

func (uc signPayloadUseCase) WithStorage(storage logical.Storage) usecases.KeysSignUseCase {
	uc.getKeyUC = uc.getKeyUC.WithStorage(storage)
	return &uc
}

func (uc *signPayloadUseCase) Execute(ctx context.Context, id, namespace, data string) (string, error) {
	logger := log.FromContext(ctx).With("namespace", namespace).With("id", id)
	logger.Debug("signing message")

	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		errMessage := "data must be a base64 string"
		logger.With("error", err).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}

	key, err := uc.getKeyUC.Execute(ctx, id, namespace)
	if err != nil {
		return "", err
	}

	switch {
	case key.Algorithm == entities.EDDSA && key.Curve == entities.BN254:
		return uc.signEDDSA(logger, key.PrivateKey, dataBytes)
	case key.Algorithm == entities.ECDSA && key.Curve == entities.Secp256k1:
		return uc.signECDSA(logger, key.PrivateKey, dataBytes)
	default:
		errMessage := "invalid signing algorithm/elliptic curve combination"
		logger.Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}
}

func (uc *signPayloadUseCase) signECDSA(logger hclog.Logger, privKeyString string, data []byte) (string, error) {
	ecdsaPrivKey, err := crypto.HexToECDSA(privKeyString)
	if err != nil {
		errMessage := "failed to parse ECDSA private key"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	signatureB, err := crypto.Sign(crypto.Keccak256(data), ecdsaPrivKey)
	if err != nil {
		errMessage := "failed to sign payload with ECDSA"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	// We remove the recID from the signature (last byte).
	return base64.StdEncoding.EncodeToString(signatureB[:len(signatureB)-1]), nil
}

func (uc *signPayloadUseCase) signEDDSA(logger hclog.Logger, privKeyString string, data []byte) (string, error) {
	privKey := eddsa.PrivateKey{}
	privKeyB, _ := hexutil.Decode(privKeyString)
	_, err := privKey.SetBytes(privKeyB)
	if err != nil {
		errMessage := "failed to parse EDDSA private key"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	signatureB, err := privKey.Sign(data, hash.MIMC_BN254.New("seed"))
	if err != nil {
		errMessage := "failed to sign payload with EDDSA"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	return base64.StdEncoding.EncodeToString(signatureB), nil
}
