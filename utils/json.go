package utils

import (
	"io"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type JsonChain struct {
	data []byte
	err  error
	opt  *sjson.Options
}

func NewFromString(str string, opt *sjson.Options) *JsonChain {
	return &JsonChain{
		data: []byte(str),
		opt:  opt,
	}
}

func NewJsonChainFromBytes(data []byte, opt *sjson.Options) *JsonChain {
	return &JsonChain{
		data: data,
		opt:  opt,
	}
}

func NewJsonChainFromBytesWithCopy(data []byte, opt *sjson.Options) *JsonChain {
	json := JsonChain{
		data: make([]byte, len(data)),
		opt:  opt,
	}
	copy(json.data, data)
	return &json
}

func NewJsonChainFromReader(reader io.Reader, opt *sjson.Options) (*JsonChain, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &JsonChain{data: data, opt: opt}, nil
}

func (json *JsonChain) Set(path string, value any) *JsonChain {
	if json.err != nil {
		return json
	}

	json.data, json.err = sjson.SetBytesOptions(json.data, path, value, json.opt)
	return json
}

func (json *JsonChain) Delete(path string) *JsonChain {
	if json.err != nil {
		return json
	}

	json.data, json.err = sjson.DeleteBytes(json.data, path)
	return json
}

func (json *JsonChain) Get(path string) gjson.Result {
	return gjson.GetBytes(json.data, path)
}

func (json *JsonChain) Result() ([]byte, error) {
	return json.data, json.err
}

func (json *JsonChain) ResultString() (string, error) {
	return string(json.data), json.err
}

func (json *JsonChain) ResultToWriter(writer io.Writer) error {
	if json.err != nil {
		return json.err
	}
	_, err := writer.Write(json.data)
	return err
}
