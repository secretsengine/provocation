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

package postgresql

import (
	"context"
	"errors"
	"net/url"

	"github.com/hashicorp/vault/plugins/database/postgresql"
	"github.com/secretsengine/provocation/engines/internal/database"
)

// Engine is a PostgreSQL secrets engine that can provision and revoke
// credentials.
type Engine database.Engine

var (
	// ErrInvalidURI is returned when the PostgreSQL URI is malformed.
	ErrInvalidURI = errors.New("invalid uri")
)

func (e *Engine) client() (*database.Engine, interface{}, error) {
	engine := database.Engine(*e)
	db, err := postgresql.New()

	// Update URI to include username/password
	uri, err := url.Parse(engine.URI)
	if err != nil {
		return nil, nil, ErrInvalidURI
	}
	uri.User = url.UserPassword(e.Username, e.Password)
	engine.URI = uri.String()

	return &engine, db, err
}

// Provision provisions credentials using the engine configuration.
func (e *Engine) Provision(ctx context.Context, namespace, name string) ([]byte, map[string][]byte, error) {
	engine, db, err := e.client()
	if err != nil {
		return nil, nil, err
	}

	return engine.Provision(ctx, db, namespace, name)
}

// Revoke revokes credentials using the engine configuration.
func (e *Engine) Revoke(ctx context.Context, revocation []byte) error {
	engine, db, err := e.client()
	if err != nil {
		return err
	}

	return engine.Revoke(ctx, db, revocation)
}
