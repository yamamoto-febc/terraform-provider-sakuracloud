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
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSakuraCloudSubnetDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSubnetDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceSubnetBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSubnetDataSourceID("sakuracloud_subnet.foobar"),
					testAccCheckSakuraCloudSubnetDataSourceID("sakuracloud_subnet.foobar2"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSubnetDataSourceID("data.sakuracloud_subnet.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_subnet.foobar", "ipaddresses.#", "16"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSubnetConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSubnetDataSourceNotExists("data.sakuracloud_subnet.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudSubnetDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Subnet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Subnet data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudSubnetDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found Subnet data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudSubnetDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_subnet" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Subnet.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Subnet still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceSubnetBase = `
resource sakuracloud_internet "foobar" {
    name = "subnet_test"
}
resource "sakuracloud_subnet" "foobar" {
    internet_id = "${sakuracloud_internet.foobar.id}"
    next_hop = "${sakuracloud_internet.foobar.ipaddresses[0]}"
}
resource "sakuracloud_subnet" "foobar2" {
    internet_id = "${sakuracloud_internet.foobar.id}"
    next_hop = "${sakuracloud_internet.foobar.ipaddresses[1]}"
}
`

var testAccCheckSakuraCloudDataSourceSubnetConfig = `
resource sakuracloud_internet "foobar" {
    name = "subnet_test"
}
resource "sakuracloud_subnet" "foobar" {
    internet_id = "${sakuracloud_internet.foobar.id}"
    next_hop = "${sakuracloud_internet.foobar.ipaddresses[0]}"
}
resource "sakuracloud_subnet" "foobar2" {
    internet_id = "${sakuracloud_internet.foobar.id}"
    next_hop = "${sakuracloud_internet.foobar.ipaddresses[1]}"
}

data sakuracloud_subnet "foobar" {
    internet_id = "${sakuracloud_internet.foobar.id}"
    index = 1
}
`

var testAccCheckSakuraCloudDataSourceSubnetConfig_NotExists = `
resource sakuracloud_internet "foobar" {
    name = "subnet_test"
}
resource "sakuracloud_subnet" "foobar" {
    internet_id = "${sakuracloud_internet.foobar.id}"
    next_hop = "${sakuracloud_internet.foobar.ipaddresses[0]}"
}
resource "sakuracloud_subnet" "foobar2" {
    internet_id = "${sakuracloud_internet.foobar.id}"
    next_hop = "${sakuracloud_internet.foobar.ipaddresses[1]}"
}
data sakuracloud_subnet "foobar" {
    internet_id = "${sakuracloud_internet.foobar.id}"
    index = 2
}
`
