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
	url := getURLInput(reader)

	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Errorf("an error was encountered getting the url"))
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	ioutil.WriteFile("dat1", body, 0644)
}

func getURLInput(r *bufio.Reader) string {
	fmt.Print("Please enter the url you would like to download from then hit enter: \n")

	delimiter := '\n'
	url, err := r.ReadString(byte(delimiter))
	if err != nil {
		panic(fmt.Errorf("the string that was entered is invalid"))
	}
	url = strings.TrimSuffix(url, string(delimiter))

	return url
}
