package parser

import (
	"fmt"
	"github.com/frostyslav/gopenvswitch-db/app/xmlschema"
	"math"
	"strings"
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
)

type columnType struct {
	tableName     string
	columnName    string
	containsKey   bool
	containsValue bool
	containsMin   bool
	containsMax   bool
	keyRawData    interface{}
	valueRawData  interface{}
	rawData       interface{}
	xmlDoc        *xmlschema.Database
}

type keyOfType struct {
	rawData interface{}
}

type valueOfType struct {
	rawData interface{}
}

func (ct *columnType) parse() string {
	switch t := ct.rawData.(type) {
	case string:
		return fmt.Sprint(typeMatching[t])
	case bool:
		return fmt.Sprintf("bool")
	case map[string]interface{}:
		return fmt.Sprintf("%s", ct.parseCompound())
	default:
		return fmt.Sprintf("\n")
	}
}

func (ct *columnType) parseCompound() string {

	data := ct.rawData.(map[string]interface{})
	if _, ok := data[keyKey]; ok {
		ct.containsKey = true
	}

	if _, ok := data[valueKey]; ok {
		ct.containsValue = true
	}

	if _, ok := data[maxKey]; ok {
		ct.containsMax = true
	}

	var firstType, lastType string
	if ct.containsKey {
		ct.keyRawData = data[keyKey]
		switch ct.keyRawData.(type) {
		case map[string]interface{}:
			lastType = ct.parseCompoundKey()
		default:
			lastType = ct.parseSimpleKey()
		}

		if ct.containsValue {
			ct.valueRawData = data[valueKey]
			switch ct.valueRawData.(type) {
			case map[string]interface{}:
				firstType = ct.parseCompoundValue()
			default:
				firstType = ct.parseSimpleValue()
			}

			return fmt.Sprintf("map[%s]%s", firstType, lastType)
		}

		if ct.containsMax {
			if maxValue(ct.rawData.(map[string]interface{})) > 1 {
				return fmt.Sprintf("[]%s", lastType)
			}
		}

		return fmt.Sprintf("%s", lastType)
	}

	return ""
}

func (ct *columnType) parseSimpleKey() string {
	switch keyType := ct.keyRawData.(type) {
	case string:
		return fmt.Sprintf(typeMatching[keyType])
	case int:
		return fmt.Sprintf("int")
	case bool:
		return fmt.Sprintf("bool")
	}
	return ""
}

func (ct *columnType) parseCompoundKey() string {
	key := ct.keyRawData.(map[string]interface{})
	if v, ok := key[typeKey].(string); ok {
		return ct.parseGenericCompound(v, key)
	}

	return ""
}

func (ct *columnType) parseSimpleValue() string {
	switch valueType := ct.valueRawData.(type) {
	case string:
		return fmt.Sprintf(typeMatching[valueType])
	case int:
		return fmt.Sprintf("int")
	case bool:
		return fmt.Sprintf("bool")
	}
	return ""
}

func (ct *columnType) parseCompoundValue() string {
	value := ct.valueRawData.(map[string]interface{})
	if v, ok := value[typeKey].(string); ok {
		return ct.parseGenericCompound(v, value)
	}

	return ""
}

func (ct *columnType) parseGenericCompound(t string, data map[string]interface{}) string {
	switch t {
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

			customTypeName := ct.createCustomStringType(values)
			return customTypeName
		}

		return "string // (parse MAXLENGTH)"
	case intType:
		minVal := minIntegerValue(data)
		maxVal := maxIntegerValue(data)
		if minVal > 0 || maxVal > 0 {
			customTypeName := ct.createCustomIntType(minVal, maxVal)
			return customTypeName
		} else {
			return "int"
		}
	case uuidType:
		refTable := sanitizeName(data[refTableKey].(string))
		return fmt.Sprintf("*%s", refTable)
	default:
		return fmt.Sprintf(t)
	}
}

func (ct *columnType) createCustomIntType(min, max int64) string {
	var customTypeName string
	intSize := "int"
	if min > 0 && max > 0 {
		intSize = intSizeFromMaxValue(max)
		customTypeName = fmt.Sprintf("IntFrom%dTo%d", min, max)
	} else if min > 0 {
		intSize = "int64"
		customTypeName = fmt.Sprintf("IntFrom%d", min)
	} else if max > 0 {
		intSize = intSizeFromMaxValue(max)
		customTypeName = fmt.Sprintf("IntFrom0To%d", max)
	}

	if _, ok := customTypes[customTypeName]; !ok {
		var str strings.Builder

		str.WriteString(fmt.Sprintf("// %s custom integer type of %s size.\n", customTypeName, intSize))
		str.WriteString(fmt.Sprintf("type %s %s\n", customTypeName, intSize))

		customTypes[customTypeName] = str.String()
	}

	return customTypeName
}

func intSizeFromMaxValue(val int64) string {
	if val < math.MaxInt8 {
		return "int8"
	} else if val < math.MaxInt16 {
		return "int16"
	} else if val < math.MaxInt32 {
		return "int32"
	} else if val < math.MaxInt64 {
		return "int64"
	} else {
		return "int"
	}
}

func (ct *columnType) createCustomStringType(values []string) string {
	var str strings.Builder

	customTypeName := fmt.Sprintf("%s%s", sanitizeName(ct.tableName), sanitizeName(ct.columnName))
	if len(values) > 0 {
		str.WriteString("const (\n")
		for _, val := range values {
			sanitizedValue := sanitizeName(val)
			keyName := fmt.Sprintf("%s%s", customTypeName, sanitizedValue)
			str.WriteString(ct.addKeyComment(keyName, val))
			str.WriteString(fmt.Sprintf("%s %s = %q\n", keyName, customTypeName, val))
		}
		str.WriteString(")\n")
	}

	str.WriteString(fmt.Sprintf("type %s string\n", customTypeName))

	customTypes[customTypeName] = str.String()

	return customTypeName
}

func (ct *columnType) addKeyComment(keyName, rawKeyName string) string {
	return fmt.Sprintf("// %s %s\n", keyName, ct.xmlDoc.KeyDescription(ct.tableName, ct.columnName, rawKeyName))
}
