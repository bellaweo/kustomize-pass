// Copyright 2020 me!
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestSecretsFromPassPlugin(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t)
	defer th.Reset()

	m := th.LoadAndRunGenerator(`
apiVersion: someteam.example.com/v1
kind: SecretsFromPass
passdir: katt_test
metadata:
  name: forbiddenValues
  namespace: production
keys:
- ROCKET
- VEGETABLE
`)
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
data:
  ROCKET: cG9vcA==
  VEGETABLE: YmxhaA==
kind: Secret
metadata:
  name: forbiddenValues
  namespace: production
type: Opaque
`)
}
