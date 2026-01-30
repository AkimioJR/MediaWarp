package utils

import (
	"io"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Json struct {
	data []byte
	err  error
	opt  *sjson.Options
}

func NewFromString(str string, opt *sjson.Options) *Json {
	return &Json{
		data: []byte(str),
		opt:  opt,
	}
}

func NewFromBytes(data []byte, opt *sjson.Options) *Json {
	return &Json{
		data: data,
		opt:  opt,
	}
}

func NewFromBytesWithCopy(data []byte, opt *sjson.Options) *Json {
	json := Json{
		data: make([]byte, len(data)),
		opt:  opt,
	}
	copy(json.data, data)
	return &json
}

func NewFromReader(reader io.Reader, opt *sjson.Options) (*Json, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &Json{data: data, opt: opt}, nil
}

func (json *Json) Set(path string, value any) *Json {
	if json.err != nil {
		return json
	}

	json.data, json.err = sjson.SetBytesOptions(json.data, path, value, json.opt)
	return json
}

func (json *Json) Delete(path string) *Json {
	if json.err != nil {
		return json
	}

	json.data, json.err = sjson.DeleteBytes(json.data, path)
	return json
}

func (json *Json) Get(path string) gjson.Result {
	return gjson.GetBytes(json.data, path)
}

func (json *Json) Result() ([]byte, error) {
	return json.data, json.err
}

func (json *Json) ResultString() (string, error) {
	return string(json.data), json.err
}

func (json *Json) ResultToWriter(writer io.Writer) error {
	if json.err != nil {
		return json.err
	}
	_, err := writer.Write(json.data)
	return err
}
