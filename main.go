package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var src, lng string
var exitstingTranslationfile string

var patterns []*regexp.Regexp
var translations map[string]string

func init() {
	srcArg := flag.String("src", "", "directory with source files(absolute path)")
	lngArg := flag.String("lng", "ru", "target language(alpha-2 code)")
	flag.Parse()
	src = *srcArg
	lng = *lngArg
}

func main() {
	if src == "" {
		log.Fatal("No source specified\n")
	}

	exitstingTranslationfile = fmt.Sprintf("translations_%s.json", lng)
	log.Printf("translating into %s", lng)

	patterns = []*regexp.Regexp{
		regexp.MustCompile(`<Translate[^>]*>[\n\r\s]*(.*?)[\n\r\s]*</Translate>`),
		regexp.MustCompile(`[:=]\s*'(.+)'\s*,?\s*//\s*:translate`),
		regexp.MustCompile(`\Wt\('([^']*)'.*\)\W?`),
	}

	dataFromFile, err := ioutil.ReadFile(exitstingTranslationfile)
	if err != nil {
		log.Printf("Could not not read from file %s: "+
			"\n\t Error: %s \n"+
			"\tWill create a new one", exitstingTranslationfile, err)
		translations = make(map[string]string)
	} else {
		if err = json.Unmarshal(dataFromFile, &translations); err != nil {
			log.Fatalf("Cannot unmarshal existing translation file: %s", err)
		}
	}

	err = processDir(src)
	log.Printf("finished processing: %s", err)
	transJson, err := JSONMarshal(translations)
	if err != nil {
		log.Printf("cannot encode translations: %s", err)
	}
	err = ioutil.WriteFile(exitstingTranslationfile, transJson, 0644)
	log.Printf("wrote file %s: err-%s", exitstingTranslationfile, err)

}

func processDir(dir string) error {
	return filepath.WalkDir(dir, func(path string, dirEnt os.DirEntry, err error) error {

		_, file := filepath.Split(path)

		if dirEnt.IsDir() {
			if file == "node_modules" || strings.Contains(file, "_dist") {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(file) == ".js" {
			//fmt.Printf("test\n")
			return nil
		}

		processFile(path) //check to see if any difference between file as and argument or path as an argument

		return nil

	})
}

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

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent(" ", "  ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func translateText(lng string, text string) (string, error) {
	return Translate(text, "en", lng)
}
