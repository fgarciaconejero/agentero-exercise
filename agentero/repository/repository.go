package repository

import (
	"database/sql"
	"fmt"

	"github.com/agentero-exercise/agentero/domain/models"
	"github.com/agentero-exercise/agentero/resources/protos"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	Db sql.DB
}

func NewRepository() (*Repository, error) {
	db, err := SetDatabaseUp()
	if err != nil {
		return nil, err
	}

	return &Repository{
		Db: *db,
	}, nil
}

func (r *Repository) GetById(agentId string) (phs []*protos.PolicyHolder, err error) {
	getPolicyHoldersSQL := `SELECT * FROM policy_holders`
	rows, err := r.Db.Query(getPolicyHoldersSQL)
	if err != nil {
		fmt.Println("There was a problem while trying to get policy holders from SQLite,", err)
		return
	}

	for rows.Next() {
		ph := &protos.PolicyHolder{}
		err = rows.Scan(&ph.Name, &ph.MobileNumber)
		if err != nil {
			fmt.Println("There was an error scanning policy holders:", err.Error())
			return
		}
		phs = append(phs, ph)
	}

	getInsurancePoliciesByIdSQL := `SELECT * FROM insurance_policies WHERE agent_id = ?`

	for i, v := range phs {
		rows, err = r.Db.Query(getInsurancePoliciesByIdSQL, agentId)
		if err != nil {
			return
		}

		for rows.Next() {
			discardId := ""
			ip := &protos.InsurancePolicy{}
			err = rows.Scan(&discardId, &ip.MobileNumber, &ip.Premium, &ip.Type, &ip.AgentId)
			if err != nil {
				fmt.Println("There was an error scanning insurance policies:", err.Error())
				return
			}
			if v.MobileNumber == ip.MobileNumber {
				phs[i].InsurancePolicy = append(phs[i].InsurancePolicy, ip)
			}
		}
	}

	phs = filterOutUnmatchedPolicyHolders(phs)

	return
}

func (r *Repository) GetByMobileNumber(mobileNumber string) (ph *protos.PolicyHolder, err error) {
	getPolicyHoldersSQL := `SELECT * FROM policy_holders`
	rows, err := r.Db.Query(getPolicyHoldersSQL)
	if err != nil {
		fmt.Println("There was a problem while trying to get policy holders from SQLite,", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		phAux := &protos.PolicyHolder{}
		err = rows.Scan(&phAux.Name, &phAux.MobileNumber)
		if err != nil {
			fmt.Println("There was an error scanning policy holders:", err.Error())
			return
		}
		if phAux.MobileNumber == mobileNumber {
			ph = phAux
			break
		}
	}

	getInsurancePoliciesByMobileNumberSQL := `SELECT * FROM insurance_policies WHERE ip_mobile_number = ?`
	statement, err := r.Db.Prepare(getInsurancePoliciesByMobileNumberSQL)
	if err != nil {
		fmt.Println("There was a problem preparing the getInsurancePoliciesByMobileNumberSQL statement,", err)
		return nil, err
	}

	rows, err = statement.Query(mobileNumber)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		discardId := ""
		ip := &protos.InsurancePolicy{}
		err = rows.Scan(&discardId, &ip.MobileNumber, &ip.Premium, &ip.Type, &ip.AgentId)
		if err != nil {
			fmt.Println("There was an error scanning insurance policies:", err.Error())
			return
		}
		if mobileNumber == ip.MobileNumber {
			ph.InsurancePolicy = append(ph.InsurancePolicy, ip)
		}
	}

	return
}

