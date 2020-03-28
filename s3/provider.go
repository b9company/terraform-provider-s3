package s3

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/go-homedir"
)

func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The access key for API operations.",
			},

			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The secret key for API operations.",
			},

			"shared_credentials_file": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "The path to the shared credentials file. If " +
					"not set this defaults to `~/.aws/credentials`.",
			},

			"profile": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "The profile for API operations. If not set, " +
					"the default profile created with `aws configure` will " +
					"be used.",
			},

			"token": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "Session token. A session token is only " +
					"required if you are using temporary security credentials.",
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
				InputDefault: "us-east-1",
				Description: "The region where S3 operations will take " +
					"place. Examples are us-east-1, us-west-2, etc. Defaults " +
					"to `us-east-1`.",
			},

			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				Description: "The endpoint where S3 operations will take " +
					"place.",
			},

			"max_retries": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
				Description: "The maximum number of times an API request is " +
					"being executed. If the request still fails, an error is " +
					"thrown.",
			},

			"force_path_style": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Set this to true to force the request to use " +
					"path-style addressing, i.e., " +
					"http://s3.amazonaws.com/BUCKET/KEY. By default, the S3 " +
					"client will use virtual hosted bucket addressing when " +
					"possible (http://BUCKET.s3.amazonaws.com/KEY). Specific " +
					"to the Amazon S3 service.",
			},

			"aws_signature_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "v4",
				Description: "Signature version to use to sign API requests. " +
					"Default to `v4`.",
				ValidateFunc: validation.StringInSlice([]string{
					"v4",
					"v2",
				}, false),
			},

			"tls_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "Set this to `false` to allow the provider to " +
					"perform HTTP requests. Defaults to `true`.",
			},

			"tls_skip_verify": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Accepts any certificate presented by the " +
					"endpoint and any host name in that certificate.",
			},

			"tls_cert_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"s3_bucket":        resourceS3Bucket(),
			"s3_bucket_object": resourceS3BucketObject(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		AccessKey:           d.Get("access_key").(string),
		SecretKey:           d.Get("secret_key").(string),
		Profile:             d.Get("profile").(string),
		Region:              d.Get("region").(string),
		Endpoint:            d.Get("endpoint").(string),
		MaxRetries:          d.Get("max_retries").(int),
		ForcePathStyle:      d.Get("force_path_style").(bool),
		AwsSignatureVersion: d.Get("aws_signature_version").(string),
		TlsEnabled:          d.Get("tls_enabled").(bool),
		TlsSkipVerify:       d.Get("tls_skip_verify").(bool),
		TlsCertPath:         d.Get("tls_cert_path").(string),
		terraformVersion:    terraformVersion,
	}

	credsPath, err := homedir.Expand(d.Get("shared_credentials_file").(string))
	if err != nil {
		return nil, err
	}
	config.CredsFilename = credsPath

	return config.Client()
}
