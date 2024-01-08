package genopenapi

import (
	"encoding/json"
	"errors"
	"io"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
)

type ContentEncoder interface {
	Encode(v interface{}) (err error)
}

func (f Format) Validate() error {
	switch f {
	case FormatJSON:
		return nil
	default:
		return errors.New("unknown format: " + string(f))
	}
}

func (f Format) NewEncoder(w io.Writer) (ContentEncoder, error) {
	switch f {

	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")

		return enc, nil
	default:
		return nil, errors.New("unknown format: " + string(f))
	}
}
