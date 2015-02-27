[![](https://img.shields.io/badge/godoc-complete-blue.svg)](http://godoc.org/github.com/iamthebot/easyrss)
# easyrss
A Go library designed from the ground up to handle the complete RSS 2.0 specification with Itunes and MediaRSS extensions.

It features a very simple but powerful API with helpful logic to determine what kind of feed you're looking at and whether each field is available. Easyrss will also decode podcast episode durations, publication dates, and many other fields into appropriate Go objects like Time.Time and map[string][string].

Instead of relying on encoding/xml (which won't work on many feeds that deviate from the spec), easyrss uses libxml by way of the [gokogiri](https://github.com/moovweb/gokogiri) bindings. This ensures a high degree of robustness for feeds that include non-standard fields and have non-complaint XML. Performance is on-par or faster than using encoding/XML directly.

## installation
First, install libxml and its development headers by way of your package manager. Usually this is something like:

Ubuntu:
```bash
sudo apt-get install libxml2-dev libxml2
```

Fedora/CentOS/RHEL:
```bash
sudo apt-get install libxml2-devel
```

ArchLinux:
```bash
sudo pacman -S libxml2
```

Then install this package:
```bash
go get https://github.com/iamthebot/easyrss.git
```
Easy! Refer to [godoc](http://godoc.org/github.com/iamthebot/easyrss) for complete API documentation.
