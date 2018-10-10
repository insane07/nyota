package validation

import (
	"nyota/backend/i18n"
	"nyota/backend/model"
	"nyota/backend/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	i "github.com/nicksnyder/go-i18n/i18n"
)

func init() {
	govalidator.TagMap["custom"] = govalidator.Validator(func(str string) bool {
		return str == "duck"
	})

	govalidator.CustomTypeTagMap.Set("customNameValidator", govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
		switch v := o.(type) { // you can type switch on the context interface being validated
		case StructWithCustomField:
			// you can check and validate against some other field in the context,
			// return early or not validate against the context at all – your choice
			if v.Name1 == v.Name2 {
				return false
			}
		case User1:
			//
			return true
		default:
			// expecting some other type? Throw/panic here or continue
		}
		return true
	}))
}

//User1 struct
type User1 struct {
	Name string `valid:"length(5|10)"`
	//Email string `json:"eemail" valid:"email,custom"`
	Email string `json:"EMAIL" valid:"email~key_valid_email"`
}

//StructWithCustomField custom sample
type StructWithCustomField struct {
	Name1 string `valid:"customNameValidator"`
	Name2 string `valid:"-"`
}

//Example1 func
func Example1() {
	b, err := govalidator.ValidateStruct(&User1{"John", "Juan@abc."})
	if false == b {
		errsMap := govalidator.ErrorsByField(err)

		T := i18n.Translate(&model.SessionContext{
			Lang: "fr-FR",
			User: &model.UserContext{
				UserName: "UserA",
				TenantId: "T1"},
		})

		for k, v := range errsMap {
			// Tanslate each key from the value field
			errsMap[k] = T(v)
		}

		fmt.Println("Asaskevich Field Validation:", errsMap)
	}

	b, err = govalidator.ValidateStruct(&User1{Name: "John"})
	if false == b {
		fmt.Println("Asaskevich Validator:", err)
	}

	b, err = govalidator.ValidateStruct(&StructWithCustomField{Name1: "Mahesh", Name2: "Kumar"})
	if false == b {
		fmt.Println("Asaskevich Validator:", err)
	}

	b, err = govalidator.ValidateStruct(&StructWithCustomField{Name1: "Mahesh", Name2: "Mahesh"})
	if false == b {
		fmt.Println("Asaskevich Custom Name Validator:", err)
	}

}

//ExampleSegment struct
type ExampleSegment struct {
	ID          int    // segment identifier
	TenantID    string // Tenant ID
	CollectorID string // Foreign key to collector identifier
	Name        string `valid:"required~key_name_required,length(5|255)~key_name_length"` // Display name
	Description string `valid:"optional,length(5|255)"`                                   // Segment description
}

//UpsertSegment1 method
func UpsertSegment1(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var t ExampleSegment
	err1 := decoder.Decode(&t)

	//EDIT operation
	/*if req.Method == "PUT" {
		if 0 == t.ID {

		}
	}*/
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	valid := ValidateModel(s.TFunc, &t, w)

	if true == valid {
		//Write actual response
	}

}

// ValidateModel validates given model and writes validation errors as internationalized json error messages
func ValidateModel(tf i.TranslateFunc, st interface{}, w http.ResponseWriter) bool {

	valid, err := govalidator.ValidateStruct(st)
	if false == valid {
		errors := govalidator.ErrorsByField(err)
		for k, v := range errors {
			// Tanslate each key from the value field
			errors[k] = tf(v)
		}
		w.Header().Set(utils.HTTPContentTypeKey, utils.HTTPContentJSONValue)
		js, err := json.Marshal(errors)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		w.Write(js)
		//json.NewEncoder(w).Encode(errors)
		return false
	}
	return true
}
