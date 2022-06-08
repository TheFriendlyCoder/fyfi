package datamodel

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/TheFriendlyCoder/fyfi/ent"
	"golang.org/x/net/html"
)

func CreatePackage(ctx context.Context, client *ent.Client, anchor *html.Node, parent *ent.PythonDistro) (*ent.PythonPackage, error) {
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

	retval, err := client.PythonPackage.
		Create().
		SetURL(url).
		SetChecksum(checksum).
		SetFilename(filename).
		SetPythonVersion(pyver).
		SetDistro(parent).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating entity: %w", err)
	}
	//log.Println("entity was created: ", retval)
	return retval, nil
}
