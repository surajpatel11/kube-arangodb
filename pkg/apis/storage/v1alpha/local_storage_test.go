//
// DISCLAIMER
//
// Copyright 2018 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Jan Christoph Uhde
//

package v1alpha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test creation of arango local storage
func TestArangoLocalStorageCreation(t *testing.T) {
	// REVIEW - is there something more meaningful to test
	storage := ArangoLocalStorage{}
	list := ArangoLocalStorageList{}
	list.Items = append(list.Items, storage)
	assert.Equal(t, 1, len(list.Items))
}
