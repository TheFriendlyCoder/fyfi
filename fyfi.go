package main

// Reference: https://peps.python.org/pep-0503/
import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/TheFriendlyCoder/fyfi/ent"
	"github.com/TheFriendlyCoder/fyfi/internal/datamodel"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
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

	distro, err := datamodel.CreateDistro(context.Background(), client, string(body))
	if err != nil {
		log.Printf("Error parsing response data: %v\n", err)
		return
	}

	log.Printf("Loading distro %s\n", distro.Name)
	// log.Printf("Found %d packages for distro %s\n", len(distro.Edges.Packages), distro.Name)
	// log.Printf("First package %s has pyver %s\n", distro.Edges.Packages[0].Filename, distro.Edges.Packages[0].PythonVersion)
	fmt.Fprint(w, string(body))
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
