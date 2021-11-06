package mocks

import "github.com/agentero-exercise/agentero/resources/protos"

var Users = []protos.PolicyHolder{
	{
		Name:         "user1",
		MobileNumber: "1234567890",
	},
	{
		Name:         "user2",
		MobileNumber: "123 456 7891",
	},
	{
		Name:         "user3",
		MobileNumber: "(123) 456 7892",
	},
	{
		Name:         "user4",
		MobileNumber: "(123) 456-7893",
	},
	{
		Name:         "user5",
		MobileNumber: "123-456-7894",
	},
}
