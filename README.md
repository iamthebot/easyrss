[![](https://img.shields.io/badge/godoc-complete-blue.svg)](http://godoc.org/github.com/iamthebot/easyrss)
# easyrss
Made this after finding no barebones rss packages that simply parsed an RSS feed from an existing byte slice. Easyrss doesn't require you fetch, update, parse, or format anything. Simply pass along a byte slice containing the feed and you're in business.

Absolutely no date parsing, etc. happens inside. You simply get a struct containing the channels (and metadata) with channels containing all the items you want. The package is only 45LOC.
