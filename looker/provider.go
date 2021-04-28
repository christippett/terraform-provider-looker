package looker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	v3 "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

var (
	version           string = "dev"
	defaultAPIVersion string = "3.1"
)

// Provider -
func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_BASE_URL", nil),
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_CLIENT_SECRET", nil),
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_VERSION", defaultAPIVersion),
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_VERIFY_SSL", true),
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_TIMEOUT", nil),
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_WORKSPACE_ID", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"looker_user":             resourceUser(),
			"looker_project":          resourceProject(),
			"looker_git_deploy_key":   resourceGitDeployKey(),
			"looker_project_git_repo": resourceProjectGitRepo(),
			"looker_folder":           resourceFolder(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"looker_session": dataSourceSession(),
			"looker_project": dataSourceProject(),
		},
	}
	p.ConfigureContextFunc = configure(version, p)
	return p
}

type Config struct {
	sdk         *v3.LookerSDK
	accessToken *string
	workspaceId *string
	session     *rtl.AuthSession
	url         string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {

	return func(ctx context.Context, d *schema.ResourceData) (m interface{}, diags diag.Diagnostics) {

		workspaceId := d.Get("workspace_id").(string)
		timeout := d.Get("timeout").(int)

		clientConfig := rtl.ApiSettings{
			BaseUrl:      d.Get("base_url").(string),
			VerifySsl:    d.Get("verify_ssl").(bool),
			Timeout:      int32(timeout),
			ClientId:     d.Get("client_id").(string),
			ClientSecret: d.Get("client_secret").(string),
			ApiVersion:   d.Get("api_version").(string),
			AgentTag:     p.UserAgent("terraform-provider-scaffolding", version),
		}

		// New instance of LookerSDK
		authSession := rtl.NewAuthSession(clientConfig)
		sdk := v3.NewLookerSDK(authSession)

		// Perform initial login
		accessToken := validateLogin(authSession)

		// Update workspace for the current API session
		sessionDetail := v3.WriteApiSession{
			WorkspaceId: &workspaceId,
		}
		session, err := sdk.UpdateSession(sessionDetail, nil)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		config := Config{
			sdk:         sdk,
			accessToken: accessToken,
			workspaceId: session.WorkspaceId,
			session:     authSession,
			url:         fmt.Sprintf("%s/api/%s", authSession.Config.BaseUrl, authSession.Config.ApiVersion),
		}

		return &config, nil
	}
}

func validateLogin(authSession *rtl.AuthSession) *string {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		diag.FromErr(err)
	}
	authSession.Authenticate(req)
	return extractAuthToken(req.Header.Get("Authorization"))
}
