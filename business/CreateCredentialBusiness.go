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
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	strfmt "github.com/go-openapi/strfmt"
	uid "github.com/segmentio/ksuid"

	bl "github.com/ccamaleon5/CredentialMother/blockchain"
	"github.com/ccamaleon5/CredentialMother/models"

	qrcode "github.com/skip2/go-qrcode"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "cursos@itama.pe"

	// Specify a configuration set. If you do not want to use a configuration
	// set, comment out the following constant and the
	// ConfigurationSetName: aws.String(ConfigurationSet) argument below
	//ConfigurationSet = "ConfigSet"

	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	AwsRegion = "us-east-1"

	// The subject line for the email.
	Subject = "Bienvenido al Curso Blockchain"

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

		err, tx := client.SignCredential(address, options, credentialHash, big.NewInt(date.Unix()))
		if err != nil {
			fmt.Println("Transaction wasn't sent")
		}

		qrFile, err := generateQR("http://cursos.itama.com/", credentialHash, getNameSubject(subject.Content))
		if err != nil {
			fmt.Printf("Failed generate QR: %s", err)
		}

		err = sendCredentialByEmail(getReceiverMail(subject.Content), string(rawCredential), tx, getNameSubject(subject.Content), getLastNameSubject(subject.Content), qrFile)
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

func sendCredentialByEmail(destination, credential, tx, firstName, lastName string, qrFile []byte) error {
	jsonData := models.MailRequest{FolderId: 13233, Name: firstName + ".png", FileData: base64.StdEncoding.EncodeToString(qrFile)}
	jsonValue, _ := json.Marshal(jsonData)

	request, err := http.NewRequest("POST", "https://us7.api.mailchimp.com/3.0/file-manager/files", bytes.NewBuffer(jsonValue))
	request.Header.Add("Authorization", "/apikey 6e29f9fe593c0caaef0caf99ef8ef499-us7")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string([]byte(body)))

	//Send Email to Alumn

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	fullURL := getFullURL(result)

	mergeFields := models.MergeFields{FNAME: firstName, LNAME: lastName, MMERGE5: fullURL, MMERGE6: tx}
	emailData := models.SendMailRequest{EmailAddress: destination, Status: "subscribed", MergeFields: mergeFields}
	emailValue, _ := json.Marshal(emailData)
	request, err = http.NewRequest("POST", "https://us7.api.mailchimp.com/3.0/lists/4a956ef616/members", bytes.NewBuffer(emailValue))
	request.Header.Add("Authorization", "/apikey 6e29f9fe593c0caaef0caf99ef8ef499-us7")

	resp, err = client.Do(request)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
		return err
	}

	body, _ = ioutil.ReadAll(resp.Body)
	log.Println(string([]byte(body)))

	return nil
}

/*
func sendCredentialByEmail(destination, credential, tx string) error { // Create a new session and specify an AWS Region.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	// Create an SES client in the session.
	svc := ses.New(sess)

	// Assemble the email.
	// The HTML body for the email.
	HtmlBody := "<h1>Itama Cursos</h1><p>Tu credencial fue registrada en la Blockchain " +
		"<a href='https://ropsten.etherscan.io/tx/" + tx + "'>Ethereum Block Explorer</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>Go</a>.</p>"

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
}*/

func generateQR(url string, hashCredential [32]byte, filename string) ([]byte, error) {
	var hash = make([]byte, 32, 64)

	for i, j := range hashCredential {
		hash[i] = j
	}

	log.Println("qrrr hex:", hex.EncodeToString(hash))
	log.Println("qrrr:", hash)
	qrFile, err := qrcode.Encode(url+hex.EncodeToString(hash), qrcode.Medium, 256)
	return qrFile, err
}

func generateID() string {
	id := uid.New()
	return id.String()
}

func getNameSubject(contentSubject interface{}) string {
	content := contentSubject.(map[string]interface{})
	name := content["name"]

	return fmt.Sprintf("%v", name)
}

func getLastNameSubject(contentSubject interface{}) string {
	content := contentSubject.(map[string]interface{})
	lastname := content["lastname"]

	return fmt.Sprintf("%v", lastname)
}

func getReceiverMail(contentSubject interface{}) string {
	content := contentSubject.(map[string]interface{})
	email := content["email"]

	return fmt.Sprintf("%v", email)
}

func getFullURL(contentSubject interface{}) string {
	content := contentSubject.(map[string]interface{})
	fullURL := content["full_size_url"]

	return fmt.Sprintf("%v", fullURL)
}
