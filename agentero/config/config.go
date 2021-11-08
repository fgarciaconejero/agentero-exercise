package config

import "flag"

var (
	AmsUrlFlag    = flag.String("ams-api-url", "http://localhost:8081", "Fake AMS API's url")
	ClientUrlFlag = flag.String("client-url", "localhost:8080", "Client's url")
	ServerUrlFlag = flag.String("server-url", "localhost:50051", "Client's url")

	// Mac flag
	// 	AmsUrlFlag    = flag.String("ams-api-url", "http://127.0.0.1:8081", "Fake AMS API's url for Mac")
	// 	ClientUrlFlag = flag.String("client-url", "127.0.0.1:8080", "Client's url for Mac")
	// 	ServerUrlFlag = flag.String("server-url", "127.0.0.1:50051", "Client's url for Mac")

	SchedulePeriodFlag = flag.Int("schedule_period", 5, "Amount of time (in minutes) that the server will wait in order to import data from the Fake AMS API")
)
