package postgres

import "github.com/jackc/tern/migrate"

var schemaMigrations = []migrate.Migration{
	{
		Name: "1",
		UpSQL: `--
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
	id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
	username text NOT NULL,
	password text NOT NULL,
    UNIQUE (username)
);

INSERT INTO users (username, password) VALUES ('admin', (SELECT crypt('admin',gen_salt('bf')) ));

CREATE TYPE repotype AS ENUM ('local', 'remote', 'group');
CREATE TYPE pkgtype AS ENUM ('alpine', 'docker');

CREATE TABLE repo (
    id serial PRIMARY KEY,
    name text NULL,
    repotype repotype NULL,
    pkgtype pkgtype NULL,
    remote_url text NULL,
    remote_username text NULL,
    remote_password text NULL,
	anonymous boolean NOT NULL DEFAULT false,
    UNIQUE (name)
);

CREATE TABLE artifact_node (
    id serial PRIMARY KEY,
    repo_id integer REFERENCES repo (id),
    name text NOT NULL,
    path text NULL,
    upstream_url text NULL,
    etag text NULL,
    checksum text NULL,
    is_folder boolean NOT NULL DEFAULT false,
    parent_id integer REFERENCES artifact_node (id),
    UNIQUE(name),
    UNIQUE(path)
);

CREATE VIEW artifact_node_v AS (
    SELECT node.id, repo.name as repo_name, node.repo_id, node.name as node_name, node.path, node.upstream_url, node.etag, node.checksum, node.is_folder, node.parent_id FROM artifact_node as node
    INNER JOIN repo ON repo.id = node.repo_id
);

CREATE INDEX artifact_node_path_idx ON artifact_node(path);
CREATE INDEX artifact_node_etag_idx ON artifact_node(etag);
CREATE INDEX artifact_node_checksum_idx ON artifact_node(checksum);

CREATE TABLE artifact_property (
    id serial PRIMARY KEY,
    artifact_id integer REFERENCES artifact_node (id),
    key text NOT NULL,
    value text NOT NULL
);
`,
	},
}
