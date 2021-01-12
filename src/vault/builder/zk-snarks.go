package builder

import (
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/zk-snarks"
)

type zkSnarksUseCases struct {
	createAccount  usecases.CreateZksAccountUseCase
	getAccount     usecases.GetZksAccountUseCase
	listAccounts   usecases.ListZksAccountsUseCase
	listNamespaces usecases.ListZksNamespacesUseCase
	sign           usecases.ZksSignUseCase
}

func NewZkSnarksUseCases() usecases.ZksUseCases {
	getAccount := zksnarks.NewGetZksAccountUseCase()
	return &zkSnarksUseCases{
		createAccount:  zksnarks.NewCreateAccountUseCase(),
		getAccount:     getAccount,
		listAccounts:   zksnarks.NewListZksAccountsUseCase(),
		listNamespaces: zksnarks.NewListNamespacesUseCase(),
		sign:           zksnarks.NewSignUseCase(getAccount),
	}
}

func (z *zkSnarksUseCases) CreateAccount() usecases.CreateZksAccountUseCase {
	return z.createAccount
}

func (z *zkSnarksUseCases) GetAccount() usecases.GetZksAccountUseCase {
	return z.getAccount
}

func (z *zkSnarksUseCases) ListAccounts() usecases.ListZksAccountsUseCase {
	return z.listAccounts
}

func (z *zkSnarksUseCases) ListNamespaces() usecases.ListZksNamespacesUseCase {
	return z.listNamespaces
}

func (z *zkSnarksUseCases) SignPayload() usecases.ZksSignUseCase {
	return z.sign
}
