package exampleservice

import (
	exampleclient "employee-worklog-service/clients/example_client"
	"employee-worklog-service/models"
	"encoding/json"
)

func GetExample() string {
    r, _ := exampleclient.GetExample()

    var j models.Response
    jErr := json.NewDecoder(r.Body).Decode(&j)
    if jErr != nil {
        return "Service up, but client down?"
    } else {
        return j.Joke
    }
}
