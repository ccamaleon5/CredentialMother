/*
	Credential Mother
	version 0.9
	author: Adrian Pareja Abarca
	email: adriancc5.5@gmail.com
*/

package business

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	strfmt "github.com/go-openapi/strfmt"
	uid "github.com/segmentio/ksuid"

	bl "github.com/ccamaleon5/CredentialMother/blockchain"
	"github.com/ccamaleon5/CredentialMother/models"

	"github.com/aws/aws-sdk-go/service/ses"
)

//CreateCredential saving the hash into blockchain
func CreateCredential(subjects []*models.CredentialSubject, nodeURL string, issuer string, privateKey string, verificationContract string, repositoryContract string) ([]*models.Credential, error) {
	client := new(bl.Client)
	err := client.Connect(nodeURL)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	options, err := client.ConfigTransaction(privateKey)
	if err != nil {
		return nil, err
	}
	address := common.HexToAddress(verificationContract)
	//Verify DID validating signing

	//Iterate subject and generate Credentials
	credentials := make([]*models.Credential, 0, 50)
	for _, subject := range subjects {
		credential := new(models.Credential)
		id := generateID()
		credential.ID = &id
		types := make([]string, 0, 4)
		types = append(types, "VerifiableCredential")
		types = append(types, subject.Type)
		credential.Type = types
		credential.CredentialSubject = subject.Content
		credential.Evidence = subject.Evidence
		credential.Issuer = issuer
		credential.Proof = getProof("SmartContract", verificationContract)
		credentials = append(credentials, credential)

		rawCredential, err := json.Marshal(credential)
		if err != nil {
			return nil, errors.New("Credential isn't Json format")
		}

		fmt.Println(string(rawCredential))

		fmt.Println("date:", subject.ExpirationDate.String())

		date := time.Time(subject.ExpirationDate)

		credentialHash := sha256.Sum256(rawCredential)

		err = client.SignCredential(address, options, credentialHash, big.NewInt(date.Unix()))
		if err != nil{
			fmt.Println("Transaction wasn't sent")
		}

		options, err = client.ConfigTransaction(privateKey)
		if err != nil {
			return nil, err
		}
		idHash := sha256.Sum256([]byte(*credential.ID))

		err = client.SetCredential(common.HexToAddress(repositoryContract), options, idHash, credentialHash)
		if err!=nil{
			fmt.Println("Transaction wasn't sent")
		}
	}

	return credentials, nil

}

func getProof(typeProof string, verificationMethod string) *models.Proof {
	var proof = new(models.Proof)
	proof.Type = typeProof
	proof.VerificationMethod = verificationMethod
	proof.Created = strfmt.DateTime(time.Now())

	return proof
}

func generateID() string {
	id := uid.New()
	return id.String()
}

func sendCredentialByEmail() error {
	svc := ses.New(session.New())
	input := &ses.SendEmailInput{
    	Destination: &ses.Destination{
        	CcAddresses: []*string{
            	aws.String("adriancc5.5@gmail.com"),
        	},
        	ToAddresses: []*string{
            	aws.String("mauroleonpayano@gmail.com"),
            	aws.String("tabo5015@gmail.com"),
        	},
    	},
    	Message: &ses.Message{
        	Body: &ses.Body{
            	Html: &ses.Content{
                	Charset: aws.String("UTF-8"),
                	Data:    aws.String("This message body contains HTML formatting. It can, for example, contain links like this one: <a class=\"ulink\" href=\"http://docs.aws.amazon.com/ses/latest/DeveloperGuide\" target=\"_blank\">Amazon SES Developer Guide</a>."),
            	},
            	Text: &ses.Content{
                	Charset: aws.String("UTF-8"),
                	Data:    aws.String("This is the message body in text format."),
            	},
        	},
        	Subject: &ses.Content{
            	Charset: aws.String("UTF-8"),
            	Data:    aws.String("Registro Curso Blockchain"),
        	},
    	},
    ReturnPath:    aws.String(""),
    ReturnPathArn: aws.String(""),
    Source:        aws.String("cursos@itama.pe"),
    SourceArn:     aws.String(""),
	}


	result, err := svc.SendEmail(input)
	if err != nil {
    	if aerr, ok := err.(awserr.Error); ok {
        	switch aerr.Code() {
        	case ses.ErrCodeMessageRejected:
            	fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
        	case ses.ErrCodeMailFromDomainNotVerifiedException:
            	fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
        	case ses.ErrCodeConfigurationSetDoesNotExistException:
            	fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
        	case ses.ErrCodeConfigurationSetSendingPausedException:
            	fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
        	case ses.ErrCodeAccountSendingPausedException:
            	fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
        	default:
            	fmt.Println(aerr.Error())
        	}
    	} else {
        	// Print the error, cast err to awserr.Error to get the Code and
        	// Message from an error.
        	fmt.Println(err.Error())
    	}
    return
	}

	fmt.Println(result)
}
