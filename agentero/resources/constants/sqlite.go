package constants

const GetPolicyHoldersSQL string = `SELECT * FROM policy_holders`
const GetInsurancePoliciesByIdSQL string = `SELECT * FROM insurance_policies WHERE agent_id = ?`
const GetInsurancePoliciesByMobileNumberSQL string = `SELECT * FROM insurance_policies WHERE ip_mobile_number = ?`
const GetAllInsuranceAgentsSQL string = `SELECT * FROM insurance_agents`
const InsertPolicyHolderSQL string = `INSERT INTO policy_holders(name, ph_mobile_number) VALUES (?, ?) ON CONFLICT(ph_mobile_number) DO UPDATE SET name = ?`
const InsertInsurancePolicySQL string = `INSERT INTO insurance_policies(ip_mobile_number, premium, type, agent_id) VALUES (?, ?, ?, ?) ON CONFLICT(ip_id) DO UPDATE SET ip_mobile_number = ?, premium = ?, type = ?, agent_id = ?`
const InsertInsuranceAgentSQL string = `INSERT INTO insurance_agents(agent_id, name) VALUES (?, ?) ON CONFLICT(agent_id) DO UPDATE SET name = ?`
