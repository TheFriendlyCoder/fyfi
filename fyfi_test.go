package main

import (
	"context"
	"os"
	"testing"
)

func TestBasicParsing(t *testing.T) {
	client = setupDB(true)
	defer client.Close()

	dat, err := os.ReadFile("testdata/distro_sample.html")
	if err != nil {
		t.Fatalf("Error loading sample data: %v", err)
	}

	distro, err := CreateDistro(context.Background(), client, string(dat))
	if err != nil {
		t.Fatalf("Parsing should not have failed: %v", err)
	}

	expected_name := "pip"
	if distro.Name != expected_name {
		t.Errorf("Distro name incorrect %s != %s", distro.Name, expected_name)
	}
	expected_package_count := 215
	allPackages, _ := distro.QueryPackages().All(context.Background())
	if len(allPackages) != expected_package_count {
		t.Fatalf("Number of packages incorrect %d != %d", len(allPackages), expected_package_count)
	}
	// if len(distro.Edges.Packages) != expected_package_count {
	// 	t.Fatalf("Number of packages incorrect %d != %d", len(distro.Edges.Packages), expected_package_count)
	// }

	// <a href="https://files.pythonhosted.org/packages/3d/9d/1e313763bdfb6a48977b65829c6ce2a43eaae29ea2f907c8bbef024a7219/pip-0.2.tar.gz#sha256=88bb8d029e1bf4acd0e04d300104b7440086f94cc1ce1c5c3c31e3293aee1f81">pip-0.2.tar.gz</a><br/>
	pck := allPackages[0]
	expected_package_name := "pip-0.2.tar.gz"
	expected_package_url := "https://files.pythonhosted.org/packages/3d/9d/1e313763bdfb6a48977b65829c6ce2a43eaae29ea2f907c8bbef024a7219/pip-0.2.tar.gz#sha256=88bb8d029e1bf4acd0e04d300104b7440086f94cc1ce1c5c3c31e3293aee1f81"
	expected_py_ver := ""
	expected_package_checksum := "88bb8d029e1bf4acd0e04d300104b7440086f94cc1ce1c5c3c31e3293aee1f81"
	if pck.URL != expected_package_url {
		t.Errorf("Package url incorrect %s != %s", pck.URL, expected_package_url)
	}
	if pck.Filename != expected_package_name {
		t.Errorf("Package filename incorrect %s != %s", pck.Filename, expected_package_name)
	}
	if pck.PythonVersion != expected_py_ver {
		t.Errorf("Package Python version incorrect %s != %s", pck.PythonVersion, expected_py_ver)
	}
	if pck.Checksum != expected_package_checksum {
		t.Errorf("Package checksum incorrect %s != %s", pck.Checksum, expected_package_checksum)
	}

	// <a href="https://files.pythonhosted.org/packages/4b/b6/0fa7aa968a9fa4ef63a51b3ff0644e59f49dcd7235b3fd6cceb23f202e08/pip-22.1.2.tar.gz#sha256=6d55b27e10f506312894a87ccc59f280136bad9061719fac9101bdad5a6bce69" data-requires-python="&gt;=3.7">pip-22.1.2.tar.gz</a><br/>
	//pck = distro.Edges.Packages[len(distro.Edges.Packages)-1]
	pck = allPackages[len(allPackages)-1]
	expected_package_name = "pip-22.1.2.tar.gz"
	expected_package_url = "https://files.pythonhosted.org/packages/4b/b6/0fa7aa968a9fa4ef63a51b3ff0644e59f49dcd7235b3fd6cceb23f202e08/pip-22.1.2.tar.gz#sha256=6d55b27e10f506312894a87ccc59f280136bad9061719fac9101bdad5a6bce69"
	expected_py_ver = ">=3.7"
	expected_package_checksum = "6d55b27e10f506312894a87ccc59f280136bad9061719fac9101bdad5a6bce69"
	if pck.URL != expected_package_url {
		t.Errorf("Package url incorrect %s != %s", pck.URL, expected_package_url)
	}
	if pck.Filename != expected_package_name {
		t.Errorf("Package filename incorrect %s != %s", pck.Filename, expected_package_name)
	}
	if pck.PythonVersion != expected_py_ver {
		t.Errorf("Package Python version incorrect %s != %s", pck.PythonVersion, expected_py_ver)
	}
	if pck.Checksum != expected_package_checksum {
		t.Errorf("Package checksum incorrect %s != %s", pck.Checksum, expected_package_checksum)
	}
}
