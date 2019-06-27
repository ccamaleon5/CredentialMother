// Code generated by go-swagger; DO NOT EDIT.

package credential

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/ccamaleon5/CredentialMother/models"
)

// CreateCredentialOKCode is the HTTP code returned for type CreateCredentialOK
const CreateCredentialOKCode int = 200

/*CreateCredentialOK successful operation

swagger:response createCredentialOK
*/
type CreateCredentialOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Credential `json:"body,omitempty"`
}

// NewCreateCredentialOK creates CreateCredentialOK with default headers values
func NewCreateCredentialOK() *CreateCredentialOK {

	return &CreateCredentialOK{}
}

// WithPayload adds the payload to the create credential o k response
func (o *CreateCredentialOK) WithPayload(payload []*models.Credential) *CreateCredentialOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create credential o k response
func (o *CreateCredentialOK) SetPayload(payload []*models.Credential) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateCredentialOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.Credential, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// CreateCredentialBadRequestCode is the HTTP code returned for type CreateCredentialBadRequest
const CreateCredentialBadRequestCode int = 400

/*CreateCredentialBadRequest Invalid credential supplied

swagger:response createCredentialBadRequest
*/
type CreateCredentialBadRequest struct {

	/*
	  In: Body
	*/
	Payload []*models.CredentialSubject `json:"body,omitempty"`
}

// NewCreateCredentialBadRequest creates CreateCredentialBadRequest with default headers values
func NewCreateCredentialBadRequest() *CreateCredentialBadRequest {

	return &CreateCredentialBadRequest{}
}

// WithPayload adds the payload to the create credential bad request response
func (o *CreateCredentialBadRequest) WithPayload(payload []*models.CredentialSubject) *CreateCredentialBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create credential bad request response
func (o *CreateCredentialBadRequest) SetPayload(payload []*models.CredentialSubject) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateCredentialBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.CredentialSubject, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// CreateCredentialNotFoundCode is the HTTP code returned for type CreateCredentialNotFound
const CreateCredentialNotFoundCode int = 404

/*CreateCredentialNotFound DID not found

swagger:response createCredentialNotFound
*/
type CreateCredentialNotFound struct {
}

// NewCreateCredentialNotFound creates CreateCredentialNotFound with default headers values
func NewCreateCredentialNotFound() *CreateCredentialNotFound {

	return &CreateCredentialNotFound{}
}

// WriteResponse to the client
func (o *CreateCredentialNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// CreateCredentialInternalServerErrorCode is the HTTP code returned for type CreateCredentialInternalServerError
const CreateCredentialInternalServerErrorCode int = 500

/*CreateCredentialInternalServerError Error Internal Server

swagger:response createCredentialInternalServerError
*/
type CreateCredentialInternalServerError struct {
}

// NewCreateCredentialInternalServerError creates CreateCredentialInternalServerError with default headers values
func NewCreateCredentialInternalServerError() *CreateCredentialInternalServerError {

	return &CreateCredentialInternalServerError{}
}

// WriteResponse to the client
func (o *CreateCredentialInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}