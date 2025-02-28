package repo

import (
	"redirektor/server/pkg/psql"

	"gorm.io/gorm"
)

type PostgresClient struct {
	DB *gorm.DB
}

func NewPostgresClient() *PostgresClient {
	db, err := psql.New()
	if err != nil {
		panic(err)
	}
	ret := new(PostgresClient)
	ret.DB = db

	return ret
}

func (pc *PostgresClient) selfIfNil(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return pc.DB
}
