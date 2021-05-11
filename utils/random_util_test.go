package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetRandomNum(t *testing.T) {
	url := "http://8.131.110.43:8080/flowable-ui/app/rest/content/3f5edd40-48ee-11eb-9b43-f6dae279c576/raw/"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "FLOWABLE_REMEMBER_ME=QkRZRzB1TkpFQnNtNEFVSXIzSlRoUSUzRCUzRDpWdEpHbTRvSWlQN1RtQ29UN2YwNk1BJTNEJTNE")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
