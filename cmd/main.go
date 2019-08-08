package main

import (
	"flag"
	"fmt"
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
	defaultOutputLocation     string = "output"
)

var (
	jsonSchemaFile = flag.String("json-schema", defaultJSONSchemaLocation, "json schema location")
	xmlDocFile     = flag.String("xml-doc", defaultXMLDocLocation, "xml doc location")
	outputLocation = flag.String("output", defaultOutputLocation, "output location")
)

func main() {
	flag.Parse()

	p, err := parser.New(*jsonSchemaFile, *xmlDocFile)
	if err != nil {
		log.Errorf("new parser: %v", err)
		return
	}

	p.Parse()

	err = writeGeneratedCodeToFile("ovnnb", *outputLocation + "/types.go", mapToString(p.CustomTypes()))
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	err = writeGeneratedCodeToFile("ovnnb", *outputLocation + "/structs.go", mapToString(p.Structures()))
	if err != nil {
		log.Errorf(err.Error())
		return
	}
}

func writeGeneratedCodeToFile(packageName, fileName, contents string) error {
	var str strings.Builder

	str.WriteString("package " + packageName + "\n")
	str.WriteString(contents)

	gofmted, err := format.Source([]byte(str.String()))
	if err != nil {
		return fmt.Errorf("run go fmt: %v", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}

	defer file.Close()

	_, err = file.Write(gofmted)
	if err != nil {
		return fmt.Errorf("write to file: %v", err)
	}

	return nil
}

func mapToString(m map[string]string) string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var s strings.Builder
	for _, k := range keys {
		s.WriteString(m[k])
	}

	return s.String()
}
