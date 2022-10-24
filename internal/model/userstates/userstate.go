package userstates

import (
	"encoding/json"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
)

const (
	ExpectedCommand = iota

	ExpectedCurrency
	IncorrectCurrency

	ExpectedCategory
	IncorrectCategory

	ExpectedAmount
	IncorrectAmount
	LimitExceededAmount

	ExpectedDate
	IncorrectDate

	ExpectedSetLimitMonth
	IncorrectSetLimitMonth

	ExpectedDelLimitMonth
	IncorrectDelLimitMonth

	ExpectedSetLimitAmount
	IncorrectSetLimitAmount
)

const (
	AmountNotAdded = iota
	AmountAddedUnconverted
	AmountAddedConverted
)

type UserState struct {
	UserId      int64
	Currency    string
	status      int
	inputBuffer map[string]interface{}
}

func NewUserState(UserId int64) *UserState {
	return &UserState{
		UserId:      UserId,
		Currency:    "RUB",
		status:      ExpectedCommand,
		inputBuffer: make(map[string]interface{}),
	}
}

func CreateUserState(userId int64, Currency string, status int, inputBuffer map[string]interface{}) *UserState {
	return &UserState{
		UserId:      userId,
		Currency:    Currency,
		status:      status,
		inputBuffer: inputBuffer,
	}
}

func (s *UserState) BeforeSave() {
	if s.Currency == "" {
		s.Currency = "RUB"
	}
}

func (s *UserState) GetStatus() int {
	return s.status
}

func (s *UserState) SetStatus(newStatus int) {
	s.status = newStatus
	switch newStatus {
	case ExpectedCommand:
		s.clearInputBuffer()
	}
}

func (s *UserState) clearInputBuffer() {
	for key := range s.inputBuffer {
		delete(s.inputBuffer, key)
	}
}

func (s *UserState) GetBufferValue(key string) interface{} {
	return s.inputBuffer[key]
}

func (s *UserState) SetBufferValue(key string, value interface{}) {
	s.inputBuffer[key] = value
}

func (s *UserState) BufferValueExists(key string) bool {
	_, ok := s.inputBuffer[key]
	return ok
}

func (s *UserState) ClearBufferValue(key string) {
	delete(s.inputBuffer, key)
}

func (s *UserState) GetJSONBuffer() (string, error) {
	jsonBuffer, err := json.Marshal(s.inputBuffer)
	if err != nil {
		return "", err
	}
	return string(jsonBuffer), nil
}

func (s *UserState) IfFloatTransformToInt(key string) (int, error) {
	switch val := s.GetBufferValue(key).(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	default:
		return 0, localerr.ErrNotNumericValue
	}
}
