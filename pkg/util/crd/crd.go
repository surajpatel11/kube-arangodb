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
// Author Ewout Prangsma
//

package crd

import (
	"fmt"
	"time"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/arangodb/k8s-operator/pkg/apis/arangodb/v1alpha"
	"github.com/arangodb/k8s-operator/pkg/util/k8sutil"
	"github.com/arangodb/k8s-operator/pkg/util/retry"
)

// CreateCRD creates a custom resouce definition.
func CreateCRD(clientset apiextensionsclient.Interface, crdName, rkind, rplural string, shortName ...string) error {
	crd := &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: crdName,
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   api.SchemeGroupVersion.Group,
			Version: api.SchemeGroupVersion.Version,
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: rplural,
				Kind:   rkind,
			},
		},
	}
	if len(shortName) != 0 {
		crd.Spec.Names.ShortNames = shortName
	}
	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && !k8sutil.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// WaitCRDReady waits for a custom resource definition with given name to be ready.
func WaitCRDReady(clientset apiextensionsclient.Interface, crdName string) error {
	op := func() error {
		crd, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crdName, metav1.GetOptions{})
		if err != nil {
			return maskAny(err)
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1beta1.Established:
				if cond.Status == apiextensionsv1beta1.ConditionTrue {
					return nil
				}
			case apiextensionsv1beta1.NamesAccepted:
				if cond.Status == apiextensionsv1beta1.ConditionFalse {
					return maskAny(fmt.Errorf("Name conflict: %v", cond.Reason))
				}
			}
		}
		return maskAny(fmt.Errorf("Retry needed"))
	}
	if err := retry.Retry(op, time.Minute*5); err != nil {
		return maskAny(err)
	}
	return nil
}