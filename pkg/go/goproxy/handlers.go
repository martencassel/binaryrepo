package goproxy

import (
	"fmt"
	"net/http"
)

const ModuleList = "{module:.+}/@v/list"

func GoModuleList(w http.ResponseWriter, req *http.Request) {
	// Check if path exists in cache
	// If not, fetch from remote, cache and serve
	fmt.Println("Hello, World!")
}

const ModuleLatest = "{module:.+}/@v/latest"

func GoModuleLatest(w http.ResponseWriter, req *http.Request) {
	// Check if path exists in cache
	// If not, fetch from remote, cache and serve
	fmt.Println("Hello, World!")
}

const ModuleVersion = "{module:.+}/@v/version"

func GoModuleVersion(w http.ResponseWriter, req *http.Request) {
	// Check if path exists in cache
	// If not, fetch from remote, cache and serve
	fmt.Println("Hello, World!")
}

const ModuleInfo = "{module:.+}/@v/{version}.info"

func GoModuleInfo(w http.ResponseWriter, req *http.Request) {
	// Check if path exists in cache
	// If not, fetch from remote, cache and serve
	fmt.Println("Hello, World!")
}

const ModuleMod = "{module:.+}/@v/{version}.mod"

func GoModuleMod(w http.ResponseWriter, req *http.Request) {
	// Check if path exists in cache
	// If not, fetch from remote, cache and serve
	fmt.Println("Hello, World!")
}

const ModuleZip = "{module:.+}/@v/{version}.zip"

func GoModuleZip(w http.ResponseWriter, req *http.Request) {
	// Check if path exists in cache
	// If not, fetch from remote, cache and serve
	fmt.Println("Hello, World!")
}

const HelmIndex = "index.yaml"
const HelmPkg = "{name}-{version}.tgz"
