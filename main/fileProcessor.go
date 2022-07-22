package main

import (
	"io/ioutil"
	"log"
)

func processFile(file string) {
	log.Printf("Starting to process %s", file)

	data, err := ioutil.ReadFile(file)

	if err != nil {
		log.Printf("Error opening the file %s: %s", file, err)
	}

	for _, rgexp := range patterns {
		matches := rgexp.FindAllStringSubmatch(string(data), -1)

		if len(matches) == 0 {
			continue
		}

		for _, m := range matches {
			text := m[1]
			textLength := len(text)

			if text[0] == '{' && text[textLength-1] == '}' {
				if textLength > 4 && text[1] == '\'' && text[textLength-2] == '\'' {
					text = text[2 : textLength-2]
				} else {
					continue
				}
			}
			if _, exists := translations[text]; !exists {
				translation, e := translateText(lng, text)

				if e != nil {
					log.Printf("Cannot translate text to %s: %s", lng, e)
				}
				translations[text] = translation
				log.Printf("Added new translation text:\n %s", translation)
			}
		}
	}
}
