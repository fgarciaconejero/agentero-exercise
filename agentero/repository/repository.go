package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/agentero-exercise/agentero/domain/models"
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

func (r *Repository) GetAllInsuranceAgentsIds() (result []string, err error) {
	getAllInsuranceAgentsSQL := `SELECT * FROM insurance_agents;`
	statement, err := r.db.Prepare(getAllInsuranceAgentsSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	rows, err := statement.Query()
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer rows.Close()

	if !rows.Next() {
		fmt.Println("no rows ")
	}
	for rows.Next() {
		agentId := ""
		rows.Scan(agentId, nil)
		result = append(result, agentId)
	}

	return result, nil
}

func (r *Repository) UpsertPolicyHolder(ph *protos.PolicyHolder) error {
	insertPolicyHolderSQL := `INSERT INTO policy_holders(name, ph_mobile_number) VALUES (?, ?) ON CONFLICT(ph_mobile_number) DO UPDATE SET name = ?`

	statement, err := r.db.Prepare(insertPolicyHolderSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(ph.Name, ph.MobileNumber, ph.Name)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return nil
}

func (r *Repository) UpsertInsurancePolicy(ip *protos.InsurancePolicy) error {
	insertInsurancePolicySQL := `INSERT INTO insurance_policies(ip_mobile_number, premium, type) VALUES (?, ?, ?) ON CONFLICT(ip_mobile_number) DO UPDATE SET premium = ?, type = ?`

	statement, err := r.db.Prepare(insertInsurancePolicySQL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(ip.MobileNumber, ip.Premium, ip.Type, ip.Premium, ip.Type)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return nil
}

func (r *Repository) UpsertInsuranceAgent(agents *models.Agent) error {
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
			"(`name` TEXT, `ph_mobile_number` TEXT, PRIMARY KEY (`mobile_number`))")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE UNIQUE INDEX `ph_UNIQUE` ON `policy_holders`(`mobile_number`)")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE TABLE `insurance_policies`" +
			"(`ip_mobile_number` TEXT, `premium` integer, `type` TEXT, `agentId` TEXT, PRIMARY KEY (`mobile_number`))")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE UNIQUE INDEX `ip_UNIQUE` ON `policy_holders`(`mobile_number`)")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE TABLE `insurance_agents`" +
			"(`agent_id` TEXT, `name` TEXT, PRIMARY KEY (`agent_id`))")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE UNIQUE INDEX `agent_id_UNIQUE` ON `insurance_agents`(`agent_id`)")
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
