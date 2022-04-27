package binaryrepodb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4"
	domain "github.com/martencassel/binaryrepo/pkg/domain"
)

type tnxSession struct {
	pgx.Tx
	ctx context.Context
}

func (t *tnxSession) Commit() error {
	return t.Tx.Commit(t.ctx)
}

func (t *tnxSession) Rollback() error {
	return t.Tx.Rollback(t.ctx)
}

func (t *tnxSession) InsertRepo(repo *domain.PackageRepository) error {
	q := build().Insert("repos")
	q = q.Columns("repos", "repo_id")

	if err := t.doInsert(q); err != nil {
		return fmt.Errorf("failed to insert repo  %w", err)
	}

	return nil
}

func scanRepos(rows pgx.Rows) (repos []domain.PackageRepository, err error) {
	defer rows.Close()

	for rows.Next() {
		var r domain.PackageRepository
		if err := rows.Scan(&r); err != nil {
			return nil, errors.Wrap(err, "scan: failed to read repo column")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func (t *tnxSession) ListRepos() (repos []domain.PackageRepository, err error) {
	q := build().Select().Columns("repo_id").From("repos").OrderBy("repo_id");
	rows, err := t.doQuery(q)
	if err != nil {
		return nil, err
	}
	return scanRepos(rows)
}

func (t *tnxSession) GetRepo(id int64) (*domain.PackageRepository, error) {
	q := build().Select().Columns("repo_id").From("repos").Where(sq.Eq{"repo_id": id});
	row, err := t.doQueryRow(q)
	if err != nil {
		return nil, err
	}

	var r domain.PackageRepository
	if err := row.Scan(&r); err != nil {
		return nil, errors.Wrap(err, "scan: failed to read repo column")
	}

	return &r, nil
}

func (t *tnxSession) doQuery(q sq.SelectBuilder) (pgx.Rows, error) {
	sqlStr, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query to sql conversion: %v", err)
	}
	return t.Query(t.ctx, sqlStr, args...)
}

func (t *tnxSession) doQueryRow(q sq.SelectBuilder) (pgx.Row, error) {
	sqlStr, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query to sql conversion: %v", err)
	}
	return t.Query(t.ctx, sqlStr, args...)
}

func (t *tnxSession) doInsert(q sq.InsertBuilder) error {
	sqlStr, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("insert statement to sql conversion failed: %v", err)
	}
	_, err = t.Exec(t.ctx, sqlStr, args...)
	return err
}

func (t *tnxSession) doUpdate( q sq.UpdateBuilder) error {
	sqlStr, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("update statement to sql conversion failed: %v", err)
	}
	_, err = t.Exec(t.ctx, sqlStr, args...)
	return err
}

type txnFunc func(tnxSession) error

func inTxn(txn *tnxSession, fn func(txn *tnxSession) error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			txn.Rollback()
			panic(r)
		}
		if err != nil {
			txn.Rollback()
			return
		}
		err = txn.Commit()
	}()

	if err := fn(txn); err != nil {
		return err
	}
	return txn.Commit()
}

func build() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}