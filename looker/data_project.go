package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "Looker project data source.",

		ReadContext: dataSourceProjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Project display name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"uses_git": {
				Description: "True if the project is configured with a git repository.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"git_remote_url": {
				// This description is used by the documentation generator and the language server.
				Description: "Git remote repository url.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"git_service_name": {
				// This description is used by the documentation generator and the language server.
				Description: "Name of the git service provider.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"git_username": {
				Description: "Git username for HTTPS authentication. [production only]",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"git_password": {
				Description: "Git password for HTTPS authentication. [production only]",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"git_production_branch_name": {
				Description: "Git production branch name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"git_username_user_attribute": {
				Description: "User attribute name for username in per-user HTTPS authentication.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"git_password_user_attribute": {
				Description: "User attribute name for password in per-user HTTPS authentication.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"git_application_server_http_port": {
				Description: "Port that HTTP(S) application server is running on (for PRs, file browsing, etc).",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"git_application_server_http_scheme": {
				Description: "Scheme that is running on the application server (for PRs, file browsing, etc).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pull_request_mode": {
				Description: "The git pull request policy for this project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"validation_required": {
				Description: "True if the project must pass validation checks before project changes can be committed to the git repository.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"allow_warnings": {
				Description: "True if the project can be committed with warnings when `validation_required` is true. Does nothing if `validation_required` is false.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"git_release_mgmt_enabled": {
				Description: "True if advanced git release management is enabled for this project.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdk := meta.(*Config).sdk

	projectId := d.Get("name").(string)
	d.SetId(projectId)

	project, err := sdk.Project(d.Id(), "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	projectToResourceData(project, d)
	return nil
}

func projectToResourceData(project v3.Project, d *schema.ResourceData) {
	d.Set("name", project.Name)
	d.Set("git_remote_url", project.GitRemoteUrl)
	d.Set("git_service_name", project.GitServiceName)
	d.Set("git_username", project.GitUsername)
	d.Set("git_password", project.GitPassword)
	d.Set("git_production_branch_name", project.GitProductionBranchName)
	d.Set("git_username_user_attribute", project.GitUsernameUserAttribute)
	d.Set("git_password_user_attribute", project.GitPasswordUserAttribute)
	d.Set("git_application_server_http_port", project.GitApplicationServerHttpPort)
	d.Set("git_application_server_http_scheme", project.GitApplicationServerHttpScheme)
	d.Set("pull_request_mode", project.PullRequestMode)
	d.Set("validation_required", project.ValidationRequired)
	d.Set("allow_warnings", project.AllowWarnings)
	d.Set("git_release_mgmt_enabled", project.GitReleaseMgmtEnabled)
}
