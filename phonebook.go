package phonebook

import "encoding/xml"

type PhonebookContainer struct {
	Phonebooks []*Phonebook `xml:"phonebook"`
	XMLName    struct{}     `xml:"phonebooks"`
}

type Phonebook struct {
	Contacts []*Contact `xml:"contact"`
	Name     string     `xml:"name,attr,omitempty"`
}

type Contact struct {
	Person    *Person
	Category  int        `xml:"category"` // seems yet unused, always 0
	Telephony *Telephony `xml:"telephony"`
	Features  Features   `xml:"features"`
	Services  struct{}   `xml:"services"`
	ModTime   int64      `xml:"mod_time"`
	UniqueID  int        `xml:"uniqueid"`
}

type Person struct {
	RealName string `xml:"realName"`
}

type Features struct {
	Doorphone int `xml:"doorphone,attr"`
}

type Telephony struct {
	Count   int       `xml:"nid,attr"`
	Numbers []*Number `xml:"number"`
}

type Number struct {
	Type      Type     `xml:"type,attr"`
	ID        int      `xml:"id,attr"`
	Vanity    string   `xml:"vanity,attr,omitempty"`
	Priority  Priority `xml:"prio,attr"`
	Quickdial int      `xml:"quickdial,attr,omitempty"`
	Number    string   `xml:",chardata"`
}

type Type string

const (
	TypeHome   Type = "home"
	TypeMobile      = "mobile"
	TypeWork        = "work"
)

type Priority int

const (
	PriorityNormal Priority = 0
	PriorityHigh            = 1
)

func NewPhonebookContainer() *PhonebookContainer {
	p := &PhonebookContainer{}
	return p
}

func (p *PhonebookContainer) AddPhonebook(pb *Phonebook) *PhonebookContainer {
	p.Phonebooks = append(p.Phonebooks, pb)
	return p
}

// ToXMLStringIndented returns the XML representation as a string with indention to make it more readable
func (p PhonebookContainer) ToXMLStringIndented() string {
	return p.toXML(true)
}

// ToXMLString returns the XML representation as a string
func (p PhonebookContainer) ToXMLString() string {
	return p.toXML(false)
}

func (p PhonebookContainer) toXML(indented bool) string {
	var (
		myString []byte
		err      error
	)

	for _, pb := range p.Phonebooks {
		pb.optimize()
		for _, c := range pb.Contacts {
			c.optimize()
		}
	}

	if indented {
		myString, err = xml.MarshalIndent(p, "", "    ")
	} else {
		myString, err = xml.Marshal(p)
	}
	if err != nil {
		panic(err)
	}
	return xml.Header + string(myString)
}

// NewPhonebook creates a new Phonebook
func NewPhonebook(name string) *Phonebook {
	p := &Phonebook{
		Name: name,
	}
	return p
}

// AddContact adds a contact to the phonebook
func (p *Phonebook) AddContact(contact *Contact) {
	p.Contacts = append(p.Contacts, contact)
}

func (p *Phonebook) optimize() {
	for i, c := range p.Contacts {
		if c.UniqueID == 0 {
			c.UniqueID = i + 1
		}
	}
}

func NewContact(name string) *Contact {
	c := &Contact{
		Person: &Person{
			RealName: name,
		},
	}
	return c
}

func (c *Contact) AddNumber(number *Number) {
	if c.Telephony == nil {
		c.Telephony = &Telephony{}
	}
	c.Telephony.Numbers = append(c.Telephony.Numbers, number)
}

func (c *Contact) optimize() {
	if c.Telephony != nil {
		c.Telephony.Count = len(c.Telephony.Numbers)

		for i, number := range c.Telephony.Numbers {
			// set default
			if number.Type == "" {
				number.Type = TypeHome
			}
			number.ID = i
		}
	}
}

func NewNumber(value string) *Number {
	n := &Number{
		Number: value,
	}
	return n
}
