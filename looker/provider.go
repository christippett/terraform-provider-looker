package looker

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	v3 "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

var defaultLookerAPIVersion string = "3.1"

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
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
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_VERSION", defaultLookerAPIVersion),
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"looker_user": resourceUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiVersion := d.Get("api_version").(string)

	cfg := rtl.ApiSettings{
		BaseUrl:   d.Get("base_url").(string),
		VerifySsl: d.Get("verify_ssl").(bool),
		// Timeout:      d.Get("timeout").(int32),
		Timeout:      120,
		AgentTag:     "",
		FileName:     "",
		ClientId:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		ApiVersion:   apiVersion,
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var sdk *v3.LookerSDK

	// New instance of LookerSDK
	if strings.HasPrefix(apiVersion, "3.") {
		sdk = v3.NewLookerSDK(rtl.NewAuthSession(cfg))
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unsupported Looker API version",
		})
		return nil, diags
	}

	_, err := sdk.Login(v3.RequestLogin{
		ClientId:     &cfg.ClientId,
		ClientSecret: &cfg.ClientSecret,
	}, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to authenticate with Looker SDK",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return sdk, diags
}
