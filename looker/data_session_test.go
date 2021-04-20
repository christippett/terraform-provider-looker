package looker

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestDataSession(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "looker_session" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					testDataSession_exists("data.looker_session.test"),
				),
			},
		},
	})
}

func testDataSession_exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Session data not available")
		}

		if rs.Primary.Attributes["access_token"] == "" {
			return fmt.Errorf("Access token not available")
		}

		workspaceId := os.Getenv("LOOKER_WORKSPACE_ID")
		if rs.Primary.Attributes["workspace_id"] != workspaceId {
			return fmt.Errorf("Active workspace differs from provider's")
		}

		return nil
	}
}
