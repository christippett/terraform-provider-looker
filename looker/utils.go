package looker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/looker-open-source/sdk-codegen/go/rtl"
)

func Find(slice []string, val string) (int, error) {
	for i, item := range slice {
		if item == val {
			return i, nil
		}
	}
	return -1, fmt.Errorf("'%s' not a member of: %s", val, slice)
}

func convertIntToInt64(n int) *int64 {
	n64 := int64(n)
	return &n64
}

func convertIntSlice(s []interface{}) []int64 {
	s64 := make([]int64, len(s))
	for i, n := range s {
		s64[i] = int64(n.(int))
	}
	return s64
}

func extractAuthToken(s string) *string {
	re := regexp.MustCompile(`(?i)^(?:token|bearer) ([[:alnum:]]+)$`)

	m := re.FindStringSubmatch(s)
	if m != nil {
		return &m[1]
	}
	return nil
}

type ClientError interface {
	error
	Status() int
}

type HTTPError struct {
	err  error
	code int
}

func (e HTTPError) Error() string {
	return e.err.Error()
}

func (e HTTPError) Status() int {
	return e.code
}

func doRequest(method, path string, session *rtl.AuthSession) ([]byte, error) {
	url := fmt.Sprintf("%s/api/%s%s", session.Config.BaseUrl, session.Config.ApiVersion, path)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if err := session.Authenticate(req); err != nil {
		return nil, err
	}

	client := http.Client{
		Transport: session.Transport,
		Timeout:   time.Duration(session.Config.Timeout) * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 226 {
		return nil, HTTPError{err, res.StatusCode}
	}

	return body, nil
}
