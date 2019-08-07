package main

import (
	"flag"
	"go/format"
	"os"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/frostyslav/gopenvswitch-db/app/parser"
)

const (
	defaultJSONSchemaLocation string = "files/ovn-nb.json"
	defaultXMLDocLocation     string = "files/ovn-nb.xml"
	defaultOutputLocation     string = "output/generated.go"
)

var (
	jsonSchemaFile = flag.String("json-schema", defaultJSONSchemaLocation, "json schema location")
	xmlDocFile     = flag.String("xml-doc", defaultXMLDocLocation, "xml doc location")
	outputFile     = flag.String("output", defaultOutputLocation, "output location")
)

func main() {
	flag.Parse()

	p, err := parser.New(*jsonSchemaFile, *xmlDocFile)
	if err != nil {
		log.Errorf("new parser: %v", err)
		return
	}

	tables, customTypes := p.Parse()

	var formattedOutput strings.Builder
	formattedOutput.WriteString("package ovnnb\n")

	var customTypeKeys []string
	for k := range customTypes {
		customTypeKeys = append(customTypeKeys, k)
	}
	sort.Strings(customTypeKeys)

	for _, k := range customTypeKeys {
		formattedOutput.WriteString(customTypes[k])
	}

	var keys []string
	for k := range tables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		formattedOutput.WriteString(tables[k])
	}

	output, err := format.Source([]byte(formattedOutput.String()))
	if err != nil {
		log.Errorf("format output: %v", err)
		return
	}

	file, err := os.Create(*outputFile)
	if err != nil {
		log.Errorf("create file: %v", err)
		return
	}

	defer file.Close()

	file.Write(output)
}
