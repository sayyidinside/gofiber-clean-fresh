package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type TxRepository interface {
	Transaction(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(ctx context.Context, tx *gorm.DB) error
	RollbackTx(ctx context.Context, tx *gorm.DB) error
	DeferTx(ctx context.Context, tx *gorm.DB, retErr *error)
}

type txRepository struct {
	*gorm.DB
}

func NewTxRepository(db *gorm.DB) TxRepository {
	return &txRepository{DB: db}
}

// Wraper for transaction
func (r *txRepository) Transaction(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error {
	// GORM's Transaction method handles begin, commit, and rollback automatically.
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Execute the business logic function inside GORM's transaction wrapper.
		return fn(ctx, tx)
	})
}

// Manualy start transaction
func (r *txRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tx, nil
}

// Manualy commit transaction
func (r *txRepository) CommitTx(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		return errors.New("no transaction provided")
	}

	return tx.WithContext(ctx).Commit().Error
}

// Manualy rollback transaction
func (r *txRepository) RollbackTx(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		return errors.New("no transaction provided")
	}

	// GORM returns ErrInvalidTransaction if transaction already committed/rolled back.
	if err := tx.WithContext(ctx).Rollback().Error; err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
		return err
	}

	return nil
}

// DeferTx inspects the named return error pointer.
//
// Usage pattern: func (...) (err error) { tx, err := r.BeginTx(ctx); if err != nil { return err }
//
//	defer r.DeferTx(ctx, tx, &err)
//
// ... do work that may set err ...
func (r *txRepository) DeferTx(ctx context.Context, tx *gorm.DB, retErr *error) {
	if tx == nil {
		// nothing to do
		return
	}

	if p := recover(); p != nil {
		_ = r.RollbackTx(ctx, tx)
		panic(p)
	}

	if retErr != nil && *retErr != nil {
		_ = r.RollbackTx(ctx, tx)
		return
	}

	if cerr := r.CommitTx(ctx, tx); cerr != nil {
		// If commit fails, propagate into the named return error if possible
		if retErr != nil && *retErr == nil {
			*retErr = cerr
		}
	}
}
