package modules 
import (
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"fmt"
)

func CreateAccount(accountName string, accountPath string){
	
	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		panic(err)
	}

	account, mnemonic, err := registry.Create(accountName)
	if err != nil {
		panic(err)
	}

	fmt.Println("new account created: ",account, mnemonic)
}