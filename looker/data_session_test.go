package looker

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSession(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "looker_session" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.looker_session.test", "workspace_id", "dev"),
					resource.TestCheckResourceAttrSet("data.looker_session.test", "access_token"),
					resource.TestCheckResourceAttrSet("data.looker_session.test", "user_id"),
					resource.TestCheckResourceAttrSet("data.looker_session.test", "email"),
				),
			},
		},
	})
}
