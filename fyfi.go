package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/html"
	"gorm.io/gorm"
)

type Package struct {
	URL           string
	filename      string
	pythonVersion string
	checksum      string
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

	parts := strings.Split(url, "#")
	if len(parts) != 2 {
		return nil, fmt.Errorf("failed to parse out checksum: %s", url)
	}
	checksum := parts[len(parts)-1]
	if !strings.HasPrefix(checksum, "sha256=") {
		return nil, fmt.Errorf("failed to load checksum type for %s", checksum)
	}
	checksum = strings.Split(checksum, "=")[1]
	return &Package{URL: url, filename: filename, pythonVersion: pyver, checksum: checksum}, nil
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

func simple(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	DistName := ps.ByName("library")
	log.Printf("Querying for dist %s\n", DistName)

	SrcURL := fmt.Sprintf("https://pypi.org/simple/%s/", DistName)

	resp, err := http.Get(SrcURL)
	if err != nil {
		log.Println("Unable to communicate with source repo")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to communicate with source repo")
		return
	}

	distro, err := NewDistro(string(body))
	if err != nil {
		log.Printf("Error parsing response data: %v\n", err)
		return
	}
	log.Printf("Loading distro %s\n", distro.name)
	log.Printf("Found %d packages for distro %s\n", len(distro.packages), distro.name)
	log.Printf("First package %s has pyver %s\n", distro.packages[0].filename, distro.packages[0].pythonVersion)
	fmt.Fprint(w, string(body))
}

func main() {
	router := httprouter.New()
	router.GET("/simple/:library/", simple)
	router.GET("/simple/:library", simple)

	log.Println("Listing for requests at http://localhost:8000/simple")
	log.Fatal(http.ListenAndServe(":8000", router))
}
