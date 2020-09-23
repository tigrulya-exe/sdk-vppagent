// Copyright (c) 2020 Cisco Systems, Inc.
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

// +build !windows

// Package kernelvethpair provides networkservice chain elements that support the kernel Mechanism using veth pairs
package kernelvethpair

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/kernel"

	"github.com/networkservicemesh/api/pkg/api/networkservice"

	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"

	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/vppagent"
	"github.com/networkservicemesh/sdk-vppagent/pkg/tools/kernelctx"
)

type kernelVethPairServer struct{}

// NewServer provides NetworkServiceServer chain elements that support the kernel Mechanism using veth pairs
func NewServer() networkservice.NetworkServiceServer {
	return &kernelVethPairServer{}
}

func (k *kernelVethPairServer) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*networkservice.Connection, error) {
	if mechanism := kernel.ToMechanism(request.GetConnection().GetMechanism()); mechanism != nil {
		err := appendInterfaceConfig(ctx, request.GetConnection(), "server")
		if err != nil {
			return nil, err
		}
		linuxIfaces := vppagent.Config(ctx).GetLinuxConfig().GetInterfaces()
		ctx = kernelctx.WithServerInterface(ctx, linuxIfaces[len(linuxIfaces)-1])
	}
	return next.Server(ctx).Request(ctx, request)
}

func (k *kernelVethPairServer) Close(ctx context.Context, conn *networkservice.Connection) (*empty.Empty, error) {
	if mechanism := kernel.ToMechanism(conn.GetMechanism()); mechanism != nil {
		err := appendInterfaceConfig(ctx, conn, "server")
		if err != nil {
			return nil, err
		}
		linuxIfaces := vppagent.Config(ctx).GetLinuxConfig().GetInterfaces()
		ctx = kernelctx.WithServerInterface(ctx, linuxIfaces[len(linuxIfaces)-1])
	}
	return next.Server(ctx).Close(ctx, conn)
}
