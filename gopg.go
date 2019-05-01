package hansip

import (
	"github.com/go-pg/pg"
)

type gopgSQL struct {
	db *pg.DB
}

func (s *gopgSQL) query(dest interface{}, query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	_, err := s.db.Query(dest, query, args...)
	return err
}

func (s *gopgSQL) exec(query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *gopgSQL) newTransaction() (Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return &gopgTransaction{db: tx}, nil
}

type gopgTransaction struct {
	db       *pg.Tx
	finished bool
}

func (tx *gopgTransaction) Query(dest interface{}, query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	_, err := tx.db.Query(dest, query, args...)
	return err
}

func (tx *gopgTransaction) Exec(query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	_, err := tx.db.Exec(query, args...)
	return err
}

func (tx *gopgTransaction) Commit() error {
	if tx.finished {
		return ErrTxFinished
	}
	err := tx.db.Commit()
	tx.finished = true
	return err
}

func (tx *gopgTransaction) Rollback() error {
	if tx.finished {
		return ErrTxFinished
	}
	err := tx.db.Rollback()
	tx.finished = true
	return err
}
