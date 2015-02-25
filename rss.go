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
	XMLname         xml.Name    `xml:"channel"`
	Title           string      `xml:"title"`           //Channel title
	Link            string      `xml:"link"`            //Channel link
	Generator       string      `xmlL"generator"`       //Channel Generator
	Description     string      `xml:"description"`     //Channel description
	Language        string      `xml:"language"`        //Channel language
	Copyright       string      `xml:"copyright"`       //Channel Copyright
	Categories      []string    `xml:"category"`        //Channel Categories
	Items           []*Item     `xml:"item"`            //Slice of the items in the channel
	MediaCategories []string    `xml:"media:category"`  //Slice of media tag categories
	MediaCopyright  string      `xml:"media:copyright"` //Media Copyright
	MediaRating     string      `xml:"media:rating"`    //Media Rating
	MediaThumbnail  Image       `xml:"media:thumbnail"` //Media Thumbnail
	ItunesCategory  string      `xml:"itunes:category"` //Itunes Podcast Category
	ItunesOwner     ItunesOwner `xml:"itunes:owner"`    //Itunes Podcast Owner Contact Info
	ItunesKeywords  string      `xml:"itunes:keywords"` //Itunes Podcast Keywords

}

type ItunesOwner struct {
	Name  string `xml:"itunes:name"`
	Email string `xml:"itunes:email"`
}
type RSSEnclosure struct {
	URL       string `xml:"url,attr"`
	MediaType string `xml:"type,attr"`
	Size      string `xml:"fileSize,attr"`
}
type Image struct {
	URL    string `xml:"url,attr"`
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
}
type Item struct {
	XMLname        xml.Name     `xml:"item"`
	Title          string       `xml:"title"`       //Item title
	Link           string       `xml:"link"`        //Item link
	Content        string       `xml:"encoded"`     //Item content
	Description    string       `xml:"description"` //Item description
	Category       []string     `xml:"category"`    //Item categories
	Enclosure      RSSEnclosure `xml"enclosure"`    //Optional RSS Media Enclosure
	Date           string       `xml:"pubDate"`     //Last date of item publication
	GUID           string       `xml:"guid"`
	ItunesAuthor   string       `xml:"itunes:author"`   //Itunes Episode Author
	ItunesImage    Image        `xml:"itunes:image"`    //Itunes Episode Th
	ItunesSubtitle string       `xml:"itunes:subtitle"` //Itunes Episode Subtitle
	ItunesSummary  string       `xml"itunes:summary"`   //Itunes Episode Summary
	MediaContent   RSSEnclosure `xml:"media:content"`   //Media Payload
	MediaThumbnail Image        `xml:"media:thumbnail"` //Media Thumbnail
}

//Pass in a byte slice containing the feed, get an RSS struct back, with stuff populated.
func Decode(data []byte) (*RSS, error) {
	r := bytes.NewReader(data)
	decoder := xml.NewDecoder(r)
	rss := &RSS{}
	return rss, decoder.Decode(rss)
}
