package s3

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/minio/minio-go/v6"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

func resourceS3BucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceS3BucketObjectPut,
		Read:   resourceS3BucketObjectRead,
		Update: resourceS3BucketObjectUpdate,
		Delete: resourceS3BucketObjectDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"STANDARD",
					"REDUCED_REDUNDANCY",
				}, false),
			},

			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},

			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
			},

			"cache_control": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_language": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_disposition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),

			"object_lock_legal_hold_status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(minio.LegalHoldEnabled),
					string(minio.LegalHoldDisabled),
				}, false),
			},

			"object_lock_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(minio.Governance),
					string(minio.Compliance),
				}, false),
			},

			"object_lock_retain_until_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
		},
	}
}

func resourceS3BucketObjectPut(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*S3Client).client

	var body io.ReadSeeker
	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error expanding homedir in source (%s)", source))
		}
		file, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error opening S3 bucket object source (%s)", path))
		}

		body = file
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Error closing S3 bucket object source (%s): %s", path, err)
			}
		}()
	} else if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		body = bytes.NewReader([]byte(content))
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	options := minio.PutObjectOptions{}

	if v, ok := d.GetOk("storage_class"); ok {
		options.StorageClass = v.(string)
	}

	if v, ok := d.GetOk("cache_control"); ok {
		options.CacheControl = v.(string)
	}

	if v, ok := d.GetOk("content_type"); ok {
		options.ContentType = v.(string)
	}

	if v, ok := d.GetOk("content_encoding"); ok {
		options.ContentEncoding = v.(string)
	}

	if v, ok := d.GetOk("content_language"); ok {
		options.ContentLanguage = v.(string)
	}

	if v, ok := d.GetOk("content_disposition"); ok {
		options.ContentDisposition = v.(string)
	}

	if v := d.Get("tags").(map[string]interface{}); len(v) > 0 {
		options.UserTags = MakeUserTags(v)
	}

	if v, ok := d.GetOk("website_redirect"); ok {
		options.WebsiteRedirectLocation = v.(string)
	}

	if v, ok := d.GetOk("object_lock_legal_hold_status"); ok {
		options.LegalHold = parseObjectLockLegalHoldStatus(v.(string))
	}

	if v, ok := d.GetOk("object_lock_mode"); ok {
		options.Mode = parseObjectLockMode(v.(string))
	}

	if v, ok := d.GetOk("object_lock_retain_until_date"); ok {
		options.RetainUntilDate = parseObjectLockDeadline(v.(string))
	}

	if _, err := client.PutObject(bucket, key, body, -1, options); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to put object in S3 bucket (%s): %s", bucket, key))
	}

	d.SetId(key)

	return resourceS3BucketObjectRead(d, meta)
}

func resourceS3BucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*S3Client).client

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	log.Printf("[DEBUG] Reading S3 Bucket Object info")

	stat, err := client.StatObject(bucket, key, minio.StatObjectOptions{})
	if err != nil {
		// TODO: find out what caused the error, e.g. if the object doesn't
		// exist, mark it as destroyed.
		return err
	}
	d.Set("content_type", stat.ContentType)

	mode, deadline, err := client.GetObjectRetention(bucket, key, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to retrieve retention info for S3 Bucket (%s) Object (%s)", bucket, key))
	}
	d.Set("object_lock_mode", string(*mode))
	if !deadline.IsZero() {
		d.Set("object_lock_retain_until_date", deadline.String())
	}

	tags, err := client.GetObjectTagging(bucket, key)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to list tags for S3 Bucket (%s) Object (%s)", bucket, key))
	}
	d.Set("tags", tags)

	return nil
}
func resourceS3BucketObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	// Changes to any of these attributes requires creation of a new object
	// version (if bucket is versioned).
	for _, key := range []string{
		"storage_class",
		"source",
		"content",
		"cache_control",
		"content_type",
		"content_encoding",
		"content_language",
		"content_disposition",
	} {
		if d.HasChange(key) {
			return resourceS3BucketObjectPut(d, meta)
		}
	}

	client := meta.(*S3Client).client

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	if d.HasChange("object_lock_mode") || d.HasChange("object_lock_retain_until_date") {
		opts := minio.PutObjectRetentionOptions{}
		opts.RetainUntilDate = parseObjectLockDeadline(d.Get("object_lock_retain_until_date").(string))
		opts.Mode = parseObjectLockMode(d.Get("object_lock_mode").(string))

		// Bypass required to lower or clear retain-until date.
		if d.HasChange("object_lock_retain_until_date") {
			oraw, nraw := d.GetChange("object_lock_retain_until_date")
			o := parseObjectLockDeadline(oraw.(string))
			n := parseObjectLockDeadline(nraw.(string))
			if n == nil || (o != nil && n.Before(*o)) {
				opts.GovernanceBypass = true
			}
		}

		if err := client.PutObjectRetention(bucket, key, opts); err != nil {
			return errors.Wrap(err, "Failed putting S3 object lock retention")
		}
	}

	if d.HasChange("tags") {
		_, n := d.GetChange("tags")
		log.Printf("[DEBUG] Updating S3 Bucket user tags")
		if err := client.PutObjectTagging(bucket, key, MakeUserTags(n)); err != nil {
			return errors.Wrap(err, "Failed to update tags")
		}
	}

	return resourceS3BucketObjectRead(d, meta)
}

func resourceS3BucketObjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*S3Client).client

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	if err := client.RemoveObject(bucket, key); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to delete s3://%s/%s", bucket, key))
	}

	return nil
}

func parseObjectLockLegalHoldStatus(v string) minio.LegalHoldStatus {
	var status minio.LegalHoldStatus
	status = minio.LegalHoldStatus(v)
	if !status.IsValid() {
		return minio.LegalHoldStatus("")
	}

	return status
}

func parseObjectLockMode(v string) *minio.RetentionMode {
	var mode minio.RetentionMode
	mode = minio.RetentionMode(v)
	if !mode.IsValid() {
		return nil
	}

	return &mode
}

func parseObjectLockDeadline(v string) *time.Time {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil
	}

	return &t
}
