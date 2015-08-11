//A high performance, robust library for parsing RSS1.0/2.0 Including Itunes and MediaRSS Extensions
package easyrss

import (
	"errors"
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"strconv"
	"time"
)

type RSS struct {
	channel Channel
}

type xmlRSS struct {
	channel xml.Node
}

type xmlChannel struct {
	items []xml.Node
}

type Channel struct {
	title       string           //Channel title
	link        string           //Channel link
	generator   string           //channel Generator
	description string           //Channel description
	language    string           //Channel language
	copyright   string           //Channel Copyright
	categories  []string         //Channel Categories
	items       []Item           //Slice of the items in the channel
	itunes      ItunesMeta       //Itunes Podcast Category
	media       MediaChannelMeta //MediaRSS Channel Metadata

	isItunes bool
	isMRSS   bool
}

type RSSEnclosure struct {
	url       string
	mediaType string
	size      uint64
}

type GUIDField struct {
	IsPermaLink bool
	Content     string
}

type Item struct {
	title       string       //Item title
	link        string       //Item link
	date        *time.Time   //Item publication time
	media       MediaMeta    //MediaRSS Fields
	description string       //Item description
	enclosure   RSSEnclosure //Optional RSS Media Enclosure
	guid        GUIDField    //Item GUID Info
	itunes      ItunesMeta   //ITunes Podcast RSS Fields

	isItunes     bool //Whether item contains ITunes RSS Extensions
	isMRSS       bool //Whether item contains MediaRSS Extensions
	hasEnclosure bool //Whether item contains an enclosure
}

//Pass in a byte slice containing the feed, get an *RSS back. You can then explore the feed easily.
func Decode(data []byte) (*RSS, error) {
	rssObj := RSS{}
	xmlrssObj := xmlRSS{} //For quick access to key nodes
	xmlChanObj := xmlChannel{}
	xmlDoc, err := gokogiri.ParseXml(data)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	rootNode := xmlDoc.Root()
	getChannel(&xmlrssObj, rootNode)
	getChannelElem(&rssObj, xmlrssObj.channel)
	getItems(&xmlChanObj, xmlrssObj.channel)
	rssObj.channel.items = make([]Item, len(xmlChanObj.items))
	for itemID := 0; itemID < len(xmlChanObj.items); itemID++ {
		getItemMeta(&rssObj, itemID, xmlChanObj.items[itemID])
	}
	xmlDoc.Free()
	return &rssObj, nil
}

//Searches for Channel Metadata and Populates it
func getChannel(x *xmlRSS, root *xml.ElementNode) {
	for activeChannel := root.FirstChild(); activeChannel != nil; activeChannel = activeChannel.NextSibling() {
		tag := activeChannel.Name()
		if tag == "channel" {
			x.channel = activeChannel
			return //We should only have one channel element
		}
	}
}

//Get Metadata and Fetch Item nodes for Channel
func getChannelElem(r *RSS, c xml.Node) xmlChannel {
	x := xmlChannel{}
	for activeElem := c.FirstChild(); activeElem != nil; activeElem = activeElem.NextSibling() {
		tag := activeElem.Name()
		namespace := activeElem.Namespace()
		tagContent := activeElem.Content()
		if tag == "text" {
			continue
		}
		switch namespace {
		case "http://www.itunes.com/dtds/podcast-1.0.dtd":
			r.channel.isItunes = true
			setItunesMetaField(activeElem, &r.channel.itunes)
		case "http://search.yahoo.com/mrss/":
			r.channel.isMRSS = true
			setMediaChannelMetaField(activeElem, &r.channel.media)
		case "":
			switch tag {
			case "title":
				r.channel.title = tagContent
			case "link":
				r.channel.link = tagContent
			case "generator":
				r.channel.generator = tagContent
			case "description":
				r.channel.description = tagContent
			case "language":
				r.channel.language = tagContent
			case "copyright":
				r.channel.copyright = tagContent
			case "category":
				r.channel.categories = append(r.channel.categories, tagContent)
			case "item":
				x.items = append(x.items, activeElem)
			}
		}
	}
	return x
}

//Filter out stuff that's not an item
func getItems(x *xmlChannel, c xml.Node) {
	for activeItem := c.FirstChild(); activeItem != nil; activeItem = activeItem.NextSibling() {
		tag := activeItem.Name()
		if tag == "item" {
			x.items = append(x.items, activeItem)
		}
	}
}

