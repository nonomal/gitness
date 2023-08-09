// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by the Polyform Free Trial License
// that can be found in the LICENSE.md file for this repository.

package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/harness/gitness/encrypt"
	"github.com/harness/gitness/internal/store"
	gitness_store "github.com/harness/gitness/store"
	"github.com/harness/gitness/store/database"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.SecretStore = (*secretStore)(nil)

const (
	secretQueryBase = `
		SELECT
		secret_id,
		secret_description,
		secret_space_id,
		secret_uid,
		secret_data,
		secret_created,
		secret_updated,
		secret_version
		FROM secrets
		`

	secretColumns = `
	secret_id,
	secret_description,
	secret_space_id,
	secret_uid,
	secret_data,
	secret_created,
	secret_updated,
	secret_version
	`

	secretInsertStmt = `
	INSERT INTO secrets (
		secret_description,
		secret_space_id,
		secret_uid,
		secret_data,
		secret_created,
		secret_updated,
		secret_version
	) VALUES (
		:secret_description,
		:secret_space_id,
		:secret_uid,
		:secret_data,
		:secret_created,
		:secret_updated,
		:secret_version
	) RETURNING secret_id`

	secretUpdateStmt = `
	UPDATE secrets
	SET
		secret_description = :secret_description,
		secret_space_id = :secret_space_id,
		secret_uid = :secret_uid,
		secret_data = :secret_data,
		secret_created = :secret_created,
		secret_updated = :secret_updated,
		secret_version = :secret_version
	WHERE secret_id = :secret_id AND secret_version = :secret_version - 1`
)

// NewSecretStore returns a new SecretStore.
func NewSecretStore(enc encrypt.Encrypter, db *sqlx.DB) *secretStore {
	return &secretStore{
		db:  db,
		enc: enc,
	}
}

type secretStore struct {
	db  *sqlx.DB
	enc encrypt.Encrypter
}

// Find returns a secret given a secret ID
func (s *secretStore) Find(ctx context.Context, id int64) (*types.Secret, error) {
	const findQueryStmt = secretQueryBase + `
		WHERE secret_id = $1`
	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(types.Secret)
	if err := db.GetContext(ctx, dst, findQueryStmt, id); err != nil {
		return nil, database.ProcessSQLErrorf(err, "Failed to find secret")
	}
	return dec(s.enc, dst)
}

// FindByUID returns a secret in a given space with a given UID
func (s *secretStore) FindByUID(ctx context.Context, spaceID int64, uid string) (*types.Secret, error) {
	const findQueryStmt = secretQueryBase + `
		WHERE secret_space_id = $1 AND secret_uid = $2`
	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(types.Secret)
	if err := db.GetContext(ctx, dst, findQueryStmt, spaceID, uid); err != nil {
		return nil, database.ProcessSQLErrorf(err, "Failed to find secret")
	}
	return dec(s.enc, dst)
}

// Create creates a secret
func (s *secretStore) Create(ctx context.Context, secret *types.Secret) error {
	db := dbtx.GetAccessor(ctx, s.db)

	secret, err := enc(s.enc, secret)
	if err != nil {
		return err
	}

	query, arg, err := db.BindNamed(secretInsertStmt, secret)
	if err != nil {
		return database.ProcessSQLErrorf(err, "Failed to bind secret object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&secret.ID); err != nil {
		return database.ProcessSQLErrorf(err, "secret query failed")
	}

	return nil
}

func (s *secretStore) Update(ctx context.Context, secret *types.Secret) (*types.Secret, error) {
	updatedAt := time.Now()

	secret.Version++
	secret.Updated = updatedAt.UnixMilli()

	db := dbtx.GetAccessor(ctx, s.db)

	secret, err := enc(s.enc, secret)
	if err != nil {
		return nil, err
	}

	query, arg, err := db.BindNamed(secretUpdateStmt, secret)
	if err != nil {
		return nil, database.ProcessSQLErrorf(err, "Failed to bind secret object")
	}

	result, err := db.ExecContext(ctx, query, arg...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(err, "Failed to update secret")
	}

	count, err := result.RowsAffected()
	if err != nil {
		return nil, database.ProcessSQLErrorf(err, "Failed to get number of updated rows")
	}

	if count == 0 {
		return nil, gitness_store.ErrVersionConflict
	}

	return secret, nil

}

// List lists all the secrets present in a space
func (s *secretStore) List(ctx context.Context, parentID int64, opts *types.SecretFilter) ([]types.Secret, error) {
	stmt := database.Builder.
		Select(secretColumns).
		From("secrets").
		Where("secret_space_id = ?", fmt.Sprint(parentID))

	if opts.Query != "" {
		stmt = stmt.Where("LOWER(secret_uid) LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(opts.Query)))
	}

	stmt = stmt.Limit(database.Limit(opts.Size))
	stmt = stmt.Offset(database.Offset(opts.Page, opts.Size))

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	dst := []types.Secret{}
	if err = db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(err, "Failed executing custom list query")
	}

	return dst, nil
}

// Delete deletes a secret given a secret ID
func (s *secretStore) Delete(ctx context.Context, id int64) error {
	const secretDeleteStmt = `
		DELETE FROM secrets
		WHERE secret_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, secretDeleteStmt, id); err != nil {
		return database.ProcessSQLErrorf(err, "Could not delete secret")
	}

	return nil
}

// DeleteByUID deletes a secret with a given UID in a space
func (s *secretStore) DeleteByUID(ctx context.Context, spaceID int64, uid string) error {
	const secretDeleteStmt = `
	DELETE FROM secrets
	WHERE secret_space_id = $1 AND secret_uid = $2`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, secretDeleteStmt, spaceID, uid); err != nil {
		return database.ProcessSQLErrorf(err, "Could not delete secret")
	}

	return nil
}

// Count of secrets in a space.
func (s *secretStore) Count(ctx context.Context, parentID int64, opts *types.SecretFilter) (int64, error) {
	stmt := database.Builder.
		Select("count(*)").
		From("secrets").
		Where("secret_space_id = ?", parentID)

	if opts.Query != "" {
		stmt = stmt.Where("secret_uid LIKE ?", fmt.Sprintf("%%%s%%", opts.Query))
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var count int64
	err = db.QueryRowContext(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, database.ProcessSQLErrorf(err, "Failed executing count query")
	}
	return count, nil
}

// helper function returns the same secret with encrypted data
func enc(encrypt encrypt.Encrypter, secret *types.Secret) (*types.Secret, error) {
	s := *secret
	ciphertext, err := encrypt.Encrypt(secret.Data)
	if err != nil {
		return nil, err
	}
	s.Data = string(ciphertext)
	return &s, nil
}

// helper function returns the same secret with decrypted data
func dec(encrypt encrypt.Encrypter, secret *types.Secret) (*types.Secret, error) {
	s := *secret
	plaintext, err := encrypt.Decrypt([]byte(secret.Data))
	if err != nil {
		return nil, err
	}
	s.Data = string(plaintext)
	return &s, nil
}
