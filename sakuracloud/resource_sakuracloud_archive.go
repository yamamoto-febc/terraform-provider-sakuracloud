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
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sacloud/ftps"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

var allowArchiveSizes = []string{"20", "40", "60", "80", "100", "250", "500", "750", "1024"}

func resourceSakuraCloudArchive() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudArchiveCreate,
		Read:   resourceSakuraCloudArchiveRead,
		Update: resourceSakuraCloudArchiveUpdate,
		Delete: resourceSakuraCloudArchiveDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      20,
				ValidateFunc: validateIntInWord(allowArchiveSizes),
			},
			"archive_file": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validation.StringInSlice([]string{"is1a", "is1b", "tk1a", "tk1v"}, false),
			},
		},
	}
}

func resourceSakuraCloudArchiveCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	opts := client.Archive.New()

	opts.Name = d.Get("name").(string)

	source := d.Get("archive_file").(string)
	path, err := homedir.Expand(source)
	if err != nil {
		return fmt.Errorf("Error expanding homedir in source (%s): %s", source, err)
	}
	// file exists?
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("Error opening archive_file (%s): %s", source, err)
	}

	opts.SizeMB = toSizeMB(d.Get("size").(int))
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandTags(client, rawTags)
	}

	archive, err := client.Archive.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Archive resource: %s", err)
	}

	// upload
	ftpServer, err := client.Archive.OpenFTP(archive.ID)
	if err != nil {
		return fmt.Errorf("Failed to Open FTPS Connection: %s", err)
	}

	ftpClient := ftps.NewClient(ftpServer.User, ftpServer.Password, ftpServer.HostName)
	if err := ftpClient.Upload(path); err != nil {
		return fmt.Errorf("Failed to upload SakuraCloud Archive resource: %s", err)
	}

	// close
	if _, err := client.Archive.CloseFTP(archive.ID); err != nil {
		return fmt.Errorf("Failed to Close FTPS Connection from Archive resource: %s", err)

	}

	d.SetId(archive.GetStrID())
	return resourceSakuraCloudArchiveRead(d, meta)
}

func resourceSakuraCloudArchiveRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	archive, err := client.Archive.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud Archive resource: %s", err)
	}

	return setArchiveResourceData(d, client, archive)
}

func resourceSakuraCloudArchiveUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	archive, err := client.Archive.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Archive resource: %s", err)
	}
	if d.HasChange("name") {
		archive.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			archive.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			archive.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			archive.Description = description.(string)
		} else {
			archive.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			archive.Tags = expandTags(client, rawTags)
		} else {
			archive.Tags = expandTags(client, []interface{}{})
		}
	}
	archive, err = client.Archive.Update(archive.ID, archive)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Archive resource: %s", err)
	}

	contentAttrs := []string{"iso_image_file", "hash"}
	isContentChanged := false
	for _, attr := range contentAttrs {
		if d.HasChange(attr) {
			isContentChanged = true
			break
		}
	}
	if isContentChanged {

		source := d.Get("archive_file").(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("Error expanding homedir in source (%s): %s", source, err)
		}
		// file exists?
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("Error opening archive_file (%s): %s", source, err)
		}

		// upload
		ftpServer, err := client.Archive.OpenFTP(archive.ID)
		if err != nil {
			return fmt.Errorf("Failed to Open FTPS Connection: %s", err)
		}

		ftpClient := ftps.NewClient(ftpServer.User, ftpServer.Password, ftpServer.HostName)
		if err := ftpClient.Upload(path); err != nil {
			return fmt.Errorf("Failed to upload SakuraCloud Archive resource: %s", err)
		}

		// close
		if _, err := client.Archive.CloseFTP(archive.ID); err != nil {
			return fmt.Errorf("Failed to Close FTPS Connection from Archive resource: %s", err)

		}

	}

	return resourceSakuraCloudArchiveRead(d, meta)
}

func resourceSakuraCloudArchiveDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	_, err := client.Archive.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Archive resource: %s", err)
	}

	_, err = client.Archive.Delete(toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Archive resource: %s", err)
	}

	return nil
}

func setArchiveResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Archive) error {

	d.Set("name", data.Name)
	d.Set("size", toSizeGB(data.SizeMB))
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	// NOTE 本来はAPIにてmd5ハッシュを取得できるのが望ましい
	if v, ok := d.GetOk("archive_file"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("Error expanding homedir in source (%s): %s", source, err)
		}
		// file exists?
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("Error opening archive_file (%s): %s", source, err)
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Error opening archive_file(%s): %s", source, err)
		}
		defer f.Close()

		b := base64.NewEncoder(base64.StdEncoding, f)
		defer b.Close()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, f); err != nil {
			return fmt.Errorf("Error encoding to base64 from archive_file (%s): %s", source, err)
		}

		h := md5.New()
		if _, err := io.Copy(h, &buf); err != nil {
			return fmt.Errorf("Error calculate md5 from archive_file (%s): %s", source, err)
		}

		d.Set("hash", h.Sum(nil))
	}

	d.Set("zone", client.Zone)
	return nil
}
