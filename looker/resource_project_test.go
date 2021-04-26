package looker

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					resource.TestCheckResourceAttr(
						"looker_project.test",
						"id",
						"test_project",
					),
					resource.TestCheckResourceAttr(
						"looker_project.test",
						"name",
						"test_project",
					),
				),
			},
		},
	})
}
