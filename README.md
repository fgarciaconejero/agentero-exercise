# agentero-exercise

## Hi! My name is Facu, thank you for taking a look at this exercise 

This exercise consists of **two** services, the Agentero service and the AMS (Fake) service.

In order to try it, you'll need to run the AMS client, then the Agentero client and server.

It's recommended to run the AMS client **first** given that the server will retrieve information from it at start up. However, as the server retrieves the information every 5 minutes, it's not necessary, but, if you don't do it in that order, you may have to wait for the server to retrieve the information before you receive anything when making calls to the Agentero client.

## Prerequisites
You must have Go installed and it's advised that the Go extension for VS Code is installed too. The rest of the dependecies should be installed by themselves when running the services for the first time.

1° **If you are in Mac**, please go to the file `./agentero/config/config.go`, **comment** the url-related variables that are **uncommented** and then **uncomment** the ones that were previously **commented**

2° - Use the following command to run the Fake AMS API client: `go run .\ams\client.go` (or `go run ./ams/client.go` if you are using MacOS)

3° - Use the following command to run the Agentero server: `go run .\agentero\server\server.go` (or `go run ./agentero/server/server.go` if you are using MacOS)

4° - Use the following command to run the Agentero client: `go run .\agentero\client\client.go` (or `go run ./agentero/client/client.go` if you are using MacOS)

5° You are ready to go! You should be able to send requests to the different endpoints successfully now!

To get the policy holders (and their insurance policies) related to an agent, send a request to `/getById/:agentId`  
__Although you can try with whatever agent id you like, the available agents' ids are: "1", "2", "3" (unless you were to add more in `./ams/resources/mocks/agents_mock.go`)__

To get a specific policy holder (and its insurance policies) send a request to `/getByMobileNumber/:mobileNumber`  
__Although you can try whatever mobile number you like, the ones that will retrieve data are: "1234567890", "1234567891", "1234567892", "1234567893", "1234567894" (unless you were to add more in `./ams/resources/mocks/policies_mock.go` and `./ams/resources/mocks/users_mock.go`)__