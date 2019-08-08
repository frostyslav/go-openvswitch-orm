package parser

import (
	"fmt"
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
		return fmt.Sprintf("%s", p.parseCompoundType(i, d))
	default:
		return fmt.Sprintf("%s", p.parseSimpleType(d))
	}
}

func (p *parser) parseSimpleType(data interface{}) string {
	switch keyType := data.(type) {
	case string:
		return fmt.Sprintf(typeMatching[keyType])
	case int:
		return fmt.Sprintf("int")
	case bool:
		return fmt.Sprintf("bool")
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

	var firstType, lastType string
	if containsKey {
		switch d := data[keyKey].(type) {
		case map[string]interface{}:
			lastType = p.parseKeyValueCompound(i, d)
		default:
			lastType = p.parseSimpleType(d)
		}

		if containsValue {
			switch d := data[valueKey].(type) {
			case map[string]interface{}:
				firstType = p.parseKeyValueCompound(i, d)
			default:
				firstType = p.parseSimpleType(d)
			}

			return fmt.Sprintf("map[%s]%s", firstType, lastType)
		}

		if containsMaxKey {
			if maxValue(data) > 1 {
				return fmt.Sprintf("[]%s", lastType)
			}
		}

		return lastType
	}

	return ""
}

func (p *parser) parseKeyValueCompound(i *info, data map[string]interface{}) string {
	if v, ok := data[typeKey].(string); ok {
		switch v {
		case stringType:
			if enum, ok := data[enumKey]; ok {
				var values []string
				e, ok := enum.([]interface{})
				if ok {
					if len(e) == 2 {
						if eValues, ok := e[1].([]interface{}); ok {
							for _, v := range eValues {
								values = append(values, v.(string))
							}
						}
					}
				}

				return p.createAcceptableValuesStringType(i, values)
			}

			if maxLength, ok := data[maxLengthKey]; ok {
				if val, ok := maxLength.(float64); ok {
					return p.createNamedStringType(fmt.Sprintf("StringOf%dSymbols", int(val)), fmt.Sprintf("is a string with a maximum length of %d symbols.", int(val)))
				}
			}
			return "string"
		case intType:
			minVal := minIntegerValue(data)
			maxVal := maxIntegerValue(data)
			if minVal > 0 || maxVal > 0 {
				return p.createCustomIntType(minVal, maxVal)
			} else {
				return "int"
			}
		case uuidType:
			refTable := sanitizeName(data[refTableKey].(string))
			return fmt.Sprintf("*%s", refTable)
		default:
			return fmt.Sprintf(v)
		}
	}

	return ""
}
