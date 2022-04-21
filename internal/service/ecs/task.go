package ecs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_ecs_task", name="Task")
// @Tags(identifierAttribute="arn")
func ResourceTask() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAwsEcsTaskUpdate,
		ReadWithoutTimeout:   resourceAwsEcsTaskRead,
		UpdateWithoutTimeout: resourceAwsEcsTaskUpdate,
		DeleteWithoutTimeout: resourceAwsEcsTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cluster_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"task_definition_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"launch_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"containers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interfaces": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attachment_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ipv4_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"health_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAwsEcsTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*conns.AWSClient).ECSClient(ctx)

	taskArn := d.Get("arn").(string)
	clusterArn := d.Get("cluster_arn").(string)
	log.Printf("[DEBUG] Reading task %s", d.Id())
	out, err := client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Tasks:   []string{taskArn},
		Cluster: aws.String(clusterArn),
	})
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading ECS Task (%s): %s", d.Id(), err)
	}

	task := out.Tasks[0] // DescribeTasks returns a list but we're only passing in a single task ARN

	d.SetId(taskArn)
	d.Set("arn", &taskArn)
	d.Set("cluster_arn", task.ClusterArn)
	d.Set("task_definition_arn", task.TaskDefinitionArn)
	d.Set("group", task.Group)
	d.Set("launch_type", task.LaunchType)
	d.Set("health_status", task.HealthStatus)

	if err := d.Set("containers", flattenEcsTaskContainers(task.Containers)); err != nil {
		return sdkdiag.AppendErrorf(diags, "reading ECS Task (%s): failed to set flattened ECS task containers: %s", d.Id(), err)
	}

	if err := d.Set("attachments", flattenEcsTaskAttachments(task.Attachments)); err != nil {
		return sdkdiag.AppendErrorf(diags, "reading ECS Task (%s): failed to set flattened ECS task attachments: %s", d.Id(), err)
	}

	return nil
}

func resourceAwsEcsTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceAwsEcsTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func flattenEcsTaskContainers(config []types.Container) []map[string]interface{} {
	containers := make([]map[string]interface{}, 0, len(config))

	for _, raw := range config {
		item := make(map[string]interface{})
		item["container_arn"] = *raw.ContainerArn
		item["task_arn"] = *raw.TaskArn
		item["network_interfaces"] = flattenEcsTaskContainerNetworkInterfaces(raw.NetworkInterfaces)

		containers = append(containers, item)
	}

	return containers
}

func flattenEcsTaskContainerNetworkInterfaces(config []types.NetworkInterface) []map[string]interface{} {
	networkInterfaces := make([]map[string]interface{}, 0, len(config))

	for _, raw := range config {
		item := make(map[string]interface{})
		item["attachment_id"] = *raw.AttachmentId
		item["private_ipv4_address"] = *raw.PrivateIpv4Address

		networkInterfaces = append(networkInterfaces, item)
	}

	return networkInterfaces
}

func flattenEcsTaskAttachments(config []types.Attachment) []map[string]interface{} {
	attachments := make([]map[string]interface{}, 0, len(config))

	for _, raw := range config {
		item := make(map[string]interface{})
		item["id"] = *raw.Id
		item["type"] = *raw.Type
		item["details"] = flattenEcsTaskAttachmentDetails(raw.Details)

		attachments = append(attachments, item)
	}

	return attachments
}

func flattenEcsTaskAttachmentDetails(config []types.KeyValuePair) []map[string]interface{} {
	details := make([]map[string]interface{}, 0, len(config))

	for _, raw := range config {
		item := make(map[string]interface{})
		item["name"] = *raw.Name
		item["value"] = *raw.Value

		details = append(details, item)
	}

	return details
}
