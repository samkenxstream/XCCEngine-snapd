// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2020 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package sysconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/snapcore/snapd/dirs"
	"github.com/snapcore/snapd/osutil"
)

func DisableCloudInit() error {
	ubuntuDataCloud := filepath.Join(dirs.RunMnt, "ubuntu-data/system-data/etc/cloud/")
	if err := os.MkdirAll(ubuntuDataCloud, 0755); err != nil {
		return fmt.Errorf("cannot make cloud config dir: %v", err)
	}
	if err := ioutil.WriteFile(filepath.Join(ubuntuDataCloud, "cloud-init.disabled"), nil, 0644); err != nil {
		return fmt.Errorf("cannot disable cloud-init: %v", err)
	}

	return nil
}

func installCloudInitCfg(src string) error {
	ccl, err := filepath.Glob(filepath.Join(src, "*.cfg"))
	if err != nil {
		return err
	}
	if len(ccl) == 0 {
		return nil
	}

	ubuntuDataCloudCfgDir := filepath.Join(dirs.RunMnt, "ubuntu-data/system-data/etc/cloud/cloud.cfg.d/")
	if err := os.MkdirAll(ubuntuDataCloudCfgDir, 0755); err != nil {
		return fmt.Errorf("cannot make cloud config dir: %v", err)
	}

	for _, cc := range ccl {
		if err := osutil.CopyFile(cc, filepath.Join(ubuntuDataCloudCfgDir, filepath.Base(cc)), 0); err != nil {
			return err
		}
	}
	return nil
}

// disable cloud-init by default (as it's not confined)
// TODO:UC20: - allow gadget cloud.cfg.d/* (with whitelisted keys?)
//            - allow cloud.cfg.d (with whitelisted keys) for non
//               grade dangerous systems
func configureCloudInit(opts *Options) (err error) {
	switch opts.CloudInitSrcDir {
	case "":
		err = DisableCloudInit()
	default:
		err = installCloudInitCfg(opts.CloudInitSrcDir)
	}
	return err
}
