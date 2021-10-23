package parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/devbackend/crutch-db/internal/operation"
)

type Parser interface {
	Parse(str string) (*operation.Operation, error)
}

func New() Parser {
	return new(parser)
}

type parser struct{}

func (p *parser) Parse(str string) (*operation.Operation, error) {
	str = strings.TrimSpace(str)

	if len(str) == 0 {
		return nil, errors.New("empty request")
	}

	parts := strings.Split(str, " ")

	parts[0] = strings.ToUpper(parts[0])

	op := new(operation.Operation)

	switch parts[0] {
	case "GET":
		if len(parts) < 2 {
			return nil, errors.New("empty key mame")
		}

		op.Type = operation.Get
		op.Key = parts[1]
	case "SET":
		if len(parts) < 2 {
			return nil, errors.New("empty key mame")
		}

		if len(parts) < 3 {
			return nil, errors.New("empty value for set")
		}

		op.Type = operation.Set
		op.Key = parts[1]

		value, err := p.parseValue(strings.Join(parts[2:], " "))
		if err != nil {
			return nil, err
		}

		op.Value = value
	case "DELETE":
		if len(parts) < 2 {
			return nil, errors.New("empty key mame")
		}

		op.Type = operation.Delete
		op.Key = parts[1]
	case "KEYS":
		op.Type = operation.Keys
	default:
		return nil, errors.New("unknown operation " + parts[0])
	}

	return op, nil
}

func (p *parser) parseValue(str string) (interface{}, error) {
	str = strings.TrimSpace(str)

	if str[0] == '\'' {
		if len(str) == 1 || str[len(str)-1:][0] != '\'' {
			return nil, errors.New("incomplete string")
		}

		val := str[1 : len(str)-1]

		ixQuote := strings.Index(val, "'")
		if ixQuote != -1 && val[ixQuote-1] != '\\' {
			return nil, errors.New("bad string")
		}

		return strings.ReplaceAll(val, "\\'", "'"), nil
	}

	if str[0] == '[' {
		if len(str) == 1 || str[len(str)-1:][0] != ']' {
			return nil, errors.New("bad array syntax")
		}

		if len(str) == 2 {
			return []interface{}{}, nil
		}

		elems := strings.Split(str[1:len(str)-1], ",")
		res := make([]interface{}, len(elems))

		for k, v := range elems {
			val, err := p.parseValue(v)
			if err != nil {
				return nil, errors.New("bad array syntax: " + err.Error())
			}

			res[k] = val
		}

		return res, nil
	}

	var val interface{}

	val, err := strconv.Atoi(str)
	if err != nil {
		val, err = strconv.ParseFloat(str, 64)
	}

	return val, err
}
