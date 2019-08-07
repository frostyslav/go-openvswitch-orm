package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/frostyslav/gopenvswitch-db/app/xmlschema"
)

// database keys
const (
	tablesKey string = "tables"
)

// table keys
const (
	columnsKey string = "columns"
)

// database types
const (
	stringType string = "string"
	intType    string = "integer"
	boolType   string = "boolean"
	uuidType   string = "uuid"
)

var (
	customTypes = map[string]string{}
)

var (
	typeMatching = map[string]string{
		stringType: "string",
		intType:    "int",
		boolType:   "bool",
	}
)

type parser struct {
	DBSchema map[string]interface{}
	XMLDoc   *xmlschema.Database
}

func New(jsonSchemaFile, xmlDocFile string) (*parser, error) {
	ovnNB, err := ioutil.ReadFile(jsonSchemaFile)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	var dbSchema map[string]interface{}
	err = json.Unmarshal([]byte(ovnNB), &dbSchema)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	xmlDoc, err := xmlschema.NewXML(xmlDocFile)
	if err != nil {
		return nil, fmt.Errorf("new xml: %v", err)
	}

	return &parser{DBSchema: dbSchema, XMLDoc: xmlDoc}, nil
}

func (p *parser) Parse() (map[string]string, map[string]string) {
	rawTables := p.DBSchema[tablesKey].(map[string]interface{})
	tables := make(map[string]string, len(rawTables))

	for rawTableName, data := range rawTables {
		var str strings.Builder

		tableName := sanitizeName(rawTableName)

		str.WriteString(p.addTableComment(rawTableName))
		str.WriteString(fmt.Sprintf("type %s struct {\n", tableName))

		table, ok := data.(map[string]interface{})
		if !ok {
			log.Warnf("Can't parse table %q", rawTableName)
			continue
		}

		columns, ok := table[columnsKey].(map[string]interface{})
		if !ok {
			log.Warnf("Can't parse columns in table %q", rawTableName)
			continue
		}

		for rawColumnName, data := range columns {
			columnName := sanitizeName(rawColumnName)
			str.WriteString(p.addColumnComment(rawTableName, rawColumnName))
			str.WriteString(fmt.Sprintf("%s ", columnName))

			col := column{tableName: tableName, name: columnName, rawName: rawColumnName, rawData: data}
			str.WriteString(fmt.Sprintf("%s", col.parse()))
		}
		str.WriteString("}\n")
		tables[tableName] = str.String()
	}

	return tables, customTypes
}

func (p *parser) addTableComment(rawTableName string) string {
	tableName := sanitizeName(rawTableName)

	return fmt.Sprintf("// %s %s\n", tableName, p.XMLDoc.TableDescription(rawTableName))
}

func (p *parser) addColumnComment(rawTableName, rawColumnName string) string {
	columnName := sanitizeName(rawColumnName)

	return fmt.Sprintf("// %s %s\n", columnName, p.XMLDoc.ColumnDescription(rawTableName, rawColumnName))
}

func minValue(data map[string]interface{}) int64 {
	if min, ok := data[minKey]; ok {
		val, ok := min.(float64)
		if !ok {
			return 0
		}
		return int64(val)
	}

	return 0
}

func maxValue(data map[string]interface{}) int64 {
	if max, ok := data[maxKey]; ok {
		if max == "unlimited" {
			return math.MaxInt64
		}
		val, ok := max.(float64)
		if ok {
			return int64(val)
		}
	}

	return 0
}

func minIntegerValue(data map[string]interface{}) int64 {
	if min, ok := data[minIntegerKey]; ok {
		val, ok := min.(float64)
		if !ok {
			return -1
		}
		return int64(val)
	}

	return -1
}

func maxIntegerValue(data map[string]interface{}) int64 {
	if max, ok := data[maxIntegerKey]; ok {
		val, ok := max.(float64)
		if ok {
			return int64(val)
		}
	}

	return 0
}

func sanitizeName(name string) string {
	splittedName := strings.FieldsFunc(name, split)
	var sanitizedName string
	for _, splitted := range splittedName {
		sanitizedName = fmt.Sprintf("%s%s", sanitizedName, strings.Title(splitted))
	}

	return strings.TrimSpace(sanitizedName)
}

func split(r rune) bool {
	return r == '-' || r == '_'
}