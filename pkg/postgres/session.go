package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	binaryrepo "github.com/martencassel/binaryrepo"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

const (
	MigrationTable = "testdb_schema_migrations"
)

type txnSession struct {
	pgx.Tx
	ctx context.Context
}

func (d *Database) Begin(ctx context.Context) (txn *txnSession, err error) {
	tx, err := d.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &txnSession{ Tx: tx, ctx: ctx}, nil
}

func (t *txnSession) ListRepos() (repos []binaryrepo.Repo, err error) {
	log.Println("ListRepos")
	q := build().Select().Columns("name", "repotype", "pkgtype", "remote_url", "remote_username", "remote_password", "anonymous").From("repo");
	log.Println(q)
	rows, err := t.doQuery(q)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return scanRepos(rows)
}

func scanRepos(rows pgx.Rows) (repos []binaryrepo.Repo, err error) {
	defer rows.Close()

	for rows.Next() {
		var r binaryrepo.Repo

		var repoType sql.NullString
		var pkgType sql.NullString
		var remoteUrl sql.NullString
		var remote_username sql.NullString
		var remote_password sql.NullString
		var anonymous sql.NullBool
		if err := rows.Scan(&r.Name, &repoType, &pkgType, &remoteUrl, &remote_username, &remote_password, &anonymous); err != nil {
			return nil, errors.Wrap(err, "scan failed to read repo column")
		}
		if repoType.Valid {
			r.Repotype = repoType.String
		}
		if pkgType.Valid {
			r.Pkgtype = pkgType.String
		}
		if remoteUrl.Valid {
			r.Remoteurl = remoteUrl.String
		}
		if remote_username.Valid {
			r.RemoteUsername = remote_username.String
		}
		if remote_password.Valid {
			r.RemotePassword = remote_password.String
		}
		if anonymous.Valid {
			r.Anonymous = anonymous.Bool
		}
		repos = append(repos, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return repos, nil
}

func (t *txnSession) InsertRepo(repo binaryrepo.Repo) error {
	log.Printf("InsertRepo\n")
	log.Println(repo)
	q := build().Insert("repo")
	q = q.Columns("name", "repotype", "pkgtype", "remote_url", "remote_username", "remote_password", "anonymous")
	q = q.Values(repo.Name, repo.Repotype, repo.Pkgtype, repo.Remoteurl, repo.RemoteUsername, repo.RemotePassword, repo.Anonymous)
 		q = q.Suffix("ON CONFLICT DO NOTHING")
	if err := t.doInsert(q); err != nil {
		return fmt.Errorf("failed to insert repo: %w", err)
	}
	return nil
}

func (t *txnSession) LookupRepo(repoName string) (bool, error) {
	q := build().Select().Columns("*").From("artifact_node_v").Where("repo_name = ?", repoName)
	row, err := t.doQueryRow(q)
	if err != nil {
		return false, err
	}
	var repo string
	err = row.Scan(&repo)
	if err != pgx.ErrNoRows {
		return true, nil
	}
	return false, nil
}


func (t *txnSession) UpdateRepo(repo binaryrepo.Repo) error {
	log.Printf("updateRepo\n")
	log.Println(repo)
	q := build().Update("repo")
	q = q.Set("repotype", repo.Repotype)
	q = q.Set("pkgtype", repo.Pkgtype)
	q = q.Set("remote_url", repo.Remoteurl)
	q = q.Set("remote_username", repo.RemoteUsername)
	q = q.Set("remote_password", repo.RemotePassword)
	q = q.Set("anonymous", repo.Anonymous)
	q = q.Where("name = ?", repo.Name)
	if err := t.doUpdate(q); err != nil {
		return fmt.Errorf("failed to update repo: %w", err)
	}
	return nil
}

func (t *txnSession) CreateNode(node binaryrepo.Node) error {
	log.Printf("CreateNode\n")
	q := build().Insert("artifact_node")
	q = q.Columns("repo_id", "name", "path", "upstream_url", "etag", "checksum", "is_folder", "parent_id")
	q = q.Values(node.RepoID, node.NodeName, node.Path, node.UpstreamUrl, node.ETag, node.Checksum, node.IsFolder, node.ParentID)
	q = q.Suffix("ON CONFLICT DO NOTHING")
	if err := t.doInsert(q); err != nil {
		return fmt.Errorf("failed to insert node: %w", err)
	}
	return nil
}

func (t *txnSession) GetNodeByDigest(digest digest.Digest) (node *binaryrepo.Node, err error) {
	// Build select query and execute it.
	return nil, nil
}

func (t *txnSession) UpdateNode(node binaryrepo.Node) error {
	log.Printf("UpdateNode\n")
	log.Println(node)
	q := build().Update("repo")
	q = q.Set("repo_id", node.RepoID)
	q = q.Set("name", node.NodeName)
	q = q.Set("path", node.Path)
	q = q.Set("upstream_url", node.UpstreamUrl)
	q = q.Set("etag", node.ETag)
	q = q.Set("checksum", node.Checksum)
	q = q.Set("is_folder", node.IsFolder)
	q = q.Set("parent_id", node.ParentID)
	q = q.Where("id = ?", node.ID)
	if err := t.doUpdate(q); err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}
	return nil
}

func (t *txnSession) DeleteNode(node binaryrepo.Node) error {
	log.Printf("DeleteNode\n")
	log.Println(node)
	q := build().Delete("artifact_node").Where("id = ?", node.ID)
	if err := t.doDelete(q); err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}
	return nil
}

