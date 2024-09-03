package mysql

import (
	"database/sql"

	"github.com/AIGPTku/api-aigptku.id/app/repository"
)

type repoMysql struct {
	master *sql.DB
	trx *sql.DB
}

func New(dbMaster, dbTrx *sql.DB) repository.MysqlInterface {
	return &repoMysql{
		master: dbMaster,
		trx: dbTrx,
	}
}