package translators

import "log"

var Translations map[string]string

var Lng string

func TranslateText(lng string, text string) {

	if _, exists := Translations[text]; !exists {
		translation, e := translate(text, "en", lng)

		if e != nil {
			log.Printf("Cannot translate text to %s: %s", lng, e)
		}
		Translations[text] = translation
		log.Printf("Added new translation text:\n %s", translation)
	}
}
