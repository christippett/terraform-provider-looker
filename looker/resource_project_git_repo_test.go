package looker

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceProjectGitRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "looker_project" "test" {
						name = "%s"
					}

					resource "looker_project_git_repo" "test" {
						project = looker_project.test.name
						git_remote_url = "git@source.servian.com:looker/demo/terraform-looker-test.git"
						git_service_name = "gitlab"
					}
				`, testAccProjectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"looker_project_git_repo.test", "id", testAccProjectName),
					resource.TestCheckResourceAttr(
						"looker_project.test", "uses_git", "true"),
				),
			},
		},
	})
}
