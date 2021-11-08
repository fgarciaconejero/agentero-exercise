# agentero-exercise

## Hi! My name is Facu, thank you for taking a look at this exercise 

This exercise consists of **two** services, the Agentero service and the AMS (Fake) service.    
In order to try it, you'll need to run the AMS client, then the Agentero client and server.    
It's recommended to run the AMS client **first** given that the server will retrieve information from it at start up. However, as the server retrieves the information every 5 minutes, it's not necessary, but, if you don't do it in that order, you may have to wait for the server to retrieve the information before you receive anything when making calls to the Agentero client.

1° - Use the following command to run the Fake AMS API client: `go run .\ams\client`
2° - Use the following command to run the Agentero server: `go run .\agentero\server\server.go`
3° - Use the following command to run the Agentero client: `go run .\agentero\client\client.go`
