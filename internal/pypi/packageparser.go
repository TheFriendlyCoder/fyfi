// Primitives for interacting with data loaded from a PEP-503 compatible pypi
// service. https://peps.python.org/pep-0503/
package pypi

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

// TODO: add unit tests for parsing code
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

// getAttributes converts the HTML attributes associated with a specific
// node / tag into a hash map format to simplify access
func getAttributes(node *html.Node) map[string]string {
	retval := map[string]string{}

	for _, attr := range node.Attr {
		retval[attr.Key] = attr.Val
	}
	return retval
}

// ParseDistribution parses Python distribution and package information from
// HTML content retrieved from pypi.org or a cmpatible mirror that implements
// the PEP-503 standard: https://peps.python.org/pep-0503/. Returns the heading
// node containing descriptive information about the distribution, and a list
// of 0 or more anchor nodes containing information about individual package
// releases of the distribution
func ParseDistribution(htmlData string) (*html.Node, []*html.Node, error) {
	node, err := html.Parse(strings.NewReader(htmlData))
	if err != nil {
		return nil, nil, err
	}

	node = findFirstChild(node, "html")
	if node == nil {
		return nil, nil, errors.New("unable to parse HTML content")
	}
	node = findFirstChild(node, "body")
	if node == nil {
		return nil, nil, errors.New("unable to parse BODY content")
	}

	heading := findFirstChild(node, "h1")
	if heading == nil {
		return nil, nil, errors.New("unable to parse H1 heading content")
	}
	anchors := findAllChildren(node, "a")

	return heading, anchors, nil
}
