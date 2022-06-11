package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicParsing(t *testing.T) {
	dat, err := NewConfigData("testdata/sample.yml")
	assert.Nil(t, err)
	assert.Equal(t, dat.CacheFolder, "/path/to/cache")
}

func TestMissingCacheFolder(t *testing.T) {
	dat, err := NewConfigData("testdata/missing_cache_folder.yml")
	assert.NotNil(t, err)
	assert.Nil(t, dat)
	assert.Contains(t, err.Error(), "cache_folder")
}
