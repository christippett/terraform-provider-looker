package looker

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFolder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "looker_folder" "parent" {
					name = "Terraform Test Folder (parent)"

					content_metadata {
						inherits = false
					}
				}

				resource "looker_folder" "child" {
					name = "Terraform Test Folder (child)"
					parent_id = looker_folder.parent.id
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_folder.parent", "name", "Terraform Test Folder (parent)"),
					resource.TestCheckResourceAttr("looker_folder.parent", "parent_id", "1"),
					resource.TestCheckResourceAttr("looker_folder.child", "name", "Terraform Test Folder (child)"),
					resource.TestCheckResourceAttrPair("looker_folder.parent", "id", "looker_folder.child", "parent_id"),
				),
			},
			{
				Config: `
				resource "looker_folder" "parent" {
					name = "Terraform Test Folder (parent)"
				}

				resource "looker_folder" "child" {
					name = "Terraform Test Directory (child)"
					parent_id = "1"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_folder.child", "name", "Terraform Test Directory (child)"),
					resource.TestCheckResourceAttr("looker_folder.child", "parent_id", "1"),
				),
			},
		},
	})
}
