package parser

import (
	"fmt"

	"github.com/frostyslav/gopenvswitch-db/app/wrapper"
)

func (p *parser) addStructComment(i *info) string {
	s := fmt.Sprintf("%s %s\n", sanitizeName(i.dbTableName), p.XMLDoc.TableDescription(i.dbTableName))
	return wrapper.WrapAsComment(s)
}

func (p *parser) addFieldComment(i *info) string {
	s := fmt.Sprintf("%s %s\n", sanitizeName(i.dbColumnName), p.XMLDoc.ColumnDescription(i.dbTableName, i.dbColumnName))
	return wrapper.WrapAsComment(s)
}

func (p *parser) addVariableComment(i *info, variableName, dbKeyName string) string {
	s := fmt.Sprintf("%s %s\n", variableName, p.RawXMLDoc.KeyDescription(i.dbTableName, i.dbColumnName, dbKeyName))
	return wrapper.WrapAsComment(s)
}

func (p *parser) addComment(s string) string {
	return wrapper.WrapAsComment(s)
}
