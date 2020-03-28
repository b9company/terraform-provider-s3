package s3

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/minio/minio-go/v6"
)

type Config struct {
	AccessKey     string
	SecretKey     string
	CredsFilename string
	Profile       string
	Token         string
	Region        string
	Endpoint      string
	MaxRetries    int

	ForcePathStyle      bool
	AwsSignatureVersion string

	TlsEnabled    bool
	TlsSkipVerify bool
	TlsCertPath   string

	terraformVersion string
}

type S3Client struct {
	client           *minio.Client
	region           string
	terraformVersion string
}

func (c *Config) Client() (interface{}, error) {
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.StaticProvider{Value: credentials.Value{
				AccessKeyID:     c.AccessKey,
				SecretAccessKey: c.SecretKey,
				SessionToken:    c.Token,
			}},
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{
				Filename: c.CredsFilename,
				Profile:  c.Profile,
			},
		})
	cred, err := creds.Get()
	if err != nil {
		return nil, err
	}

	var client *minio.Client
	switch c.AwsSignatureVersion {
	case "v4":
		client, err = minio.NewV4(c.Endpoint, cred.AccessKeyID, cred.SecretAccessKey, c.TlsEnabled)
	case "v2":
		client, err = minio.NewV2(c.Endpoint, cred.AccessKeyID, cred.SecretAccessKey, c.TlsEnabled)
	default:
		client, err = minio.New(c.Endpoint, cred.AccessKeyID, cred.SecretAccessKey, c.TlsEnabled)
	}
	if err != nil {
		return nil, err
	}

	client.SetCustomTransport(c.getTlsConfig())

	return &S3Client{
		region:           c.Region,
		client:           client,
		terraformVersion: c.terraformVersion,
	}, nil
}

func (c *Config) getTlsConfig() http.RoundTripper {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	certs, err := ioutil.ReadFile(c.TlsCertPath)
	if err != nil {
		log.Println("[DEBUG] No certs file, ignoring")
	}

	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("[DEBUG] No certs appended, using system certs only")
	}

	return &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.TlsSkipVerify,
			RootCAs:            rootCAs,
		},
	}
}
