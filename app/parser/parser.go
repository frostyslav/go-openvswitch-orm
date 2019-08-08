package parser

import (
	"encoding/json"
	"fmt"
	"go/types"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/frostyslav/gopenvswitch-db/app/sanitize"
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
		stringType: types.Typ[types.String].String(),
		intType:    types.Typ[types.Int].String(),
		boolType:   types.Typ[types.Bool].String(),
	}
)

type parser struct {
	DBSchema       map[string]interface{}
	XMLDoc         *xmlschema.Database
	ModifiedXMLDoc *xmlschema.Database
	structures     map[string]string
	customTypes    map[string]string
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

	modifiedXMLDoc, err := xmlschema.NewXML("files/ovn-nb-with-keys.xml")
	if err != nil {
		return nil, fmt.Errorf("new xml: %v", err)
	}

	return &parser{DBSchema: dbSchema, XMLDoc: xmlDoc, ModifiedXMLDoc: modifiedXMLDoc}, nil
}

func (p *parser) Parse() {
	rawTables := p.DBSchema[tablesKey].(map[string]interface{})
	p.structures = make(map[string]string, len(rawTables))
	p.customTypes = make(map[string]string)
	i := &info{}

	for rawTableName, data := range rawTables {
		i.dbTableName = rawTableName
		i.structName = sanitize.Name(i.dbTableName)
		var str strings.Builder

		str.WriteString(p.structComment(i))
		str.WriteString("type " + i.structName + " struct {\n")

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
			i.fieldName = sanitize.Name(i.dbColumnName)
			str.WriteString(p.fieldComment(i))
			str.WriteString(i.fieldName + " " + p.parseColumn(i, data))
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
