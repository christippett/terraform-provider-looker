package looker

import (
	"regexp"
	"strings"

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
