package destinations

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	jsoniter "github.com/json-iterator/go"
)

type Marshaler interface {
	Marshal(message types.Message) ([]byte, error)
}

type JSONMarshaler struct {
	json   jsoniter.API
	pretty bool
}

func NewJSONMarshaler(opts ...func(*JSONMarshaler)) *JSONMarshaler {
	m := &JSONMarshaler{
		json: jsoniter.ConfigCompatibleWithStandardLibrary,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func WithPrettyJSON(pretty bool) func(*JSONMarshaler) {
	return func(m *JSONMarshaler) {
		m.pretty = pretty
	}
}

func (m *JSONMarshaler) Marshal(message types.Message) ([]byte, error) {
	if m.pretty {
		return m.json.MarshalIndent(message, "", "  ")
	}
	return m.json.Marshal(message)
}
