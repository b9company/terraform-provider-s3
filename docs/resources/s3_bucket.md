# s3_bucket Resource

## Example Usage

```hcl
resource "s3_bucket" "example" {
  bucket = "example"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Optional, Forces new resource) The name of the bucket. If omitted,
  Terraform will assign a random, unique name.

* `region` - (Optional) If specified, the region this bucket should reside in.
  Otherwise, the region used by the callee.

* `versioning` - (Optional) (Documented below).

The `versioning` object supports the following:

* `enabled` - (Optional) Enable versioning. Once you version-enable a bucket, it
  can never return to an unversioned state.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the bucket.

## Import

S3 bucket can be imported using the `bucket`, e.g.

```
$ terraform import s3_bucket.bucket bucket-name
```

