package looker

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLookerUserBasics(t *testing.T) {
	firstName := "John"
	lastName := "Smith"
	email := "john.smith@example.com"

	resource.Test(t, resource.TestCase{
		IsUnitTest:   true,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: generateLookerUserConfig(firstName, lastName, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("looker_user.new"),
				),
			},
		},
	})
}

func generateLookerUserConfig(firstName, lastName, email string) string {
	return fmt.Sprintf(`
	resource "looker_user" "new" {
		first_name = "%s"
		last_name = "%s"
		group_ids = [1]
		role_ids = [2]

		credentials_email {
			email = "%s"
		}
	}
	`, firstName, lastName, email)
}

func testAccCheckUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No UserID set")
		}

		return nil
	}
}

func testAccCheckUserDestroy(s *terraform.State) error {
	sdk := testAccProvider.Meta().(*Config).sdk

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_user" {
			continue
		}

		userID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}

		_, err = sdk.DeleteUser(userID, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
