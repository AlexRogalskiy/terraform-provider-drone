package drone

import (
	"fmt"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testUserConfig = `
resource "drone_user" "octocat" {
	login = "octocat"
	admin = false
	active = true
}
`

func TestUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_user.octocat",
						"login",
						"octocat",
					),
				),
			},
		},
	})
}

func testUserDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_user" {
			continue
		}

		err := client.UserDelete(resource.Primary.Attributes["login"])

		if err == nil {
			return fmt.Errorf("User still exists: %s", resource.Primary.Attributes["login"])
		}
	}

	return nil
}