//Sets Appropriate Item Metadata
func getItemMeta(r *RSS, itemID int, i xml.Node) {
	r.channel.items[itemID].media.credits = make(map[string]string)
	for activeElem := i.FirstChild(); activeElem != nil; activeElem = activeElem.NextSibling() {
		tag := activeElem.Name()
		if tag == "text" {
			continue
		}
		tagContent := activeElem.Content()
		namespace := activeElem.Namespace()
		switch namespace {
		case "http://www.itunes.com/dtds/podcast-1.0.dtd": //iTunes Podcast RSS Namespace
			r.channel.items[itemID].isItunes = true
			setItunesMetaField(activeElem, &r.channel.items[itemID].itunes)
		case "http://search.yahoo.com/mrss/": //MediaRSS Namespace
			setMediaMetaField(activeElem, &r.channel.items[itemID].media)
		case "":
			switch tag {
			case "title":
				r.channel.items[itemID].title = tagContent
			case "link":
				r.channel.items[itemID].link = tagContent
			case "pubDate":
				if parsedDate, perr := time.Parse("Mon, 02 Feb 2006 15:04:05 -0700", tagContent); perr == nil {
					r.channel.items[itemID].date = &parsedDate
				}
			case "description":
				r.channel.items[itemID].description = tagContent
			case "enclosure":
				r.channel.items[itemID].hasEnclosure = true
				if urlAttr := activeElem.Attribute("url"); urlAttr != nil {
					r.channel.items[itemID].enclosure.url = urlAttr.Value()
				}
				if mediaType := activeElem.Attribute("type"); mediaType != nil {
					r.channel.items[itemID].enclosure.mediaType = mediaType.Value()
				}
				if fileSize := activeElem.Attribute("length"); fileSize != nil {
					r.channel.items[itemID].enclosure.size, _ = strconv.ParseUint(fileSize.Value(), 10, 64)
				}
			}
		}
	}
}

//Returns available feed items. Will return an error if the feed is empty.
func (r *RSS) Items() ([]Item, error) {
	if len(r.channel.items) == 0 {
		return nil, errors.New("Feed contains no items")
	}
	return r.channel.items, nil
}

//Returns the feed title. If the field is not populated, will return an empty string and an error.
func (r *RSS) Title() (string, error) {
	if r.channel.title == "" {
		return "", errors.New("Feed title is not populated")
	}
	return r.channel.title, nil
}

//Returns the feed generator. If the field is not populated, will return an empty string and an error.
func (r *RSS) Generator() (string, error) {
	if r.channel.generator == "" {
		return "", errors.New("Feed generator is not populated")
	}
	return r.channel.generator, nil
}

//Returns the feed description. If the field is not populated, will return an empty string and an error.
func (r *RSS) Description() (string, error) {
	if r.channel.description == "" {
		return "", errors.New("Feed description is not populated")
	}
	return r.channel.description, nil
}

//Returns the feed language. This isn't always that standardized so be careful while parsing the field. An empty string and error will be returned if the field is not populated.
func (r *RSS) Language() (string, error) {
	if r.channel.language == "" {
		return "", errors.New("Feed language is not populated")
	}
	return r.channel.language, nil
}

//Returns the feed categories. If no category tags were found for the channel, you'll get a nil result with an accompanying error.
func (r *RSS) Categories() ([]string, error) {
	if len(r.channel.categories) == 0 {
		return nil, errors.New("Feed categories not populated")
	}
	return r.channel.categories, nil
}

//Returns the item title. If the item title is not populated, you'll get an empty string and an error.
func (i Item) Title() (string, error) {
	if i.title == "" {
		return "", errors.New("Item title is not populated")
	}
	return i.title, nil
}

//Returns the item link. If the item link is not populated, you'll get an empty string and an error.
func (i Item) Link() (string, error) {
	if i.link == "" {
		return "", errors.New("Item link is not populated")
	}
	return i.link, nil
}

//Returns the item date. If the item date is not populated, you'll get nil and an error.
func (i Item) Date() (*time.Time, error) {
	if i.date == nil {
		return nil, errors.New("Item date not populated")
	}
	return i.date, nil
}

//Whether or not the item has a media enclosure.
func (i Item) HasEnclosure() bool {
	return i.hasEnclosure
}

//The media enclosure url.
func (i Item) GetEnclosureURL() string {
	return i.enclosure.url
}

//Returns the item . If the item title is not populated, you'll get an empty string and an error.
func (i Item) Description() (string, error) {
	if i.description == "" {
		return "", errors.New("Item description is not populated")
	}
	return i.description, nil
}
