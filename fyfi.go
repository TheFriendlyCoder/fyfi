package main

// Reference: https://peps.python.org/pep-0503/
import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/TheFriendlyCoder/fyfi/ent"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/html"
)

var client *ent.Client

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

	distro, err := CreateDistro(context.Background(), client, string(body))
	if err != nil {
		log.Printf("Error parsing response data: %v\n", err)
		return
	}

	log.Printf("Loading distro %s\n", distro.Name)
	// log.Printf("Found %d packages for distro %s\n", len(distro.Edges.Packages), distro.Name)
	// log.Printf("First package %s has pyver %s\n", distro.Edges.Packages[0].Filename, distro.Edges.Packages[0].PythonVersion)
	fmt.Fprint(w, string(body))
}

func CreatePackage(ctx context.Context, client *ent.Client, anchor *html.Node) (*ent.PythonPackage, error) {
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
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating entity: %w", err)
	}
	//log.Println("entity was created: ", retval)
	return retval, nil
}

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

	anchors := find_all_children(node, "a")
	packages := make([]*ent.PythonPackage, len(anchors))
	for i, a := range anchors {
		temp, err := CreatePackage(context.Background(), client, a)
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

	//return &Distro{name: temp[len(temp)-1], packages: packages}, nil

	pd := client.PythonDistro.
		Create().
		SetName(temp[len(temp)-1])

	for _, p := range packages {
		pd.AddPackages(p)
	}
	retval, err := pd.Save(ctx)

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

func setupDB(memory bool) *ent.Client {
	var connection_string string
	if memory {
		connection_string = "file:ent?mode=memory&cache=shared&_fk=1"
	} else {
		connection_string = "file:metadata.db?cache=shared&_fk=1"
	}
	retval, err := ent.Open("sqlite3", connection_string)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	// Run the auto migration tool.
	if err := retval.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return retval
}
func main() {
	client = setupDB(false)
	defer client.Close()
	router := httprouter.New()
	router.GET("/simple/:library/", simple)
	router.GET("/simple/:library", simple)

	log.Println("Listing for requests at http://localhost:8000/simple")
	log.Fatal(http.ListenAndServe(":8000", router))

	log.Println("Done")
}
