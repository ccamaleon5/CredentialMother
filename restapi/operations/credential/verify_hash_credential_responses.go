// Code generated by go-swagger; DO NOT EDIT.

package credential

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/ccamaleon5/CredentialMother/models"
)

// VerifyHashCredentialOKCode is the HTTP code returned for type VerifyHashCredentialOK
const VerifyHashCredentialOKCode int = 200

/*VerifyHashCredentialOK successful operation

swagger:response verifyHashCredentialOK
*/
type VerifyHashCredentialOK struct {

	/*
	  In: Body
	*/
	Payload *models.VerifyResponse `json:"body,omitempty"`
}

// NewVerifyHashCredentialOK creates VerifyHashCredentialOK with default headers values
func NewVerifyHashCredentialOK() *VerifyHashCredentialOK {

	return &VerifyHashCredentialOK{}
}

// WithPayload adds the payload to the verify hash credential o k response
func (o *VerifyHashCredentialOK) WithPayload(payload *models.VerifyResponse) *VerifyHashCredentialOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the verify hash credential o k response
func (o *VerifyHashCredentialOK) SetPayload(payload *models.VerifyResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *VerifyHashCredentialOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// VerifyHashCredentialBadRequestCode is the HTTP code returned for type VerifyHashCredentialBadRequest
const VerifyHashCredentialBadRequestCode int = 400

/*VerifyHashCredentialBadRequest Invalid credential supplied

swagger:response verifyHashCredentialBadRequest
*/
type VerifyHashCredentialBadRequest struct {
}

// NewVerifyHashCredentialBadRequest creates VerifyHashCredentialBadRequest with default headers values
func NewVerifyHashCredentialBadRequest() *VerifyHashCredentialBadRequest {

	return &VerifyHashCredentialBadRequest{}
}

// WriteResponse to the client
func (o *VerifyHashCredentialBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// VerifyHashCredentialNotFoundCode is the HTTP code returned for type VerifyHashCredentialNotFound
const VerifyHashCredentialNotFoundCode int = 404

/*VerifyHashCredentialNotFound DID not found

swagger:response verifyHashCredentialNotFound
*/
type VerifyHashCredentialNotFound struct {
}

// NewVerifyHashCredentialNotFound creates VerifyHashCredentialNotFound with default headers values
func NewVerifyHashCredentialNotFound() *VerifyHashCredentialNotFound {

	return &VerifyHashCredentialNotFound{}
}

// WriteResponse to the client
func (o *VerifyHashCredentialNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// VerifyHashCredentialInternalServerErrorCode is the HTTP code returned for type VerifyHashCredentialInternalServerError
const VerifyHashCredentialInternalServerErrorCode int = 500

/*VerifyHashCredentialInternalServerError Error Internal Server

swagger:response verifyHashCredentialInternalServerError
*/
type VerifyHashCredentialInternalServerError struct {
}

// NewVerifyHashCredentialInternalServerError creates VerifyHashCredentialInternalServerError with default headers values
func NewVerifyHashCredentialInternalServerError() *VerifyHashCredentialInternalServerError {

	return &VerifyHashCredentialInternalServerError{}
}

// WriteResponse to the client
func (o *VerifyHashCredentialInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}