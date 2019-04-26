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

// DELIMITER is the delimiter needed to delimit user input
var DELIMITER = '\n'

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

	chunks, err := promptNumberOfChunks(reader)
	if err != nil {
		panic(err)
	}

	chunkSize, err := promptChunkSize(reader)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	body, err := getURLResponseBody(client, url, chunks, chunkSize)
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

	url, err := r.ReadString(byte(DELIMITER))
	if err != nil {
		return "", fmt.Errorf("Unable to read the string")
	}

	url = trimWhiteSpaceAndDelimiter(url)

	return url, nil
}

func promptFilenameInput(r *bufio.Reader) (string, error) {
	fmt.Print("Please enter a name for your file (if you skip this we will randomly generate one): \n")

	filename, err := r.ReadString(byte(DELIMITER))
	if err != nil {
		return "", fmt.Errorf("Unable to read the string")
	}

	filename = trimWhiteSpaceAndDelimiter(filename)

	if filename == "" {
		filename = randstr.Base62(10)
	}

	return filename, nil
}

func promptNumberOfChunks(r *bufio.Reader) (int, error) {
	defaultChunks := 4
	fmt.Printf("Please enter the number of chunks you would like to download the file in (leave blank for default=%v): \n",
		defaultChunks,
	)

	chunksString, err := r.ReadString(byte(DELIMITER))
	if err != nil {
		return 0, fmt.Errorf("Unable to read the string")
	}

	chunksString = trimWhiteSpaceAndDelimiter(chunksString)

	if chunksString == "" {
		return defaultChunks, nil
	}

	chunks, err := strconv.Atoi(chunksString)
	if err != nil {
		return 0, fmt.Errorf("Unable to convert the value provided for the number of chunks to an integer")
	}

	return chunks, nil
}

func promptChunkSize(r *bufio.Reader) (int, error) {
	defaultChunkSize := 1048576
	fmt.Printf("Please enter, in bytes, the size of chunks you would like to use to download the file in (leave blank for default=%v): \n",
		defaultChunkSize,
	)

	chunkSizeString, err := r.ReadString(byte(DELIMITER))
	if err != nil {
		return 0, fmt.Errorf("Unable to read the string")
	}

	chunkSizeString = trimWhiteSpaceAndDelimiter(chunkSizeString)

	if chunkSizeString == "" {
		return defaultChunkSize, nil
	}

	chunkSize, err := strconv.Atoi(chunkSizeString)
	if err != nil {
		return 0, fmt.Errorf("Unable to convert the value provided for the number of chunks to an integer")
	}

	return chunkSize, nil
}

func trimWhiteSpaceAndDelimiter(s string) string {
	s = strings.TrimSuffix(s, string(DELIMITER))
	s = strings.Trim(s, " ")
	return s
}

func getURLResponseBody(client *http.Client, url string, chunks int, chunkSize int) ([]byte, error) {
	var wg sync.WaitGroup

	// Create a list of channels for the expected responses
	channels := make([]chan []byte, 0, chunks)

	// Create a channel to handle any errors that may occur during the go func calls
	errChannel := make(chan error, chunks)

	for i := 0; i < chunks; i++ {
		byteRange := getByteRange(i, chunkSize)

		// Create a channel for each expected response
		channel := make(chan []byte, 1)
		channels = append(channels, channel)

		wg.Add(1)
		go func(errChannel chan<- error) {
			defer wg.Done()
			responseBody, err := getURLResponseBodyAsync(client, url, byteRange)
			channel <- responseBody
			errChannel <- err
		}(errChannel)
	}
	wg.Wait()

	// Close all the channels
	close(errChannel)
	for _, channel := range channels {
		close(channel)
	}

	// Check if any errors occurred, if so return the earliest error
	for err := range errChannel {
		if err != nil {
			return nil, err
		}
	}

	// Iterate over the list of channels and combine their contents
	body := make([]byte, 0, chunks*chunkSize)
	for _, channel := range channels {
		for byteList := range channel {
			body = append(body, byteList...)

			//Exit early if you've reached the end of the range prematurely
			if len(byteList) < chunkSize {
				return body, nil
			}
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
		return nil, fmt.Errorf("An error was encountered creating the request: %v", err)
	}

	req.Header.Add("Range", byteRange)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("An error was encountered getting the url: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("An error was encountered while reading the body: %v", err)
	}

	return body, nil
}
