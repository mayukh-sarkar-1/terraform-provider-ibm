---
layout: "ibm"
page_title: "IBM : reserved_ip"
description: |-
  Lists all the info in reserved IP for Subnet.
---

# ibm\_is_subnet_reserved_ips

Import the details of all the Reserved IPs in a Subnet as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "ibm_is_subnet_reserved_ips" "data_reserved_ips" {
  subnet = ibm_is_subnet.test_subnet.id
}
```

## Argument Reference

The following arguments are supported as inputs/request params:

* `subnet` - (Required, string) The id for the Subnet.


## Attribute Reference

The following attributes are exported as output/response:

* `id` - The id for the all the reserved ID (current timestamp)
* `limit` - The number of reserved IPs to list
* `reserved_ips` - The unique reference for the reserved IP
* `sort` - The keyword on which all the reserved IPs are sorted
* `subnet` - The id for the subnet for the reserved IP
* `total_count` - The number of reserved IP in the subnet
