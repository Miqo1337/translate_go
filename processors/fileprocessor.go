package processors

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"translator/translators"
)

var patterns []*regexp.Regexp
var textsTobeTranslated []string

func processFile(file string) {
	log.Printf("Starting to process %s", file)

	data, err := ioutil.ReadFile(file)

	if err != nil {
		log.Printf("Error opening the file %s: %s", file, err)
	}

	patterns = []*regexp.Regexp{
		//regexp.MustCompile(`<Translate[^>]*>[\n\r\s]*(.*?)[\n\r\s]*</Translate>`), ///matches only a single line
		regexp.MustCompile(`<Translate[^>]*>[\n\r\s]*((.|\n)*)[\n\r\s]*</Translate>`),
		regexp.MustCompile(`[:=]\s*'(.+)'\s*,?\s*//\s*:translate`),
		regexp.MustCompile(`\Wt\('([^']*)'.*\)\W?`),
	}

	for _, rgexp := range patterns {

		matches := rgexp.FindAllStringSubmatch(string(data), -1)

		if len(matches) == 0 {
			continue
		}

		for _, m := range matches {
			text := m[1]
			textLength := len(text)

			text = strings.Trim(text, "\n")
			if text[0] == '{' && text[textLength-1] == '}' {
				if textLength > 4 && text[1] == '\'' && text[textLength-2] == '\'' {
					text = text[2 : textLength-2]
				} else {
					continue
				}
			}
			translators.TranslateText(translators.Lng, text)
		}
	}
}
