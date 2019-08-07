package parser

import (
	"fmt"
	"github.com/frostyslav/gopenvswitch-db/app/xmlschema"

	log "github.com/sirupsen/logrus"
)

// column keys
const (
	typeKey string = "type"
)

type column struct {
	tableName string
	name      string
	rawData   interface{}
	xmlDoc    *xmlschema.Database
}

func (c *column) parse() string {
	column, ok := c.rawData.(map[string]interface{})
	if !ok {
		log.Warnf("Can't parse column %q", c.name)
		return ""
	}

	if t, ok := column[typeKey]; ok {
		colType := &columnType{rawData: t, columnName: c.name, tableName: c.tableName, xmlDoc: c.xmlDoc}
		return fmt.Sprintf("%s\n", colType.parse())
	}

	return ""
}
