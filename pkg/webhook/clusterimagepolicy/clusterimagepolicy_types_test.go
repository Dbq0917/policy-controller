// Copyright 2026 The Sigstore Authors.
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

package clusterimagepolicy

import (
	"encoding/json"
	"testing"

	ociremote "github.com/sigstore/cosign/v3/pkg/oci/remote"
)

func TestAuthorityUnmarshalSource(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		wantRemoteOpts int
		wantRepository string
	}{{
		name:           "registry only",
		source:         "registry.example.com",
		wantRemoteOpts: 0,
	}, {
		name:           "registry only with port",
		source:         "localhost:5000",
		wantRemoteOpts: 0,
	}, {
		name:           "full repository",
		source:         "build.example.com/team/app",
		wantRemoteOpts: 1,
		wantRepository: "build.example.com/team/app",
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(map[string]any{
				"source": []map[string]string{{"oci": tc.source}},
			})
			if err != nil {
				t.Fatal(err)
			}

			var authority Authority
			if err := json.Unmarshal(data, &authority); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got := len(authority.RemoteOpts); got != tc.wantRemoteOpts {
				t.Fatalf("RemoteOpts length = %d, want %d", got, tc.wantRemoteOpts)
			}
			if tc.wantRepository != "" {
				repository := ociremote.TargetRepositoryFromOptions(authority.RemoteOpts...)
				if got := repository.Name(); got != tc.wantRepository {
					t.Errorf("target repository = %q, want %q", got, tc.wantRepository)
				}
			}
		})
	}
}
