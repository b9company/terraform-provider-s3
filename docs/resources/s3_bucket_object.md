# s3_bucket_object Resource

## Example Usage

```hcl
resource "s3_bucket" "example" {
  bucket = "example"
}

resource "s3_bucket_object" "example" {
  bucket = s3_bucket.example.id
  key = "example.txt"

  content = <<EOF
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
EOF
}
```

## Argument Reference

The following arguments are supported:

-> **Note:** If you specify `content_encoding` you are responsible for encoding
the body appropriately. `source`, `content`, and `content_base64` all expect
already encoded/compressed bytes.

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to put the file in.

* `key` - (Required) The name of the object once it is in the bucket.

* `source` - (Optional, conflicts with `content`) The path to a file that will
  be read and uploaded as raw bytes for the object content.

* `content` - (Optional, conflicts with `source`) Literal string value to use as
  the object content, which will be uploaded as UTF-8-encoded text.

* `content_disposition` - (Optional) Specifies presentational information for
  the object.

* `content_encoding` - (Optional) Specifies what content encodings have been
  applied to the object and thus what decoding mechanisms must be applied to
  obtain the media-type referenced by the Content-Type header field.

* `content_language` - (Optional) The language the content is in e.g. en-US or
  en-GB.

* `content_type` - (Optional) A standard MIME type describing the format of the
  object data, e.g. application/octet-stream. All Valid MIME Types are valid for
  this input.

* `storage_class` - (Optional) Specifies the desired Storage Class. Supported
  storage classes are: "`STANDARD`", "`REDUCED_REDUNDANCY`". Defaults to
  "`STANDARD`".

* `tags` - (Optional) A mapping of tags to assign to the object.

* `object_lock_legal_hold_status` - (Optional) The legal hold status that you
  want to apply to the specified object. Valid values are `ON` and `OFF`.

* `object_lock_mode` - (Optional) The object lock retention mode that you want
  to apply to this object. Valid values are `GOVERNANCE` and `COMPLIANCE`.

* `object_lock_retain_until_date` - (Optional) The date and time, in RFC3339
  format, when this object's object lock will expire.

If no content is provided through `source` or `content`, then the object will be
empty.

-> **Note:** Terraform ignores all leading `/`s in the object's `key` and treats
multiple `/`s in the rest of the object's `key` as a single `/`, so values of
`/index.html` and `index.html` correspond to the same S3 object as do
`first//second///third//` and `first/second/third/`.

## Attributes Reference

The following attributes are exported

* `id` - the `key` of the resource supplied above
