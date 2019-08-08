package parser

import (
	"fmt"
	"go/types"
	"math"
	"strings"

	"github.com/frostyslav/gopenvswitch-db/app/sanitize"
)

func (p *parser) createNamedStringType(name, comment string) string {
	var str strings.Builder

	str.WriteString(p.customComment(name + " " + comment))
	str.WriteString(p.typeDeclaration(name, types.Typ[types.String].String()))

	p.customTypes[name] = str.String()

	return name
}

func (p *parser) createAcceptableValuesStringType(i *info, values []string) string {
	var str strings.Builder

	typeName := i.structName + i.fieldName
	if len(values) > 0 {
		str.WriteString("const (\n")
		for _, v := range values {
			val := sanitize.Name(v)
			varName := typeName + val
			str.WriteString(p.variableComment(i, varName, v))
			str.WriteString(p.variableDeclaration(varName, typeName, v))
		}
		str.WriteString(")\n")
	}

	comment := typeName + " contains acceptable values for '" + i.dbColumnName + "' column of '" + i.dbTableName + "' table."
	str.WriteString(p.customComment(comment))
	str.WriteString(p.typeDeclaration(typeName, types.Typ[types.String].String()))

	p.customTypes[typeName] = str.String()

	return typeName
}

func (p *parser) createCustomIntType(min, max int64) string {
	var customTypeName string
	intSize := types.Typ[types.Int].String()
	if min > 0 && max > 0 {
		intSize = p.intSizeFromMaxValue(max)
		customTypeName = fmt.Sprintf("IntFrom%dTo%d", min, max)
	} else if min > 0 {
		intSize = types.Typ[types.Int64].String()
		customTypeName = fmt.Sprintf("IntFrom%d", min)
	} else if max > 0 {
		intSize = p.intSizeFromMaxValue(max)
		customTypeName = fmt.Sprintf("IntFrom0To%d", max)
	}

	if _, ok := p.customTypes[customTypeName]; !ok {
		var str strings.Builder

		str.WriteString(p.customComment(customTypeName + " custom integer type of " + intSize + " size."))
		str.WriteString(p.typeDeclaration(customTypeName, intSize))

		p.customTypes[customTypeName] = str.String()
	}

	return customTypeName
}

func (p *parser) intSizeFromMaxValue(val int64) string {
	if val < math.MaxInt8 {
		return types.Typ[types.Int8].String()
	} else if val < math.MaxInt16 {
		return types.Typ[types.Int16].String()
	} else if val < math.MaxInt32 {
		return types.Typ[types.Int32].String()
	} else if val < math.MaxInt64 {
		return types.Typ[types.Int64].String()
	} else {
		return types.Typ[types.Int].String()
	}
}

func (p *parser) typeDeclaration(name, typ string) string {
	return "type " + name + " " + typ + "\n"
}

func (p *parser) variableDeclaration(name, typ, val string) string {
	return name + " " + typ + " = \"" + val + "\"\n"
}
