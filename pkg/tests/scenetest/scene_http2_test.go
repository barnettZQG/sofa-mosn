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
package tests

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/alipay/sofa-mosn/pkg/mosn"
	"github.com/alipay/sofa-mosn/pkg/protocol"
	"golang.org/x/net/http2"
)

func TestHttp2(t *testing.T) {
	meshAddr := "127.0.0.1:2045"
	http2Addr := "127.0.0.1:8080"
	server := NewUpstreamHttp2(t, http2Addr)
	server.GoServe()
	defer server.Close()
	mesh_config := CreateSimpleMeshConfig(meshAddr, []string{http2Addr}, protocol.Http2, protocol.Http2)
	mesh := mosn.NewMosn(mesh_config)
	go mesh.Start()
	defer mesh.Close()
	time.Sleep(5 * time.Second) //wait mesh and server start
	//Client Run
	tr := &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(netw, addr)
		},
	}

	httpClient := http.Client{Transport: tr}
	for i := 0; i < 20; i++ {
		requestId := fmt.Sprintf("%d", i)
		request, err := http.NewRequest("GET", fmt.Sprintf("http://%s", meshAddr), nil)
		if err != nil {
			t.Fatalf("create request error:%v\n", err)
		}
		request.Header.Add("service", "testhttp2")
		request.Header.Add("Requestid", requestId)
		resp, err := httpClient.Do(request)
		if err != nil {
			t.Errorf("request %s response error: %v\n", requestId, err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("request %s read body error: %v\n", requestId, err)
			continue
		}
		t.Logf("request %s get data: %s\n", requestId, body)
	}

}
