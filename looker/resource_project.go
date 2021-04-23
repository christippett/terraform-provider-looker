package looker

import (
	"context"
	"log"
	"net/http"
	"strings"

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
				Required: true,
			},
			"validation_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allow_warnings": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	sdk := m.(*Config).sdk

	// Create project
	writeProject := createWriteProject(d)
	project, err := sdk.CreateProject(writeProject, nil)

	if err == nil {
		d.SetId(*project.Id)

	} else if err != nil && strings.Contains(err.Error(), http.StatusText(http.StatusUnprocessableEntity)) {
		log.Printf("A project named '%s' already exists and will be linked to this resource.", *writeProject.Name)
		d.SetId(*writeProject.Name)

	} else {
		return diag.FromErr(err)
	}

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	sdk := m.(*Config).sdk

	project, err := sdk.Project(d.Id(), "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", project.Name)
	d.Set("validation_required", project.ValidationRequired)
	d.Set("allow_warnings", project.AllowWarnings)

	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	sdk := m.(*Config).sdk

	project, err := sdk.Project(d.Id(), "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	writeProject := createWriteProject(d)
	_, err = sdk.UpdateProject(*project.Id, writeProject, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	d.SetId("")

	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Terraform is unable to delete Looker project.",
			Detail:   "Looker projects cannot be deleted programmatically and must be removed manually via the web console.",
		},
	}
}

func createWriteProject(d *schema.ResourceData) v3.WriteProject {
	name := d.Get("name").(string)
	validationRequired := d.Get("validation_required").(bool)
	allowWarnings := d.Get("allow_warnings").(bool)

	// Remove invalid characters from name
	name = formatName(name)

	return v3.WriteProject{
		Name:               &name,
		ValidationRequired: &validationRequired,
		AllowWarnings:      &allowWarnings,
	}
}
