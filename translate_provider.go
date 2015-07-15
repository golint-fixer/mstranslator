package mstranslator

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/st3v/tracerr"
)

type TranslateProvider struct {
	authenicator *Authenicator
}

func NewTranslateProvider(auth *Authenicator) *TranslateProvider {
	return &TranslateProvider{authenicator: auth}
}

func (t *TranslateProvider) Translate(text, from, to string) (string, error) {
	token := t.authenicator.GetToken()

	uri := fmt.Sprintf(
		"%s?text=%s&from=%s&to=%s",
		TranslationURL,
		url.QueryEscape(text),
		url.QueryEscape(from),
		url.QueryEscape(to))

	client := &http.Client{}
	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(request)
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	translation := &ResponseXML{}
	err = xml.Unmarshal(body, &translation)
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	return translation.Value, nil
}

func (t *TranslateProvider) TransformText(lang, category, text string) (string, error) {
	token := t.authenicator.GetToken()

	uri := fmt.Sprintf(
		"%s?sentence=%s&category=%s&language=%s",
		TransformTextURL,
		url.QueryEscape(text),
		url.QueryEscape(category),
		url.QueryEscape(lang))

	client := &http.Client{}
	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(request)
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	// Microsoft Server json response contain BOM, need to trim.
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	transTransform := TransformTextResponse{}
	err = json.Unmarshal(body, &transTransform)
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	return transTransform.Sentence, nil
}
