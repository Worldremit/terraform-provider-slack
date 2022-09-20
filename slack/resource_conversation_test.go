package slack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/slack-go/slack"
	"os"
	"testing"
)

func newWithError(version string, commit string) func() (*schema.Provider, error) {
	return func() (*schema.Provider, error) {
		p := New(version, commit)
		return p(), nil
	}
}

func getGithubToken() {

}

func TestAccResourceConversation(t *testing.T) {

	rName := fmt.Sprintf("acc-test-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	rNameRenamed := "acc-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	slackToken := os.Getenv("SLACK_TOKEN")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if slackToken == "" {
				t.Skipf("SLACK_TOKEN not setup skipping")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"slack": newWithError("test", "commit"),
		},
		ProtoV5ProviderFactories:  nil,
		ProtoV6ProviderFactories:  nil,
		ExternalProviders:         nil,
		PreventPostDestroyRefresh: false,
		CheckDestroy:              nil,
		ErrorCheck:                nil,
		Steps: []resource.TestStep{
			{
				Config: tfCode(slackToken, rName, "The topic for test", "purpose", "archive", "false", "false"),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("slack_conversation.this", "name", rName),
					resource.TestCheckResourceAttr("slack_conversation.this", "topic", "The topic for test"),
					resource.TestCheckResourceAttr("slack_conversation.this", "purpose", "purpose"),
					resource.TestCheckResourceAttr("slack_conversation.this", "action_on_destroy", "archive"),
					resource.TestCheckResourceAttr("slack_conversation.this", "is_archived", "false"),
					resource.TestCheckResourceAttr("slack_conversation.this", "is_private", "false"),
					func(state *terraform.State) error {
						res, ok := state.RootModule().Resources["slack_conversation.this"]
						if !ok {
							fmt.Errorf("provider didn't create provider-unit-test")
						}
						client := slack.New(os.Getenv("SLACK_TOKEN"))
						_, err := client.GetConversationInfo(res.Primary.ID, false)

						return err
					},
				),
			},
			{ //renaming all fields
				Config: tfCode(slackToken, rNameRenamed, "The topic for test2", "purpose2", "none", "true", "false"),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("slack_conversation.this", "name", rNameRenamed),
					resource.TestCheckResourceAttr("slack_conversation.this", "topic", "The topic for test2"),
					resource.TestCheckResourceAttr("slack_conversation.this", "purpose", "purpose2"),
					resource.TestCheckResourceAttr("slack_conversation.this", "action_on_destroy", "none"),
					resource.TestCheckResourceAttr("slack_conversation.this", "is_archived", "true"),
					resource.TestCheckResourceAttr("slack_conversation.this", "is_private", "false"),
					func(state *terraform.State) error {
						res, ok := state.RootModule().Resources["slack_conversation.this"]
						if !ok {
							fmt.Errorf("provider didn't create provider-unit-test")
						}
						client := slack.New(os.Getenv("SLACK_TOKEN"))
						_, err := client.GetConversationInfo(res.Primary.ID, false)

						return err
					},
				),
			},
		},
		IDRefreshName:   "",
		IDRefreshIgnore: nil,
	})
}

func tfCode(token, name, topic, purpose, actionOnDestroy, isArchived, isPrivate string) string {
	tf := fmt.Sprintf(`
provider "slack" {
	token="%s"
}

resource slack_conversation this {
  name              = "%s"
  topic             = "%s"
  purpose           = "%s"
  action_on_destroy = "%s" 
  is_archived       = %s
  is_private        = %s
}`, token, name, topic, purpose, actionOnDestroy, isArchived, isPrivate)

	return tf
}
