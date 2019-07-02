/*
	Credential Mother
	version 0.9
	author: Adrian Pareja Abarca
	email: adriancc5.5@gmail.com
*/

package business

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"mime/multipart"
	"net/textproto"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	strfmt "github.com/go-openapi/strfmt"
	uid "github.com/segmentio/ksuid"

	bl "github.com/ccamaleon5/CredentialMother/blockchain"
	"github.com/ccamaleon5/CredentialMother/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
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
	Subject = "Bienvenido al Curso Blockchain"

	// The HTML body for the email.
	HtmlBody = "<h1>Itama Cursos</h1><p>Tu credencial fue registrada en la Blockchain " +
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
		if err != nil {
			fmt.Println("Transaction wasn't sent")
		}

		err = sendCredentialByEmail(getReceiverMail(subject.Content), string(rawCredential))
		if err != nil {
			fmt.Printf("Failed to send email: %s", err)
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

func getReceiverMail(contentSubject interface{}) string {
	content := contentSubject.(map[string]interface{})
	email := content["email"]

	return fmt.Sprintf("%v", email)
}

func sendCredentialByEmail(destination, credential string) error { // Create a new session and specify an AWS Region.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	// Create an SES client in the session.
	svc := ses.New(sess)

	// Assemble the email.
	input, err := buildEmailInput(Sender, destination, Subject, HtmlBody,
		[]byte(credential))
	if err != nil {
		fmt.Println("Error al enviar", err)
	}

	// Attempt to send the email.
	result, err := svc.SendRawEmail(input)

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

func buildEmailInput(source, destination, subject, message string,
	credentialFile []byte) (*ses.SendRawEmailInput, error) {

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// email main header:
	h := make(textproto.MIMEHeader)
	h.Set("From", source)
	h.Set("To", destination)
	h.Set("Return-Path", source)
	h.Set("Subject", subject)
	h.Set("Content-Language", "en-US")
	h.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")
	h.Set("MIME-Version", "1.0")
	_, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}

	// body:
	h = make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "7bit")
	h.Set("Content-Type", "text/plain; charset=us-ascii")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte(message))
	if err != nil {
		return nil, err
	}

	// file attachment:
	fn := "credential.json"
	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition", "attachment; filename="+fn)
	h.Set("Content-Type", "application/json; x-unix-mode=0644; name=\""+fn+"\"")
	h.Set("Content-Transfer-Encoding", "7bit")
	part, err = writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(credentialFile)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Strip boundary line before header (doesn't work with it present)
	s := buf.String()
	if strings.Count(s, "\n") < 2 {
		return nil, fmt.Errorf("invalid e-mail content")
	}
	s = strings.SplitN(s, "\n", 2)[1]

	raw := ses.RawMessage{
		Data: []byte(s),
	}
	input := &ses.SendRawEmailInput{
		Destinations: []*string{aws.String(destination)},
		Source:       aws.String(source),
		RawMessage:   &raw,
	}

	return input, nil
}
