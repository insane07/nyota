package i18n

import (
	"fmt"
	"nyota/backend/model"
	"testing"

	"github.com/nicksnyder/go-i18n/i18n"
)

func TestI18N(t *testing.T) {
	var user = &model.UserContext{
		UserName: "aaa",
		TenantId: "T1",
	}
	var session = &model.SessionContext{
		User: user,
		Lang: "en-US",
	}

	if "Welcome to Nyota Backend Services" != Translate(session)("program_greeting") {
		t.Errorf("Testcase failed")
	}

	session = &model.SessionContext{
		User: user,
		Lang: "xh",
	}

	if "英語 - Welcome to Nyota Backend Services" != Translate(session)("program_greeting") {
		t.Errorf("XH language Testcase failed")
	}
}

func example() {
	ENG, _ := i18n.Tfunc("en-US")
	FR, _ := i18n.Tfunc("fr-FR")
	fmt.Println(ENG)
	fmt.Println(ENG("program_greeting"))
	pMap := map[string]interface{}{"Person": "mahesh"}
	pStruct := struct{ Person string }{Person: "kumar"}
	fmt.Println(ENG("person_greeting", pMap))
	fmt.Println(ENG("person_greeting", pStruct))

	type P struct{ Person string }
	fmt.Println(ENG("person_greeting", P{"Person"}))

	fmt.Println(ENG("d_days", 1))
	fmt.Println(ENG("d_days", 5))
	fmt.Println(ENG("new_msg", 5))

	fmt.Println(FR("person_address", struct {
		Name   string
		Street string
		City   string
	}{Name: "P1", Street: "St1", City: "C1"}))

	fmt.Println(ENG("person_unread_email_count_timeframe", 3, map[string]interface{}{
		"Person":    "Bob",
		"Count":     10,
		"Timeframe": ENG("d_days", 1),
	}))
}
