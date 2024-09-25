// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crio

import (
	"k8s.io/klog/v2"

	"github.com/swarnimarun/cadvisor/container"
	"github.com/swarnimarun/cadvisor/fs"
	info "github.com/swarnimarun/cadvisor/info/v1"
	"github.com/swarnimarun/cadvisor/watcher"
)

// NewPlugin returns an implementation of container.Plugin suitable for passing to container.RegisterPlugin()
func NewPlugin(success *bool) container.Plugin {
	return &plugin{success: success}
}

type plugin struct{ success *bool }

func (p *plugin) InitializeFSContext(context *fs.Context) error {
	crioClient, err := Client()
	if err != nil {
		return err
	}

	crioInfo, err := crioClient.Info()
	if err != nil {
		klog.V(5).Infof("CRI-O not connected: %v", err)
	} else {
		context.Crio = fs.CrioContext{Root: crioInfo.StorageRoot, ImageStore: crioInfo.StorageImage, Driver: crioInfo.StorageDriver}
	}
	return nil
}

func (p *plugin) Register(factory info.MachineInfoFactory, fsInfo fs.FsInfo, includedMetrics container.MetricSet) (watcher.ContainerWatcher, error) {
	err := Register(factory, fsInfo, includedMetrics)
	*p.success = err == nil
	return nil, err
}
