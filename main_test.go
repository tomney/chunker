package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/suite"
)

type promptURLInputTestSuite struct {
	suite.Suite
	reader *bufio.Reader
}

func (s *promptURLInputTestSuite) Test_RemovesDelimiter() {
	s.reader = bufio.NewReader(strings.NewReader("someurl\n"))
	url, _ := promptURLInput(s.reader)
	s.Equal("someurl", url)
}

func (s *promptURLInputTestSuite) Test_RemovesPrecedingAndTrailingWhitespace() {
	s.reader = bufio.NewReader(strings.NewReader("  someurl    \n"))
	url, _ := promptURLInput(s.reader)
	s.Equal("someurl", url)
}

func (s *promptURLInputTestSuite) Test_ReturnsErrIfInputDoesNotContainDelimiter() {
	s.reader = bufio.NewReader(strings.NewReader("someurl"))
	url, err := promptURLInput(s.reader)
	s.Equal("", url)
	s.Equal(fmt.Errorf("Unable to read the string"), err)
}

func TestPromptURLInputTests(t *testing.T) {
	suite.Run(t, new(promptURLInputTestSuite))
}

type promptFilenameInputTestSuite struct {
	suite.Suite
	reader *bufio.Reader
}

func (s *promptFilenameInputTestSuite) Test_RemovesDelimiter() {
	s.reader = bufio.NewReader(strings.NewReader("somefilename\n"))
	filename, _ := promptFilenameInput(s.reader)
	s.Equal("somefilename", filename)
}

func (s *promptFilenameInputTestSuite) Test_RemovesPrecedingAndTrailingWhitespace() {
	s.reader = bufio.NewReader(strings.NewReader("  somefilename    \n"))
	filename, _ := promptFilenameInput(s.reader)
	s.Equal("somefilename", filename)
}

func (s *promptFilenameInputTestSuite) Test_ReturnsErrIfInputDoesNotContainDelimiter() {
	s.reader = bufio.NewReader(strings.NewReader("somefilename"))
	filename, err := promptFilenameInput(s.reader)
	s.Equal("", filename)
	s.Equal(fmt.Errorf("Unable to read the string"), err)
}

func (s *promptFilenameInputTestSuite) Test_ReturnsTenCharacterRandomStringIfNoStringIsProvided() {
	s.reader = bufio.NewReader(strings.NewReader("\n"))
	filename, _ := promptFilenameInput(s.reader)
	s.Equal(10, len(filename))
}

func TestPromptFilenameInputTests(t *testing.T) {
	suite.Run(t, new(promptFilenameInputTestSuite))
}

type getURLResponseBodyTestSuite struct {
	suite.Suite
	server *httptest.Server
	client *http.Client
}

func (s *getURLResponseBodyTestSuite) SetupTest() {
	s.server = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "1")
			},
		),
	)
	s.client = s.server.Client()
}

func (s *getURLResponseBodyTestSuite) Test_ReturnsFourConcatenatedURLResponses() {
	body, _ := getURLResponseBody(s.client, s.server.URL, 4, 1048576)
	s.Equal("1111", string(body))
}

func (s *getURLResponseBodyTestSuite) TearDownTest() {
	s.server.Close()
}

func TestGetURLResponseBodyTests(t *testing.T) {
	suite.Run(t, new(getURLResponseBodyTestSuite))
}
