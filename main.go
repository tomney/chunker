package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/thanhpk/randstr"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	url, err := promptURLInput(reader)
	if err != nil {
		panic(err)
	}

	filename, err := promptFilenameInput(reader)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	body, err := getURLResponseBody(client, url, 4, 1048576)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		panic(err)
	}
}

func promptURLInput(r *bufio.Reader) (string, error) {
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

func promptFilenameInput(r *bufio.Reader) (string, error) {
	fmt.Print("Please enter a name for your file (if you skip this we will randomly generate one): \n")

	delimiter := '\n'

	filename, err := r.ReadString(byte(delimiter))
	if err != nil {
		return "", fmt.Errorf("Unable to read the string")
	}

	filename = strings.TrimSuffix(filename, string(delimiter))
	filename = strings.Trim(filename, " ")

	if filename == "" {
		filename = randstr.Base62(10)
	}

	return filename, nil
}

func getURLResponseBody(client *http.Client, url string, chunks int, chunkSize int) ([]byte, error) {
	var wg sync.WaitGroup
	channels := make([]chan []byte, 0, chunks)
	errChannel := make(chan error, chunks)

	for i := 0; i < chunks; i++ {
		byteRange := getByteRange(i, chunkSize)

		wg.Add(1)
		channel := make(chan []byte, 1)
		channels = append(channels, channel)
		go func(errChannel chan<- error) {
			defer wg.Done()
			responseBody, err := getURLResponseBodyAsync(client, url, byteRange)
			channel <- responseBody
			errChannel <- err
		}(errChannel)
	}

	wg.Wait()
	close(errChannel)
	for _, channel := range channels {
		close(channel)
	}

	for err := range errChannel {
		if err != nil {
			return nil, err
		}
	}

	body := make([]byte, 0, chunks*chunkSize)

	for _, channel := range channels {
		for byteList := range channel {
			body = append(body, byteList...)
		}
	}

	return body, nil
}

func getByteRange(i int, chunkSize int) string {
	return "bytes=" + strconv.Itoa(chunkSize*i) + "-" + strconv.Itoa(chunkSize*(i+1)-1)
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
