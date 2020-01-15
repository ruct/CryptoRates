package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetBody(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("%v %v", url, resp.StatusCode))
		log.Println(err)
		return "", err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	bodyString := string(bodyBytes)

	return bodyString, nil
}

func GetJson(url string) (map[string]interface{}, error) {
	body, err := GetBody(url)
	if err != nil {
		return nil, err
	}

	bytes := []byte(body)
	var fullData map[string]interface{}
	if err := json.Unmarshal(bytes, &fullData); err != nil {
		log.Println(err)
		return nil, err
	}
	return fullData, nil
}

