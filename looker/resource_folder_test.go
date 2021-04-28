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
				resource "looker_folder" "test" {
					name = "Test Looker Folder"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_folder.test", "name", "Test Looker Folder"),
					resource.TestCheckResourceAttr("looker_folder.test", "parent_id", "1"),
				),
			},
		},
	})
}
