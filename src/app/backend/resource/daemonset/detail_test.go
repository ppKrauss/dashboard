// Copyright 2017 The Kubernetes Dashboard Authors.
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

package daemonset

import (
	"testing"

	"k8s.io/client-go/kubernetes/fake"
	api "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func TestDeleteDaemonSetServices(t *testing.T) {
	cases := []struct {
		namespace, name string
		DaemonSetList   *extensions.DaemonSetList
		serviceList     *api.ServiceList
		expectedActions []string
	}{
		{
			TestNamespace, "ds-1",
			&extensions.DaemonSetList{
				Items: []extensions.DaemonSet{
					CreateDaemonSet("ds-1", TestNamespace, TestLabel),
				},
			},
			&api.ServiceList{
				Items: []api.Service{
					CreateService("svc-1", TestNamespace, TestLabel),
				},
			},
			[]string{"get", "list", "list", "delete"},
		},
		{
			TestNamespace, "ds-1",
			&extensions.DaemonSetList{
				Items: []extensions.DaemonSet{
					CreateDaemonSet("ds-1", TestNamespace, TestLabel),
					CreateDaemonSet("ds-2", TestNamespace, TestLabel),
				},
			},
			&api.ServiceList{
				Items: []api.Service{
					CreateService("svc-1", TestNamespace, TestLabel),
				},
			},
			[]string{"get", "list"},
		},
		{
			TestNamespace, "ds-1",
			&extensions.DaemonSetList{
				Items: []extensions.DaemonSet{
					CreateDaemonSet("ds-1", TestNamespace, TestLabel),
				},
			},
			&api.ServiceList{},
			[]string{"get", "list", "list"},
		},
	}

	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.DaemonSetList, c.serviceList)

		DeleteDaemonSetServices(fakeClient, c.namespace, c.name)

		actions := fakeClient.Actions()
		if len(actions) != len(c.expectedActions) {
			t.Errorf("Unexpected actions: %v, expected %d actions got %d", actions,
				len(c.expectedActions), len(actions))
			continue
		}

		for i, verb := range c.expectedActions {
			if actions[i].GetVerb() != verb {
				t.Errorf("Unexpected action: %+v, expected %s",
					actions[i], verb)
			}
		}
	}
}
