package looker

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"credentials_email": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"forced_password_reset_at_next_login": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"home_folder_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"models_dir_validated": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ui_state": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	config := m.(*Config)
	sdk := config.sdk

	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	locale := d.Get("locale").(string)
	isDisabled := d.Get("is_disabled").(bool)
	homeFolderID := d.Get("home_folder_id").(string)
	modelsDirValidated := d.Get("models_dir_validated").(bool)
	uiState := d.Get("ui_state").(map[string]interface{})

	userDetails := v3.WriteUser{
		FirstName:          &firstName,
		LastName:           &lastName,
		Locale:             &locale,
		IsDisabled:         &isDisabled,
		HomeFolderId:       &homeFolderID,
		ModelsDirValidated: &modelsDirValidated,
		UiState:            &uiState,
	}

	user, err := sdk.CreateUser(userDetails, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	userID := *user.Id

	creds := d.Get("credentials_email").([]interface{})[0]
	_, err = sdk.CreateUserCredentialsEmail(userID, makeCredentialsEmail(creds), "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(userID)))

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdk := m.(*Config).sdk

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := sdk.User(userID, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ui_state", user.UiState); err != nil {
		return diag.FromErr(err)
	}

	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("locale", user.Locale)
	d.Set("is_disabled", user.IsDisabled)
	d.Set("home_folder_id", user.HomeFolderId)
	d.Set("models_dir_validated", user.ModelsDirValidated)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdk := m.(*Config).sdk

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var writeUser v3.WriteUser

	if d.HasChange("first_name") {
		writeUser.FirstName = d.Get("first_name").(*string)
	}

	if d.HasChange("last_name") {
		writeUser.LastName = d.Get("last_name").(*string)
	}

	_, err = sdk.UpdateUser(userID, writeUser, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdk := m.(*Config).sdk

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = sdk.DeleteUser(userID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func makeCredentialsEmail(creds interface{}) v3.WriteCredentialsEmail {
	credentials := creds.(map[string]interface{})
	email := credentials["email"].(string)
	forcedReset := credentials["forced_password_reset_at_next_login"].(bool)
	return v3.WriteCredentialsEmail{
		Email:                          &email,
		ForcedPasswordResetAtNextLogin: &forcedReset,
	}
}

func flattenCredentialsEmail(creds *v3.CredentialsEmail) []interface{} {
	c := make(map[string]interface{})
	c["email"] = creds.Email
	c["forced_password_reset_at_next_login"] = creds.ForcedPasswordResetAtNextLogin

	return []interface{}{c}
}
