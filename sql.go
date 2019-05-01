package hansip

import (
	"fmt"
	"runtime"
)

// sql exposes methods needed to execute query
type sql interface {
	query(dest interface{}, query string, args ...interface{}) error
	exec(query string, args ...interface{}) error
	newTransaction() (Transaction, error)
}

// Transaction represents an sql transaction.
// Transactions are always guaranteed to run in master connection.
type Transaction interface {
	Query(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) error
	Commit() error
	Rollback() error
}

func injectCallerInfo(sql string) string {
	pc, file, line, ok := runtime.Caller(3)
	details := runtime.FuncForPC(pc)
	if !ok || details == nil {
		return sql
	}

	msg := fmt.Sprintf("/* %s at %s:%d */", details.Name(), file, line)
	return fmt.Sprintf("%s\n%s", msg, sql)
}
