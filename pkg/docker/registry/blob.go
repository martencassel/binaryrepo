package registry

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//	"github.com/google/uuid"

	//	"github.com/google/uuid"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// PathInitBlobUpload URL.
const PathInitBlobUpload = "/repo/{repo-name}/v2/{name}/blobs/upload"

// Initiate a resumable blob upload.
// POST /repo/{repo-name}/v2/<name>/blobs/uploads.
func (registry *DockerRegistry) InitBlobUpload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	repoName := vars["repo-name"]
	if registry.index.FindRepo(repoName) == nil || name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("%s /v2/%s/blobs/uploads", http.MethodPost, name)
	uuid, _ := uuid.NewUUID()
	ioutil.WriteFile(fmt.Sprintf("%s/uploads/%s", registry.fs.BasePath, uuid.String()), []byte{}, 0644)
	w.Header().Set("Content-Length", "0")
	w.Header().Set("docker-distribution-api-version", "registry/2.0")
	w.Header().Set("Docker-Upload-UUID", uuid.String())
	w.Header().Set("Range", "0-0")
	loc := fmt.Sprintf("/repo/%s/v2/%s/blobs/uploads/%s", repoName, name, uuid)
	w.Header().Set("Location", loc)
	w.Header().Set("Connection", "close")
	w.WriteHeader(http.StatusAccepted)
}

/* PathVersion URL.
const PathGetBlob = "/repo/{repo-name}/v2/{name}/blobs/{digest}"


	GET /repo/{repo-name}/v2/<name>/blobs/<digest>
	returns a blob with a digest from a repo.

func (registry *DockerRegistry) GetBlob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	d := vars["digest"]
	repoName := vars["repo-name"]
	if registry.index.FindRepo(repoName) == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("GET /repo/%s/v2/%s/blobs/%s", repoName, name, d)
	dgst, err := digest.Parse(d)
	if err != nil {
		log.Printf("Digest is invalid %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if registry.fs.Exists(dgst) {
		w.WriteHeader(http.StatusOK)
		b, err := registry.fs.ReadFile(dgst)
		if err != nil {
			log.Printf("Error reading file %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
		w.Header().Set("Docker-Content-Digest", dgst.String())
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// PathHeadBlob URL.
const PathHeadBlob = "/repo/{repo-name}/v2/{name}/blobs/{digest}"

	HEAD /repo/{repo-name}/v2/<name>/blobs/<digest> should return
	blob length and digest of blob exists, otherwise not found.
func (registry *DockerRegistry) HeadBlob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	d := vars["digest"]
	repoName := vars["repo-name"]
	if registry.index.FindRepo(repoName) == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("HEAD /repo/%s/v2/%s/blobs/%s", repoName, name, d)
	dgst, err := digest.Parse(d)
	if err != nil {
		log.Printf("Digest is invalid %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exists := registry.fs.Exists(dgst)
	if exists {
		w.WriteHeader(http.StatusOK)
		b, err := registry.fs.ReadFile(dgst)
		if err != nil {
			log.Printf("Error reading file %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
		w.Header().Set("Docker-Content-Digest", dgst.String())
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

*/

/*
// PathInitBlobUpload URL.
const PathChunkedBlobUpload = "/repo/{repo-name}/v2/{name}/blobs/uploads/{uuid}"

	Upload chunk. Append chunk into the uploads directory in the <uuid> file.

	PATCH /repo/<repo-name>/v2/<name>/blobs/uploads/<uuid>
func (registry *DockerRegistry) UploadBlobChunk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	repoName := vars["repo-name"]
	uuid := vars["uuid"]
	if registry.index.FindRepo(repoName) == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if uuid == "" {
		log.Printf("UUID is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if name == "" {
		log.Printf("Name is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("%s /v2/%s/blobs/uploads", http.MethodPatch, name)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	start := r.Header.Get("Content-Range")
	if start == "" {
		log.Printf("Content-Range header is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	start = strings.Split(start, "-")[0]
	start_offset, err := strconv.Atoi(start)
	if err != nil {
		log.Printf("Content-Range header is invalid %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path := fmt.Sprintf("%s/uploads/%s", registry.fs.BasePath, uuid)
	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ret, err := f.Seek(int64(start_offset), io.SeekStart)
	if err != nil {
		log.Printf("Error seeking file %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ret != int64(start_offset) {
		log.Printf("Error seeking file %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = f.Write(b)
	if err != nil {
		log.Printf("Error writing file %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Location", fmt.Sprintf("/repo/%s/v2/%s/blobs/uploads/%s", repoName, name, uuid))
	w.Header().Set("Range", fmt.Sprintf("%d-%d", 0, len(b)))
	w.Header().Set("Content-Length", "0")
	w.Header().Set("Docker-Upload-UUID", uuid)
}

	Completed upload using upload uuid, returns status accepted
	PUT /v2/<name>/blobs/uploads/<uuid>
func (registry *DockerRegistry) CompleteUpload(w http.ResponseWriter, r *http.Request) {

}
*/
