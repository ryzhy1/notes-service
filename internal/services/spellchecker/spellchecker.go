package spellcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const spellcheckURL = "https://speller.yandex.net/services/spellservice.json/checkText"

func CheckSpelling(text string) ([]map[string]interface{}, error) {
	const op = "spellchecker.CheckSpelling"

	resp, err := http.PostForm(spellcheckURL, url.Values{"text": {text}})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
