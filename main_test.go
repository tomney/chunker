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
	s.Equal("someurl", url)
}

func (s *getURLInputTestSuite) Test_RemovesPrecedingAndTrailingWhitespace() {
	s.reader = bufio.NewReader(strings.NewReader("  someurl    \n"))
	url, _ := getURLInput(s.reader)
	s.Equal("someurl", url)
}

func (s *getURLInputTestSuite) Test_ReturnsErrIfInputDoesNotContainDelimiter() {
	s.reader = bufio.NewReader(strings.NewReader("someurl"))
	url, err := getURLInput(s.reader)
	s.Equal("", url)
	s.Equal(fmt.Errorf("Unable to read the string"), err)
}

func TestGetUrlInputTestSuite(t *testing.T) {
	suite.Run(t, new(getURLInputTestSuite))
}

type getFilenameInputTestSuite struct {
	suite.Suite
	reader *bufio.Reader
}

func (s *getFilenameInputTestSuite) SetupTest() {
	s.reader = bufio.NewReader(strings.NewReader("somefilename\n"))
}

func (s *getFilenameInputTestSuite) Test_RemovesDelimiter() {
	filename, _ := getFilenameInput(s.reader)
	s.Equal("somefilename", filename)
}

func (s *getFilenameInputTestSuite) Test_RemovesPrecedingAndTrailingWhitespace() {
	s.reader = bufio.NewReader(strings.NewReader("  somefilename    \n"))
	filename, _ := getFilenameInput(s.reader)
	s.Equal("somefilename", filename)
}

func (s *getFilenameInputTestSuite) Test_ReturnsErrIfInputDoesNotContainDelimiter() {
	s.reader = bufio.NewReader(strings.NewReader("somefilename"))
	filename, err := getFilenameInput(s.reader)
	s.Equal("", filename)
	s.Equal(fmt.Errorf("Unable to read the string"), err)
}

func (s *getFilenameInputTestSuite) Test_ReturnsTenCharacterRandomStringIfNoStringIsProvided() {
	s.reader = bufio.NewReader(strings.NewReader("\n"))
	filename, _ := getFilenameInput(s.reader)
	s.Equal(10, len(filename))
}

func TestGetFilenameInputTests(t *testing.T) {
	suite.Run(t, new(getFilenameInputTestSuite))
}
