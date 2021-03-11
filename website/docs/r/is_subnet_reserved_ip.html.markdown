---
layout: "ibm"
page_title: "IBM : ibm_is_subnet_reserved_ip"
description: |-
  Manages IBM Subnet reserved IP
---

# ibm_is_virtual_endpoint_gateway

Provides a subnet reserved IP resource. This allows Subnet reserved IP to be created, updated, and deleted.

## Example Usage

In the following example, you can create a Reserved IP:

```hcl
    // Create a VPC
    resource "ibm_is_vpc" "vpc1" {
        name = "my-vpc"
    }

    // Create a subnet
    resource "ibm_is_subnet" "subnet1" {
        name                     = "my-subnet"
        vpc                      = ibm_is_vpc.vpc1.id
        zone                     = "us-south-1"
        total_ipv4_address_count = 256
    }

    // Create the resrved IP in the following ways

    // Only with Subnet ID
    resource "ibm_is_subnet_reserved_ip" "res_ip" {
        subnet = ibm_is_subnet.subnet1.id
    }

    // Subnet ID with a given name
    resource "ibm_is_subnet_reserved_ip" "res_ip_name" {
        subnet = ibm_is_subnet.subnet1.id
        name = "my-subnet"
    }

    // Subnet ID with auto_delete
    resource "ibm_is_subnet_reserved_ip" "res_ip_auto_delete" {
        subnet = ibm_is_subnet.subnet1.id
        auto_delete = true
    }

    // Subnet ID with both name and auto_delete
    resource "ibm_is_subnet_reserved_ip" "res_ip_auto_delete_name" {
        subnet = ibm_is_subnet.subnet1.id
        name = "my-subnet"
        auto_delete = true
    }
```

## Timeouts
`ibm_is_subnet_reserved_ip` provides the following Timeouts. However there is no wait performed as `status` was present in the response

* `create` - (Default 10 minutes) Used for creating reserved IP.
* `delete` - (Default 10 minutes) Used for deleting reserved IP.

## Argument Reference

The following arguments are supported:

* `subnet` - (Required, Forces new resource, string) The subnet id for the reserved IP.
* `name` - (Optional, string) The name of the reserved IP.
    **NOTE**: Raise error if name is given with a prefix `ibm-`.
* `auto_delete` - (Optional, boolean) If reserved IP is auto deleted.