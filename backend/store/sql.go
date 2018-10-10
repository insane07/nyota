package store

import (
	"nyota/backend/logutil"
	"nyota/backend/model"

	gorp "gopkg.in/gorp.v2"
)

// txExecFunc can perform all SQLs to be executed in transaction.
type TxExecFunc func(*gorp.Transaction) error

// sqlExecTx creates a transaction and invokes exec under it.
func execTx(s *model.SessionContext, db SqlDB, exec TxExecFunc) error {
	tx, err := db.Begin()
	if err != nil {
		logutil.Errorf(s, "tx begin(%v)", err)
		return err
	}

	if err := exec(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		logutil.Errorf(s, "tx commit(%v)", err)
		return err
	}

	return nil
}
