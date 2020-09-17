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
	"errors"
	"fmt"
	"sync"
)

// DirectMemifMetrics is struct, containing information about directmemif connection
type DirectMemifMetrics struct {
	// RxBytes is a total number of bytes received from source
	RxBytes uint
	// TxBytes is a total number of bytes transmitted to target
	TxBytes uint
}

// Collector aggregates metrics
type Collector interface {
	Update(interface{}) error
	Metrics() map[string]string
}

// Collector aggregates DirectMemifMetrics
type directMemifCollector struct {
	lock    sync.Mutex
	metrics DirectMemifMetrics
}

// NewDirectMemifCollector creates directMemifCollector
func NewDirectMemifCollector() Collector {
	return &directMemifCollector{}
}

// Update updates directMemifCollector's metrics with incoming DirectMemifMetrics
func (d *directMemifCollector) Update(incomingMetrics interface{}) error {
	memifMetrics, ok := incomingMetrics.(DirectMemifMetrics)
	if !ok {
		return errors.New("wrong metrics type, should be DirectMemifMetrics")
	}

	d.lock.Lock()
	defer d.lock.Unlock()
	d.metrics.RxBytes += memifMetrics.RxBytes
	d.metrics.TxBytes += memifMetrics.TxBytes

	return nil
}

// Metrics returns DirectMemifMetrics' map representation
func (d *directMemifCollector) Metrics() map[string]string {
	d.lock.Lock()
	defer d.lock.Unlock()
	return map[string]string{
		"rx_bytes": fmt.Sprint(d.metrics.RxBytes),
		"tx_bytes": fmt.Sprint(d.metrics.TxBytes),
	}
}
