package huobiapi

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func (A *HuobiApi) httpRequest(request *http.Request) ([]byte, error) {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	A.debugData.lastResponseData = data

	//fmt.Println("SERVER ANSWER: ", string(data)) //debug

	return data, nil
}

func (A *HuobiApi) doRequestGET(uri string) ([]byte, error) {
	A.debugData.lastSentData = []byte(uri)

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	var response []byte
	response, err = A.httpRequest(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (A *HuobiApi) doRequestPOST(uri string, jsonBytes []byte) ([]byte, error) {
	A.debugData.lastSentData = jsonBytes

	body := bytes.NewReader(jsonBytes)
	//fmt.Println("\nBODY PARAMETERS: ", string(jsonBytes))
	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	//request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	//request.Header.Set("Accept", "application/json")

	var response []byte
	response, err = A.httpRequest(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
