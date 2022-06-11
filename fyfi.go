package main

// Reference: https://peps.python.org/pep-0503/
import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/TheFriendlyCoder/fyfi/ent"
	"github.com/TheFriendlyCoder/fyfi/internal/configuration"
	"github.com/TheFriendlyCoder/fyfi/internal/datamodel"
	"github.com/TheFriendlyCoder/fyfi/internal/pypi"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

const configFilePath string = "sample.yml"

type AppSettings struct {
	Client *ent.Client
	Config *configuration.ConfigData
}

func simple(w http.ResponseWriter, r *http.Request) {
	settings, ok := r.Context().Value("settings").(AppSettings)

	if !ok {
		fmt.Println("Settings not correct type")
	}

	ps, ok := r.Context().Value("params").(httprouter.Params)
	if !ok {
		fmt.Println("Params not correct type")
	}

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

	temp, err := pypi.ParseDistribution(string(body))
	if err != nil {
		log.Println("Failed to parse distro data")
		return
	}
	err = datamodel.SaveDistro(context.Background(), settings.Client, temp)
	if err != nil {
		log.Printf("Error parsing response data: %v\n", err)
		return
	}

	log.Printf("Loading distro %s\n", temp.Name)
	// log.Printf("Found %d packages for distro %s\n", len(distro.Edges.Packages), distro.Name)
	// log.Printf("First package %s has pyver %s\n", distro.Edges.Packages[0].Filename, distro.Edges.Packages[0].PythonVersion)
	fmt.Fprint(w, string(body))
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Take the context out from the request
		ctx := r.Context()
		ctx = context.WithValue(ctx, "params", ps)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}

func settingsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config, err := configure()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(config.CacheFolder)

		var settings AppSettings
		settings.Client = setupDB(false)
		settings.Config = config
		defer settings.Client.Close()

		ctx := r.Context()
		ctx = context.WithValue(ctx, "settings", settings)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
func setupDB(memory bool) *ent.Client {
	var connection_string string
	if memory {
		connection_string = "file:ent?mode=memory&cache=shared&_fk=1"
	} else {
		os.Remove("metadata.db")
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

func configure() (*configuration.ConfigData, error) {

	config, err := configuration.NewConfigData(configFilePath)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(config.CacheFolder, 0700)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	router := httprouter.New()
	simpleReceiver := settingsMiddleware(http.HandlerFunc(simple))
	router.GET("/simple/:library/", wrapHandler(simpleReceiver))
	router.GET("/simple/:library", wrapHandler(simpleReceiver))

	log.Println("Listing for requests at http://localhost:8000/simple")
	log.Fatal(http.ListenAndServe(":8000", router))

	log.Println("Done")
}
