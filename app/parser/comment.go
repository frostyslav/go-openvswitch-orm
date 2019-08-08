package parser

import (
	"github.com/frostyslav/gopenvswitch-db/app/wrapper"
)

var (
	wrapLimit = 80
)

func (p *parser) structComment(i *info) string {
	s := i.structName + " " + p.XMLDoc.TableDescription(i.dbTableName) + "\n"
	return wrapper.WrapAsComment(s, wrapLimit)
}

func (p *parser) fieldComment(i *info) string {
	s := i.fieldName + " " + p.XMLDoc.ColumnDescription(i.dbTableName, i.dbColumnName) + "\n"
	return wrapper.WrapAsComment(s, wrapLimit)
}

func (p *parser) variableComment(i *info, variableName, dbKeyName string) string {
	s := variableName + " " + p.ModifiedXMLDoc.KeyDescription(i.dbTableName, i.dbColumnName, dbKeyName) + "\n"
	return wrapper.WrapAsComment(s, wrapLimit)
}

func (p *parser) customComment(s string) string {
	return wrapper.WrapAsComment(s + "\n", wrapLimit)
}
