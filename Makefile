.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run: ## runs generator
	go run -v cmd/main.go

prepare: install-prerequisites download generate-xml-struct ## run make step

download: download-schema download-description ## download pre-requisites

download-schema: ## download ovn schema
	curl -s -L https://raw.githubusercontent.com/openvswitch/ovs/master/ovn/ovn-nb.ovsschema -o files/ovn-nb.json

download-description: ## download ovn description
	curl -s -L https://raw.githubusercontent.com/openvswitch/ovs/master/ovn/ovn-nb.xml -o files/ovn-nb.xml

sanitize-xml: ## sanitize xml
	sed -i 's%<code>%"%g' files/ovn-nb.xml
	sed -i 's%</code>%"%g' files/ovn-nb.xml
	sed -i 's%<var>%"%g' files/ovn-nb.xml
	sed -i 's%</var>%"%g' files/ovn-nb.xml
	sed -i 's%<em>%"%g' files/ovn-nb.xml
	sed -i 's%</em>%"%g' files/ovn-nb.xml
	sed -i 's/<ref\(.*\)=\(.*\)\/>/\2/g' files/ovn-nb.xml

generate-xml-struct: sanitize-xml ## generate go struct from xml file
	echo "package xmlschema" > app/xmlschema/xml-schema.go
	echo 'import "encoding/xml"' >> app/xmlschema/xml-schema.go
	zek -e files/ovn-nb.xml >> app/xmlschema/xml-schema.go
	gofmt -w app/xmlschema/xml-schema.go

install-prerequisites: ## installs required binaries
	go get -v github.com/miku/zek/cmd/...
	go get github.com/eidolon/wordwrap
