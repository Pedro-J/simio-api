package resource

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"simio-api/service"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	dnaSimianHorizontal = []string{"CCCG", "AAAT", "GGGA", "TTTT"}
	dnaHuman            = []string{"CGAT", "GTCA", "TACG", "TCGA"}
	dnaDifCols          = []string{"CTGZ", "CT", "TATTTGGAATTT"}
	dnaInvalidFirstChar = []string{"ZGTCCCTA", "GCAGGAAT", "TTCCAAGG", "TCAATTGC", "GGTTCCAG", "CCTAGGCC", "TTGCGCAA", "AAACCGTA"}
	dnaEmpty            = []string{}
)

//Mocking simio service
type SimioServiceMock struct {
	mock.Mock
	service.SimioService
}

func (sm *SimioServiceMock) ProcessDNA(dna []string) (bool, error) {
	args := sm.Called(dna)
	return args.Bool(0), args.Error(1)
}

func (sm *SimioServiceMock) GetSimiansProportion() service.Stats {
	args := sm.Called()
	return args.Get(0).(service.Stats)
}

func doRequest(url string, reqBody string, method string) (string, int) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var req *http.Request

	if method == http.MethodGet {
		req, _ = http.NewRequest(method, url, nil)
	} else {
		req, _ = http.NewRequest(method, url, strings.NewReader(string(reqBody)))
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error on http request")
	}

	var respBody string
	respBody, err = getBody(resp.Body)

	if err != nil {
		log.Printf("Error on reading resp body")
	}

	defer resp.Body.Close()

	return respBody, resp.StatusCode
}

//Test Resources
func TestCheckSimian(t *testing.T) {
	assert := assert.New(t)

	type Case struct {
		request            SimioRequest
		procesResultErr    error
		processResult      bool
		expectedStatusCode int
	}

	cases := []Case{
		Case{request: SimioRequest{DNA: dnaSimianHorizontal}, processResult: true, procesResultErr: nil, expectedStatusCode: http.StatusOK},
		Case{request: SimioRequest{DNA: dnaHuman}, processResult: false, procesResultErr: nil, expectedStatusCode: http.StatusForbidden},
		Case{request: SimioRequest{DNA: dnaDifCols}, processResult: false, procesResultErr: fmt.Errorf(""), expectedStatusCode: http.StatusBadRequest},
		Case{request: SimioRequest{DNA: dnaInvalidFirstChar}, processResult: false, procesResultErr: fmt.Errorf(""), expectedStatusCode: http.StatusBadRequest},
		Case{request: SimioRequest{DNA: dnaEmpty}, processResult: false, procesResultErr: fmt.Errorf(""), expectedStatusCode: http.StatusBadRequest},
	}

	for _, currentCase := range cases {
		simioServiceMocked := new(SimioServiceMock)
		simioServiceMocked.On("ProcessDNA", currentCase.request.DNA).
			Return(currentCase.processResult, currentCase.procesResultErr)

		simioResource := NewSimioResource(simioServiceMocked)

		// Start a local HTTP server
		server := httptest.NewServer(http.HandlerFunc(simioResource.CheckSimian))

		body, _ := json.Marshal(currentCase.request)
		_, resultStatusCode := doRequest(server.URL, string(body), http.MethodPost)

		assert.Equal(currentCase.expectedStatusCode, resultStatusCode)

		server.Close()
	}
}

func TestGetSimiansProportion(t *testing.T) {
	assert := assert.New(t)

	type Case struct {
		stats service.Stats

		expectedStatusCode int
		expectedRespBody   string
	}

	stats := service.Stats{CountHumanDNA: 2500, CountMutantDNA: 1000, Ratio: float64(1000 / 2500)}
	expectedRespBody, _ := json.Marshal(stats)

	cases := []Case{
		Case{
			stats:              stats,
			expectedRespBody:   string(expectedRespBody),
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, currentCase := range cases {
		simioServiceMocked := new(SimioServiceMock)
		simioServiceMocked.On("GetSimiansProportion").Return(currentCase.stats)

		simioResource := NewSimioResource(simioServiceMocked)

		// Start a local HTTP server
		server := httptest.NewServer(http.HandlerFunc(simioResource.GetSimiansProportion))

		_, resultStatusCode := doRequest(server.URL, "", http.MethodGet)

		assert.Equal(currentCase.expectedStatusCode, resultStatusCode)

		server.Close()
	}
}

func TestMapToSimioRequest(t *testing.T) {
	assert := assert.New(t)

	invalidReq, _ := http.NewRequest(http.MethodPost, "url.test.com", strings.NewReader(string("invalid")))
	validReq, _ := http.NewRequest(http.MethodPost, "url.test.com", strings.NewReader(string(`{"dna": ["GTCA"]}`)))
	_, err := mapToSimioRequest(invalidReq)
	res, _ := mapToSimioRequest(validReq)

	assert.NotNil(err)
	assert.NotNil(res)
}

func TestBuildResource(t *testing.T) {
	assert := assert.New(t)
	resource := BuildSimioResource()

	assert.NotNil(resource)
	assert.NotNil(resource.simioService)
}
