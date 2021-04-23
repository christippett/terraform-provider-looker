package looker

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSession() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSession,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func readDataSession(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	config := m.(*Config)
	sdk := config.sdk

	session, err := sdk.Session(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	me, err := sdk.Me("", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("access_token", config.accessToken)
	d.Set("workspace_id", session.WorkspaceId)
	d.Set("user_id", me.Id)
	d.Set("email", me.Email)

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
