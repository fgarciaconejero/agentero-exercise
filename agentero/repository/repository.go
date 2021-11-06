package repository

import (
	"database/sql"

	"github.com/agentero-exercise/agentero/resources/protos"
)

type Repository struct {
	db sql.DB
}

func NewRepository() (*Repository, error) {
	db, err := SetDatabaseUp()
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: *db,
	}, nil
}

func (r *Repository) GetById(agentId string) ([]*protos.PolicyHolder, error) {
	return nil, nil
}

func (r *Repository) GetByMobileNumber(agentId string) (*protos.PolicyHolder, error) {
	return nil, nil
}

func (r *Repository) Upsert(phs []*protos.PolicyHolder) error {
	return nil
}
func SetDatabaseUp() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}

	// Asking if the tables are already created so that we don't have a duplicate table error
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", "policy_holders", "insurance_policies").Scan(&count)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		_, err = db.Exec("CREATE TABLE `policy_holders`" +
			"(`name` string, `mobile_number` string, PRIMARY KEY `mobile_number`)")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE UNIQUE INDEX `mobile_number_UNIQUE` ON `users`(`mobile_number`)")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE TABLE `insurance_policies`" +
			"(`mobile_number` string, `premium` integer, `type` string, PRIMARY KEY `mobile_number`)")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE UNIQUE INDEX `mobile_number_UNIQUE` ON `users`(`mobile_number`)")
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
