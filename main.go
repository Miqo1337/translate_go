package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"translator/jsonMarshaller"
	"translator/processors"
	"translator/translators"
)

var src string
var exitstingTranslationfile string

func init() {
	srcArg := flag.String("src", "", "directory with source files(absolute path)")
	lngArg := flag.String("lng", "ru", "target language(alpha-2 code)")
	flag.Parse()
	src = *srcArg
	translators.Lng = *lngArg

	//logger.LogfileInit()

}

func main() {

	file, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("unable to set logfile: %v\n", err)
	}
	defer file.Close()

	log.SetOutput(file)

	if src == "" {
		log.Fatal("No source specified\n")
	}

	exitstingTranslationfile = fmt.Sprintf("translations_%s.json", translators.Lng)
	log.Printf("translating into %s", translators.Lng)

	dataFromFile, err := ioutil.ReadFile(exitstingTranslationfile)
	if err != nil {
		log.Printf("Could not not read from file %s: "+
			"\n\t Error: %s \n"+
			"\tWill create a new one", exitstingTranslationfile, err)
		translators.Translations = make(map[string]string)
	} else {
		if err = json.Unmarshal(dataFromFile, &(translators.Translations)); err != nil {
			log.Fatalf("Cannot unmarshal existing translation file: %s", err)
		}
	}

	err = processors.ProcessDir(src)
	log.Printf("finished processing: %s", err)
	transJson, err := jsonMarshaller.JSONMarshal(translators.Translations)
	if err != nil {
		log.Printf("cannot encode translations: %s", err)
	}
	err = ioutil.WriteFile(exitstingTranslationfile, transJson, 0644)
	log.Printf("wrote file %s: err-%s", exitstingTranslationfile, err)

}
