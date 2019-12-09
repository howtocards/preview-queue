// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// UserRender user render
// swagger:model UserRender
type UserRender struct {

	// callback
	// Required: true
	Callback URL `json:"callback"`

	// extra
	Extra Extra `json:"extra,omitempty"`

	// user
	// Required: true
	User UserPath `json:"user"`
}

// Validate validates this user render
func (m *UserRender) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCallback(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUser(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UserRender) validateCallback(formats strfmt.Registry) error {

	if err := m.Callback.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("callback")
		}
		return err
	}

	return nil
}

func (m *UserRender) validateUser(formats strfmt.Registry) error {

	if err := m.User.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("user")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *UserRender) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UserRender) UnmarshalBinary(b []byte) error {
	var res UserRender
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}