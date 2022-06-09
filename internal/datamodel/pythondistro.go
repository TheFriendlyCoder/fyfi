package datamodel

import (
	"context"
	"fmt"

	"github.com/TheFriendlyCoder/fyfi/ent"
	"github.com/TheFriendlyCoder/fyfi/internal/pypi"
)

// SaveDistro serializes the data from a PythonDistribution struct so
// it can be stored in a metadata database
func SaveDistro(ctx context.Context, client *ent.Client, distro *pypi.PythonDistribution) error {

	for _, p := range distro.Packages {
		_, err := client.PythonPackage.
			Create().
			SetChecksum(p.Checksum).
			SetURL(p.URL).
			SetFilename(p.Filename).
			SetPythonVersion(p.PythonVersion).
			SetDistro(distro.Name).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed creating entity: %w", err)
		}
	}
	return nil
}
