package looker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/looker-open-source/sdk-codegen/go/rtl"
	v3 "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

var namePattern = regexp.MustCompile(`[^a-z0-9-_]`)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func updateSession(sdk *v3.LookerSDK, workspaceId string) (v3.ApiSession, error) {
	sessionDetail := v3.WriteApiSession{
		WorkspaceId: &workspaceId,
	}
	return sdk.UpdateSession(sessionDetail, nil)
}

func formatName(name string) string {
	return strings.ToLower(namePattern.ReplaceAllString(name, "_"))
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
