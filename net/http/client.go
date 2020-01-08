package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetContent(url string) (content string, err error) {
	fmt.Printf("Request url `%s`\n", url)

	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", err
		}

		return string(body), nil
	} else {
		return "", fmt.Errorf("response code is '%s'", resp.Status)
	}
}
