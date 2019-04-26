package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/thanhpk/randstr"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	url, err := getURLInput(reader)
	if err != nil {
		panic(err)
	}

	filename, err := getFilenameInput(reader)
	if err != nil {
		panic(err)
	}

	body, err := getURLResponseBody(url)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, body, 0644)
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

	url = strings.Trim(url, " ")

	return url, nil
}

func getFilenameInput(r *bufio.Reader) (string, error) {
	fmt.Print("Please enter a name for your file (if you skip this we will randomly generate one): \n")
	delimiter := '\n'
	filename, err := r.ReadString(byte(delimiter))
	if err != nil {
		return "", fmt.Errorf("Unable to read the string")
	}
	filename = strings.TrimSuffix(filename, string(delimiter))

	filename = strings.Trim(filename, " ")

	if filename == "" {
		filename = randstr.Hex(10)
	}

	return filename, nil
}

func getURLResponseBody(url string) ([]byte, error) {
	client := &http.Client{}
	var wg sync.WaitGroup
	wg.Add(4)
	c1 := make(chan []byte, 1)
	c2 := make(chan []byte, 1)
	c3 := make(chan []byte, 1)
	c4 := make(chan []byte, 1)
	// TODO: Come up with a system for getting errors as well
	// errChan := make(chan error)
	go func() {
		defer wg.Done()
		responseBody, _ := getURLResponseBodyAsync(client, url, "bytes=0-1048575")
		c1 <- responseBody
	}()
	go func() {
		defer wg.Done()
		responseBody, _ := getURLResponseBodyAsync(client, url, "bytes=1048576-2097151")
		c2 <- responseBody
	}()
	go func() {
		defer wg.Done()
		responseBody, _ := getURLResponseBodyAsync(client, url, "bytes=2097152-3145727")
		c3 <- responseBody
	}()
	go func() {
		defer wg.Done()
		responseBody, _ := getURLResponseBodyAsync(client, url, "bytes=3145728-4194303")
		c4 <- responseBody
	}()
	wg.Wait()
	close(c1)
	close(c2)
	close(c3)
	close(c4)

	body := make([]byte, 400000, 4000000)
	for byteList := range c1 {
		body = append(body, byteList...)
	}

	for byteList := range c2 {
		body = append(body, byteList...)
	}

	for byteList := range c3 {
		body = append(body, byteList...)
	}

	for byteList := range c4 {
		body = append(body, byteList...)
	}

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
