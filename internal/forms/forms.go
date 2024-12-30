package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}

	}
}

// Has check if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	validate := r.Form.Get(field)
	if validate == "" {
		return false
	}
	return true
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int, req *http.Request) bool {
	fieldFromForm := req.Form.Get(field)
	if len(fieldFromForm) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks forvalid email addresses
func (f *Form) IsEmail(field string){
	if !govalidator.IsEmail(f.Get(field)){
		f.Errors.Add(field,"Invalid email address")
	}
}
