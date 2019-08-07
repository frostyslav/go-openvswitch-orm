package parser

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// column keys
const (
	typeKey string = "type"
)

type column struct {
	tableName string
	name      string
	rawName   string
	rawData   interface{}
}

func (c *column) parse() string {
	column, ok := c.rawData.(map[string]interface{})
	if !ok {
		log.Warnf("Can't parse column %q", c.rawName)
		return ""
	}

	if t, ok := column[typeKey]; ok {
		colType := &columnType{rawData: t, columnName: c.name, tableName: c.tableName}
		return fmt.Sprintf("%s\n", colType.parse())
	}

	return ""
}
