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
	fmt.Print("Please enter the url you would like to download from then hit enter: \n")

	delimiter := '\n'
	url, err := reader.ReadString(byte(delimiter))
	url = strings.TrimSuffix(url, string(delimiter))
	if err != nil {
		panic(fmt.Errorf("the string that was entered is invalid"))
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Errorf("an error was encountered getting the url"))
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	ioutil.WriteFile("dat1", body, 0644)
}
