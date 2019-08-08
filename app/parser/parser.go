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
	typeMatching = map[string]string{
		stringType: "string",
		intType:    "int",
		boolType:   "bool",
	}
)

type parser struct {
	DBSchema    map[string]interface{}
	XMLDoc      *xmlschema.Database
	RawXMLDoc   *xmlschema.Database
	structures  map[string]string
	customTypes map[string]string
}

type info struct {
	dbTableName  string
	dbColumnName string
	structName   string
	fieldName    string
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

	rawXMLDoc, err := xmlschema.NewXML("files/ovn-nb-with-keys.xml")
	if err != nil {
		return nil, fmt.Errorf("new xml: %v", err)
	}

	return &parser{DBSchema: dbSchema, XMLDoc: xmlDoc, RawXMLDoc: rawXMLDoc}, nil
}

func (p *parser) Parse() {
	rawTables := p.DBSchema[tablesKey].(map[string]interface{})
	p.structures = make(map[string]string, len(rawTables))
	p.customTypes = make(map[string]string)
	i := &info{}

	for rawTableName, data := range rawTables {
		i.dbTableName = rawTableName
		i.structName = sanitizeName(i.dbTableName)
		var str strings.Builder

		str.WriteString(p.addStructComment(i))
		str.WriteString(fmt.Sprintf("type %s struct {\n", i.structName))

		table, ok := data.(map[string]interface{})
		if !ok {
			log.Warnf("Can't parse table %q", i.dbTableName)
			continue
		}

		columns, ok := table[columnsKey].(map[string]interface{})
		if !ok {
			log.Warnf("Can't parse columns in table %q", i.dbTableName)
			continue
		}

		for rawColumnName, data := range columns {
			i.dbColumnName = rawColumnName
			i.fieldName = sanitizeName(i.dbColumnName)
			str.WriteString(p.addFieldComment(i))
			str.WriteString(fmt.Sprintf("%s ", i.fieldName))
			str.WriteString(fmt.Sprintf("%s", p.parseColumn(i, data)))
		}
		str.WriteString("}\n")
		p.structures[i.structName] = str.String()
	}
}

func (p *parser) Structures() map[string]string {
	return p.structures
}

func (p *parser) CustomTypes() map[string]string {
	return p.customTypes
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
