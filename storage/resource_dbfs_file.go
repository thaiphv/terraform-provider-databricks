package storage

import (
	"context"

	"github.com/databrickslabs/terraform-provider-databricks/common"

	"github.com/databrickslabs/terraform-provider-databricks/workspace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceDBFSFile manages files on DBFS
func ResourceDBFSFile() *schema.Resource {
	return common.Resource{
		SchemaVersion: 1,
		Schema: workspace.FileContentSchema(map[string]*schema.Schema{
			"file_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		}),
		Create: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			// TODO: make mandatory DBFS prefix or something to facilitate use for DBFS libraries?...
			path := d.Get("path").(string) // fmt.Sprintf("dbfs:%s", d.Get("path"))
			content, err := workspace.ReadContent(d)
			if err != nil {
				return err
			}
			if err = NewDbfsAPI(ctx, c).Create(path, content, true); err != nil {
				return err
			}
			d.SetId(path)
			return nil
		},
		Read: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			dbfsAPI := NewDbfsAPI(ctx, c)
			fileInfo, err := dbfsAPI.Status(d.Id())
			if err != nil {
				return err
			}
			d.Set("path", fileInfo.Path)
			d.Set("file_size", fileInfo.FileSize)
			return nil
		},
		Delete: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			return NewDbfsAPI(ctx, c).Delete(d.Id(), false)
		},
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    DbfsFileV0(),
				Upgrade: workspace.MigrateV0,
			},
		},
	}.ToResource()
}
