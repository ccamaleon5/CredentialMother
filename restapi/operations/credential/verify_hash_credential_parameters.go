// Code generated by go-swagger; DO NOT EDIT.

package credential

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewVerifyHashCredentialParams creates a new VerifyHashCredentialParams object
// no default values defined in spec.
func NewVerifyHashCredentialParams() VerifyHashCredentialParams {

	return VerifyHashCredentialParams{}
}

// VerifyHashCredentialParams contains all the bound params for the verify hash credential operation
// typically these are obtained from a http.Request
//
// swagger:parameters verifyHashCredential
type VerifyHashCredentialParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*the hash credential that need be validated. Example: A2354B54566756EF12132DC233E2BB33122
	  Required: true
	  In: path
	*/
	CredentialHash string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewVerifyHashCredentialParams() beforehand.
func (o *VerifyHashCredentialParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rCredentialHash, rhkCredentialHash, _ := route.Params.GetOK("credentialHash")
	if err := o.bindCredentialHash(rCredentialHash, rhkCredentialHash, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCredentialHash binds and validates parameter CredentialHash from path.
func (o *VerifyHashCredentialParams) bindCredentialHash(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.CredentialHash = raw

	return nil
}
