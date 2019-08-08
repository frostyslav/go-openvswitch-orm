package xmlschema

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func NewXML(filePath string) (*Database, error) {
	// Open our xmlFile
	xmlFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	result := Database{}
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	err = xml.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, fmt.Errorf("xml unmarshal: %v", err)
	}

	return &result, nil
}

func (d *Database) TableDescription(tableName string) string {
	var text string

	for _, t := range d.Table {
		if t.Name == tableName {
			text = text + t.Text
			for _, g := range t.Group {
				for _, p := range g.P {
					text = text + p
				}
			}
			for _, p := range t.P {
				text = text + p
			}
		}
	}

	text = strings.TrimSpace(text)
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")

	return text
}

func (d *Database) ColumnDescription(tableName, columnName string) string {
	var text string
	for _, t := range d.Table {
		if t.Name == tableName {
			for _, c := range t.Column {
				if c.Name == columnName {
					text = text + c.Text
					for _, p := range c.P {
						text = text + p.Text
					}
				}
			}
			for _, g := range t.Group {
				for _, c := range g.Column {
					if c.Name == columnName {
						text = text + c.Text
					}
					for _, p := range c.P {
						if c.Name == columnName {
							text = text + p
						}
					}
				}
			}
		}
	}

	text = strings.TrimSpace(text)
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")

	return text
}

func (d *Database) KeyDescription(tableName, columnName, keyName string) string {
	var text string
	for _, t := range d.Table {
		if t.Name == tableName {
			for _, c := range t.Column {
				if c.Name == columnName {
					for _, l := range c.Ul.Li {
						if l.Key == keyName {
							text = text + l.Text
						}
					}
				}
			}
		}
	}

	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ": ", " ")
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")

	return text
}
