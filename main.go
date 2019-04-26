package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	url, err := getURLInput(reader)
	if err != nil {
		panic(err)
	}

	body, err := getURLResponseBody(url)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("dat1", body, 0644)
	if err != nil {
		panic(err)
	}
}

func getURLInput(r *bufio.Reader) (string, error) {
	fmt.Print("Please enter the url you would like to download from then hit enter: \n")

	delimiter := '\n'
	url, err := r.ReadString(byte(delimiter))
	if err != nil {
		return "", fmt.Errorf("Unable to read the string")
	}
	url = strings.TrimSuffix(url, string(delimiter))

	return url, nil
}

func getURLResponseBody(url string) ([]byte, error) {
	client := &http.Client{}
	var wg sync.WaitGroup
	wg.Add(1)
	c1 := make(chan []byte, 1000000)
	// c2 := make(chan []byte)
	// c3 := make(chan []byte)
	// c4 := make(chan []byte)
	// Come up with a system for getting errors as well
	// errChan := make(chan error)
	go func() {
		defer wg.Done()
		fmt.Printf("Getting first byteList")
		responseBody, _ := getURLResponseBodyAsync(client, url, "bytes=0-999999")
		fmt.Printf("And we are assigning the responseBody to a channel")
		c1 <- responseBody
	}()
	// go func() {
	// 	defer wg.Done()
	// 	fmt.Printf("Getting second byteList")
	// 	responseBody, _ := getURLResponseBodyAsync(client, url, "bytes=1000000-1999999")
	// 	c2 <- responseBody
	// 	fmt.Printf("Got second byteList")
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	fmt.Printf("Getting third byteList")
	// 	responseBody, _ := getURLResponseBodyAsync(client, url, "bytes=2000000-2999999")
	// 	c3 <- responseBody
	// 	fmt.Printf("Got third byteList")
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	fmt.Printf("Getting fourth byteList")
	// 	responseBody, _ := getURLResponseBodyAsync(client, url, "bytes=3000000-3999999")
	// 	c4 <- responseBody
	// 	fmt.Printf("Got fourth byteList")
	// }()
	wg.Wait()
	close(c1)

	// for err := range errChan {
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	var body []byte
	for byteList := range c1 {
		fmt.Printf("Appending first byteList")
		body = append(body, byteList...)
	}

	// for byteList := range c2 {
	// 	fmt.Printf("Appending second byteList")
	// 	body = append(body, byteList...)
	// }

	// for byteList := range c3 {
	// 	fmt.Printf("Appending third byteList")
	// 	body = append(body, byteList...)
	// }

	// for byteList := range c4 {
	// 	fmt.Printf("Appending fourth byteList")
	// 	body = append(body, byteList...)
	// }

	return body, nil
}

func getURLResponseBodyAsync(client *http.Client, url string, byteRange string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("an error was encountered creating the request: %v", err)
	}
	req.Header.Add("Range", byteRange)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("an error was encountered getting the url: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("an error was encountered while reading the body: %v", err)
	}
	return body, nil
}
