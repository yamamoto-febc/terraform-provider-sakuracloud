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
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/cleanup"
)

func resourceSakuraCloudInternet() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudInternetCreate,
		Read:   resourceSakuraCloudInternetRead,
		Update: resourceSakuraCloudInternetUpdate,
		Delete: resourceSakuraCloudInternetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"netmask": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(types.AllowInternetNetworkMaskLen()),
				Default:      28,
			},
			"band_width": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(types.AllowInternetBandWidth()),
				Default:      100,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "target SakuraCloud zone",
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"network_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_prefix_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipv6_network_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudInternetCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	builder := expandInternetBuilder(d, client)

	internet, err := builder.Build(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Internet is failed: %s", err)
	}

	d.SetId(internet.ID.String())
	return resourceSakuraCloudInternetRead(d, meta)
}

func resourceSakuraCloudInternetRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	internetOp := sacloud.NewInternetOp(client)

	internet, err := internetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Internet[%s]: %s", d.Id(), err)
	}
	return setInternetResourceData(ctx, d, client, internet)
}

func resourceSakuraCloudInternetUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	internetOp := sacloud.NewInternetOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	internet, err := internetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Internet[%s]: %s", d.Id(), err)
	}

	builder := expandInternetBuilder(d, client)
	internet, err = builder.Update(ctx, zone, internet.ID)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Internet[%s] is failed: %s", d.Id(), err)
	}

	d.SetId(internet.ID.String()) // 帯域変更後はIDが変更になるため
	return resourceSakuraCloudInternetRead(d, meta)
}

func resourceSakuraCloudInternetDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	internetOp := sacloud.NewInternetOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	internet, err := internetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Internet[%s]: %s", d.Id(), err)
	}

	if err := waitForDeletionBySwitchID(ctx, client, zone, internet.Switch.ID); err != nil {
		return fmt.Errorf("waiting deletion is failed: Internet[%s] still used by others: %s", internet.ID, err)
	}

	if err := cleanup.DeleteInternet(ctx, internetOp, zone, internet.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Internet[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setInternetResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Internet) error {
	swOp := sacloud.NewSwitchOp(client)
	zone := getZone(d, client)

	sw, err := swOp.Read(ctx, zone, data.Switch.ID)
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Switch[%s]: %s", data.Switch.ID, err)
	}

	var serverIDs []string
	if sw.ServerCount > 0 {
		servers, err := swOp.GetServers(ctx, zone, sw.ID)
		if err != nil {
			return fmt.Errorf("could not find SakuraCloud Servers: %s", err)
		}
		for _, s := range servers.Servers {
			serverIDs = append(serverIDs, s.ID.String())
		}
	}

	var enableIPv6 bool
	var ipv6Prefix, ipv6NetworkAddress string
	var ipv6PrefixLen int
	if len(data.Switch.IPv6Nets) > 0 {
		enableIPv6 = true
		ipv6Prefix = data.Switch.IPv6Nets[0].IPv6Prefix
		ipv6PrefixLen = data.Switch.IPv6Nets[0].IPv6PrefixLen
		ipv6NetworkAddress = fmt.Sprintf("%s/%d", ipv6Prefix, ipv6PrefixLen)
	}

	d.Set("name", data.Name)                                    // nolint
	d.Set("icon_id", data.IconID.String())                      // nolint
	d.Set("description", data.Description)                      // nolint
	d.Set("netmask", data.NetworkMaskLen)                       // nolint
	d.Set("band_width", data.BandWidthMbps)                     // nolint
	d.Set("switch_id", sw.ID.String())                          // nolint
	d.Set("network_address", sw.Subnets[0].NetworkAddress)      // nolint
	d.Set("gateway", sw.Subnets[0].DefaultRoute)                // nolint
	d.Set("min_ip_address", sw.Subnets[0].AssignedIPAddressMin) // nolint
	d.Set("max_ip_address", sw.Subnets[0].AssignedIPAddressMax) // nolint
	d.Set("enable_ipv6", enableIPv6)                            // nolint
	d.Set("ipv6_prefix", ipv6Prefix)                            // nolint
	d.Set("ipv6_prefix_len", ipv6PrefixLen)                     // nolint
	d.Set("ipv6_network_address", ipv6NetworkAddress)           // nolint
	d.Set("zone", zone)                                         // nolint
	if err := d.Set("ip_addresses", sw.Subnets[0].GetAssignedIPAddresses()); err != nil {
		return err
	}
	if err := d.Set("server_ids", serverIDs); err != nil {
		return err
	}
	return d.Set("tags", flattenTags(data.Tags))
}
