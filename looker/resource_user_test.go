package looker

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	v3 "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

func TestAccCLookerUserBasics(t *testing.T) {
	firstName := "John"
	lastName := "Smith"
	email := "john.smith@example.com"

	resource.Test(t, resource.TestCase{
		IsUnitTest:   true,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLookerUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLookerUserBasic(firstName, lastName, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLookerExists("looker_user.new"),
				),
			},
		},
	})
}

func testAccCheckLookerUserDestroy(s *terraform.State) error {
	sdk := testAccProvider.Meta().(*v3.LookerSDK)

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

func testAccCheckLookerUserBasic(firstName, lastName, email string) string {
	return fmt.Sprintf(`
	resource "looker_user" "new" {
		first_name = "%s"
		last_name = "%s"

		credentials_email {
			email = "%s"
		}
	}
	`, firstName, lastName, email)
}

func testAccCheckLookerExists(n string) resource.TestCheckFunc {
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
