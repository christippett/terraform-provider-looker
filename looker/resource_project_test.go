package looker

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccProjectName = "test_project"

func TestResourceProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "looker_project" "test" {
					name = "%s"
				}
				`, testAccProjectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"looker_project.test",
						"id",
						testAccProjectName,
					),
					resource.TestCheckResourceAttr(
						"looker_project.test",
						"name",
						testAccProjectName,
					),
				),
			},
		},
	})
}
