/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package v2

import (
	//"time"
	//"google.golang.org/grpc"
	envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	ads "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/alipay/sofa-mosn/pkg/log"
	//"golang.org/x/net/context"
	//google_rpc "github.com/gogo/googleapis/google/rpc"
	"errors"
	envoy_api_v2_core1 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
)

func (c *V2Client) GetClusters(streamClient ads.AggregatedDiscoveryService_StreamAggregatedResourcesClient) []*envoy_api_v2.Cluster {
	err := c.ReqClusters(streamClient)
	if err != nil {
		log.DefaultLogger.Fatalf("get clusters fail: %v", err)
		return nil
	}
	r, err := streamClient.Recv()
	if err != nil {
		log.DefaultLogger.Fatalf("get clusters fail: %v", err)
		return nil
	}
	return c.HandleClustersResp(r)
}

func (c *V2Client) ReqClusters(streamClient ads.AggregatedDiscoveryService_StreamAggregatedResourcesClient) error {
	if streamClient == nil {
		return errors.New("stream client is nil")
	}
	err := streamClient.Send(&envoy_api_v2.DiscoveryRequest{
		VersionInfo:   "",
		ResourceNames: []string{},
		TypeUrl:       "type.googleapis.com/envoy.api.v2.Cluster",
		ResponseNonce: "",
		ErrorDetail:   nil,
		Node: &envoy_api_v2_core1.Node{
			Id: c.ServiceNode,
		},
	})
	if err != nil {
		log.DefaultLogger.Fatalf("get clusters fail: %v", err)
		return err
	}
	return nil
}

func (c *V2Client) HandleClustersResp(resp *envoy_api_v2.DiscoveryResponse) []*envoy_api_v2.Cluster {
	clusters := make([]*envoy_api_v2.Cluster, 0)
	for _, res := range resp.Resources {
		cluster := envoy_api_v2.Cluster{}
		cluster.Unmarshal(res.GetValue())
		clusters = append(clusters, &cluster)
	}
	return clusters
}
