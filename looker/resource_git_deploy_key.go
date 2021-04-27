package looker

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var deployKeyPattern = regexp.MustCompile(`^(?P<key>ssh-rsa AAAA[0-9A-Za-z+/]+[=]{0,3})(?: (?P<comment>[^\s]+))?`)

func resourceGitDeployKey() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceGitDeployKeyCreate,
		ReadContext:   resourceGitDeployKeyRead,
		DeleteContext: resourceGitDeployKeyDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "Looker project ID.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"public_key": {
				// This description is used by the documentation generator and the language server.
				Description: "Public key of the generated ssh key pair for the project's git repository.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceGitDeployKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	sdk := config.sdk
	projectId := d.Get("project_id").(string)

	_, err := sdk.CreateGitDeployKey(projectId, nil)
	if err != nil && !regexp.MustCompile("^response error: 409").MatchString(err.Error()) {
		return diag.FromErr(err)
	}

	d.SetId(projectId)

	return resourceGitDeployKeyRead(ctx, d, meta)
}

func resourceGitDeployKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	path := fmt.Sprintf("/projects/%s/git/deploy_key", d.Id())

	key, err := doRequest("GET", path, config.session)
	if err != nil {
		return diag.FromErr(err)
	}

	key = deployKeyPattern.Find(key)
	d.Set("public_key", string(key))

	return diags
}

func resourceGitDeployKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
