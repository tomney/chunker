package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	url, err := getURLInput(reader)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Errorf("an error was encountered getting the url"))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("an error was encountered while reading the body"))
	}

	ioutil.WriteFile("dat1", body, 0644)
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
