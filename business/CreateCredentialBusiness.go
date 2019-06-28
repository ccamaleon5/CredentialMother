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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
    // Replace sender@example.com with your "From" address. 
    // This address must be verified with Amazon SES.
    Sender = "cursos@itama.pe"
    
    // Replace recipient@example.com with a "To" address. If your account 
    // is still in the sandbox, this address must be verified.
    Recipient = "adriancc5.5@gmail.com"

    // Specify a configuration set. If you do not want to use a configuration
    // set, comment out the following constant and the 
    // ConfigurationSetName: aws.String(ConfigurationSet) argument below
    //ConfigurationSet = "ConfigSet"
    
    // Replace us-west-2 with the AWS Region you're using for Amazon SES.
    AwsRegion = "us-east-1"
    
    // The subject line for the email.
    Subject = "Bienvenido Curso Blockchain"
    
    // The HTML body for the email.
    HtmlBody =  "<h1>Itama Cursos</h1><p>Tu credencial fue registrada en la Blockchain " +
                "<a href='https://ropsten.etherscan.io/tx/0x1860823efe27a0be0bb48375db3c0c94ab9e8435d836113eb3481b7586a3180f'>Ethereum Block Explorer</a> using the " +
                "<a href='https://aws.amazon.com/sdk-for-go/'>Go</a>.</p>"
    
    //The email body for recipients with non-HTML email clients.
    TextBody = "Este email contiene tu Credential de acceso al curso de Blockchain, se te solicitará el JSON adjunto, por favor no lo borres o pierdas."
    
    // The character encoding for the email.
    CharSet = "UTF-8"
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

		err = sendCredentialByEmail()
		if err != nil{
			fmt.Println("Failed to send email: %s", err)
		}

		//Deprecated code to save credential into blockchain
		/*options, err = client.ConfigTransaction(privateKey)
		if err != nil {
			return nil, err
		}
		idHash := sha256.Sum256([]byte(*credential.ID))

		err = client.SetCredential(common.HexToAddress(repositoryContract), options, idHash, credentialHash)
		if err!=nil{
			fmt.Println("Transaction wasn't sent")
		}*/
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

func sendCredentialByEmail() error {// Create a new session and specify an AWS Region.
    sess, err := session.NewSession(&aws.Config{
        Region:aws.String(AwsRegion)},
    )
    
    // Create an SES client in the session.
    svc := ses.New(sess)
    
    // Assemble the email.
    input := &ses.SendEmailInput{
        Destination: &ses.Destination{
            CcAddresses: []*string{
				aws.String("mauroleonpayano@gmail.com"),
				aws.String("tabo5015@gmail.com"),
            },
            ToAddresses: []*string{
                aws.String(Recipient),
            },
        },
        Message: &ses.Message{
            Body: &ses.Body{
                Html: &ses.Content{
                    Charset: aws.String(CharSet),
                    Data:    aws.String(HtmlBody),
                },
                Text: &ses.Content{
                    Charset: aws.String(CharSet),
                    Data:    aws.String(TextBody),
                },
            },
            Subject: &ses.Content{
                Charset: aws.String(CharSet),
                Data:    aws.String(Subject),
            },
        },
        Source: aws.String(Sender),
            // Comment or remove the following line if you are not using a configuration set
        //    ConfigurationSetName: aws.String(ConfigurationSet),
    }

    // Attempt to send the email.
    result, err := svc.SendEmail(input)
    
    // Display error messages if they occur.
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case ses.ErrCodeMessageRejected:
                fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
            case ses.ErrCodeMailFromDomainNotVerifiedException:
                fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
            case ses.ErrCodeConfigurationSetDoesNotExistException:
                fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
        return err
    }
    
    fmt.Println("Email Sent!")
	fmt.Println(result)
	
	return nil
}
