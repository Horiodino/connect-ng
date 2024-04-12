package zypper

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SUSE/connect-ng/internal/util"
)

func TestParseProductsXML(t *testing.T) {
	products, err := parseProductsXML(util.ReadTestFile("products.xml", t))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(products) != 2 {
		t.Errorf("Expected len()==2. Got %d", len(products))
	}
	if products[0].ToTriplet() != "SUSE-MicroOS/5.0/x86_64" {
		t.Errorf("Expected SUSE-MicroOS/5.0/x86_64 Got %s", products[0].ToTriplet())
	}
}

func TestParseProductsXML_ReleaseType(t *testing.T) {
	zypperFilesystemRoot = t.TempDir()
	// write test OEM file
	testRT := "SLES-OEM-TEST"
	oemDir := filepath.Join(zypperFilesystemRoot, oemPath)
	if err := os.MkdirAll(oemDir, 0755); err != nil {
		t.Fatalf("Error creating OEM dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(oemDir, "sles"), []byte(testRT), 0600); err != nil {
		t.Fatalf("Error writing OEM file: %v", err)
	}

	products, err := parseProductsXML(util.ReadTestFile("products-rt.xml", t))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if products[0].ReleaseType != "" {
		t.Errorf("Expected empty ReleaseType Got %v", products[0].ReleaseType)
	}
	if products[1].ReleaseType != testRT {
		t.Errorf("Expected ReleaseType=%v Got %v", testRT, products[1].ReleaseType)
	}
	if products[2].ReleaseType != testRT {
		t.Errorf("Expected ReleaseType=%v Got %v", testRT, products[2].ReleaseType)
	}
	if products[3].ReleaseType != "rel2" {
		t.Errorf("Expected ReleaseType=rel2 Got %v", products[3].ReleaseType)
	}
	if products[4].ReleaseType != "" {
		t.Errorf("Expected empty ReleaseType Got %v", products[4].ReleaseType)
	}
}

func TestParseServicesXML(t *testing.T) {
	services, err := parseServicesXML(util.ReadTestFile("services.xml", t))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(services) != 1 {
		t.Errorf("Expected len()==1. Got %d", len(services))
	}
	if services[0].Name != "SUSE_Linux_Enterprise_Micro_5.0_x86_64" {
		t.Errorf("Expected SUSE_Linux_Enterprise_Micro_5.0_x86_64 Got %s", services[0].Name)
	}
}

func TestParseReposXML(t *testing.T) {
	repos, err := parseReposXML(util.ReadTestFile("repos.xml", t))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(repos) != 3 {
		t.Errorf("Expected len()==3. Got %v", len(repos))
	}
	if repos[0].Name != "SLE-Module-Basesystem15-SP2-Pool" {
		t.Errorf("Expected SLE-Module-Basesystem15-SP2-Pool. Got %v", repos[0].Name)
	}
	if repos[0].Priority != 99 {
		t.Errorf("Expected priority 99. Got %v", repos[0].Priority)
	}
	if !repos[0].Enabled {
		t.Errorf("Expected Enabled. Got %v", repos[0].Enabled)
	}
	if repos[1].Priority != 50 {
		t.Errorf("Expected priority 99. Got %v", repos[1].Priority)
	}
	if repos[1].Enabled {
		t.Errorf("Expected not Enabled Got %v", repos[1].Enabled)
	}
}

func TestInstalledProducts(t *testing.T) {
	util.Execute = func(_ []string, _ []int) ([]byte, error) {
		return util.ReadTestFile("products.xml", t), nil
	}

	products, err := InstalledProducts()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(products) != 2 {
		t.Errorf("Expected len()==2. Got %d", len(products))
	}
	if products[0].ToTriplet() != "SUSE-MicroOS/5.0/x86_64" {
		t.Errorf("Expected SUSE-MicroOS/5.0/x86_64 Got %s", products[0].ToTriplet())
	}
}

func TestBaseProduct(t *testing.T) {
	util.Execute = func(_ []string, _ []int) ([]byte, error) {
		return util.ReadTestFile("products.xml", t), nil
	}

	base, err := BaseProduct()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if base.ToTriplet() != "SUSE-MicroOS/5.0/x86_64" {
		t.Errorf("Expected SUSE-MicroOS/5.0/x86_64 Got %s", base.ToTriplet())
	}
}

func TestBaseProductError(t *testing.T) {
	util.Execute = func(_ []string, _ []int) ([]byte, error) {
		return util.ReadTestFile("products-no-base.xml", t), nil
	}
	_, err := BaseProduct()
	if err != ErrCannotDetectBaseProduct {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestParseSearchResultXML(t *testing.T) {
	packages, err := parseSearchResultXML(util.ReadTestFile("product-search.xml", t))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(packages) != 2 {
		t.Errorf("Expected len()==2. Got %v", len(packages))
	}
	if packages[0].Name != "SLES" {
		t.Errorf("Expected SLES. Got %v", packages[0].Name)
	}
	if packages[0].Edition != "15.2-0" {
		t.Errorf("Expected edition 15.2-0. Got %v", packages[0].Edition)
	}
	if packages[0].Repo != "SLE-Product-SLES15-SP2-Updates" {
		t.Errorf("Expected SLE-Product-SLES15-SP2-Updates. Got %v", packages[0].Repo)
	}
	if packages[1].Edition != "15.2-0" {
		t.Errorf("Expected edition 15.2-0. Got %v", packages[1].Edition)
	}
	if packages[1].Repo != "SLE-Product-SLES15-SP2-Pool" {
		t.Errorf("Expected SLE-Product-SLES15-SP2-Pool. Got %v", packages[1].Repo)
	}
}
