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

package rabbitmq

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/michaelklishin/rabbit-hole"
)

// Engine is a RabbitMQ secrets engine that can provision and revoke
// credentials.
type Engine struct {
	URI      string
	Username string
	Password string

	Tags   []string
	VHosts map[string]VHost
}

// VHost defines permissions to be set on a vhost.
type VHost struct {
	Configure string
	Write     string
	Read      string
}

// Provision provisions credentials using the engine configuration.
func (e *Engine) Provision(ctx context.Context, namespace, name string) ([]byte, map[string][]byte, error) {
	client, err := rabbithole.NewClient(e.URI, e.Username, e.Password)
	if err != nil {
		return nil, nil, err
	}

	username := fmt.Sprintf("%s-%s", namespace, name)
	password, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, err
	}

	_, err = client.PutUser(username, rabbithole.UserSettings{
		Password: password,
		Tags:     strings.Join(e.Tags, ","),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user %q: %v", username, err)
	}

	for vhost, permissions := range e.VHosts {
		_, err := client.UpdatePermissionsIn(vhost, username, rabbithole.Permissions(permissions))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to set permissions on user %q: %v", username, err)
		}
	}

	return []byte(username), map[string][]byte{
		"username": []byte(username),
		"password": []byte(password),
	}, nil
}

// Revoke revokes credentials using the engine configuration.
func (e *Engine) Revoke(ctx context.Context, revocation []byte) error {
	client, err := rabbithole.NewClient(e.URI, e.Username, e.Password)
	if err != nil {
		return err
	}

	username := string(revocation)
	if _, err = client.DeleteUser(username); err != nil {
		return fmt.Errorf("failed to revoke user %q: %v", username, err)
	}

	return nil
}
