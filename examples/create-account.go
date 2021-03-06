package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
	"time"
)

func main() {
	//
	// Generate keys
	//

	// Read and decode the operator secret key
	operatorSecret, err := hedera.SecretKeyFromString(os.Getenv("OPERATOR_SECRET"))
	if err != nil {
		panic(err)
	}

	// Generate a new keypair for the new account
	secret := hedera.GenerateSecretKey()
	public := secret.Public()

	fmt.Printf("secret = %v\n", secret)
	fmt.Printf("public = %v\n", public)

	//
	// Connect to Hedera
	//

	client := hedera.Dial("testnet.hedera.com:50001")
	// TODO: client.SetRetryOnFailure(0) // default: 5
	defer client.Close()

	//
	// Send transaction to create account
	//

	nodeAccountId := hedera.NewAccountID(0, 0, 3)
	operatorAccountID := hedera.NewAccountID(0, 0, 2)
	response, err := client.CreateAccount().
		Operator(operatorAccountID).
		Node(nodeAccountId).
		Key(public).
		InitialBalance(0).
		Memo("[test] hedera-sdk-go v2").
		Sign(operatorSecret).
		Execute()

	if err != nil {
		panic(err)
	}

	transactionID := response.ID
	fmt.Printf("created account; transaction = %v\n", transactionID)

	//
	// Get receipt to prove we created it ok
	//

	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	receipt, err := client.GetTransactionReceipt(transactionID).Answer()
	if err != nil {
		panic(err)
	}

	fmt.Printf("account = %v\n", *receipt.AccountID)
}
