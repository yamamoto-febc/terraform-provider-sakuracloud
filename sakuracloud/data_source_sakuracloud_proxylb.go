// Copyright 2016-2019 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/sacloud"
)

func dataSourceSakuraCloudProxyLB() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudProxyLBRead,

		Schema: map[string]*schema.Schema{
			"name_selectors": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tag_selectors": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vip_failover": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sticky_session": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bind_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"redirect_to_https": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"support_http2": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"response_header": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"health_check": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_loop": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"sorry_server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intermediate_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"additional_certificates": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_cert": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"intermediate_cert": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"icon_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy_networks": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceSakuraCloudProxyLBRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.ProxyLB.FilterBy(key, f)
		}
	}

	res, err := client.ProxyLB.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud ProxyLB resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}
	var data *sacloud.ProxyLB
	targets := res.CommonServiceProxyLBItems

	if rawNameSelector, ok := d.GetOk("name_selectors"); ok {
		selectors := expandStringList(rawNameSelector.([]interface{}))
		var filtered []sacloud.ProxyLB
		for _, a := range targets {
			if hasNames(&a, selectors) {
				filtered = append(filtered, a)
			}
		}
		targets = filtered
	}
	if rawTagSelector, ok := d.GetOk("tag_selectors"); ok {
		selectors := expandStringList(rawTagSelector.([]interface{}))
		var filtered []sacloud.ProxyLB
		for _, a := range targets {
			if hasTags(&a, selectors) {
				filtered = append(filtered, a)
			}
		}
		targets = filtered
	}

	if len(targets) == 0 {
		return filterNoResultErr()
	}
	data = &targets[0]

	d.SetId(data.GetStrID())
	return setProxyLBResourceData(d, client, data)
}
