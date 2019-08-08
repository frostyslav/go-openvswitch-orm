package parser

import (
	"fmt"
	"go/types"

	"github.com/frostyslav/gopenvswitch-db/app/sanitize"
)

// type keys
const (
	keyKey   string = "key"
	valueKey string = "value"
	maxKey   string = "max"
	minKey   string = "min"
)

// key keys
const (
	maxIntegerKey string = "maxInteger"
	minIntegerKey string = "minInteger"
	refTableKey   string = "refTable"
	refTypeKey    string = "refTypes"
	enumKey       string = "enum"
	maxLengthKey  string = "maxLength"
)

func (p *parser) parseColumnType(i *info, data interface{}) string {
	switch d := data.(type) {
	case map[string]interface{}:
		return p.parseCompoundType(i, d)
	default:
		return p.parseSimpleType(d)
	}
}

func (p *parser) parseSimpleType(data interface{}) string {
	switch keyType := data.(type) {
	case string:
		return typeMatching[keyType]
	case int:
		return types.Typ[types.Int].String()
	case bool:
		return types.Typ[types.Bool].String()
	}
	return ""
}

func (p *parser) parseCompoundType(i *info, data map[string]interface{}) string {
	var containsKey, containsValue, containsMaxKey bool
	if _, ok := data[keyKey]; ok {
		containsKey = true
	}

	if _, ok := data[valueKey]; ok {
		containsValue = true
	}

	if _, ok := data[maxKey]; ok {
		containsMaxKey = true
	}

	var key, value string
	if containsKey {
		switch d := data[keyKey].(type) {
		case map[string]interface{}:
			value = p.parseKeyValueCompound(i, d)
		default:
			value = p.parseSimpleType(d)
		}

		if containsValue {
			switch d := data[valueKey].(type) {
			case map[string]interface{}:
				key = p.parseKeyValueCompound(i, d)
			default:
				key = p.parseSimpleType(d)
			}

			return "map[" + key + "]" + value
		}

		if containsMaxKey {
			if maxValue(data) > 1 {
				return "[]" + value
			}
		}

		return value
	}

	return ""
}

func (p *parser) parseKeyValueCompound(i *info, data map[string]interface{}) string {
	if v, ok := data[typeKey].(string); ok {
		switch v {
		case stringType:
			if enum, ok := data[enumKey]; ok {
				values := enumValues(enum)
				return p.createAcceptableValuesStringType(i, values)
			}

			if _, ok := data[maxLengthKey]; ok {
				val := maxLengthValue(data)
				return p.createNamedStringType(fmt.Sprintf("StringOf%dSymbols", val), fmt.Sprintf("is a string with a maximum length of %d symbols.", val))
			}
			return types.Typ[types.String].String()
		case intType:
			minVal := minIntegerValue(data)
			maxVal := maxIntegerValue(data)
			if minVal > 0 || maxVal > 0 {
				return p.createCustomIntType(minVal, maxVal)
			} else {
				return types.Typ[types.Int].String()
			}
		case uuidType:
			refTable := sanitize.Name(data[refTableKey].(string))
			return "*" + refTable
		default:
			return v
		}
	}

	return ""
}
