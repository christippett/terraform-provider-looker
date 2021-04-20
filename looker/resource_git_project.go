package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGitProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitProjectCreate,
		ReadContext:   resourceGitProjectRead,
		UpdateContext: resourceGitProjectUpdate,
		DeleteContext: resourceGitProjectDelete,
		Schema: map[string]*schema.Schema{
			"remote_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"production_branch_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_cookie": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username_user_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password_user_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_server_http_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_server_http_scheme": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pull_request_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGitProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceGitProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceGitProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceGitProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}
