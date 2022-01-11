package main

import (
	"os"

	"github.com/martencassel/binaryrepo/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting.....")

}
func main() {
	cmd.Execute()
}
