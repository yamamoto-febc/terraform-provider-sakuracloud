// Copyright 2016-2019 The Libsacloud Authors
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

// Code generated by "stringer -type=ArchiveOSType"; DO NOT EDIT.

package ostype

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CentOS-0]
	_ = x[CentOS8-1]
	_ = x[CentOS7-2]
	_ = x[CentOS6-3]
	_ = x[Ubuntu-4]
	_ = x[Ubuntu1804-5]
	_ = x[Ubuntu1604-6]
	_ = x[Debian-7]
	_ = x[Debian10-8]
	_ = x[Debian9-9]
	_ = x[CoreOS-10]
	_ = x[RancherOS-11]
	_ = x[K3OS-12]
	_ = x[Kusanagi-13]
	_ = x[SophosUTM-14]
	_ = x[FreeBSD-15]
	_ = x[Netwiser-16]
	_ = x[OPNsense-17]
	_ = x[Windows2016-18]
	_ = x[Windows2016RDS-19]
	_ = x[Windows2016RDSOffice-20]
	_ = x[Windows2016SQLServerWeb-21]
	_ = x[Windows2016SQLServerStandard-22]
	_ = x[Windows2016SQLServer2017Standard-23]
	_ = x[Windows2016SQLServerStandardAll-24]
	_ = x[Windows2016SQLServer2017StandardAll-25]
	_ = x[Windows2019-26]
	_ = x[Custom-27]
}

const _ArchiveOSType_name = "CentOSCentOS8CentOS7CentOS6UbuntuUbuntu1804Ubuntu1604DebianDebian10Debian9CoreOSRancherOSK3OSKusanagiSophosUTMFreeBSDNetwiserOPNsenseWindows2016Windows2016RDSWindows2016RDSOfficeWindows2016SQLServerWebWindows2016SQLServerStandardWindows2016SQLServer2017StandardWindows2016SQLServerStandardAllWindows2016SQLServer2017StandardAllWindows2019Custom"

var _ArchiveOSType_index = [...]uint16{0, 6, 13, 20, 27, 33, 43, 53, 59, 67, 74, 80, 89, 93, 101, 110, 117, 125, 133, 144, 158, 178, 201, 229, 261, 292, 327, 338, 344}

func (i ArchiveOSType) String() string {
	if i < 0 || i >= ArchiveOSType(len(_ArchiveOSType_index)-1) {
		return "ArchiveOSType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ArchiveOSType_name[_ArchiveOSType_index[i]:_ArchiveOSType_index[i+1]]
}