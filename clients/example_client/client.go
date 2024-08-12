package exampleclient

import (
	"net/http"
)

func GetExample() (http.Response, error) {
    uri := "https://icanhazdadjoke.com"

    client := &http.Client{}

    req, _ := http.NewRequest(http.MethodGet, uri, nil)
    req.Header.Add("Accept", "application/json")

    res, err := client.Do(req)

    return *res, err
}
