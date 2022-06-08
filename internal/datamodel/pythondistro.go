package datamodel

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/TheFriendlyCoder/fyfi/ent"
	"golang.org/x/net/html"
)

func CreateDistro(ctx context.Context, client *ent.Client, htmlData string) (*ent.PythonDistro, error) {
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

	header, err := find_first_child(node, "h1")
	if err != nil {
		return nil, err
	}
	temp := strings.Split(header.FirstChild.Data, " ")

	//return &Distro{name: temp[len(temp)-1], packages: packages}, nil

	pd := client.PythonDistro.
		Create().
		SetName(temp[len(temp)-1])
	retval, err := pd.Save(ctx)

	anchors := find_all_children(node, "a")
	packages := make([]*ent.PythonPackage, len(anchors))
	for i, a := range anchors {
		temp, err := CreatePackage(context.Background(), client, a, retval)
		if err != nil {
			return nil, err
		}
		packages[i] = temp
	}
	for _, p := range packages {
		pd.AddPackages(p)
	}

	// retval, err := client.PythonDistro.
	// 	Create().
	// 	SetName(temp[len(temp)-1]).
	// 	//AddPackageIDs(packages[0].ID, packages[1].ID).
	// 	//AddPackages(packages[0], packages[1]).
	// 	AddPackages(packages).
	// 	Save(ctx)
	//t := asdf.QueryPackages().WithDistro()
	// retval, err := client.PythonDistro.
	// 	Query().
	// 	Where(ent.PythonDistro.ID(asdf.ID)).
	// 	WithPackages().
	// 	All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating entity: %w", err)
	}
	//log.Println("entity was created: ", retval)
	return retval, nil
}
