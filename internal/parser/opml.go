package parser

import (
	"encoding/xml"
)

// OPML represents the root of the XML file
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Body    struct {
		Outlines []Outline `xml:"outline"`
	} `xml:"body"`
}

// Outline represents an individual feed (or a folder of feeds)
type Outline struct {
	XMLURL   string    `xml:"xmlUrl,attr"`
	Outlines []Outline `xml:"outline"` // OPML allows folders inside folders!
}

// ExtractFeeds parses an OPML byte slice and returns a flat list of feed URLs
func ExtractFeeds(data []byte) []string {
	var opml OPML
	err := xml.Unmarshal(data, &opml)
	if err != nil {
		return nil
	}

	var urls []string
	
	// We use a recursive function because users might have feeds nested deep inside folders
	var extract func(outlines []Outline)
	extract = func(outlines []Outline) {
		for _, o := range outlines {
			if o.XMLURL != "" {
				urls = append(urls, o.XMLURL)
			}
			// Dive into any sub-folders
			extract(o.Outlines)
		}
	}

	extract(opml.Body.Outlines)
	return urls
}