// Copyright (c) 2020 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"errors"
	"fmt"
	"sync"
)

type Collector interface {
	Update(interface{}) error
	Metrics() map[string]string
}

type DirectMemifMetrics map[string]uint

type directMemifCollector struct {
	lock    sync.Mutex
	metrics DirectMemifMetrics
}

func NewDirectMemifCollector() Collector {
	return &directMemifCollector{
		metrics: make(map[string]uint),
	}
}

func (d *directMemifCollector) Update(toMerge interface{}) error {
	memifMetrics, ok := toMerge.(DirectMemifMetrics)
	if !ok {
		return errors.New("Wrong metrics type, should be DirectMemifMetrics")
	}

	d.lock.Lock()
	defer d.lock.Unlock()
	for k, v := range memifMetrics {
		d.metrics[k] += v
	}

	return nil
}

func (d *directMemifCollector) Metrics() map[string]string {
	d.lock.Lock()
	defer d.lock.Unlock()
	result := make(map[string]string)
	for k, v := range d.metrics {
		result[k] = fmt.Sprint(v)
	}

	return result
}
