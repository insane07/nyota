package i18n

import (
	"nyota/backend/logutil"
	"nyota/backend/model"

	"github.com/gobs/simplejson"

	"github.com/nicksnyder/go-i18n/i18n"
)

const (
	//DefaultLanguage to be used
	DefaultLanguage = "en-US"
)

func init() {
	logutil.Printf(nil, "Initializing I18N handler")

	//i18n.MustLoadTranslationFile("/Users/mrathinasamy/Dev/prizm/src/nyota/backend/resources/i18n/en-us.all.json")
	//i18n.MustLoadTranslationFile("/Users/mrathinasamy/Dev/prizm/src/nyota/backend/resources/i18n/fr-fr.all.json")
	//i18n.ParseTranslationFileBytes("en-us.all.json", []byte(enTranslation))

	i18n.ParseTranslationFileBytes("en-us.all.json", []byte(EnUs))
	//i18n.ParseTranslationFileBytes("fr-fr.all.json", []byte(FrFr))
	i18n.ParseTranslationFileBytes("xh.all.json", []byte(XhXh))

	// Check if i18n support is available
	T, _ := i18n.Tfunc(DefaultLanguage)
	logutil.Printf(nil, "I18N Check (en-US) (%s)", T("program_greeting"))

}

//Translate provides I18N function based on Language set on UserContext
func Translate(s *model.SessionContext) i18n.TranslateFunc {
	// Set it to US English if there are no "Acccept-Language" value passed in header
	var lang string
	if nil == s {
		lang = DefaultLanguage
	} else {
		lang = s.Lang
	}

	if lang == "" {
		lang = DefaultLanguage
	}

	//logutil.Printf(nil,"Accept-Language=%s is set for the TenantId=%s and User=%s", lang, s.User.TenantId, s.User.UserName)
	// en-US is the default language
	T, _ := i18n.Tfunc(lang, DefaultLanguage)
	return T
}

//GetAllMessages returns all tags and its associated translations
func GetAllMessages() string {

	// If we have access to bundle we can get all the messages in one shot
	var all = make(map[string]map[string]string)

	for _, lang := range i18n.LanguageTags() {
		ln := make(map[string]string)
		T, _ := i18n.Tfunc(lang)
		for _, elem := range i18n.LanguageTranslationIDs(lang) {
			ln[elem] = T(elem)
		}
		all[lang] = ln
	}

	str, err := simplejson.DumpString(all)
	if nil != err {
		return ""
	}
	return str
}
