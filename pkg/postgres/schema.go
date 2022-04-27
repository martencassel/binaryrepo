package binaryrepodb

import (
	"github.com/jackc/tern/migrate"
)

var schemaMigrations = []migrate.Migration{
	{
		Name: "1",
		UpSQL: `
CREATE TABLE repos (
	repo_id SERIAL,

	PRIMARY KEY(id)

);`,
	},
}

