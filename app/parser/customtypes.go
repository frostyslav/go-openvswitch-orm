package parser

import (
	"fmt"
	"math"
	"strings"
)

func (p *parser) createNamedStringType(name, comment string) string {
	var str strings.Builder

	str.WriteString(p.addComment(fmt.Sprintf("%s %s\n", name, comment)))
	str.WriteString(fmt.Sprintf("type %s string\n", name))

	p.customTypes[name] = str.String()

	return name
}

func (p *parser) createAcceptableValuesStringType(i *info, values []string) string {
	var str strings.Builder

	customTypeName := fmt.Sprintf("%s%s", i.structName, i.fieldName)
	if len(values) > 0 {
		str.WriteString("const (\n")
		for _, val := range values {
			sanitizedDbKeyValue := sanitizeName(val)
			variableName := fmt.Sprintf("%s%s", customTypeName, sanitizedDbKeyValue)
			str.WriteString(p.addVariableComment(i, variableName, val))
			str.WriteString(fmt.Sprintf("%s %s = %q\n", variableName, customTypeName, val))
		}
		str.WriteString(")\n")
	}

	comment := fmt.Sprintf("%s contains acceptable values for %q column of %q table.\n", customTypeName, i.dbColumnName, i.dbTableName)
	str.WriteString(p.addComment(comment))
	str.WriteString(fmt.Sprintf("type %s string\n", customTypeName))

	p.customTypes[customTypeName] = str.String()

	return customTypeName
}

func (p *parser) createCustomIntType(min, max int64) string {
	var customTypeName string
	intSize := "int"
	if min > 0 && max > 0 {
		intSize = p.intSizeFromMaxValue(max)
		customTypeName = fmt.Sprintf("IntFrom%dTo%d", min, max)
	} else if min > 0 {
		intSize = "int64"
		customTypeName = fmt.Sprintf("IntFrom%d", min)
	} else if max > 0 {
		intSize = p.intSizeFromMaxValue(max)
		customTypeName = fmt.Sprintf("IntFrom0To%d", max)
	}

	if _, ok := p.customTypes[customTypeName]; !ok {
		var str strings.Builder

		str.WriteString(fmt.Sprintf("// %s custom integer type of %s size.\n", customTypeName, intSize))
		str.WriteString(fmt.Sprintf("type %s %s\n", customTypeName, intSize))

		p.customTypes[customTypeName] = str.String()
	}

	return customTypeName
}

func (p *parser) intSizeFromMaxValue(val int64) string {
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
