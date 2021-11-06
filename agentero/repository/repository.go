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

func (r *Repository) GetById(agentId string) (phs []*protos.PolicyHolder, err error) {
	getPolicyHoldersSQL := `SELECT * FROM policy_holders`
	statement, err := r.db.Prepare(getPolicyHoldersSQL)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	rows, err := statement.Query()
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		fmt.Println("no policy holders ")
		return
	}

	for rows.Next() {
		ph := &protos.PolicyHolder{}
		rows.Scan(ph.Name, ph.MobileNumber, nil)
		phs = append(phs, ph)
	}

	getInsurancePoliciesByIdSQL := `SELECT * FROM insurance_policies WHERE agent_id = ?`
	statement, err = r.db.Prepare(getInsurancePoliciesByIdSQL)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	for i, v := range phs {
		rows, err = statement.Query(agentId)
		if err != nil {
			log.Fatalln(err.Error())
			return nil, err
		}
		defer rows.Close()

		if !rows.Next() {
			fmt.Println("no insurance policies with mobile number: ", v.MobileNumber)
			return
		}

		for rows.Next() {
			ip := &protos.InsurancePolicy{}
			rows.Scan(ip.MobileNumber, ip.Premium, ip.Type)
			phs[i].InsurancePolicy = append(phs[i].InsurancePolicy, ip)
		}
	}

	return
}

// TODO: Implement this
func (r *Repository) GetByMobileNumber(mobileNumber string) (ph *protos.PolicyHolder, err error) {
	// TODO: Duplicated code from here to, at least, the if after the defer. Extract to a new helper function
	getPolicyHoldersSQL := `SELECT * FROM policy_holders`
	statement, err := r.db.Prepare(getPolicyHoldersSQL)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	rows, err := statement.Query()
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		fmt.Println("no policy holders ")
		return nil, nil
	}

	for rows.Next() {
		phAux := &protos.PolicyHolder{}
		rows.Scan(phAux.Name, phAux.MobileNumber, nil)
		if phAux.MobileNumber == mobileNumber {
			ph = phAux
			break
		}
	}

	getInsurancePoliciesByMobileNumberSQL := `SELECT * FROM insurance_policies WHERE mobile_number = ?`
	statement, err = r.db.Prepare(getInsurancePoliciesByMobileNumberSQL)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	rows, err = statement.Query(mobileNumber)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		fmt.Println("no insurance policies with mobile number: ", v.MobileNumber)
		return
	}

	for rows.Next() {
		ip := &protos.InsurancePolicy{}
		rows.Scan(ip.MobileNumber, ip.Premium, ip.Type)
		ph.InsurancePolicy = append(ph.InsurancePolicy, ip)
	}

	return
}

func (r *Repository) GetAllInsuranceAgentsIds() (result []string, err error) {
	getAllInsuranceAgentsSQL := `SELECT * FROM insurance_agents`
	statement, err := r.db.Prepare(getAllInsuranceAgentsSQL)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	rows, err := statement.Query()
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		fmt.Println("no insurance agents ")
		return result, nil
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
		return err
	}

	_, err = statement.Exec(ph.Name, ph.MobileNumber, ph.Name)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	return nil
}

func (r *Repository) UpsertInsurancePolicy(ip *protos.InsurancePolicy, agentId string) error {
	insertInsurancePolicySQL := `INSERT INTO insurance_policies(ip_mobile_number, premium, type, agentId) VALUES (?, ?, ?, ?) ON CONFLICT(ip_mobile_number) DO UPDATE SET premium = ?, type = ?, agentId = ?`

	statement, err := r.db.Prepare(insertInsurancePolicySQL)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	_, err = statement.Exec(ip.MobileNumber, ip.Premium, ip.Type, agentId, ip.Premium, ip.Type, agentId)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	return nil
}

func (r *Repository) UpsertInsuranceAgent(agent *models.Agent) error {
	insertInsuranceAgentSQL := `INSERT INTO insurance_agents(agent_id, name) VALUES (?, ?) ON CONFLICT(agent_id) DO UPDATE SET name = ?`
	statement, err := r.db.Prepare(insertInsuranceAgentSQL)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	_, err = statement.Exec(agent.Id, agent.Name, agent.Id)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	return nil
}

func SetDatabaseUp() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}

	// Asking if the tables are already created so that we don't have a duplicate table error
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", "policy_holders", "insurance_policies", "insurance_agents").Scan(&count)
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

		_, err = db.Exec("CREATE UNIQUE INDEX `ip_UNIQUE` ON `insurance_policies`(`mobile_number`)")
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
