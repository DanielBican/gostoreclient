package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"sync"
)

const (
	addr              = "http://localhost:8080/v1"
	loginBodyTemplate = `{"Email": "%s","Password": "%s"}`
	product1          = `{"ProductName": "Product 1","Quantity": 1}`
	product2          = `{"ProductName": "Product 2","Quantity": 1}`
)

func main() {

	var wg sync.WaitGroup

	// Needed for saving response Set-Cookie header value
	jar1, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client1 := &http.Client{Jar: jar1}
	jar2, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client2 := &http.Client{Jar: jar2}

	wg.Add(2)
	go executeFlow(client1, "jon.doe@company.com", "1234", &wg)
	go executeFlow(client2, "jon.doe2@company.com", "1234", &wg)
	wg.Wait()
	fmt.Println("finished executing flows")
}

func executeFlow(client *http.Client, username string, password string, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Printf("executing flow for %s\n", username)

	post(client, "/login", fmt.Sprintf(loginBodyTemplate, username, password), "", "")
	post(client, "/cart", product1, username, password)
	post(client, "/cart", product1, username, password)
	post(client, "/cart", product2, username, password)
	post(client, "/checkout", "", username, password)
	// For adding logout to the flow enable following line
	// post(client, "/logout", "", username, password)
}

func post(client *http.Client, urlSuffix string, body string, username string, password string) {
	req, err := http.NewRequest("POST", addr+urlSuffix, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Panic(err)
	}
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s post to %s with body %s returned %d with body %s\n",
		username, addr+urlSuffix, body, resp.StatusCode, respBody)
}
