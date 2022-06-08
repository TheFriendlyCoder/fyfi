package datamodel

import (
	"context"
	"fmt"
	"strings"

	"github.com/TheFriendlyCoder/fyfi/ent"
	"github.com/TheFriendlyCoder/fyfi/internal/pypi"
)

func CreateDistro(ctx context.Context, client *ent.Client, htmlData string) (*ent.PythonDistro, error) {
	heading, anchors, err := pypi.ParseDistribution(htmlData)
	if err != nil {
		return nil, err
	}
	temp := strings.Split(heading.FirstChild.Data, " ")
	//return &Distro{name: temp[len(temp)-1], packages: packages}, nil

	pd := client.PythonDistro.
		Create().
		SetName(temp[len(temp)-1])
	retval, err := pd.Save(ctx)

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
