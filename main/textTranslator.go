package main

func translateText(lng string, text string) (string, error) {
	return Translate(text, "en", lng)
}
