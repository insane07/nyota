package validation

import (
	"nyota/backend/i18n"
	"nyota/backend/model"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

//Customer struct
type Customer struct {
	Name    string   `json:"name"`
	Gender  string   `json:"gender"`
	Email   string   `json:"email"`
	Address Address  `json:"address"`
	Ips     []string `json:"ips"`
}

//Address struct
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

//NameRule 2 rules
var NameRule = []validation.Rule{
	validation.Required.Error("key_name_required"),
	validation.Length(5, 15).Error("key_name_length"),
}

//Validate method
func (c Customer) Validate() error {
	//var rules []validation.Rule

	// Add to existing Rule[]
	//rules = append(NameRule, is.Email)

	return validation.ValidateStruct(&c,
		// Name cannot be empty, and the length must be between 5 and 20.
		//validation.Field(&c.Name, rules...),
		validation.Field(&c.Name, NameRule...),
		// Gender is optional, and should be either "Female" or "Male".
		validation.Field(&c.Gender, validation.In("Female", "Male")),
		// Email cannot be empty and should be in a valid email format.
		validation.Field(&c.Email, validation.Required, is.Email),
		// Validate Address using its own validation rules
		validation.Field(&c.Address),
	)
}

//Validate method
func (a Address) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Street, validation.Required, validation.Length(5, 20)),
		validation.Field(&a.City, validation.Required),
		validation.Field(&a.State, validation.Required),
		validation.Field(&a.Zip, validation.Required),
	)
}

func ipValidator(value interface{}) error {
	s, _ := value.([]string)
	if s != nil {
		for _, ip := range s {
			e := validation.Validate(ip, is.IP)
			if nil != e {
				return errors.New(ip + " is not a valid IP")
			}
		}

	}
	return nil
}

func ipArrayValidation(ips []string) {
	e1 := validation.Validate(ips,
		validation.Required,        // not empty
		validation.By(ipValidator), // is a valid URL
	)
	if e1 != nil {
		fmt.Printf("IPValidation %s : error (%s)\n", ips, e1)
	}
}

// Example validation
func Example() {

	ipArrayValidation([]string{"1.2.3.4", "abcd"})
	ipArrayValidation(nil)

	data := "example@a.com"
	err := validation.Validate(data,
		validation.Required,       // not empty
		validation.Length(5, 100), // length between 5 and 100
		is.URL, // is a valid URL
	)
	if err != nil {
		fmt.Println("Email Validation:", err)
	}

	c := Customer{
		Name:  "abcd",
		Email: "a@b.com",
		Ips:   []string{"1.1.1.1", "abcde"},
		Address: Address{
			Street: "25th Street",
			City:   "Unknown",
			State:  "Virginia",
			Zip:    "12345",
		},
	}

	err1 := c.Validate()
	if nil != err1 {
		m := make(map[string]string)
		b, _ := json.Marshal(err1)
		json.Unmarshal(b, &m)

		T := i18n.Translate(&model.SessionContext{
			Lang: "fr-FR",
			User: &model.UserContext{
				UserName: "UserA",
				TenantId: "T1"},
		})

		for k, v := range m {
			// Tanslate each key from the value field
			m[k] = T(v)
		}
		data, _ := json.Marshal(m)
		fmt.Println(string(data))
	}

}

//Example3 method
func Example3() {

	s := Customer{
		Name:  "abcd",
		Email: "a@b.com",
		Address: Address{
			Street: "25th Street",
			City:   "Unknown",
			State:  "Virginia",
			Zip:    "12345",
		},
	}

	var fieldRules []*validation.FieldRules

	fieldRules = append(fieldRules, validation.Field(&s.Name, validation.Required, validation.Length(5, 20)))
	fieldRules = append(fieldRules, validation.Field(&s.Gender, validation.In("Female", "Male")))
	fieldRules = append(fieldRules, validation.Field(&s.Email, is.Email))

	if len(strings.TrimSpace(s.Email)) > 0 {
		fieldRules = append(fieldRules, validation.Field(&s.Address, validation.Required))
		fieldRules = append(fieldRules, validation.Field(&s.Address))
	}

	err1 := validation.ValidateStruct(&s, fieldRules...)
	if nil != err1 {
		m := make(map[string]string)
		b, _ := json.Marshal(err1)
		json.Unmarshal(b, &m)

		T := i18n.Translate(&model.SessionContext{
			Lang: "fr-FR",
			User: &model.UserContext{
				UserName: "UserA",
				TenantId: "T1"},
		})

		for k, v := range m {
			// Tanslate each key from the value field
			m[k] = T(v)
		}
		data, _ := json.Marshal(m)
		fmt.Println(string(data))
	}
}
