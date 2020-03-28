package s3

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/minio/minio-go/v6"
	"github.com/pkg/errors"
)

func resourceS3Bucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceS3BucketCreate,
		Read:   resourceS3BucketRead,
		Update: resourceS3BucketUpdate,
		Delete: resourceS3BucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 63),
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"versioning": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"object_lock_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
		},
	}
}

func resourceS3BucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*S3Client).client

	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		bucket = resource.UniqueId()
	}
	d.Set("bucket", bucket)

	log.Printf("[DEBUG] S3 bucket create: %s", bucket)

	region := meta.(*S3Client).region

	log.Printf("[DEBUG] S3 bucket create: %s, using region: %s", bucket, region)

	var err error
	if d.Get("object_lock_enabled").(bool) {
		err = client.MakeBucketWithObjectLock(bucket, region)
	} else {
		err = client.MakeBucket(bucket, region)
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Unable to create bucket [%s] in region [%s]", bucket, region))
	}

	log.Printf("[DEBUG] Created bucket: [%s] in region: [%s]", bucket, region)

	// Assign the bucket name as the resource ID
	d.SetId(bucket)
	return resourceS3BucketUpdate(d, meta)
}

func resourceS3BucketRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*S3Client).client
	bucket := d.Id()

	exists, _ := client.BucketExists(bucket)
	if !exists {
		log.Printf("[WARN] S3 Bucket (%s) not found, removing from state", bucket)
		d.SetId("")
		return nil
	}

	if _, ok := d.GetOk(bucket); !ok {
		d.Set(bucket, d.Id())
	}

	return nil
}

func resourceS3BucketUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*S3Client).client

	if d.HasChange("versioning") {
		if err := resourceS3BucketVersioningUpdate(client, d); err != nil {
			return err
		}
	}

	return resourceS3BucketRead(d, meta)
}

func resourceS3BucketDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*S3Client).client
	bucket := d.Id()

	log.Printf("[DEBUG] S3 Delete Bucket: %s", bucket)

	// TODO: Implement force_destroy feature, which removes the bucket even if
	// it's not empty.
	err := client.RemoveBucket(bucket)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error deleting S3 Bucket (%s)", d.Id()))
	}

	return nil
}

func resourceS3BucketVersioningUpdate(client *minio.Client, d *schema.ResourceData) (err error) {
	bucket := d.Get("bucket").(string)
	v := d.Get("versioning").([]interface{})

	setVersioning := false

	if len(v) > 0 {
		c := v[0].(map[string]interface{})

		if c["enabled"].(bool) {
			setVersioning = true
		}
	}

	log.Printf("[DEBUG] S3 set bucket versioning: %#v", setVersioning)

	if setVersioning {
		err = client.EnableVersioning(bucket)
	} else {
		err = client.DisableVersioning(bucket)
	}

	return
}
