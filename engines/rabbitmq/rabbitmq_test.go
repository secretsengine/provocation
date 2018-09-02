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

package rabbitmq_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/secretsengine/provocation/engines/rabbitmq"

	"github.com/michaelklishin/rabbit-hole"
	"github.com/ory/dockertest"
)

func inContainer(t *testing.T, fn func(resource *dockertest.Resource)) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not connect to docker: %v", err)
	}

	resource, err := pool.Run("rabbitmq", "3-management", []string{})
	if err != nil {
		t.Fatalf("could not start rabbitmq container: %v", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("could not purge rabbitmq container: %v", err)
		}
	}()

	if err = pool.Retry(func() error {
		client, err := rabbithole.NewClient(fmt.Sprintf("http://localhost:%v", resource.GetPort("15672/tcp")), "guest", "guest")
		if err != nil {
			return err
		}

		_, err = client.Overview()

		return err
	}); err != nil {
		t.Fatalf("could not connect to docker: %v", err)
	}

	fn(resource)
}

func TestBasic(t *testing.T) {
	inContainer(t, func(resource *dockertest.Resource) {
		engine := rabbitmq.Engine{
			URI:      fmt.Sprintf("http://localhost:%v", resource.GetPort("15672/tcp")),
			Username: "guest",
			Password: "guest",
			Tags:     []string{"foo", "bar"},
			VHosts: map[string]rabbitmq.VHost{
				"/": {Configure: ".*", Write: ".*", Read: ".*"},
			},
		}

		revocation, credentials, err := engine.Provision(context.TODO(), "foo", "bar")
		if err != nil {
			t.Fatalf("error provisioning consul credentials: %v", err)
		}

		if len(credentials["username"]) == 0 {
			t.Error("expected username to have length > 0")
		}
		if len(credentials["password"]) == 0 {
			t.Error("expected password to have length > 0")
		}

		err = engine.Revoke(context.TODO(), revocation)
		if err != nil {
			t.Fatalf("error revoking rabbitmq credentials: %v", err)
		}
	})
}
