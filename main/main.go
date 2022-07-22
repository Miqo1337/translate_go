package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
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
		//regexp.MustCompile(`<Translate[^>]*>[\n\r\s]*(.*?)[\n\r\s]*</Translate>`), ///matches only a single line
		regexp.MustCompile(`<Translate[^>]*>[\n\r\s]*((.|\n)*)[\n\r\s]*</Translate>`),
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
