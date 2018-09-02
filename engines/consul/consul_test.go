/*
Copyright 2018 The SecretsEngine Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package consul_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/secretsengine/provocation/engines/consul"

	"github.com/hashicorp/consul/api"
	"github.com/ory/dockertest"
)

func inContainer(t *testing.T, fn func(resource *dockertest.Resource)) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not connect to docker: %v", err)
	}

	config := struct {
		Datacenter       string `json:"datacenter,omitempty"`
		ACLDatacenter    string `json:"acl_datacenter,omitempty"`
		ACLDefaultPolicy string `json:"acl_default_policy,omitempty"`
		ACLMasterToken   string `json:"acl_master_token,omitempty"`
	}{
		Datacenter:       "test",
		ACLDatacenter:    "test",
		ACLDefaultPolicy: "deny",
		ACLMasterToken:   "test",
	}

	encodedConfig, _ := json.Marshal(config)
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "consul",
		Cmd:        []string{"agent", "-dev", "-client", "0.0.0.0"},
		Env:        []string{fmt.Sprintf("CONSUL_LOCAL_CONFIG=%s", encodedConfig)},
	})

	defer func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("could not purge consul container: %v", err)
		}
	}()

	if err = pool.Retry(func() error {
		config := api.DefaultConfig()
		config.Address = fmt.Sprintf("localhost:%s", resource.GetPort("8500/tcp"))
		config.Token = "test"

		client, err := api.NewClient(config)
		if err != nil {
			return err
		}

		_, err = client.KV().Put(&api.KVPair{
			Key:   "ready",
			Value: []byte("ready"),
		}, nil)

		return err
	}); err != nil {
		t.Fatalf("could not connect to docker: %v", err)
	}

	fn(resource)
}

func TestBasic(t *testing.T) {
	inContainer(t, func(resource *dockertest.Resource) {
		engine := consul.Engine{
			Address:   fmt.Sprintf("localhost:%s", resource.GetPort("8500/tcp")),
			Token:     "test",
			TokenType: "management",
		}

		revocation, credentials, err := engine.Provision(context.TODO(), "foo", "bar")
		if err != nil {
			t.Fatalf("error provisioning consul credentials: %v", err)
		}

		if len(credentials["token"]) == 0 {
			t.Error("expected token to have length > 0")
		}

		err = engine.Revoke(context.TODO(), revocation)
		if err != nil {
			t.Fatalf("error revoking consul credentials: %v", err)
		}
	})
}
