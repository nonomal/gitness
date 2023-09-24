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

package pullreq

import (
	"github.com/harness/gitness/gitrpc"
	"github.com/harness/gitness/internal/auth/authz"
	pullreqevents "github.com/harness/gitness/internal/events/pullreq"
	"github.com/harness/gitness/internal/services/codecomments"
	"github.com/harness/gitness/internal/services/pullreq"
	"github.com/harness/gitness/internal/sse"
	"github.com/harness/gitness/internal/store"
	"github.com/harness/gitness/internal/url"
	"github.com/harness/gitness/lock"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(db *sqlx.DB, urlProvider *url.Provider, authorizer authz.Authorizer,
	pullReqStore store.PullReqStore, pullReqActivityStore store.PullReqActivityStore,
	codeCommentsView store.CodeCommentView,
	pullReqReviewStore store.PullReqReviewStore, pullReqReviewerStore store.PullReqReviewerStore,
	repoStore store.RepoStore, principalStore store.PrincipalStore, fileViewStore store.PullReqFileViewStore,
	rpcClient gitrpc.Interface, eventReporter *pullreqevents.Reporter,
	mtxManager lock.MutexManager, codeCommentMigrator *codecomments.Migrator,
	pullreqService *pullreq.Service, sseStreamer sse.Streamer,
) *Controller {
	return NewController(db, urlProvider, authorizer,
		pullReqStore, pullReqActivityStore,
		codeCommentsView,
		pullReqReviewStore, pullReqReviewerStore,
		repoStore, principalStore, fileViewStore,
		rpcClient, eventReporter,
		mtxManager, codeCommentMigrator, pullreqService, sseStreamer)
}
