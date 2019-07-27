package utils

import (
	"bytes"
	"encoding/xml"
	"errors"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
)

var (
	dataParser    = regexp.MustCompile(`^sips?:.+?@(.+:\d+)`)
	runningParser = regexp.MustCompile(`^RUNNING \((\d+)\)`)
)

type Profile struct {
	Name  string `xml:"name"`
	Type  string `xml:"type"`
	Data  string `xml:"data"`
	State string `xml:"state"`
}
type Profiles struct {
	Profiles []*Profile `xml:"profile"`
}

type SofiaProfile struct {
	Name    string
	Address string
	Running string
}

func (sp *SofiaProfile) loadXMLProfile(p *Profile) error {
	sp.Name = p.Name
	dataGroupSlice := dataParser.FindStringSubmatch(p.Data)
	if len(dataGroupSlice) != 2 {
		return errors.New("could not find address info in: " + p.Data)
	}
	sp.Address = dataGroupSlice[1]
	runningGroupSlice := runningParser.FindStringSubmatch(p.State)
	if len(runningGroupSlice) != 2 {
		return errors.New("cannot find running info in: " + p.State)
	}
	sp.Running = runningGroupSlice[1]
	return nil
}

func ParseSofiaStatus(data string) ([]*SofiaProfile, error) {
	data = strings.TrimSpace(data)
	dec := xml.NewDecoder(bytes.NewBufferString(data))
	dec.CharsetReader = charset.NewReaderLabel

	profiles := &Profiles{}
	err := dec.Decode(profiles)
	if err != nil {
		return nil, err
	}
	sofiaProfiles := make([]*SofiaProfile, len(profiles.Profiles))
	for i, p := range profiles.Profiles {
		sp := &SofiaProfile{}
		if err := sp.loadXMLProfile(p); err != nil {
			return nil, err
		}
		sofiaProfiles[i] = sp
	}
	return sofiaProfiles, nil
}