func (r *Repository) GetAllInsuranceAgentsIds() (result []string, err error) {
	getAllInsuranceAgentsSQL := `SELECT * FROM insurance_agents`
	rows, err := r.Db.Query(getAllInsuranceAgentsSQL)
	if err != nil {
		fmt.Println("There was an error getting insurance agents", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		agent := &models.Agent{}
		err = rows.Scan(&agent.Id, &agent.Name)
		if err != nil {
			fmt.Println("There was an error scanning insurance agents:", err.Error())
			return
		}
		result = append(result, agent.Id)
	}

	return
}

func (r *Repository) UpsertPolicyHolder(ph *protos.PolicyHolder) (err error) {
	insertPolicyHolderSQL := `INSERT INTO policy_holders(name, ph_mobile_number) VALUES (?, ?) ON CONFLICT(ph_mobile_number) DO UPDATE SET name = ?`
	_, err = r.Db.Exec(insertPolicyHolderSQL, ph.Name, ph.MobileNumber, ph.Name)
	if err != nil {
		fmt.Println("There was a problem executing the insertPolicyHolderSQL statement,", err)
		return err
	}

	return nil
}

func (r *Repository) UpsertInsurancePolicy(ip *protos.InsurancePolicy, agentId string) (err error) {
	insertInsurancePolicySQL := `INSERT INTO insurance_policies(ip_mobile_number, premium, type, agent_id) VALUES (?, ?, ?, ?) ON CONFLICT(ip_id) DO UPDATE SET ip_mobile_number = ?, premium = ?, type = ?, agent_id = ?`
	_, err = r.Db.Exec(insertInsurancePolicySQL, ip.MobileNumber, ip.Premium, ip.Type, agentId, ip.MobileNumber, ip.Premium, ip.Type, agentId)
	if err != nil {
		fmt.Println("There was a problem executing the insertInsurancePolicySQL statement,", err)
		return err
	}

	return nil
}

func (r *Repository) UpsertInsuranceAgent(agent *models.Agent) (err error) {
	insertInsuranceAgentSQL := `INSERT INTO insurance_agents(agent_id, name) VALUES (?, ?) ON CONFLICT(agent_id) DO UPDATE SET name = ?`
	_, err = r.Db.Exec(insertInsuranceAgentSQL, agent.Id, agent.Name, agent.Id)
	if err != nil {
		fmt.Println("There was a problem executing the insertInsuranceAgentSQL statement,", err)
		return
	}

	return
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
			"(`name` TEXT, `ph_mobile_number` TEXT, PRIMARY KEY (`ph_mobile_number`))")
		if err != nil {
			fmt.Println("There was a problem while creating policy_holders table, ", err)
			return nil, err
		}

		_, err = db.Exec("CREATE UNIQUE INDEX `ph_UNIQUE` ON `policy_holders`(`ph_mobile_number`)")
		if err != nil {
			fmt.Println("There was a problem while creating ph_UNIQUE index, ", err)
			return nil, err
		}

		_, err = db.Exec("CREATE TABLE `insurance_policies`" +
			"(`ip_id` integer, `ip_mobile_number` TEXT, `premium` integer, `type` TEXT, `agent_id` TEXT, PRIMARY KEY (`ip_id`))")
		if err != nil {
			fmt.Println("There was a problem while creating insurance_policies table, ", err)
			return nil, err
		}

		_, err = db.Exec("CREATE UNIQUE INDEX `ip_UNIQUE` ON `insurance_policies`(`ip_id`)")
		if err != nil {
			fmt.Println("There was a problem while creating ip_UNIQUE index, ", err)
			return nil, err
		}

		_, err = db.Exec("CREATE TABLE `insurance_agents`" +
			"(`agent_id` TEXT, `name` TEXT, PRIMARY KEY (`agent_id`))")
		if err != nil {
			fmt.Println("There was a problem while creating agent_id table, ", err)
			return nil, err
		}

		_, err = db.Exec("CREATE UNIQUE INDEX `agent_id_UNIQUE` ON `insurance_agents`(`agent_id`)")
		if err != nil {
			fmt.Println("There was a problem while creating agent_id_UNIQUE index, ", err)
			return nil, err
		}
	}

	return db, nil
}

func filterOutUnmatchedPolicyHolders(phs []*protos.PolicyHolder) (result []*protos.PolicyHolder) {
	for _, v := range phs {
		if len(v.InsurancePolicy) != 0 {
			result = append(result, v)
		}
	}
	return
}
