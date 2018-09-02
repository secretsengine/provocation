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

package password_test

import (
	"context"
	"testing"

	"github.com/secretsengine/provocation/engines/password"
)

func TestBasic(t *testing.T) {
	engine := password.Engine{
		Length: 16,
	}

	revocation, credentials, err := engine.Provision(context.TODO(), "foo", "bar")
	if err != nil {
		t.Fatalf("error provisioning password credentials: %v", err)
	}

	if len(credentials["password"]) != 16 {
		t.Error("expected password to have length == 16")
	}

	err = engine.Revoke(context.TODO(), revocation)
	if err != nil {
		t.Fatalf("error revoking password credentials: %v", err)
	}
}
