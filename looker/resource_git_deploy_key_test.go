package looker

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceGitDeployKey(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "looker_project" "test" {
					name = "%s"
				}

				resource "looker_git_deploy_key" "test" {
					project_id = looker_project.test.name
				}
				`, testAccProjectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"looker_git_deploy_key.test", "public_key", regexp.MustCompile("^ssh-rsa ")),
				),
			},
		},
	})
}
