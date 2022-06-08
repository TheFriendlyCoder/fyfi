package pypi

import (
	"os"
	"testing"
)

func TestBasicParsing(t *testing.T) {
	dat, err := os.ReadFile("testdata/distro_sample.html")
	if err != nil {
		t.Fatalf("Error loading sample data: %v", err)
	}

	heading, anchors, err := ParseDistribution(string(dat))
	if err != nil {
		t.Fatalf("Parsing should not have failed: %v", err)
	}

	expected_heading := "Links for pip"
	if heading.FirstChild.Data != expected_heading {
		t.Errorf("Parsed heading incorrect %s != %s", heading.FirstChild.Data, expected_heading)
	}

	expected_package_count := 215
	if len(anchors) != expected_package_count {
		t.Fatalf("Number of packages incorrect %d != %d", len(anchors), expected_package_count)
	}

	// <a href="https://files.pythonhosted.org/packages/3d/9d/1e313763bdfb6a48977b65829c6ce2a43eaae29ea2f907c8bbef024a7219/pip-0.2.tar.gz#sha256=88bb8d029e1bf4acd0e04d300104b7440086f94cc1ce1c5c3c31e3293aee1f81">pip-0.2.tar.gz</a><br/>
	attributes := getAttributes(anchors[0])
	expected_package_name := "pip-0.2.tar.gz"
	expected_package_url := "https://files.pythonhosted.org/packages/3d/9d/1e313763bdfb6a48977b65829c6ce2a43eaae29ea2f907c8bbef024a7219/pip-0.2.tar.gz#sha256=88bb8d029e1bf4acd0e04d300104b7440086f94cc1ce1c5c3c31e3293aee1f81"
	expected_py_ver := ""
	if attributes["href"] != expected_package_url {
		t.Errorf("Package url incorrect %s != %s", attributes["href"], expected_package_url)
	}
	if anchors[0].FirstChild.Data != expected_package_name {
		t.Errorf("Package filename incorrect %s != %s", anchors[0].FirstChild.Data, expected_package_name)
	}
	if attributes["data-requires-python"] != expected_py_ver {
		t.Errorf("Package Python version incorrect %s != %s", attributes["data-requires-python"], expected_py_ver)
	}

	// <a href="https://files.pythonhosted.org/packages/4b/b6/0fa7aa968a9fa4ef63a51b3ff0644e59f49dcd7235b3fd6cceb23f202e08/pip-22.1.2.tar.gz#sha256=6d55b27e10f506312894a87ccc59f280136bad9061719fac9101bdad5a6bce69" data-requires-python="&gt;=3.7">pip-22.1.2.tar.gz</a><br/>
	anchor := anchors[len(anchors)-1]
	attributes = getAttributes(anchor)
	expected_package_name = "pip-22.1.2.tar.gz"
	expected_package_url = "https://files.pythonhosted.org/packages/4b/b6/0fa7aa968a9fa4ef63a51b3ff0644e59f49dcd7235b3fd6cceb23f202e08/pip-22.1.2.tar.gz#sha256=6d55b27e10f506312894a87ccc59f280136bad9061719fac9101bdad5a6bce69"
	expected_py_ver = ">=3.7"
	if attributes["href"] != expected_package_url {
		t.Errorf("Package url incorrect %s != %s", attributes["href"], expected_package_url)
	}
	if anchor.FirstChild.Data != expected_package_name {
		t.Errorf("Package filename incorrect %s != %s", anchor.FirstChild.Data, expected_package_name)
	}
	if attributes["data-requires-python"] != expected_py_ver {
		t.Errorf("Package Python version incorrect %s != %s", attributes["data-requires-python"], expected_py_ver)
	}
}
