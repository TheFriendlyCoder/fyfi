// Primitives for interacting with data loaded from a PEP-503 compatible pypi
// service. https://peps.python.org/pep-0503/
package pypi

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type PythonPackage struct {
	URL           string
	Filename      string
	PythonVersion string
	Checksum      string
}

type PythonDistribution struct {
	Name     string
	Packages []PythonPackage
}

// findFirstChild locates the first matching child HTML node under a parent
// node with the specified tag name. Returns nil if no matching child node
// can be found.
// NOTE: Search is non-recursive
func findFirstChild(node *html.Node, name string) *html.Node {
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		if n.Type != html.ElementNode {
			continue
		}
		if n.Data == name {
			return n
		}
	}
	return nil
}

//findAllChildren looks for multiple child HTML nodes under a parent node
// with a specified tag name. Returns the list of 0 or more nodes with
// tags that match.
// NOTE: Search is non-recursive
func findAllChildren(node *html.Node, name string) []*html.Node {
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

func createPackage(anchor *html.Node) (*PythonPackage, error) {
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

	return &PythonPackage{URL: url, Filename: filename, PythonVersion: pyver, Checksum: checksum}, nil
}

func createDistro(node *html.Node) (*PythonDistribution, error) {
	node = findFirstChild(node, "html")
	if node == nil {
		return nil, errors.New("unable to parse HTML content")
	}
	node = findFirstChild(node, "body")
	if node == nil {
		return nil, errors.New("unable to parse BODY content")
	}

	heading := findFirstChild(node, "h1")
	if heading == nil {
		return nil, errors.New("unable to parse H1 heading content")
	}

	parts := strings.Split(heading.FirstChild.Data, " ")
	name := parts[len(parts)-1]

	anchors := findAllChildren(node, "a")

	packages := make([]PythonPackage, len(anchors))
	for i, a := range anchors {
		temp, err := createPackage(a)
		if err != nil {
			return nil, err
		}
		packages[i] = *temp
	}

	return &PythonDistribution{Name: name, Packages: packages}, nil
}

// ParseDistribution parses Python distribution and package information from
// HTML content retrieved from pypi.org or a cmpatible mirror that implements
// the PEP-503 standard: https://peps.python.org/pep-0503/. Returns the heading
// node containing descriptive information about the distribution, and a list
// of 0 or more anchor nodes containing information about individual package
// releases of the distribution
func ParseDistribution(htmlData string) (*PythonDistribution, error) {
	node, err := html.Parse(strings.NewReader(htmlData))
	if err != nil {
		return nil, err
	}
	return createDistro(node)
}
