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

	p.Parse()

	var str strings.Builder
	str.WriteString("package ovnnb\n")
	str.WriteString(sortedMapOutput(p.CustomTypes()))
	str.WriteString(sortedMapOutput(p.Structures()))

	gofmted, err := format.Source([]byte(str.String()))
	if err != nil {
		log.Errorf("run go fmt: %v", err)
		return
	}

	file, err := os.Create(*outputFile)
	if err != nil {
		log.Errorf("create file: %v", err)
		return
	}

	defer file.Close()

	file.Write(gofmted)
}

func sortedMapOutput(m map[string]string) string {
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