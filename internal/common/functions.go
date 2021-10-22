package common

import (
	"WebApplication/internal/logger"
	"crypto/md5"
	"encoding"
	"encoding/json"
	"github.com/jmoiron/sqlx"
)

func CreateSimpleHash(value interface{}) (string, error) {
	marshal, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	h := md5.New()
	h.Write(marshal)
	marshaller := h.(encoding.BinaryMarshaler)
	data, err := marshaller.MarshalBinary()
	return string(data), err
}

func FinalizeTransaction(tx *sqlx.Tx, hasErr error) (err error) {
	if hasErr == nil {
		err = tx.Commit()

	} else {
		err = tx.Rollback()
	}
	if err != nil {
		logger.Error(err)
	}
	return
}
