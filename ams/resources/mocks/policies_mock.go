package mocks

import "github.com/agentero-exercise/agentero/resources/protos"

var Policies = []protos.InsurancePolicy{
	{
		MobileNumber: "1234567890",
		Premium:      2000,
		Type:         "homeowner",
		AgentId:      "1",
	},
	{
		MobileNumber: "123 456 7891",
		Premium:      200,
		Type:         "renter",
		AgentId:      "2",
	},
	{
		MobileNumber: "123-456-7892",
		Premium:      1500,
		Type:         "homeowner",
		AgentId:      "1",
	},
	{
		MobileNumber: "(123) 456-7893",
		Premium:      155,
		Type:         "personal_auto",
		AgentId:      "3",
	},
	{
		MobileNumber: "123-456-7894",
		Premium:      1000,
		Type:         "homeowner",
		AgentId:      "1",
	},
	{
		MobileNumber: "123-456-7890",
		Premium:      500,
		Type:         "personal_auto",
		AgentId:      "1",
	},
	{
		MobileNumber: "1234567892",
		Premium:      100,
		Type:         "life",
		AgentId:      "2",
	},
	{
		MobileNumber: "(123)456-7892",
		Premium:      200,
		Type:         "homeowner",
		AgentId:      "1",
	},
}
