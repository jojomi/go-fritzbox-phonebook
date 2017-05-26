package phonebook

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasicXMLGeneration(t *testing.T) {
	p := NewPhonebookContainer().AddPhonebook(NewPhonebook("testbook"))
	str := p.ToXMLString()
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<phonebooks><phonebook name="testbook"></phonebook></phonebooks>`

	assert.Equal(t, expected, str, "they should be equal")
}

func TestBasicXMLGenerationIndented(t *testing.T) {
	p := NewPhonebookContainer().AddPhonebook(NewPhonebook("testbook"))
	str := p.ToXMLStringIndented()
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<phonebooks>
    <phonebook name="testbook"></phonebook>
</phonebooks>`

	assert.Equal(t, expected, str, "they should be equal")
}

func TestPhonebookEmptyName(t *testing.T) {
	p := NewPhonebookContainer().AddPhonebook(NewPhonebook(""))
	str := p.ToXMLStringIndented()
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<phonebooks>
    <phonebook></phonebook>
</phonebooks>`

	assert.Equal(t, expected, str, "they should be equal")
}

func TestAddContact(t *testing.T) {
	p := NewPhonebook("testbook")
	pc := NewPhonebookContainer().AddPhonebook(p)

	c := NewContact("Gerd Müller")
	l, err := time.LoadLocation("UTC")
	assert.Nil(t, err)
	c.ModTime = time.Date(2015, 01, 02, 03, 04, 05, 0, l).Unix()

	n := NewNumber("+491237612315674")
	n.Quickdial = 3
	c.AddNumber(n)

	n2 := NewNumber("01785451213689")
	n2.Type = TypeWork
	c.AddNumber(n2)

	p.AddContact(c)

	str := pc.ToXMLStringIndented()
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<phonebooks>
    <phonebook name="testbook">
        <contact>
            <Person>
                <realName>Gerd Müller</realName>
            </Person>
            <category>0</category>
            <telephony nid="2">
                <number type="home" id="0" prio="0" quickdial="3">+491237612315674</number>
                <number type="work" id="1" prio="0">01785451213689</number>
            </telephony>
            <features doorphone="0"></features>
            <services></services>
            <mod_time>1420167845</mod_time>
            <uniqueid>1</uniqueid>
        </contact>
    </phonebook>
</phonebooks>`

	assert.Equal(t, expected, str, "they should be equal")
}
