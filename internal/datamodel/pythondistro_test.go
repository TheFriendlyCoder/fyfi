package datamodel

import (
	"context"
	"testing"

	"github.com/TheFriendlyCoder/fyfi/ent"
	"github.com/TheFriendlyCoder/fyfi/internal/pypi"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestBasicSerialization(t *testing.T) {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed creating schema resources: %v", err)
	}
	pkg := pypi.PythonPackage{
		Checksum:      "sha256=abcd1234",
		PythonVersion: ">=2.7",
		Filename:      "sample.whl",
		URL:           "http://fubar.com/sample",
	}

	distro := pypi.PythonDistribution{
		Name:     "sample",
		Packages: make([]pypi.PythonPackage, 1),
	}
	distro.Packages[0] = pkg
	SaveDistro(context.Background(), client, &distro)

	items, err := client.PythonPackage.Query().All(context.Background())
	if err != nil {
		t.Fatalf("failed querying todos: %v", err)
	}
	assert.Equal(t, 1, len(items))
	assert.Equal(t, pkg.Filename, items[0].Filename)
	assert.Equal(t, pkg.Checksum, items[0].Checksum)
	assert.Equal(t, pkg.PythonVersion, items[0].PythonVersion)
	assert.Equal(t, pkg.URL, items[0].URL)

}
