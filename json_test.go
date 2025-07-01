package main

import (
	"testing"
	"time"
)

func TestInterfaceToData(t *testing.T) {
	data := Intercom{
		Address: "Ул. Черниковская 75/2, под. 3",
		Aparts:  []int64{101, 102},
		Vendor:  "esp32",
		Status:  "online",
		MAC:     "00:11:22:33:44:55",
	}
	jsonStr, err := interfaceToData(data)
	if err != nil {
		t.Fatalf("interfaceToData returned error: %v", err)
	}
	expectedSubstring := `"mac":"00:11:22:33:44:55"`
	if !contains(jsonStr, expectedSubstring) {
		t.Errorf("Expected JSON to contain %s, got %s", expectedSubstring, jsonStr)
	}
}

func TestUnmarshallIntercom(t *testing.T) {
	jsonStr := `{
		"address": "Ул. Черниковская 75/2, под. 3",
		"aparts": [101, 102],
		"vendor": "esp32",
		"status": "online",
		"mac": "00:11:22:33:44:55"
	}`
	intercom, err := UnmarshallIntercom(jsonStr)
	if err != nil {
		t.Fatalf("UnmarshallIntercom returned error: %v", err)
	}
	if intercom.Address != "Ул. Черниковская 75/2, под. 3" || intercom.MAC != "00:11:22:33:44:55" {
		t.Errorf("Unexpected intercom data: %+v", intercom)
	}
}

func TestUnmarshallEvent(t *testing.T) {
	jsonStr := `{
		"event": "open",
		"arg": 1,
		"desc": "door opened",
		"mac": "00:11:22:33:44:55",
		"time": "2025-07-01T15:04:05Z"
	}`
	event, err := UnmarshallEvent(jsonStr)
	if err != nil {
		t.Fatalf("UnmarshallEvent returned error: %v", err)
	}
	expectedTime, _ := time.Parse(time.RFC3339, "2025-07-01T15:04:05Z")
	if event.Name != "open" || event.Args != 1 || !event.Time.Equal(expectedTime) {
		t.Errorf("Unexpected event data: %+v", event)
	}
}

// contains — простая подстрочная проверка
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
