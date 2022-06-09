package pypi

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicParsing(t *testing.T) {
	dat, err := os.ReadFile("testdata/distro_sample.html")
	if err != nil {
		t.Fatalf("Error loading sample data: %v", err)
	}

	distro, err := ParseDistribution(string(dat))
	if err != nil {
		t.Fatalf("Parsing should not have failed: %v", err)
	}

	assert.Equal(t, "pip", distro.Name)
	assert.Equal(t, 215, len(distro.Packages))

	// <a href="https://files.pythonhosted.org/packages/3d/9d/1e313763bdfb6a48977b65829c6ce2a43eaae29ea2f907c8bbef024a7219/pip-0.2.tar.gz#sha256=88bb8d029e1bf4acd0e04d300104b7440086f94cc1ce1c5c3c31e3293aee1f81">pip-0.2.tar.gz</a><br/>
	pck := distro.Packages[0]
	expected_package_name := "pip-0.2.tar.gz"
	expected_package_url := "https://files.pythonhosted.org/packages/3d/9d/1e313763bdfb6a48977b65829c6ce2a43eaae29ea2f907c8bbef024a7219/pip-0.2.tar.gz#sha256=88bb8d029e1bf4acd0e04d300104b7440086f94cc1ce1c5c3c31e3293aee1f81"
	expected_py_ver := ""
	expected_checksum := "sha256=88bb8d029e1bf4acd0e04d300104b7440086f94cc1ce1c5c3c31e3293aee1f81"
	assert.Equal(t, expected_package_url, pck.URL)
	assert.Equal(t, expected_package_name, pck.Filename)
	assert.Equal(t, expected_py_ver, pck.PythonVersion)
	assert.Equal(t, expected_checksum, pck.Checksum)

	// <a href="https://files.pythonhosted.org/packages/4b/b6/0fa7aa968a9fa4ef63a51b3ff0644e59f49dcd7235b3fd6cceb23f202e08/pip-22.1.2.tar.gz#sha256=6d55b27e10f506312894a87ccc59f280136bad9061719fac9101bdad5a6bce69" data-requires-python="&gt;=3.7">pip-22.1.2.tar.gz</a><br/>
	pck = distro.Packages[len(distro.Packages)-1]
	expected_package_name = "pip-22.1.2.tar.gz"
	expected_package_url = "https://files.pythonhosted.org/packages/4b/b6/0fa7aa968a9fa4ef63a51b3ff0644e59f49dcd7235b3fd6cceb23f202e08/pip-22.1.2.tar.gz#sha256=6d55b27e10f506312894a87ccc59f280136bad9061719fac9101bdad5a6bce69"
	expected_py_ver = ">=3.7"
	expected_checksum = "sha256=6d55b27e10f506312894a87ccc59f280136bad9061719fac9101bdad5a6bce69"
	assert.Equal(t, expected_package_url, pck.URL)
	assert.Equal(t, expected_package_name, pck.Filename)
	assert.Equal(t, expected_py_ver, pck.PythonVersion)
	assert.Equal(t, expected_checksum, pck.Checksum)

}
