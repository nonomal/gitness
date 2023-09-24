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

package gitrpc

import "regexp"

var matchCommitSHA = regexp.MustCompile("^[0-9a-f]+$")

func ValidateCommitSHA(commitSHA string) bool {
	if len(commitSHA) != 40 && len(commitSHA) != 64 {
		return false
	}

	return matchCommitSHA.MatchString(commitSHA)
}
