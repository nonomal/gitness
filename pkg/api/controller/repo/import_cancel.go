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

package repo

import (
	"context"
	"fmt"

	apiauth "github.com/harness/gitness/pkg/api/auth"
	"github.com/harness/gitness/pkg/api/usererror"
	"github.com/harness/gitness/pkg/auth"
	"github.com/harness/gitness/types/enum"
)

// ImportCancel cancels a repository import.
func (c *Controller) ImportCancel(ctx context.Context,
	session *auth.Session,
	repoRef string,
) error {
	// note: can't use c.getRepoCheckAccess because this needs to fetch a repo being imported.
	repo, err := c.repoStore.FindByRef(ctx, repoRef)
	if err != nil {
		return err
	}

	if err = apiauth.CheckRepo(ctx, c.authorizer, session, repo, enum.PermissionRepoDelete, false); err != nil {
		return err
	}

	if !repo.Importing {
		return usererror.BadRequest("repository is not being imported")
	}

	if err = c.importer.Cancel(ctx, repo); err != nil {
		return fmt.Errorf("failed to cancel repository import")
	}

	return c.DeleteNoAuth(ctx, session, repo)
}