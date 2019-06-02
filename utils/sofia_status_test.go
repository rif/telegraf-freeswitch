package utils

import "testing"

var (
	sofiaData = `
<?xml version="1.0" encoding="ISO-8859-1"?>
<profiles>
<profile>
<name>dot142</name>
<type>profile</type>
<data>sip:mod_sofia@173.244.183.142:5060</data>
<state>RUNNING (0)</state>
</profile>
<gateway>
<name>example.com</name>
<type>gateway</type>
<data>sip:joeuser@example.com</data>
<state>NOREG</state>
</gateway>

<profile>
<name>dot139</name>
<type>profile</type>
<data>sip:mod_sofia@173.244.183.139:5060</data>
<state>RUNNING (0)</state>
</profile>
<gateway>
<name>example.com</name>
<type>gateway</type>
<data>sip:joeuser@example.com</data>
<state>NOREG</state>
</gateway>

<profile>
<name>dot141</name>
<type>profile</type>
<data>sip:mod_sofia@173.244.183.141:5060</data>
<state>RUNNING (0)</state>
</profile>
<gateway>
<name>example.com</name>
<type>gateway</type>
<data>sip:joeuser@example.com</data>
<state>NOREG</state>
</gateway>

<profile>
<name>external</name>
<type>profile</type>
<data>sip:mod_sofia@173.244.183.138:5060</data>
<state>RUNNING (63)</state>
</profile>
<gateway>
<name>example.com</name>
<type>gateway</type>
<data>sip:joeuser@example.com</data>
<state>NOREG</state>
</gateway>

<profile>
<name>dot140</name>
<type>profile</type>
<data>sip:mod_sofia@173.244.183.140:5060</data>
<state>RUNNING (0)</state>
</profile>
<gateway>
<name>example.com</name>
<type>gateway</type>
<data>sip:joeuser@example.com</data>
<state>NOREG</state>
</gateway>

</profiles>
`
)

func TestSofiaStatusParse(t *testing.T) {
	sofiaProfiles, err := ParseSofiaStatus(sofiaData)
	if err != nil {
		t.Fatal("error parsing data: ", err)
	}
	if len(sofiaProfiles) != 5 {
		t.Errorf("error extracting profiles: %+v", sofiaProfiles)
	}
	sp := sofiaProfiles[3]
	if sp.Name != "external" || sp.Address != "173.244.183.138:5060" || sp.Running != "63" {
		t.Error("error on sofia profiles correction: ", sp)
	}
}

func TestSofiaStatusParseEmpty(t *testing.T) {
	_, err := ParseSofiaStatus("")
	if err == nil {
		t.Error("failed to return parsing error: ", err)
	}
}

func TestSofiaStatusParseBadData(t *testing.T) {
	_, err := ParseSofiaStatus("bad data")
	if err == nil {
		t.Error("failed to return parsing error: ", err)
	}
}
