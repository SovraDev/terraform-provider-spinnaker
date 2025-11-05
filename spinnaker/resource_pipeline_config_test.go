package spinnaker

import (
	"testing"
	"fmt"
	"net/http"
	"time"


	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSpinnakerPipelineConfig_basic(t *testing.T) {
	resourceName := "spinnaker_pipeline_config.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSpinnakerPipelineConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "application", rName),
					resource.TestCheckResourceAttr(resourceName, "pipeline", "deploy-my-app"),
					resource.TestCheckResourceAttr(resourceName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.0.type", "webhook"),
					resource.TestCheckResourceAttr(resourceName, "triggers.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "triggers.0.source", "start"),
				),
			},
		},
	})
}

func testAccCheckApplicationConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("PipelineConfig Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PipelineConfig ID is set")
		}

		client := testAccProvider.Meta().(gateConfig).client
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, resp, err := client.ApplicationControllerApi.GetPipelineConfigUsingGET(client.Context, rs.Primary.ID, "deploy-my-app")
			if resp != nil {
				if resp.StatusCode == http.StatusNotFound {
					return resource.RetryableError(fmt.Errorf("pipeline config does not exit"))
				} else if resp.StatusCode != http.StatusOK {
					return resource.NonRetryableError(fmt.Errorf("encountered an error getting pipeline config, status code: %d", resp.StatusCode))
				}
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Unable to find pipeline config after retries: %s", err)
		}
		return nil
	}
}


func testAccSpinnakerPipelineConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "spinnaker_pipeline_config" "test" {
	application = %q
	pipeline    = "deploy-my-app"
	
	triggers = [
		{
			type               = "webhook"
			enabled            = true
			source             = "start"
			payload_constraints = {}
		}
	]
}
`, rName)
}

func testAccSpinnakerPipelineConfig_nondefault(rName string) string {
	return fmt.Sprintf(`
resource "spinnaker_pipeline_config" "test" {
	application = %q
	pipeline    = "deploy-my-app"
	
	triggers = [
		{
			type    = "webhook"
			enabled = true
			source  = "custom-source"
			payload_constraints = {
				"key1" = "value1"
			}
		}
	]
}
`, rName)
}
