CREATE DATABASE binaryrepo;

\c binaryrepo;

DROP TYPE IF EXISTS repotype CASCADE;
CREATE TYPE repotype AS ENUM ('local', 'remote', 'group');

DROP TYPE IF EXISTS pkgtype CASCADE;
CREATE TYPE pkgtype AS ENUM ('alpine', 'docker');

DROP TABLE IF EXISTS repo;
CREATE TABLE repo (
    id serial PRIMARY KEY,
    name text NULL,
    repotype repotype NULL,
    pkgtype pkgtype NULL,
    remote_url text NULL,
    username text NULL,
    password text NULL,
    UNIQUE (name),
    CONSTRAINT remote_must_have_url_constraint
        CHECK (CASE repotype WHEN 'remote' THEN remote_url IS NOT NULL
                             ELSE false END)
);

INSERT INTO repo (name, repotype, pkgtype, remote_url, username, password) VALUES ('alpine-local', 'local', 'alpine', NULL, NULL, NULL);

INSERT INTO repo (name, repotype, pkgtype, remote_url, username, password) VALUES ('alpine-remote', 'remote', 'alpine', NULL, NULL, NULL);

SELECT * FROM repo;