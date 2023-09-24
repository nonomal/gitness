// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secret

import (
	"github.com/harness/gitness/encrypt"
	"github.com/harness/gitness/internal/auth/authz"
	"github.com/harness/gitness/internal/store"
	"github.com/harness/gitness/types/check"

	"github.com/jmoiron/sqlx"
)

type Controller struct {
	db          *sqlx.DB
	uidCheck    check.PathUID
	encrypter   encrypt.Encrypter
	secretStore store.SecretStore
	authorizer  authz.Authorizer
	spaceStore  store.SpaceStore
}

func NewController(
	db *sqlx.DB,
	uidCheck check.PathUID,
	authorizer authz.Authorizer,
	encrypter encrypt.Encrypter,
	secretStore store.SecretStore,
	spaceStore store.SpaceStore,
) *Controller {
	return &Controller{
		db:          db,
		uidCheck:    uidCheck,
		encrypter:   encrypter,
		secretStore: secretStore,
		authorizer:  authorizer,
		spaceStore:  spaceStore,
	}
}
