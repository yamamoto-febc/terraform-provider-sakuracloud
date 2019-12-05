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
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudLoadBalancerVIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudLoadBalancerVIPCreate,
		Read:   resourceSakuraCloudLoadBalancerVIPRead,
		Delete: resourceSakuraCloudLoadBalancerVIPDelete,
		Update: resourceSakuraCloudLoadBalancerVIPUpdate,
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"vip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				ForceNew:     true,
			},
			"delay_loop": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(10, 2147483647),
				Default:      10,
			},
			"sorry_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"servers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 40,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Required: true,
						},
						"check_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(types.LoadBalancerHealthCheckProtocolsStrings(), false),
						},
						"check_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"check_status": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceSakuraCloudLoadBalancerVIPCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)
	lbID := d.Get("load_balancer_id").(string)

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	lb, err := lbOp.Read(ctx, zone, sakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud LoadBalancer resource: %s", err)
	}

	vip := expandLoadBalancerVIP(d)
	if r := findLoadBalancerVIPMatch(lb, vip); r != nil {
		return fmt.Errorf("already exists: LoadBalancer VIP: %s:%d", r.VirtualIPAddress, r.Port)
	}
	vips := append(lb.VirtualIPAddresses, vip)

	lb, err = lbOp.Update(ctx, zone, lb.ID, &sacloud.LoadBalancerUpdateRequest{
		Name:               lb.Name,
		Description:        lb.Description,
		Tags:               lb.Tags,
		IconID:             lb.IconID,
		VirtualIPAddresses: vips,
		SettingsHash:       lb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud LoadBalancerVIP is failed: %s", err)
	}

	d.SetId(loadBalancerVIPIDHash(lbID, vip))
	return resourceSakuraCloudLoadBalancerVIPRead(d, meta)
}

func resourceSakuraCloudLoadBalancerVIPRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)
	lbID := d.Get("load_balancer_id").(string)

	lb, err := lbOp.Read(ctx, zone, sakuraCloudID(lbID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer: %s", err)
	}

	src := expandLoadBalancerVIP(d)
	vip := findLoadBalancerVIPMatch(lb, src)
	if vip == nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	d.Set("vip", vip.VirtualIPAddress)
	d.Set("port", vip.Port.Int())
	if err := d.Set("servers", flattenLoadBalancerServers(vip)); err != nil {
		return err
	}
	d.Set("delay_loop", vip.DelayLoop.Int())
	d.Set("sorry_server", vip.SorryServer)
	d.Set("zone", getZone(d, client))
	return nil
}

func resourceSakuraCloudLoadBalancerVIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)
	lbID := d.Get("load_balancer_id").(string)

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	lb, err := lbOp.Read(ctx, zone, sakuraCloudID(lbID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer: %s", err)
	}

	src := expandLoadBalancerVIP(d)
	vip := findLoadBalancerVIPMatch(lb, src)
	if vip == nil {
		d.SetId("")
		return nil
	}

	vip.DelayLoop = src.DelayLoop
	vip.SorryServer = src.SorryServer

	lb, err = lbOp.Update(ctx, zone, lb.ID, &sacloud.LoadBalancerUpdateRequest{
		Name:               lb.Name,
		Description:        lb.Description,
		Tags:               lb.Tags,
		IconID:             lb.IconID,
		VirtualIPAddresses: lb.VirtualIPAddresses,
		SettingsHash:       lb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud LoadBalancerVIP is failed: %s", err)
	}

	return resourceSakuraCloudLoadBalancerVIPRead(d, meta)
}

func resourceSakuraCloudLoadBalancerVIPDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)
	lbID := d.Get("load_balancer_id").(string)

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	lb, err := lbOp.Read(ctx, zone, sakuraCloudID(lbID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer: %s", err)
	}

	src := expandLoadBalancerVIP(d)
	var vips []*sacloud.LoadBalancerVirtualIPAddress
	for _, v := range lb.VirtualIPAddresses {
		if !isSameLoadBalancerVIP(src, v) {
			vips = append(vips, v)
		}
	}

	lb, err = lbOp.Update(ctx, zone, lb.ID, &sacloud.LoadBalancerUpdateRequest{
		Name:               lb.Name,
		Description:        lb.Description,
		Tags:               lb.Tags,
		IconID:             lb.IconID,
		VirtualIPAddresses: vips,
		SettingsHash:       lb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("deleting SakuraCloud LoadBalancerVIP is failed: %s", err)
	}
	return nil
}

func findLoadBalancerVIPMatch(lb *sacloud.LoadBalancer, vip *sacloud.LoadBalancerVirtualIPAddress) *sacloud.LoadBalancerVirtualIPAddress {
	for _, v := range lb.VirtualIPAddresses {
		if isSameLoadBalancerVIP(v, vip) {
			return v
		}
	}
	return nil
}

func isSameLoadBalancerVIP(v1, v2 *sacloud.LoadBalancerVirtualIPAddress) bool {
	return v1.VirtualIPAddress == v2.VirtualIPAddress && v1.Port == v2.Port
}

func loadBalancerVIPIDHash(loadBalancerID string, s *sacloud.LoadBalancerVirtualIPAddress) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", loadBalancerID))
	buf.WriteString(fmt.Sprintf("%s-", s.VirtualIPAddress))
	buf.WriteString(fmt.Sprintf("%s", s.Port.String()))
	return buf.String()
}

func flattenLoadBalancerServers(vip *sacloud.LoadBalancerVirtualIPAddress) []interface{} {
	var servers []interface{}
	for _, s := range vip.Servers {
		servers = append(servers, flattenLoadBalancerServer(s))
	}
	return servers
}

func flattenLoadBalancerServer(s *sacloud.LoadBalancerServer) interface{} {
	return map[string]interface{}{
		"ipaddress":      s.IPAddress,
		"check_protocol": s.HealthCheck.Protocol.String(),
		"check_path":     s.HealthCheck.Path,
		"check_status":   s.HealthCheck.ResponseCode.String(),
		"enabled":        s.Enabled.Bool(),
	}
}
