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

package postgresql_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/secretsengine/provocation/engines/postgresql"

	"github.com/ory/dockertest"
)

func inContainer(t *testing.T, fn func(resource *dockertest.Resource)) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not connect to docker: %v", err)
	}

	resource, err := pool.Run("postgres", "latest", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=database"})
	if err != nil {
		t.Fatalf("could not start postgres container: %v", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("could not purge postgres container: %v", err)
		}
	}()

	if err = pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/database?sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatalf("could not connect to docker: %v", err)
	}

	fn(resource)
}
func TestBasic(t *testing.T) {
	inContainer(t, func(resource *dockertest.Resource) {
		engine := postgresql.Engine{
			URI:      fmt.Sprintf("postgres://localhost:%s/database?sslmode=disable", resource.GetPort("5432/tcp")),
			Username: "postgres",
			Password: "secret",
			Creation: []string{
				`CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}'`,
				`GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}"`,
			},
		}

		revocation, credentials, err := engine.Provision(context.TODO(), "foo", "bar")
		if err != nil {
			t.Fatalf("error provisioning postgresql credentials: %v", err)
		}

		if len(credentials["username"]) == 0 {
			t.Error("expected username to have length > 0")
		}
		if len(credentials["password"]) == 0 {
			t.Error("expected password to have length > 0")
		}

		err = engine.Revoke(context.TODO(), revocation)
		if err != nil {
			t.Fatalf("error revoking postgresql credentials: %v", err)
		}
	})
}
