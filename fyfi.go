package main

import (
	"errors"
	"log"
	"strings"

	"golang.org/x/net/html"
	"gorm.io/gorm"
)

type Package struct {
	URL           string
	filename      string
	pythonVersion string
	// TODO: parse out checksum from URL and store separately here
	// checksum string
}

func NewPackage(anchor *html.Node) (*Package, error) {
	url := ""
	var pyver string
	for _, attr := range anchor.Attr {
		switch {

		case attr.Key == "href":
			url = attr.Val
		case attr.Key == "data-requires-python":
			pyver = html.UnescapeString(attr.Val)
		}
	}
	filename := anchor.FirstChild.Data
	if url == "" {
		return nil, errors.New("failed to parse package URL from pypi response")
	}
	return &Package{URL: url, filename: filename, pythonVersion: pyver}, nil
}

type Distro struct {
	gorm.Model
	name     string
	packages []*Package
}

func NewDistro(htmlData string) (*Distro, error) {
	node, err := html.Parse(strings.NewReader(htmlData))
	if err != nil {
		return nil, err
	}

	var find_first_child = func(node *html.Node, name string) (*html.Node, error) {
		for n := node.FirstChild; n != nil; n = n.NextSibling {
			if n.Type != html.ElementNode {
				continue
			}
			if n.Data == name {
				return n, nil
			}
		}
		return nil, errors.New("unable to find node")
	}

	var find_all_children = func(node *html.Node, name string) []*html.Node {
		retval := make([]*html.Node, 0)
		for n := node.FirstChild; n != nil; n = n.NextSibling {
			if n.Type != html.ElementNode {
				continue
			}
			if n.Data == name {
				retval = append(retval, n)
			}
		}
		return retval
	}
	node, err = find_first_child(node, "html")
	if err != nil {
		return nil, err
	}
	node, err = find_first_child(node, "body")
	if err != nil {
		return nil, err
	}

	anchors := find_all_children(node, "a")
	packages := make([]*Package, len(anchors))
	for i, a := range anchors {
		temp, err := NewPackage(a)
		if err != nil {
			return nil, err
		}
		packages[i] = temp

	}

	header, err := find_first_child(node, "h1")
	if err != nil {
		return nil, err
	}
	temp := strings.Split(header.FirstChild.Data, " ")

	return &Distro{name: temp[len(temp)-1], packages: packages}, nil
}

func main() {
	log.Fatal("Error from KSP")
	log.Println("Shouldn't see me")
}
