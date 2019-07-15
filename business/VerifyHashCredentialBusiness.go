package business

import (
	"encoding/hex"
	"errors"
	"log"

	bl "github.com/ccamaleon5/CredentialMother/blockchain"
	"github.com/ccamaleon5/CredentialMother/models"

	"github.com/ethereum/go-ethereum/common"
)

//VerifyCredential saved into blockchain
func VerifyHashCredential(hexaCredential, nodeURL, publicAddress, verificationContract string) (*models.VerifyResponse, error) {
	client := new(bl.Client)
	err := client.Connect(nodeURL)

	if err != nil {
		return nil, err
	}

	defer client.Close()

	response := new(models.VerifyResponse)
	response.Valid = true
	errorResponse := new(models.Error)
	errorResponse.Code = "200"
	errorResponse.Message = "OK"
	response.Error = errorResponse

	contractAddress := common.HexToAddress(verificationContract)
	address := common.HexToAddress(publicAddress)

	var fail error

	//Convert Credential Hexa to Bytes and verify
	log.Printf("Verifying credential ID: %s", hexaCredential)
	
	hashCredential, err := hex.DecodeString(hexaCredential)

	var hash [32]byte 

	for i, j := range hashCredential {
		hash[i] = j
	}

	if err != nil {
		log.Println("Credential isn't a json format")
		return nil, err
	}
	result, err := client.VerifyHashCredential(contractAddress, hash, address)
	if err != nil {
		return nil, err
	}
	if !result {
		fail = errors.New("Failed verifying credentials")
	}

	if fail != nil {
		errorResponse.Code = "501"
		errorResponse.Message = "Credential isn't valid"
		return response, fail
	}

	return response, nil
}
