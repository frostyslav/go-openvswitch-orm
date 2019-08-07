package parser

import (
	"fmt"
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

			customTypeName := fmt.Sprintf("%s%s", ct.tableName, ct.columnName)
			createCustomStringType(customTypeName, values...)
			return customTypeName
		}

		return fmt.Sprintf("string // (parse MAXLENGTH)")
	case intType:
		minVal := minIntegerValue(data)
		maxVal := maxIntegerValue(data)
		if minVal >= 0 && maxVal > 0 {
			return fmt.Sprintf("IntMin%dMax%d", minVal, maxVal)
		} else if minVal >= 0 {
			return fmt.Sprintf("IntMin%d", minVal)
		} else if maxVal > 0 {
			return fmt.Sprintf("IntMax%d", minVal)
		} else {
			return fmt.Sprintf("int")
		}
	case uuidType:
		refTable := sanitizeName(data[refTableKey].(string))
		return fmt.Sprintf("*%s", refTable)
	default:
		return fmt.Sprintf(t)
	}
}

func createCustomIntType(name string) string {
	var str strings.Builder

	str.WriteString(fmt.Sprintf("type %s int64", name))

	return str.String()
}

func createCustomStringType(name string, values ...string) {
	var str strings.Builder

	sanitizedName := sanitizeName(name)
	if len(values) > 0 {
		str.WriteString("const (\n")
		for _, val := range values {
			sanitizedValue := sanitizeName(val)
			str.WriteString(fmt.Sprintf("%s %s = 	%q\n", sanitizedValue, sanitizedName, val))
		}
		str.WriteString(")\n")
	}

	str.WriteString(fmt.Sprintf("type %s string\n", name))

	customTypes[name] = str.String()
}
