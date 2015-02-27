package easyrss

import "errors"

type Image struct {
	title  string
	url    string
	link   string
	width  int
	height int
}

//Returns image title, if available. If not, an empty string is returned along with an error.
func (i *Image) Title() (string, error) {
	if i.title == "" {
		return "", errors.New("Image title field is not populated")
	}
	return i.title, nil
}

//Returns the url where the image is located, if available. If not, an empty string is returned along with an error.
func (i *Image) URL() (string, error) {
	if i.url == "" {
		return "", errors.New("Image URL is not populated")
	}
	return i.url, nil
}

//Returns the url where the image links to, if available. If not, an empty string is returned along with an error.
func (i *Image) Link() (string, error) {
	if i.link == "" {
		return "", errors.New("Image link is not populated")
	}
	return i.link, nil
}

//Returns the url image width in pixels, if available. If not, a width of 0 is returned along with an error.
func (i *Image) Width() (int, error) {
	if i.width == 0 {
		return 0, errors.New("Image width is not populated")
	}
	return i.width, nil
}

//Returns the url image width in pixels, if available. If not, a width of 0 is returned along with an error.
func (i *Image) Height() (int, error) {
	if i.height == 0 {
		return 0, errors.New("Image height is not populated")
	}
	return i.height, nil
}
