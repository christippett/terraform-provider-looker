package looker

import (
	"context"
	"fmt"
	"strconv"

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
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"home_folder_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"personal_folder_id": {
				Type:     schema.TypeInt,
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
			"group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	config := m.(*Config)
	sdk := config.sdk

	user, err := sdk.CreateUser(makeWriteUser(d), "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	userId := *user.Id
	d.SetId(strconv.Itoa(int(userId)))

	// create email credentials
	creds := d.Get("credentials_email").([]interface{})[0]
	_, err = sdk.CreateUserCredentialsEmail(userId, makeCredentialsEmail(creds), "", nil)
	if err != nil {
		// delete user if unable to create email credential
		sdk.DeleteUser(userId, nil)
		return diag.FromErr(err)
	}

	// add user to group(s)
	groups := d.Get("group_ids").(*schema.Set)
	for _, g := range groups.List() {
		u := v3.GroupIdForGroupUserInclusion{
			UserId: &userId,
		}
		sdk.AddGroupUser(int64(g.(int)), u, nil)
	}

	// add user to role(s)
	roles := convertIntSlice(d.Get("role_ids").(*schema.Set).List())
	sdk.SetUserRoles(userId, roles, "", nil)

	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdk := m.(*Config).sdk

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := sdk.User(userId, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ui_state", user.UiState); err != nil {
		return diag.FromErr(err)
	}

	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("display_name", user.DisplayName)
	d.Set("email", user.Email)
	d.Set("locale", user.Locale)
	d.Set("is_disabled", user.IsDisabled)
	d.Set("home_folder_id", user.HomeFolderId)
	d.Set("personal_folder_id", user.PersonalFolderId)
	d.Set("models_dir_validated", user.ModelsDirValidated)
	d.Set("group_ids", user.GroupIds)
	d.Set("role_ids", user.RoleIds)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdk := m.(*Config).sdk

	userId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	// update user
	user := makeWriteUser(d)
	_, err = sdk.UpdateUser(userId, user, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// update email credentials
	creds := d.Get("credentials_email").([]interface{})[0]
	_, err = sdk.UpdateUserCredentialsEmail(userId, makeCredentialsEmail(creds), "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// add user to role(s)
	roles := convertIntSlice(d.Get("role_ids").([]interface{}))
	sdk.SetUserRoles(userId, roles, "", nil)

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	sdk := m.(*Config).sdk

	userId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = sdk.DeleteUser(userId, nil)
	if err != nil {
		return diag.FromErr(err)
	}

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

func makeWriteUser(d *schema.ResourceData) v3.WriteUser {
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	locale := d.Get("locale").(string)
	isDisabled := d.Get("is_disabled").(bool)
	homeFolderID := fmt.Sprint(d.Get("home_folder_id").(int))
	modelsDirValidated := d.Get("models_dir_validated").(bool)
	uiState := d.Get("ui_state").(map[string]interface{})

	user := v3.WriteUser{
		FirstName:          &firstName,
		LastName:           &lastName,
		Locale:             &locale,
		IsDisabled:         &isDisabled,
		HomeFolderId:       &homeFolderID,
		ModelsDirValidated: &modelsDirValidated,
		UiState:            &uiState,
	}
	return user
}

func flattenCredentialsEmail(creds *v3.CredentialsEmail) []interface{} {
	c := make(map[string]interface{})
	c["email"] = creds.Email
	c["forced_password_reset_at_next_login"] = creds.ForcedPasswordResetAtNextLogin

	return []interface{}{c}
}
