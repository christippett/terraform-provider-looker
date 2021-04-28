package looker

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v4 "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceProjectGitRepo() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceProjectGitRepoCreate,
		ReadContext:   resourceProjectGitRepoRead,
		DeleteContext: resourceProjectGitRepoDelete,
		UpdateContext: resourceProjectGitRepoUpdate,

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "Looker project ID/name.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"git_remote_url": {
				// This description is used by the documentation generator and the language server.
				Description: "Git remote repository url.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"git_service_name": {
				// This description is used by the documentation generator and the language server.
				Description: "Name of the git service provider.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"git_username": {
				Description:  "Git username for HTTPS authentication. [production only]",
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"git_password"},
			},
			"git_password": {
				Description:  "Git password for HTTPS authentication. [production only]",
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"git_username"},
			},
			"git_production_branch_name": {
				Description: "Git production branch name.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"git_username_user_attribute": {
				Description:  "User attribute name for username in per-user HTTPS authentication.",
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"git_password_user_attribute"},
			},
			"git_password_user_attribute": {
				Description:  "User attribute name for password in per-user HTTPS authentication.",
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"git_username_user_attribute"},
			},
			"git_application_server_http_port": {
				Description: "Port that HTTP(S) application server is running on (for PRs, file browsing, etc).",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"git_application_server_http_scheme": {
				Description: "Scheme that is running on the application server (for PRs, file browsing, etc).",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateDiagFunc: func(v interface{}, p cty.Path) (diags diag.Diagnostics) {
					val := v.(string)
					opts := []string{
						string(v4.GitApplicationServerHttpScheme_Http),
						string(v4.GitApplicationServerHttpScheme_Https),
					}
					if _, err := Find(opts, val); err != nil {
						return diag.FromErr(err)
					}
					return diags
				},
			},
			"pull_request_mode": {
				Description: "The git pull request policy for this project.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ValidateDiagFunc: func(v interface{}, p cty.Path) (diags diag.Diagnostics) {
					val := v.(string)
					opts := []string{
						string(v4.PullRequestMode_Links),
						string(v4.PullRequestMode_Off),
						string(v4.PullRequestMode_Recommended),
						string(v4.PullRequestMode_Required),
					}
					if _, err := Find(opts, val); err != nil {
						return diag.FromErr(err)
					}
					return diags
				},
			},
			"validation_required": {
				Description: "True if the project must pass validation checks before project changes can be committed to the git repository.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"allow_warnings": {
				Description: "True if the project can be committed with warnings when `validation_required` is true. Does nothing if `validation_required` is false.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"git_release_mgmt_enabled": {
				Description: "True if advanced git release management is enabled for this project.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceProjectGitRepoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdk := meta.(*Config).sdk

	projectId := d.Get("project").(string)
	d.SetId(projectId)

	project, err := sdk.Project(projectId, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	writeProject := makeWriteProjectGitRepo(d)
	_, err = sdk.UpdateProject(*project.Id, writeProject, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectGitRepoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	return dataSourceProjectRead(ctx, d, meta)
}

func resourceProjectGitRepoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceProjectGitRepoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceProjectGitRepoCreate(ctx, d, meta)
}

func makeWriteProjectGitRepo(d *schema.ResourceData) v4.WriteProject {
	var validationRequired, allowWarnings bool
	var git struct {
		remoteUrl                   string
		serviceName                 string
		username                    string
		password                    string
		productionBranchName        string
		usernameUserAttribute       string
		passwordUserAttribute       string
		applicationServerHttpPort   int64
		applicationServerHttpScheme v4.GitApplicationServerHttpScheme
		releaseMgmtEnabled          bool
		pullRequestMode             v4.PullRequestMode
	}
	var writeProject v4.WriteProject
	if d.HasChange("git_remote_url") {
		git.remoteUrl = d.Get("git_remote_url").(string)
		git.serviceName = d.Get("git_service_name").(string)
		writeProject.GitRemoteUrl = &git.remoteUrl
		writeProject.GitServiceName = &git.serviceName
	}
	if d.HasChanges("git_username", "git_password") {
		git.username = d.Get("git_username").(string)
		git.password = d.Get("git_password").(string)
		writeProject.GitUsername = &git.username
		writeProject.GitPassword = &git.password
	}
	if d.HasChanges("git_username_user_attribute", "git_password_user_attribute") {
		git.usernameUserAttribute = d.Get("git_username_user_attribute").(string)
		git.passwordUserAttribute = d.Get("git_password_user_attribute").(string)
		writeProject.GitUsernameUserAttribute = &git.usernameUserAttribute
		writeProject.GitPasswordUserAttribute = &git.passwordUserAttribute
	}
	if d.HasChanges("git_application_server_http_scheme", "git_application_server_http_port") {
		git.applicationServerHttpPort = int64(d.Get("git_application_server_http_port").(int))
		git.applicationServerHttpScheme = v4.GitApplicationServerHttpScheme(d.Get("git_application_server_http_scheme").(string))
		writeProject.GitApplicationServerHttpPort = &git.applicationServerHttpPort
		writeProject.GitApplicationServerHttpScheme = &git.applicationServerHttpScheme
	}
	if d.HasChange("git_production_branch_name") {
		git.productionBranchName = d.Get("git_production_branch_name").(string)
		writeProject.GitProductionBranchName = &git.productionBranchName
	}
	if d.HasChange("git_release_mgmt_enabled") {
		git.releaseMgmtEnabled = d.Get("git_release_mgmt_enabled").(bool)
		writeProject.GitReleaseMgmtEnabled = &git.releaseMgmtEnabled
	}
	if d.HasChange("pull_request_mode") {
		git.pullRequestMode = v4.PullRequestMode(d.Get("pull_request_mode").(string))
		writeProject.PullRequestMode = &git.pullRequestMode
	}
	if d.HasChanges("validation_required", "allow_warnings") {
		validationRequired = d.Get("validation_required").(bool)
		allowWarnings = d.Get("allow_warnings").(bool)
		writeProject.ValidationRequired = &validationRequired
		writeProject.AllowWarnings = &allowWarnings
	}

	return writeProject
}
