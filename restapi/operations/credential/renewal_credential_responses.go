// Code generated by go-swagger; DO NOT EDIT.

package credential

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/ccamaleon5/CredentialMother/models"
)

// RenewalCredentialOKCode is the HTTP code returned for type RenewalCredentialOK
const RenewalCredentialOKCode int = 200

/*RenewalCredentialOK successful operation

swagger:response renewalCredentialOK
*/
type RenewalCredentialOK struct {

	/*
	  In: Body
	*/
	Payload *models.Credential `json:"body,omitempty"`
}

// NewRenewalCredentialOK creates RenewalCredentialOK with default headers values
func NewRenewalCredentialOK() *RenewalCredentialOK {

	return &RenewalCredentialOK{}
}

// WithPayload adds the payload to the renewal credential o k response
func (o *RenewalCredentialOK) WithPayload(payload *models.Credential) *RenewalCredentialOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the renewal credential o k response
func (o *RenewalCredentialOK) SetPayload(payload *models.Credential) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RenewalCredentialOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RenewalCredentialBadRequestCode is the HTTP code returned for type RenewalCredentialBadRequest
const RenewalCredentialBadRequestCode int = 400

/*RenewalCredentialBadRequest Invalid credential supplied

swagger:response renewalCredentialBadRequest
*/
type RenewalCredentialBadRequest struct {
}

// NewRenewalCredentialBadRequest creates RenewalCredentialBadRequest with default headers values
func NewRenewalCredentialBadRequest() *RenewalCredentialBadRequest {

	return &RenewalCredentialBadRequest{}
}

// WriteResponse to the client
func (o *RenewalCredentialBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// RenewalCredentialNotFoundCode is the HTTP code returned for type RenewalCredentialNotFound
const RenewalCredentialNotFoundCode int = 404

/*RenewalCredentialNotFound Credential ID not found

swagger:response renewalCredentialNotFound
*/
type RenewalCredentialNotFound struct {
}

// NewRenewalCredentialNotFound creates RenewalCredentialNotFound with default headers values
func NewRenewalCredentialNotFound() *RenewalCredentialNotFound {

	return &RenewalCredentialNotFound{}
}

// WriteResponse to the client
func (o *RenewalCredentialNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// RenewalCredentialInternalServerErrorCode is the HTTP code returned for type RenewalCredentialInternalServerError
const RenewalCredentialInternalServerErrorCode int = 500

/*RenewalCredentialInternalServerError Error Internal Server

swagger:response renewalCredentialInternalServerError
*/
type RenewalCredentialInternalServerError struct {
}

// NewRenewalCredentialInternalServerError creates RenewalCredentialInternalServerError with default headers values
func NewRenewalCredentialInternalServerError() *RenewalCredentialInternalServerError {

	return &RenewalCredentialInternalServerError{}
}

// WriteResponse to the client
func (o *RenewalCredentialInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
