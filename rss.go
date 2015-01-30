//Stupidly simple package for parsing XML RSS from a byte slice.
package easyrss

import (
	"bytes"
	"encoding/xml"
)

type RSS struct {
	XMLname  xml.Name   `xml:"rss"`
	Channels []*Channel `xml:"channel"` //Slice of the channels
}

type Channel struct {
	XMLname     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`       //Channel title
	Link        string   `xml:"link"`        //Channel link
	Description string   `xml:"description"` //Channel description
	Language    string   `xml:"language"`    //Channel language
	Items       []*Item  `xml:"item"`        //Slice of the items in the channel
}

type RSSEnclosure struct {
	URL       string `xml:"url,atrr"`
	MediaType string `xml:"type,attr"`
}
type Item struct {
	XMLname     xml.Name     `xml:"item"`
	Title       string       `xml:"title"`       //Item title
	Link        string       `xml:"link"`        //Item link
	Content     string       `xml:"encoded"`     //Item content
	Description string       `xml:"description"` //Item description
	Category    []string     `xml:"category"`    //Item categories
	Enclosure   RSSEnclosure `xml"enclosure"`    //Optional RSS Media Enclosure
	Date        string       `xml:"pubDate"`     //Last date of item publication
	GUID        string       `xml:"guid"`
}

//Pass in a byte slice containing the feed, get an RSS struct back, with stuff populated.
func Decode(data []byte) (*RSS, error) {
	r := bytes.NewReader(data)
	decoder := xml.NewDecoder(r)
	rss := &RSS{}
	return rss, decoder.Decode(rss)
}
