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

type DirectMemifMetrics struct {
	Rx_bytes uint
	Tx_bytes uint
}

type MetricsCollector interface {
	Update(interface{}) error
	Metrics() map[string]string
}

type directMemifCollector struct {
	lock    sync.Mutex
	metrics DirectMemifMetrics
}

func NewDirectMemifCollector() MetricsCollector {
	return &directMemifCollector{}
}

func (d *directMemifCollector) Update(toMerge interface{}) error {
	memifMetrics, ok := toMerge.(DirectMemifMetrics)
	if !ok {
		return errors.New("Wrong metrics type, should be DirectMemifMetrics")
	}

	d.lock.Lock()
	defer d.lock.Unlock()
	d.metrics.Rx_bytes += memifMetrics.Rx_bytes
	d.metrics.Tx_bytes += memifMetrics.Tx_bytes

	return nil
}

func (d *directMemifCollector) Metrics() map[string]string {
	d.lock.Lock()
	defer d.lock.Unlock()
	return map[string]string{
		"rx_bytes": fmt.Sprint(d.metrics.Rx_bytes),
		"tx_bytes": fmt.Sprint(d.metrics.Tx_bytes),
	}
}
