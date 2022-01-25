package registry

import (
	"log"
	"net/http"
)

// PathPutManifest URL.
const RegistryPathVersion = "/repo/{repo-name}/v2"

// VersionHandler implements GET baseURL/repo/v2/
func (registry *DockerRegistry) VersionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello")
}
