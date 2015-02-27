package easyrss

import (
	"errors"
	"github.com/moovweb/gokogiri/xml"
	"strconv"
	"strings"
)

//MediaRSS Item Metadata
type MediaMeta struct {
	content   RSSEnclosure
	credits   map[string]string //Media Credits... Usually role->photographer
	thumbnail Image
}

//MediaRSS Channel Metadata
type MediaChannelMeta struct {
	rating     string   //Age Rating
	copyright  string   //Feed Copyright
	thumbnail  Image    //Feed Thumbnail
	keywords   []string //Feed Keywords
	categories []string //Feed Categories
}

func setMediaMetaField(n xml.Node, m *MediaMeta) {
	tag := n.Name()
	switch tag {
	case "content":
		if urlAttr := n.Attribute("url"); urlAttr != nil {
			m.content.url = urlAttr.Value()
		}
		if typeAttr := n.Attribute("type"); typeAttr != nil {
			m.content.url = typeAttr.Value()
		}
		if sizeAttr := n.Attribute("fileSize"); sizeAttr != nil {
			fSize := sizeAttr.Value()
			m.content.size, _ = strconv.ParseUint(fSize, 10, 64)
		}
	case "thumbnail":
		if urlAttr := n.Attribute("url"); urlAttr != nil {
			m.thumbnail.url = urlAttr.Value()
		}
		if widthAttr := n.Attribute("width"); widthAttr != nil {
			width, _ := strconv.ParseInt(widthAttr.Value(), 10, 8)
			m.thumbnail.width = int(width)
		}
		if heightAttr := n.Attribute("height"); heightAttr != nil {
			height, _ := strconv.ParseInt(heightAttr.Value(), 10, 8)
			m.thumbnail.height = int(height)
		}
	case "credits":
		if roleAttr := n.Attribute("role"); roleAttr != nil {
			m.credits[roleAttr.Value()] = n.Content()
		}
	}
}

func setMediaChannelMetaField(n xml.Node, m *MediaChannelMeta) {
	tag := n.Name()
	tagContent := n.Content()
	switch tag {
	case "rating":
		m.rating = tagContent
	case "copyright":
		m.copyright = tagContent
	case "thumbnail":
		if urlAttr := n.Attribute("url"); urlAttr != nil {
			m.thumbnail.url = urlAttr.Value()
		}
		if widthAttr := n.Attribute("width"); widthAttr != nil {
			width, _ := strconv.ParseInt(widthAttr.Value(), 10, 8)
			m.thumbnail.width = int(width)
		}
		if heightAttr := n.Attribute("height"); heightAttr != nil {
			height, _ := strconv.ParseInt(heightAttr.Value(), 10, 8)
			m.thumbnail.height = int(height)
		}
	case "keywords":
		m.keywords = strings.Split(tagContent, ", ")
	case "category":
		m.categories = append(m.categories, tagContent)
	}
}

//Whether or not this feed implements MediaRSS Extensions
func (r *RSS) IsMRSS() bool {
	return r.channel.isMRSS
}

//MediaRSS Feed Rating. If the MRSS feed "rating" field is not populated or if the feed doesn't implement MediaRSS extensions, you'll receive an empty string and an error.
func (r *RSS) MRSSRating() (string, error) {
	if !r.channel.isMRSS {
		return "", errors.New("Not a MediaRSS Feed")
	}
	if r.channel.media.rating == "" {
		return "", errors.New("MediaRSS feed rating field is not populated")
	}
	return r.channel.media.rating, nil
}

//MediaRSS Feed Copyright. If the MRSS feed "copyright" field is not populated or if the feed doesn't implement MediaRSS extensions, you'll receive an empty string and an error.
func (r *RSS) MRSSCopyright() (string, error) {
	if !r.channel.isMRSS {
		return "", errors.New("Not a MediaRSS Feed")
	} else if r.channel.media.copyright == "" {
		return "", errors.New("MediaRSS feed copyright field is not populated")
	}
	return r.channel.media.copyright, nil
}

//Returns the Itunes "image" field for the channel. If the channel doesn't contain ITunes Extensions or hasn't populated the channel-wide Itunes "image" field, will return nil and an error.
func (r *RSS) Thumbnail() (*Image, error) {
	if !r.channel.isMRSS {
		return nil, errors.New("Not a MediaRSS Feed")
	} else if r.channel.media.thumbnail.url == "" {
		return nil, errors.New("MediaRSS thumbnail fields not populated")
	}
	return &r.channel.itunes.image, nil
}

//MediaRSS Feed keywords. If the MRSS feed "keywords" field is not populated or if the feed doesn't implement MediaRSS extensions, this will return nil and an error.
func (r *RSS) Keywords() ([]string, error) {
	if !r.channel.isMRSS {
		return nil, errors.New("Not a MediaRSS Feed")
	} else if len(r.channel.media.keywords) == 0 {
		return nil, errors.New("No MediaRSS Feed Keywords")
	}
	return r.channel.media.keywords, nil
}

//MediaRSS Feed categories. If the MRSS feed "categories" field is not populated or if the feed doesn't implement MediaRSS extensions, this will return nil and an error.
func (r *RSS) Categories() ([]string, error) {
	if !r.channel.isMRSS {
		return nil, errors.New("Not a MediaRSS Feed")
	} else if len(r.channel.media.categories) == 0 {
		return nil, errors.New("No MediaRSS Feed Categories")
	}
	return r.channel.media.categories, nil
}
