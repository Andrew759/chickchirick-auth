package c_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Response struct {
	*PayloadContainer
	ErrorContainer []*ErrorBody `json:"errors,omitempty"`
}

type PayloadContainer struct {
	Payload interface{}   `json:"payload"`
	buf     *bytes.Reader // Внутренний буфер для чтения
}

type ErrorBody struct {
	Message string `json:"message"`
}

func NewResponse() *Response {
	return &Response{}
}

// SendSuccess перезаписывает PayloadContainer и отправляет ответ
func (r *Response) SendSuccess(w http.ResponseWriter, payload interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	r.PayloadContainer = &PayloadContainer{
		Payload: payload,
	}

	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		r.SendError(w, err.Error(), http.StatusInternalServerError)
	}
}

// SendError добавляет одну ошибку в структуру ошибок и отправляет ответ
func (r *Response) SendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errBody := ErrorBody{
		Message: message,
	}
	r.ErrorContainer = append(r.ErrorContainer, &errBody)

	_ = json.NewEncoder(w).Encode(r)
}

// AddErrorToErrorContainer добавление ошибки в контейнер ошибок
func (r *Response) AddErrorToErrorContainer(err error) {
	r.ErrorContainer = append(r.ErrorContainer, &ErrorBody{err.Error()})
}

func (r *Response) AddErrorsToErrorContainer(errors []error) {
	for _, err := range errors {
		r.AddErrorToErrorContainer(err)
	}
}

// Send отправка данных, записанных в структуру Response
func (r *Response) Send(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(r)
}

// FirstError Возвращает первую ошибку из ErrorContainer
func (r *Response) FirstError() error {
	if len(r.ErrorContainer) > 0 {
		return errors.New(r.ErrorContainer[0].Message)
	}
	return nil
}

// Read внутренняя функция для возможности декодирования payload в сущности (используется в тестах)
func (pc *PayloadContainer) Read(p []byte) (n int, err error) {
	if pc.buf == nil {
		data, err := json.Marshal(pc.Payload)
		if err != nil {
			return 0, err
		}
		pc.buf = bytes.NewReader(data)
	}
	return pc.buf.Read(p)
}

// PayloadLength функция для подсчета количества сущностей в payload
func (pc *PayloadContainer) PayloadLength() int {
	switch pc.Payload.(type) {
	case []interface{}:
		return len(pc.Payload.([]interface{}))
	case interface{}:
		return 1
	default:
		return 0
	}
}

func (pc *PayloadContainer) GetPayloadByIndex(index int) interface{} {
	switch pc.Payload.(type) {
	case []interface{}:
		return pc.Payload.([]interface{})[index]
	case interface{}:
		if index == 0 {
			return pc.Payload
		}
	}

	panic("invalid payload, or operation is unacceptable")
}
