package local

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func AlpineLocalUploadHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info().Msg(r.RequestURI)
			log.Info().Msg("Alpine Local Upload Handler")
			vars := mux.Vars(r)
			log.Info().Msgf("%s %s %s %s", vars["branch"], vars["repo"], vars["arch"], vars["pkgname"])
			log.Info().Msg("Handle PUT")
			f, err := os.Create("/tmp/uploads/" + vars["pkgname"])
			if err != nil {
					panic(err)
			}
			defer f.Close()
			if _, err := io.Copy(f, r.Body); err != nil {
					panic(err)
			}
			if err := f.Sync(); err != nil {
					panic(err)
			}
			path := fmt.Sprintf("%s/%s/%s/%s", vars["branch"], vars["repo"], vars["arch"], vars["pkgname"])
			log.Info().Msgf("Filename: %s", path)
			fmt.Fprintln(w, "upload done...")
	})
}
