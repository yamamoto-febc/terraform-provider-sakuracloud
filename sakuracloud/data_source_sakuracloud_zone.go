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

func dataSourceSakuraCloudZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudZoneRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateZone(allowZones),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceSakuraCloudZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	zoneSlug := client.Zone
	if v, ok := d.GetOk("name"); ok {
		zoneSlug = v.(string)
	}

	res, err := client.GetZoneAPI().Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Zone resource: %s", err)
	}
	if res == nil || len(res.Zones) == 0 {
		return filterNoResultErr()
	}
	var data *sacloud.Zone

	for _, z := range res.Zones {
		if z.Name == zoneSlug {
			data = &z
			break
		}
	}
	if data == nil {
		return filterNoResultErr()
	}

	d.SetId(data.GetStrID())
	d.Set("name", data.Name)
	d.Set("zone_id", data.GetStrID())
	d.Set("description", data.Description)
	d.Set("region_id", data.Region.GetStrID())
	d.Set("region_name", data.GetRegionName())
	d.Set("dns_servers", data.Region.NameServers)

	return nil
}
