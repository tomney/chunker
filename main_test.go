package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type getURLInputTestSuite struct {
	suite.Suite
	reader *bufio.Reader
}

func (s *getURLInputTestSuite) SetupTest() {
	s.reader = bufio.NewReader(strings.NewReader("someurl\n"))

}

func (s *getURLInputTestSuite) Test_RemovesDelimiter() {
	url, _ := getURLInput(s.reader)
	s.Equal(url, "someurl")
}

func (s *getURLInputTestSuite) Test_ReturnsErrIfInputDoesNotContainDelimiter() {
	s.reader = bufio.NewReader(strings.NewReader("someurl"))
	url, err := getURLInput(s.reader)
	s.Equal("", url)
	s.Equal(fmt.Errorf("Unable to read the string"), err)
}

func TestGetURLInputTestSuite(t *testing.T) {
	suite.Run(t, new(getURLInputTestSuite))
}
