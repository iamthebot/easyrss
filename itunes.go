package easyrss

import (
	"errors"
	"github.com/moovweb/gokogiri/xml"
	"strconv"
	"strings"
	"time"
)

type ItunesMeta struct {
	author   string
	subtitle string
	summary  string
	image    Image
	explicit string
	duration time.Duration
	keywords string
}

//Sets Appropriate Field Given Itunes Node
func setItunesMetaField(n xml.Node, i *ItunesMeta) {
	tag := n.Name()
	tagContent := n.Content()
	switch tag {
	case "subtitle":
		i.subtitle = tagContent
	case "author":
		i.author = tagContent
	case "summary":
		i.summary = tagContent
	case "explicit":
		i.explicit = tagContent
	case "keywords":
		i.keywords = tagContent
	case "duration":
		splitDur := strings.Split(tagContent, ":")
		if len(splitDur) != 3 && len(splitDur) != 2 { //Not H:M:S/M:S format, so probably just a seconds integer
			duration, err := time.ParseDuration(tagContent + "s")
			if err != nil {
				return
			}
			i.duration = duration
			return
		}

		if len(splitDur) == 3 {
			hours, e := strconv.ParseInt(splitDur[0], 10, 8)
			if e != nil {
				return
			}
			minutes, e := strconv.ParseInt(splitDur[1], 10, 8)
			if e != nil {
				return
			}
			seconds, e := strconv.ParseInt(splitDur[2], 10, 8)
			if e != nil {
				return
			}
			i.duration = time.Duration(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second)
		} else {
			minutes, e := strconv.ParseInt(splitDur[0], 10, 8)
			if e != nil {
				return
			}
			seconds, e := strconv.ParseInt(splitDur[1], 10, 8)
			if e != nil {
				return
			}
			i.duration = time.Duration(time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second)
		}

	case "image":
		if urlNode := n.Attribute("href"); urlNode != nil {
			i.image.url = urlNode.Value()
		}
	default:
		return
	}
}

//Whether or not this feed implements ItunesRSS Extensions
func (r *RSS) IsItunes() bool {
	return r.channel.isItunes
}

//Returns the Itunes "author" field for the channel. If the channel doesn't contain ITunes Extensions or hasn't populated the channel-wide Itunes "author" field, will return an empty string and an error
func (r *RSS) ItunesAuthor() (string, error) {
	if !r.channel.isItunes {
		return "", errors.New("Not an Itunes RSS Feed")
	}
	if r.channel.itunes.author == "" {
		return "", errors.New("Itunes author field not populated")
	}
	return r.channel.itunes.author, nil
}

//Returns the Itunes "author" field for the channel. If the channel doesn't contain ITunes Extensions or hasn't populated the channel-wide Itunes "subtitle" field, will return an empty string and an error
func (r *RSS) ItunesSubtitle() (string, error) {
	if !r.channel.isItunes {
		return "", errors.New("Not an Itunes RSS Feed")
	}
	if r.channel.itunes.subtitle == "" {
		return "", errors.New("Itunes subtitle field not populated")
	}
	return r.channel.itunes.subtitle, nil
}

//Returns the Itunes "summary" field for the channel. If the channel doesn't contain ITunes Extensions or hasn't populated the channel-wide Itunes "summary" field, will return an empty string and an error
func (r *RSS) ItunesSummary() (string, error) {
	if !r.channel.isItunes {
		return "", errors.New("Not an Itunes RSS Feed")
	}
	if r.channel.itunes.subtitle == "" {
		return "", errors.New("Itunes summary field not populated")
	}
	return r.channel.itunes.subtitle, nil
}

//Returns the Itunes "image" field for the channel. If the channel doesn't contain ITunes Extensions or hasn't populated the channel-wide Itunes "image" field, will return nil and an error.
func (r *RSS) ItunesImage() (*Image, error) {
	if !r.channel.isItunes {
		return nil, errors.New("Not an Itunes RSS Feed")
	}
	if r.channel.itunes.image.url == "" {
		return nil, errors.New("Itunes image fields not populated")
	}
	return &r.channel.itunes.image, nil
}

//Returns the Itunes "explicit" field for the channel. If the channel doesn't contain ITunes Extensions or hasn't populated the channel-wide Itunes "explicit" field, will return an empty string and an error func (r *RSS) ItunesExplicit() (string, error) {
func (r *RSS) ItunesExplicit() (string, error) {
	if !r.channel.isItunes {
		return "", errors.New("Not an Itunes RSS Feed")
	}
	if r.channel.itunes.explicit == "" {
		return "", errors.New("Itunes explicit field not populated")
	}
	return r.channel.itunes.explicit, nil
}

//Returns the Itunes "author" field for the item. If the item doesn't contain ITunes Extensions or hasn't populated the Itunes "author" field, will return an empty string and an error
func (i *Item) ItunesAuthor() (string, error) {
	if !i.isItunes {
		return "", errors.New("Not an Itunes RSS Feed")
	}
	if i.itunes.author == "" {
		return "", errors.New("Itunes author field not populated")
	}
	return i.itunes.author, nil
}

//Returns the Itunes "author" field for the item. If the item doesn't contain ITunes Extensions or hasn't populated the Itunes "subtitle" field, will return an empty string and an error
func (i *Item) ItunesSubtitle() (string, error) {
	if !i.isItunes {
		return "", errors.New("Not an Itunes RSS Feed")
	}
	if i.itunes.subtitle == "" {
		return "", errors.New("Itunes subtitle field not populated")
	}
	return i.itunes.subtitle, nil
}

//Returns the Itunes "summary" field for the item. If the item doesn't contain ITunes Extensions or hasn't populated the Itunes "summary" field, will return an empty string and an error
func (i *Item) ItunesSummary() (string, error) {
	if !i.isItunes {
		return "", errors.New("Not an Itunes RSS Feed")
	}
	if i.itunes.subtitle == "" {
		return "", errors.New("Itunes summary field not populated")
	}
	return i.itunes.subtitle, nil
}

//Returns Itunes episode duration. If this information wasn't available or the item doesn't contain Itunes Extensions then we return nil and an error.
func (i Item) ItunesDuration() (*time.Duration, error) {
	if !i.isItunes {
		return nil, errors.New("Not an Itunes RSS Feed")
	}
	if int(i.itunes.duration) == 0 {
		return nil, errors.New("Itunes duration field missing")
	}
	return &i.itunes.duration, nil
}

//Returns the Itunes "image" field for the item. If the item doesn't contain ITunes Extensions or hasn't populated the Itunes "image" field, will return nil and an error.
func (i *Item) ItunesImage() (*Image, error) {
	if !i.isItunes {
		return nil, errors.New("Not an Itunes RSS Feed")
	}
	if i.itunes.image.url == "" {
		return nil, errors.New("Itunes image fields not populated")
	}
	return &i.itunes.image, nil
}
