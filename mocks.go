package hansip

type dummySQL struct {
	queryRun, execRun, newTransactionRun bool
}

func (d *dummySQL) query(dest interface{}, query string, args ...interface{}) error {
	d.queryRun = true
	return nil
}

func (d *dummySQL) exec(query string, args ...interface{}) error {
	d.execRun = true
	return nil
}

func (d *dummySQL) newTransaction() (Transaction, error) {
	d.newTransactionRun = true
	return nil, nil
}