func (t *txnSession) GetNode(repo string, path string) (node binaryrepo.Node, err error) {
	q := build().Select().Columns("*").From("artifact_node_v").Where("repo_name = ? and path = ?", repo, path)
	row, err := t.doQueryRow(q)
	if err != nil {
		return node, err
	}
	err = row.Scan(&node.ID, &node.RepoName, &node.RepoID, &node.NodeName, &node.Path, &node.UpstreamUrl, &node.ETag, &node.Checksum, &node.IsFolder, &node.ParentID)
	return
}


func (t *txnSession) GetUser(username string) (user binaryrepo.User, err error) {
	q := build().Select().Columns("username").From("users").Where("username = ?", username)
	row, err := t.doQueryRow(q)
	if err != nil {
		return user, err
	}
	err = row.Scan(&user.Username, &user.Password)
	return
}

func (t *txnSession) LookupUser(user binaryrepo.User) (exists bool, err error) {
	log.Println(user.Username)
	q := build().Select()
	q = q.Columns("username")
	q = q.From("users")
	q = q.Where(sq.Expr("username = ? AND password = crypt(?, password)", user.Username, user.Password))
	row, err := t.doQueryRow(q)
	if err != nil {
		return false, err
	}
	var username string
	err = row.Scan(&username)
	if err != pgx.ErrNoRows {
		return true, nil
	}
	return false, nil
}

func (t *txnSession) InsertUser(user binaryrepo.User) error {
	q := build().Insert("users")
	q = q.Columns("username", "password")
	q = q.Values(user.Username, sq.Expr("crypt(?,gen_salt('bf'))", user.Password))
	q = q.Suffix("ON CONFLICT DO NOTHING")
	if err := t.doInsert(q); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (t *txnSession) GetRepo(name string) (repo binaryrepo.Repo, err error) {
	q := build().Select().Columns("name").From("repo").Where("name = ?", name)
	row, err := t.doQueryRow(q)
	if err != nil {
		return repo, err
	}
	err = row.Scan(&repo.ID, &repo.Name, &repo.Pkgtype, &repo.Repotype)
	return
}

func (t *txnSession) doQueryRow(q sq.SelectBuilder) (pgx.Row, error) {
	sqlStr, args, err := q.ToSql()
	log.Println(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("query to sql: %w", err)
	}
	return t.QueryRow(t.ctx, sqlStr, args...), nil
}

func (t *txnSession) doInsert(q sq.InsertBuilder) error {
	sqlStr, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("insert to sql: %w", err)
	}
	log.Println(sqlStr)
	_, err = t.Exec(t.ctx, sqlStr, args...);
	log.Println(err)
	return err
}

func (t *txnSession) doUpdate(q sq.UpdateBuilder) error {
	sqlStr, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("update to sql: %w", err)
	}
	log.Println(sqlStr)
	_, err = t.Exec(t.ctx, sqlStr, args...);
	log.Println(err)
	return err
}



func (t *txnSession) doDelete(q sq.DeleteBuilder) error {
	sqlStr, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("delete to sql: %w", err)
	}
	log.Println(sqlStr)
	_, err = t.Exec(t.ctx, sqlStr, args...);
	log.Println(err)
	return err
}
