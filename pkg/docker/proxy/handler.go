package dockerproxy

import (
	"net/http"

	"github.com/gorilla/mux"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
)

func RegisterHandlers(r *mux.Router, fs *filestore.FileStore, repoIndex *repo.RepoIndex) {
	////log.Info().Msgf("Registering docker proxy handlers")
	p := NewDockerProxyApp()
	p.index = repoIndex
	r.HandleFunc(PathVersionUrl1, p.VersionHandler).Methods(http.MethodGet)
	r.HandleFunc(PathVersionUrl2, p.VersionHandler).Methods(http.MethodGet)
	r.HandleFunc(PathHeadManifest1, p.HeadManifestHandler).Methods(http.MethodHead)
	r.HandleFunc(PathHeadManifest2, p.HeadManifestHandler).Methods(http.MethodHead)
	r.HandleFunc(PathGetManifest1, p.GetManifestHandler).Methods(http.MethodGet)
	r.HandleFunc(PathGetManifest2, p.GetManifestHandler).Methods(http.MethodGet)
	r.HandleFunc(PathGetBlob1, p.DownloadLayer).Methods(http.MethodGet)
	r.HandleFunc(PathGetBlob2, p.DownloadLayer).Methods(http.MethodGet)
	r.HandleFunc(PathHeadBlob1, p.HasLayer).Methods(http.MethodHead)
	r.HandleFunc(PathHeadBlob2, p.HasLayer).Methods(http.MethodHead)
	r.HandleFunc(PathHeadBlob1, p.LayerPut).Methods(http.MethodPut)
	r.HandleFunc(PathHeadBlob2, p.LayerPut).Methods(http.MethodPut)

}
