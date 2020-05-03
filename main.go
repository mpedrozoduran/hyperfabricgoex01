package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"strconv"
)

type Bank struct {
	contractapi.Contract
}

func (b *Bank) InitAccount(ctx contractapi.TransactionContextInterface, account string, initialAmountAccount int) error {
	log.Println("Calling Bank Init")
	err := ctx.GetStub().PutState(account, []byte(strconv.Itoa(initialAmountAccount)))
	if err != nil {
		return err
	}
	return nil
}

func (b *Bank) Transfer(ctx contractapi.TransactionContextInterface, accountFrom string, accountTo string, amount int) error {
	log.Println("Calling Bank Transfer")
	amountFromBytes, err := ctx.GetStub().GetState(accountFrom)
	if err != nil {
		return fmt.Errorf("Error trying to get amount from accountFrom %s", accountFrom)
	}
	if amountFromBytes == nil {
		return fmt.Errorf("Entity accountFrom %s not found", accountFrom)
	}
	amountFrom, _ := strconv.Atoi(string(amountFromBytes))

	amountToBytes, err := ctx.GetStub().GetState(accountTo)
	if err != nil {
		return fmt.Errorf("Error trying to get amount from accountTo %s", accountTo)
	}
	if amountToBytes == nil {
		return fmt.Errorf("Entity accountTo %s not found", accountTo)
	}
	amountTo, _ := strconv.Atoi(string(amountToBytes))

	amountFrom = amountFrom - amount
	amountTo = amountTo + amount

	err = ctx.GetStub().PutState(accountFrom, []byte(strconv.Itoa(amountFrom)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(accountTo, []byte(strconv.Itoa(amountTo)))
	if err != nil {
		return err
	}

	return nil
}

func (b *Bank) Query(ctx contractapi.TransactionContextInterface, account string) (string, error) {
	log.Println("Calling Bank Query")
	accAmountBytes, err := ctx.GetStub().GetState(account)
	if err != nil {
		jsonResp := "{\"Error\": \"Failed to get state for account " + account + "}"
		return "", errors.New(jsonResp)
	}
	if accAmountBytes == nil {
		jsonResp := "{\"Error\": \"Failed to get state for account " + account + "}"
		return "", errors.New(jsonResp)
	}
	jsonResp := "{\"Account\": \"" + account + ",\"Amount\": " + string(accAmountBytes) + "}"
	log.Println(jsonResp)
	return string(accAmountBytes), nil
}

func main() {
	chainCode, err := contractapi.NewChaincode(new(Bank))
	if err != nil {
		panic(err.Error())
	}
	if err := chainCode.Start(); err != nil {
		log.Fatalf("Error starting Bank chaincode: %s", err)
	}
}
