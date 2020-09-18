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

// Package metrics - contains directmemif metrics and metrics collector
package metrics

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

const (
	// RxBytes is a total number of bytes received from target
	RxBytes = "rx_bytes"
	// TxBytes is a total number of bytes transmitted to target
	TxBytes = "tx_bytes"
)

// Collector aggregates metrics
type Collector interface {
	Update(map[string]string) error
	Metrics() map[string]string
}

// Collector aggregates direct memif metrics
type directMemifCollector struct {
	lock    sync.Mutex
	metrics map[string]uint64
}

// NewDirectMemifCollector creates directMemifCollector
func NewDirectMemifCollector() Collector {
	return &directMemifCollector{
		metrics: map[string]uint64{
			RxBytes: 0,
			TxBytes: 0,
		},
	}
}

// Update updates directMemifCollector's metrics with incoming metrics
func (d *directMemifCollector) Update(incomingMetrics map[string]string) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	for k, v := range incomingMetrics {
		value, err := strconv.ParseUint(v, 10, 0)
		if err != nil {
			return err
		}

		if _, ok := d.metrics[k]; !ok {
			return errors.Errorf("Unknown key: %v", k)
		}
		d.metrics[k] += value
	}

	return nil
}

// Metrics returns direct memif metrics' map representation
func (d *directMemifCollector) Metrics() map[string]string {
	d.lock.Lock()
	defer d.lock.Unlock()
	result := make(map[string]string)
	for k, v := range d.metrics {
		result[k] = fmt.Sprint(v)
	}

	return result
}
