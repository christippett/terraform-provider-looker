package looker

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testResourceProjectConfig_basic = `
resource "looker_project" "test" {
	name = "test_project"
}
`

func TestResourceProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceProjectConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testResourceProject_exists("looker_project.test"),
				),
			},
		},
	})
}

func testResourceProject_exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Project not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Project ID not set")
		}

		if rs.Primary.Attributes["deploy_key"] == "" {
			return fmt.Errorf("Git deploy key not available for project")
		}

		return nil
	}
}
