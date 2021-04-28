package looker

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceProject(t *testing.T) {
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

					data "looker_project" "test" {
						name = looker_project.test.name
					}
				`, testAccProjectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.looker_project.test", "name", testAccProjectName),
				),
			},
		},
	})
}
