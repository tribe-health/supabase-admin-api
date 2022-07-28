package common

import (
	"io/ioutil"
	"net/http"
)

func AuthedRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("apikey", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoic3VwYWJhc2VfYWRtaW4ifQ.veeAYq7d22dUiUe7cfQDvnZULmLJwUiB2neF_zTcD94")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
