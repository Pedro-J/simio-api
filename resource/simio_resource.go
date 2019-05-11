package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"simio-api/service"
)

type SimioRequest struct {
	DNA []string `json:"dna"`
}

type SimioResource struct {
	simioService service.SimioService
}

func (sr *SimioResource) CheckSimian(rw http.ResponseWriter, req *http.Request) {

	simioRequest, err := mapToSimioRequest(req)

	if err != nil {
		buildResponse(rw, http.StatusBadRequest, err.Error())
		return
	}

	isSimian, processErr := sr.simioService.ProcessDNA(simioRequest.DNA)

	if processErr != nil {
		buildResponse(rw, http.StatusBadRequest, processErr.Error())
		return
	}

	if isSimian {
		buildResponse(rw, http.StatusOK, "")
		log.Print("DNA is simian")
	} else {
		buildResponse(rw, http.StatusForbidden, "")
		log.Print("DNA is not simian")
	}
}

func (sr *SimioResource) GetSimiansProportion(rw http.ResponseWriter, req *http.Request) {
	stats := sr.simioService.GetSimiansProportion()
	responseBody, _ := json.Marshal(stats)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(responseBody))
}

func mapToSimioRequest(req *http.Request) (*SimioRequest, error) {

	defaultInvalidPayloadError := fmt.Errorf("Invalid Request Payload")

	bodyString, err := getBody(req.Body)

	defer req.Body.Close()

	if err != nil {
		log.Printf("%+v", err)
		return nil, defaultInvalidPayloadError
	}

	var simioRequest SimioRequest
	err = json.Unmarshal([]byte(bodyString), &simioRequest)

	if err != nil {
		log.Printf("%+v", err)
		return nil, defaultInvalidPayloadError
	}

	return &simioRequest, nil
}

func getBody(body io.ReadCloser) (string, error) {
	var responseBody string

	bodyBytes, err := ioutil.ReadAll(body)
	responseBody = string(bodyBytes)

	if err != nil {
		return "", fmt.Errorf("Error on reading payload body")
	}

	return responseBody, nil
}

func buildResponse(rw http.ResponseWriter, statusCode int, errorMsg string) {
	var responseMessage string
	if errorMsg == "" {
		responseMessage = fmt.Sprintf("%s", http.StatusText(statusCode))
	} else {
		responseMessage = fmt.Sprintf("%s - %s", http.StatusText(statusCode), errorMsg)
	}

	rw.WriteHeader(statusCode)
	rw.Write([]byte(responseMessage))
}

func BuildSimioResource() *SimioResource {
	return &SimioResource{
		simioService: service.BuildSimioService(),
	}
}
