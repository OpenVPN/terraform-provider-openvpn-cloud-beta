---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openvpn_host Resource - terraform-provider-openvpn-cloud-beta"
subcategory: ""
description: |-
  
---

# openvpn_host (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connector` (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--connector))
- `description` (String)
- `domain` (String)
- `internet_access` (String)
- `name` (String)

### Read-Only

- `id` (String) The ID of this resource.
- `system_subnets` (List of String)

<a id="nestedblock--connector"></a>
### Nested Schema for `connector`

Required:

- `description` (String)
- `name` (String)
- `vpn_region_id` (String)

Read-Only:

- `id` (String) The ID of this resource.
- `ip_v4_address` (String)
- `ip_v6_address` (String)
- `profile` (String, Sensitive)


