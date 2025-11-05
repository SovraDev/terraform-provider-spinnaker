package spinnaker

import (
	"fmt"
	"log"

	"github.com/SovraDev/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePipelineConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline": {
				Type:     schema.TypeString,
				Required: true,
			},
			"trigger": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"payload_constraints": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		Create: resourcePipelineConfigCreate,
		Read:   resourcePipelineConfigRead,
		Update: resourcePipelineConfigUpdate,
		Delete: resourcePipelineConfigDelete,
		Exists: resourcePipelineConfigExists,
	}
}

type pipelineConfigTrigger struct {
	Enabled            bool           `json:"enabled"`
	PayloadConstraints map[string]any `json:"payloadConstraints"`
	Source             string         `json:"source"`
	Type               string         `json:"type"`
}

type pipelineConfigRead struct {
	Triggers []pipelineConfigTrigger `json:"triggers"`
}

func resourcePipelineConfigExists(data *schema.ResourceData, meta any) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("pipeline").(string)

	var p pipelineConfigRead
	if _, err := api.GetPipeline(client, applicationName, pipelineName, &p); err != nil {
		return false, err
	}

	if len(p.Triggers) == 0 {
		return false, nil
	}

	return true, nil
}

func resourcePipelineConfigRead(data *schema.ResourceData, meta any) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("pipeline").(string)

	var p pipelineConfigRead
	log.Println("[DEBUG] Making request to spinnaker")
	log.Printf("Reading pipeline %s from application %s\n", pipelineName, applicationName)
	_, err := api.GetPipeline(client, applicationName, pipelineName, &p)
	if err != nil {
		return err
	}

	// Convert triggers from Spinnaker format to Terraform schema
	triggersOut := make([]map[string]any, 0, len(p.Triggers))
	for _, t := range p.Triggers {
		trigger := map[string]any{
			"enabled": t.Enabled,
			"source":  t.Source,
			"type":    t.Type,
		}
		if t.PayloadConstraints != nil {
			trigger["payload_constraints"] = t.PayloadConstraints
		}
		triggersOut = append(triggersOut, trigger)
	}

	if err := data.Set("triggers", triggersOut); err != nil {
		return fmt.Errorf("could not set triggers for pipeline %s: %w", pipelineName, err)
	}

	return readPipelineConfig(data, applicationName, pipelineName)
}

func readPipelineConfig(data *schema.ResourceData, application string, pipeline string) error {
	data.SetId(fmt.Sprintf("%s:%s", application, pipeline))
	return nil
}

func resourcePipelineConfigCreate(data *schema.ResourceData, meta any) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("pipeline").(string)
	triggersRaw := data.Get("triggers").(*schema.Set).List()

	exist, err := resourcePipelineConfigExists(data, meta)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("pipeline %s in application %s already has triggers configured", pipelineName, applicationName)
	}

	var p pipelineRead
	log.Println("[DEBUG] Making request to spinnaker")
	log.Printf("Reading pipeline %s from application %s\n", pipelineName, applicationName)
	jsonMap, err := api.GetPipeline(client, applicationName, pipelineName, &p)
	if err != nil {
		return err
	}

	// Convert triggers from Terraform schema to Spinnaker format
	triggers := make([]map[string]any, 0, len(triggersRaw))
	for _, triggerRaw := range triggersRaw {
		triggerMap := triggerRaw.(map[string]any)
		trigger := map[string]any{
			"type":    triggerMap["type"],
			"enabled": triggerMap["enabled"],
			"source":  triggerMap["source"],
		}
		if payloadConstraints, ok := triggerMap["payload_constraints"]; ok {
			trigger["payloadConstraints"] = payloadConstraints
		} else {
			trigger["payloadConstraints"] = map[string]any{}
		}
		triggers = append(triggers, trigger)
	}

	// Add triggers to pipeline
	jsonMap["triggers"] = triggers

	if err := api.UpdatePipeline(client, p.ID, jsonMap); err != nil {
		return err
	}

	return resourcePipelineConfigRead(data, meta)
}

func resourcePipelineConfigUpdate(data *schema.ResourceData, meta any) error {
	return resourcePipelineConfigCreate(data, meta)
}

func resourcePipelineConfigDelete(data *schema.ResourceData, meta any) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("pipeline").(string)

	var p pipelineRead
	log.Println("[DEBUG] Making request to spinnaker")
	log.Printf("Reading pipeline %s from application %s\n", pipelineName, applicationName)
	jsonMap, err := api.GetPipeline(client, applicationName, pipelineName, &p)
	if err != nil {
		return err
	}

	delete(jsonMap, "triggers")

	if err := api.UpdatePipeline(client, p.ID, jsonMap); err != nil {
		return err
	}

	return nil
}
