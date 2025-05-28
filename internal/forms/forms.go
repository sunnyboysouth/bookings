package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, which embeds url.Values and includes a map for errors.
type Form struct {
	url.Values
	Errors errors
}

// Valid checks if the form is valid by returning true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0 //simple way to check if the form is valid
} //

// New initializes a new Form struct with the provided url.Values and an empty errors map.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks if the specified fields are present and not empty in the form data.
func (f *Form) Required(field ...string) {
	for _, x := range field {
		value := f.Get(x)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(x, "This field cannot be blank")
		}
	}
}

// Has checks if a specific field exists in the form data.
func (f *Form) Has(field string) bool {
	//x := r.Form.Get(field)
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

func (f *Form) MinLength(field string, length int) bool {
	//x := r.Form.Get(field)
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks if the specified field contains a valid email address.
// uses the govalidator package to validate the email format.
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
