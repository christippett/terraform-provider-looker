package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"validation_required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"allow_warnings": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdk := m.(*v3.LookerSDK)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	validationRequired := d.Get("validation_required").(bool)
	allowWarnings := d.Get("allow_warnings").(bool)

	projectDetails := v3.WriteProject{
		Name:               &name,
		ValidationRequired: &validationRequired,
		AllowWarnings:      &allowWarnings,
	}

	// Set session workspace to "dev"
	err := updateSession(sdk, "dev")
	if err != nil {
		return diag.FromErr(err)
	}

	project, err := sdk.CreateProject(projectDetails, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*project.Id)

	resourceProjectRead(ctx, d, m)

	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdk := m.(*v3.LookerSDK)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// Set session workspace to "dev"
	err := updateSession(sdk, "dev")
	if err != nil {
		return diag.FromErr(err)
	}

	project, err := sdk.Project(d.Id(), "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", project.Name)
	d.Set("validation_required", project.ValidationRequired)
	d.Set("allow_warnings", project.AllowWarnings)

	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdk := m.(*v3.LookerSDK)

	// Set session workspace to "dev"
	err := updateSession(sdk, "dev")
	if err != nil {
		return diag.FromErr(err)
	}

	project, err := sdk.Project(d.Id(), "", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	var writeProject v3.WriteProject

	writeProject.Name = d.Get("name").(*string)
	writeProject.ValidationRequired = d.Get("validation_required").(*bool)
	writeProject.AllowWarnings = d.Get("alloww_warnings").(*bool)

	_, err = sdk.UpdateProject(*project.Id, writeProject, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Terraform is unable to delete Looker project.",
			Detail:   "Looker projects cannot be deleted programmatically and must be removed manually via the web console.",
		},
	}
}

func updateSession(sdk *v3.LookerSDK, workspaceId string) error {
	sessionDetail := v3.WriteApiSession{
		WorkspaceId: &workspaceId,
	}
	_, err := sdk.UpdateSession(sessionDetail, nil)
	return err
}
