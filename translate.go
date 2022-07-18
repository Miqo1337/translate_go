package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Translate(source, sourceLng, targetLng string) (string, error) {

	sourceEncoded := url.QueryEscape(source)

	url := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=" +
		sourceLng + "&tl=" + targetLng + "&dt=t&q=" + sourceEncoded

	req, err := http.Get(url)
	if err != nil {
		return "err", errors.New("error getting translate.googleapis.com")
	}

	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "err", errors.New("error reading response body")
	}

	badReq := strings.Contains(string(body), "<title>Error 400 (Bad Request)")
	if badReq {
		return "err", errors.New("error 400 (Bad Request)")
	}

	var result []interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "err", errors.New("error unmarshalling data")
	}

	var text []string

	if len(result) > 0 {
		textSlices := result[0]
		for _, slice := range textSlices.([]interface{}) {
			for _, translatedText := range slice.([]interface{}) {
				text = append(text, fmt.Sprintf("%v", translatedText))
				break
			}
		}
		finalText := strings.Join(text, " ")

		return finalText, nil
	} else {
		return "err", errors.New("no translated data in response")
	}

}
