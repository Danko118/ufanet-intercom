package main

import "encoding/json"

func interfaceToData(data interface{}) (string, error) {
	j, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func UnmarshallIntercom(jsonData string) (*Intercom, error) {

	var intercom Intercom

	err := json.Unmarshal([]byte(jsonData), &intercom)
	if err != nil {
		return nil, err
	}

	return &intercom, nil
}

func UnmarshallEvent(jsonData string) (*Event, error) {

	var event Event

	err := json.Unmarshal([]byte(jsonData), &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}
