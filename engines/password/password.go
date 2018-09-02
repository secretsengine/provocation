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

package password

import (
	"context"

	"github.com/sethvargo/go-password/password"
)

// Engine is a password secrets engine that can provision a randomly generated
// password.
type Engine struct {
	Length                int
	NumDigits             int
	NumSymbols            int
	DisableUppercase      bool
	AllowRepeatCharacters bool
}

// Provision provisions credentials using the engine configuration.
func (e *Engine) Provision(ctx context.Context, namespace, name string) ([]byte, map[string][]byte, error) {
	password, err := password.Generate(e.Length, e.NumDigits, e.NumSymbols, e.DisableUppercase, e.AllowRepeatCharacters)

	return nil, map[string][]byte{"password": []byte(password)}, err
}

// Revoke does nothing for this secrets engine.
func (e *Engine) Revoke(ctx context.Context, revocation []byte) error {
	return nil
}
