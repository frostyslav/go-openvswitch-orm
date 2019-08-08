package parser

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// column keys
const (
	typeKey string = "type"
)

func (p *parser) parseColumn(i *info, data interface{}) string {
	column, ok := data.(map[string]interface{})
	if !ok {
		log.Warnf("Can't parse column %q", i.dbColumnName)
		return ""
	}

	if typeData, ok := column[typeKey]; ok {
		return fmt.Sprintf("%s\n", p.parseColumnType(i, typeData))
	}

	return ""
}
