// chain.go
package modules

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
)

func CreateAccountMnemonic(mnemonic string, name string, accountPath string) {

	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		return
	}

	algos, _ := registry.Keyring.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algos)
	if err != nil {
		return
	}

	registryPath := hd.CreateHDPath(sdktypes.GetConfig().GetCoinType(), 0, 0).String()
	record, err := registry.Keyring.NewAccount(name, mnemonic, "", registryPath, algo)
	if err != nil {
		return
	}

	account := cosmosaccount.Account{
		Name:   name,
		Record: record,
	}

	fmt.Println("new account created: ", account, mnemonic)
}
