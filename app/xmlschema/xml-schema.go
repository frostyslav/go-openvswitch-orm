package xmlschema

import "encoding/xml"

// Database was generated 2019-08-07 22:45:11 by rostyslav_fridman on workstation.
type Database struct {
	XMLName xml.Name `xml:"database"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
	Title   string   `xml:"title,attr"`
	P       []string `xml:"p"`  // This database is the inte...
	H2      string   `xml:"h2"` // External IDs
	Dl      struct {
		Text string `xml:",chardata"`
		Dt   string `xml:"dt"` // "external_ids": map of st...
		Dd   string `xml:"dd"` // Key-value pairs for use b...
	} `xml:"dl"`
	Table []struct {
		Text  string   `xml:",chardata"` // SSL configuration for ovn...
		Name  string   `xml:"name,attr"`
		Title string   `xml:"title,attr"`
		P     []string `xml:"p"` // Northbound configuration ...
		Group []struct {
			Text   string `xml:",chardata"` // These columns allow a cli...
			Title  string `xml:"title,attr"`
			Column []struct {
				Text string   `xml:",chardata"` // Sequence number for clien...
				Name string   `xml:"name,attr"`
				Key  string   `xml:"key,attr"`
				Type string   `xml:"type,attr"`
				P    []string `xml:"p"` // To exclude some addresses...
				Ul   struct {
					Text string `xml:",chardata"`
					Li   []struct {
						Text string   `xml:",chardata"` // "192.168.0.2 192.168.0.10...
						P    []string `xml:"p"`         // If any IPv4 address is gi...
						Key  string   `xml:"key"`       // slaac
					} `xml:"li"`
				} `xml:"ul"`
				Dl struct {
					Text string   `xml:",chardata"`
					Dt   []string `xml:"dt"` // (empty string), "router",...
					Dd   []struct {
						Text string   `xml:",chardata"` // A VM (or VIF) interface.,...
						P    []string `xml:"p"`         // Represents a logical port...
						Ul   struct {
							Text string   `xml:",chardata"`
							Li   []string `xml:"li"` // VMs connected to SR-IOV n...
						} `xml:"ul"`
						Dl struct {
							Text string   `xml:",chardata"`
							Dt   []string `xml:"dt"` // "80:fa:5b:06:72:b7", "80:...
							Dd   []string `xml:"dd"` // This indicates that the l...
						} `xml:"dl"`
					} `xml:"dd"`
				} `xml:"dl"`
			} `xml:"column"`
			Group []struct {
				Text   string `xml:",chardata"`
				Title  string `xml:"title,attr"`
				P      string `xml:"p"` // These options apply when ...
				Column []struct {
					Text string   `xml:",chardata"` // BFD option "min-rx" value...
					Name string   `xml:"name,attr"`
					Key  string   `xml:"key,attr"`
					Type string   `xml:"type,attr"`
					P    []string `xml:"p"` // This is used to send grat...
					Dl   struct {
						Text string   `xml:",chardata"`
						Dt   []string `xml:"dt"` // "router", "Ethernet addre...
						Dd   []struct {
							Text string   `xml:",chardata"`
							P    []string `xml:"p"` // Gratuitous ARPs will be s...
						} `xml:"dd"`
					} `xml:"dl"`
				} `xml:"column"`
			} `xml:"group"`
			P  []string `xml:"p"` // These columns provide nam...
			Ul struct {
				Text string   `xml:",chardata"`
				Li   []string `xml:"li"` // To attach a logical switc...
			} `xml:"ul"`
		} `xml:"group"`
		Column []struct {
			Text string `xml:",chardata"` // Load balance a virtual ip...
			Name string `xml:"name,attr"`
			P    []struct {
				Text string `xml:",chardata"` // The logical ports connect...
				B    string `xml:"b"`         // Example:
			} `xml:"p"`
			Ul struct {
				Text string `xml:",chardata"`
				Li   []struct {
					Text string `xml:",chardata"` // : Used to implement filte...
					Key  string `xml:"key"`       // from-lport, to-lport, all...
				} `xml:"li"`
			} `xml:"ul"`
		} `xml:"column"`
		Pre string `xml:"pre"` // ovn-nbctl create Address_...
	} `xml:"table"`
}
