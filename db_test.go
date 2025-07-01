package main

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	db = mockDB
	cleanup := func() { mockDB.Close() }
	return mock, cleanup
}

func TestFecthEvents(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"mac", "event_name", "event_args", "event_desc", "event_time"}).
		AddRow("aa:bb:cc:aa:bb:cc", "open", 1, "opened by code", time.Now())

	mock.ExpectQuery("SELECT mac, event_name, event_args, event_desc, event_time").
		WillReturnRows(rows)

	events, err := FecthEvents("aa:bb:cc:aa:bb:cc")
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "open", events[0].Name)
}

func TestFecthIntercom(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"mac", "address", "vendor", "aparts", "status"}).
		AddRow("aa:bb:cc:aa:bb:cc", "Main St", "esp32", pq.Array([]int64{1, 2}), "online")

	mock.ExpectQuery("SELECT i.mac,i.address,i.vendor,i.aparts,s.status").
		WillReturnRows(rows)

	intercom, err := FecthIntercom("aa:bb:cc:aa:bb:cc")
	assert.NoError(t, err)
	assert.Equal(t, "Main St", intercom.Address)
	assert.Equal(t, "esp32", intercom.Vendor)
	assert.Equal(t, "online", intercom.Status)
	assert.Equal(t, []int64{1, 2}, intercom.Aparts)
}

func TestIntercomAppend(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	payload := `{"address":"Main St","aparts":[1,2],"vendor":"esp32"}`
	topic := "intercoms/aa:bb:cc:aa:bb:cc/info"

	mock.ExpectExec("INSERT INTO intercoms").
		WithArgs("Main St", pq.Array([]int64{1, 2}), "esp32", "aa:bb:cc:aa:bb:cc").
		WillReturnResult(sqlmock.NewResult(1, 1))

	msg := mockMessage{topic: topic, payload: []byte(payload)}
	err := IntercomAppend(msg)
	assert.NoError(t, err)
}

func TestEventAppend(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	payload := `{"event":"open","arg":1,"desc":"opened by code"}`

	mock.ExpectExec("INSERT INTO events").
		WithArgs("aa:bb:cc:aa:bb:cc", "open", 1, "opened by code").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := EventAppend(payload, "aa:bb:cc:aa:bb:cc")
	assert.NoError(t, err)
}

func TestAppendState(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	topic := "intercoms/aa:bb:cc:aa:bb:cc/status"
	payload := []byte("online")

	mock.ExpectExec("INSERT INTO intercomStatus").
		WithArgs("aa:bb:cc:aa:bb:cc", "online").
		WillReturnResult(sqlmock.NewResult(1, 1))

	msg := mockMessage{topic: topic, payload: payload}
	err := AppendState(msg)
	assert.NoError(t, err)
}

// mockMessage implements mqtt.Message

type mockMessage struct {
	topic   string
	payload []byte
}

func (m mockMessage) Duplicate() bool   { return false }
func (m mockMessage) Qos() byte         { return 0 }
func (m mockMessage) Retained() bool    { return false }
func (m mockMessage) Topic() string     { return m.topic }
func (m mockMessage) MessageID() uint16 { return 1 }
func (m mockMessage) Payload() []byte   { return m.payload }
func (m mockMessage) Ack()              {}
