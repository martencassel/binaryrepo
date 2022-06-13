package remote

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

/*
	Serve HTTP GET request for either for
     	* branch/repo/arch/APKINDEX.tar.gz
		* branch/repo/arch/pkgname.tar.gz
	* Check if the file exists in the local cache
		* Check for a newer version of the requested file in the upstream repository
		  If found, then download the new version and update the local cache
		  Serve the file from the local cache
	* If not found, then download the file from the upstream repository and update the local cache
	* If the file is not found in the upstream repository, then return a 404
*/
func AlpineGetHandler() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		log.Info().Msg(r.RequestURI)
		log.Info().Msg("Alpine remote get handler")
		vars := mux.Vars(r)
		log.Info().Msg("Handle GET")
		log.Info().Msgf("%s %s %s %s", vars["branch"], vars["repo"], vars["arch"], vars["filename"])
	})
}