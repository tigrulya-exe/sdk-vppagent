// Copyright (c) 2020 Cisco and/or its affiliates.
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

package ipaddress_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/cls"
	memif_mechanisms "github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/memif"
	"github.com/networkservicemesh/sdk/pkg/networkservice/core/chain"
	"github.com/networkservicemesh/sdk/pkg/networkservice/utils/checks/checkcontext"
	"github.com/networkservicemesh/sdk/pkg/networkservice/utils/checks/checkcontextonreturn"
	"github.com/networkservicemesh/sdk/pkg/networkservice/utils/checks/checkopts"
	"github.com/networkservicemesh/sdk/pkg/networkservice/utils/inject/injecterror"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/connectioncontext/ipcontext/ipaddress"
	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/mechanisms/memif"
	"github.com/networkservicemesh/sdk-vppagent/pkg/networkservice/vppagent"
)

const (
	IPAddress      = "1.1.1.1"
	BaseDir        = "BaseDir"
	SocketFilename = "socketfilename"
)

func clientRequest() *networkservice.NetworkServiceRequest {
	return &networkservice.NetworkServiceRequest{
		Connection: &networkservice.Connection{
			Mechanism: &networkservice.Mechanism{
				Cls:  cls.LOCAL,
				Type: memif_mechanisms.MECHANISM,
				Parameters: map[string]string{
					memif_mechanisms.SocketFilename: SocketFilename,
				},
			},
			Context: &networkservice.ConnectionContext{
				IpContext: &networkservice.IPContext{
					SrcIpAddr: IPAddress,
				},
			},
		},
	}
}

func TestSetIPVppClient(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	client := chain.NewNetworkServiceClient(
		checkcontextonreturn.NewClient(t, func(t *testing.T, ctx context.Context) {
			conf := vppagent.Config(ctx)
			numInterfaces := len(conf.GetVppConfig().GetInterfaces())
			require.Greater(t, numInterfaces, 0)
			assert.Equal(t, IPAddress, conf.GetVppConfig().GetInterfaces()[numInterfaces-1].GetIpAddresses()[0])
		}),
		ipaddress.NewClient(),
		checkcontext.NewClient(t, func(t *testing.T, ctx context.Context) {
			conf := vppagent.Config(ctx)
			expectedConf := vppagent.Config(vppagent.WithConfig(context.Background()))
			assert.Equal(t, expectedConf, conf)
		}),
		memif.NewClient(BaseDir),
	)
	conn, err := client.Request(vppagent.WithConfig(context.Background()), clientRequest())
	assert.NotNil(t, conn)
	assert.Nil(t, err)
	_, err = client.Close(vppagent.WithConfig(context.Background()), conn)
	assert.Nil(t, err)
}

func TestSetVppIPClientPropagatesError(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	client := chain.NewNetworkServiceClient(
		ipaddress.NewClient(),
		injecterror.NewClient(),
	)
	_, err := client.Request(vppagent.WithConfig(context.Background()), clientRequest())
	assert.NotNil(t, err)
	_, err = client.Close(vppagent.WithConfig(context.Background()), clientRequest().GetConnection())
	assert.NotNil(t, err)
}

func TestSetIPVppClientPropagatesOpts(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	client := checkopts.CheckPropogateOptsClient(t, ipaddress.NewClient())
	_, err := client.Request(vppagent.WithConfig(context.Background()), clientRequest())
	assert.Nil(t, err)
	_, err = client.Close(vppagent.WithConfig(context.Background()), clientRequest().GetConnection())
	assert.Nil(t, err)
}
