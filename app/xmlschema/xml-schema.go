package xmlschema

import "encoding/xml"

// Database was generated 2019-08-07 18:42:22 by rostyslav_fridman on workstation.
type Database struct {
	XMLName xml.Name `xml:"database"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
	Title   string   `xml:"title,attr"`
	P       []struct {
		Text string `xml:",chardata"` // This database is the inte...
		Ref  struct {
			Text string `xml:",chardata"`
			Db   string `xml:"db,attr"`
		} `xml:"ref"`
	} `xml:"p"`
	H2 string `xml:"h2"` // External IDs
	Dl struct {
		Text string `xml:",chardata"`
		Dt   string `xml:"dt"` // "external_ids": map of st...
		Dd   string `xml:"dd"` // Key-value pairs for use b...
	} `xml:"dl"`
	Table []struct {
		Text  string `xml:",chardata"` // SSL configuration for ovn...
		Name  string `xml:"name,attr"`
		Title string `xml:"title,attr"`
		P     []struct {
			Text string `xml:",chardata"` // Northbound configuration ...
			Ref  []struct {
				Text   string `xml:",chardata"`
				Column string `xml:"column,attr"`
				Table  string `xml:"table,attr"`
				Db     string `xml:"db,attr"`
			} `xml:"ref"`
		} `xml:"p"`
		Group []struct {
			Text   string `xml:",chardata"` // These columns allow a cli...
			Title  string `xml:"title,attr"`
			Column []struct {
				Text string `xml:",chardata"` // Sequence number for clien...
				Name string `xml:"name,attr"`
				Key  string `xml:"key,attr"`
				Type string `xml:"type,attr"`
				Ref  struct {
					Text   string `xml:",chardata"`
					Column string `xml:"column,attr"`
					Key    string `xml:"key,attr"`
				} `xml:"ref"`
				P []struct {
					Text string `xml:",chardata"` // To exclude some addresses...
					Ref  []struct {
						Text   string `xml:",chardata"`
						Key    string `xml:"key,attr"`
						Table  string `xml:"table,attr"`
						Column string `xml:"column,attr"`
						Db     string `xml:"db,attr"`
					} `xml:"ref"`
				} `xml:"p"`
				Ul struct {
					Text string `xml:",chardata"`
					Li   []struct {
						Text string   `xml:",chardata"` // "192.168.0.2 192.168.0.10...
						P    []string `xml:"p"`         // If any IPv4 address is gi...
					} `xml:"li"`
				} `xml:"ul"`
				Dl struct {
					Text string   `xml:",chardata"`
					Dt   []string `xml:"dt"` // (empty string), "router",...
					Dd   []struct {
						Text string `xml:",chardata"` // A VM (or VIF) interface.,...
						P    []struct {
							Text string `xml:",chardata"` // Represents a logical port...
							Ref  struct {
								Text   string `xml:",chardata"`
								Column string `xml:"column,attr"`
								Table  string `xml:"table,attr"`
								Db     string `xml:"db,attr"`
								Key    string `xml:"key,attr"`
							} `xml:"ref"`
						} `xml:"p"`
						Ul struct {
							Text string   `xml:",chardata"`
							Li   []string `xml:"li"` // VMs connected to SR-IOV n...
						} `xml:"ul"`
						Dl struct {
							Text string   `xml:",chardata"`
							Dt   []string `xml:"dt"` // "80:fa:5b:06:72:b7", "80:...
							Dd   []string `xml:"dd"` // This indicates that the l...
						} `xml:"dl"`
						Ref []struct {
							Text   string `xml:",chardata"`
							Column string `xml:"column,attr"`
							Table  string `xml:"table,attr"`
							Key    string `xml:"key,attr"`
						} `xml:"ref"`
					} `xml:"dd"`
				} `xml:"dl"`
			} `xml:"column"`
			Group []struct {
				Text   string `xml:",chardata"`
				Title  string `xml:"title,attr"`
				P      string `xml:"p"` // These options apply when ...
				Column []struct {
					Text string `xml:",chardata"` // BFD option "min-rx" value...
					Name string `xml:"name,attr"`
					Key  string `xml:"key,attr"`
					Type string `xml:"type,attr"`
					Ref  struct {
						Text  string `xml:",chardata"`
						Table string `xml:"table,attr"`
					} `xml:"ref"`
					P  []string `xml:"p"` // This is used to send grat...
					Dl struct {
						Text string   `xml:",chardata"`
						Dt   []string `xml:"dt"` // "router", "Ethernet addre...
						Dd   []struct {
							Text string `xml:",chardata"`
							P    []struct {
								Text string `xml:",chardata"` // Gratuitous ARPs will be s...
								Ref  struct {
									Text   string `xml:",chardata"`
									Column string `xml:"column,attr"`
									Key    string `xml:"key,attr"`
								} `xml:"ref"`
							} `xml:"p"`
						} `xml:"dd"`
					} `xml:"dl"`
				} `xml:"column"`
			} `xml:"group"`
			P []struct {
				Text string `xml:",chardata"` // These columns provide nam...
				Ref  []struct {
					Text   string `xml:",chardata"`
					Column string `xml:"column,attr"`
					Key    string `xml:"key,attr"`
					Table  string `xml:"table,attr"`
				} `xml:"ref"`
			} `xml:"p"`
			Ul struct {
				Text string `xml:",chardata"`
				Li   []struct {
					Text string `xml:",chardata"` // To attach a logical switc...
					Ref  []struct {
						Text   string `xml:",chardata"`
						Table  string `xml:"table,attr"`
						Column string `xml:"column,attr"`
					} `xml:"ref"`
				} `xml:"li"`
			} `xml:"ul"`
		} `xml:"group"`
		Column []struct {
			Text string `xml:",chardata"` // Load balance a virtual ip...
			Name string `xml:"name,attr"`
			P    []struct {
				Text string `xml:",chardata"` // The logical ports connect...
				Ref  struct {
					Text   string `xml:",chardata"`
					Column string `xml:"column,attr"`
					Table  string `xml:"table,attr"`
					Db     string `xml:"db,attr"`
					Colun  string `xml:"colun,attr"`
				} `xml:"ref"`
				B string `xml:"b"` // Example:
			} `xml:"p"`
			Ul struct {
				Text string   `xml:",chardata"`
				Li   []string `xml:"li"` // "from-lport": Used to imp...
			} `xml:"ul"`
		} `xml:"column"`
		Pre string `xml:"pre"` // ovn-nbctl create Address_...
	} `xml:"table"`
}
