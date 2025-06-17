package main

import "encoding/json"

func UnmarshallIntercom(jsonData string) (*Intercom, error) {

	var intercom Intercom

	err := json.Unmarshal([]byte(jsonData), &intercom)
	if err != nil {
		return nil, err
	}

	return &intercom, nil
}
