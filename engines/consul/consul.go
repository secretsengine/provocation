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

package consul

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
)

var (
	// ErrInvalidRevocation is returned if the revocation data is malformed
	ErrInvalidRevocation = errors.New("invalid revocation data")
)

// Engine is a Consul secrets engine that can provision and revoke
// credentials.
type Engine struct {
	Address string
	Scheme  string
	Token   string

	TokenType string
	Policy    string
}

func (e *Engine) client() (*api.Client, error) {
	conf := api.DefaultNonPooledConfig()
	conf.Address = e.Address
	conf.Scheme = e.Scheme
	conf.Token = e.Token

	return api.NewClient(conf)
}

// Provision provisions credentials using the engine configuration.
func (e *Engine) Provision(ctx context.Context, namespace, name string) ([]byte, map[string][]byte, error) {
	client, err := e.client()
	if err != nil {
		return nil, nil, err
	}

	writeOpts := &api.WriteOptions{}
	writeOpts = writeOpts.WithContext(ctx)

	username := fmt.Sprintf("%v-%v", namespace, name)
	token, _, err := client.ACL().Create(&api.ACLEntry{
		Name:  username,
		Type:  e.TokenType,
		Rules: e.Policy,
	}, (&api.WriteOptions{}).WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}

	revocation, _ := json.Marshal([]string{username, token})
	return revocation, map[string][]byte{
		"token": []byte(token),
	}, nil
}

// Revoke revokes credentials using the engine configuration.
func (e *Engine) Revoke(ctx context.Context, revocation []byte) error {
	client, err := e.client()
	if err != nil {
		return err
	}

	var tokenInfo []string
	if err := json.Unmarshal(revocation, &tokenInfo); err != nil {
		return ErrInvalidRevocation
	}
	if len(tokenInfo) != 2 {
		return ErrInvalidRevocation
	}

	if _, err = client.ACL().Destroy(tokenInfo[1], (&api.WriteOptions{}).WithContext(ctx)); err != nil {
		return fmt.Errorf("failed to revoke token %q: %v", string(tokenInfo[0]), err)
	}

	return nil
}
