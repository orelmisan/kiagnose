/*
 * This file is part of the kiagnose project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2022 Red Hat, Inc.
 *
 */

package vmi

import (
	"context"
	"fmt"
	"log"
	"time"

	k8scorev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/util/wait"

	kvcorev1 "kubevirt.io/api/core/v1"
)

type KubevirtVmisClient interface {
	GetVirtualMachineInstance(namespace, name string) (*kvcorev1.VirtualMachineInstance, error)
	CreateVirtualMachineInstance(namespace string, vmi *kvcorev1.VirtualMachineInstance) (*kvcorev1.VirtualMachineInstance, error)
}

func Start(c KubevirtVmisClient, namespace string, vmi *kvcorev1.VirtualMachineInstance) error {
	log.Printf("starting VMI %s/%s..", namespace, vmi.Name)
	if _, err := c.CreateVirtualMachineInstance(namespace, vmi); err != nil {
		return fmt.Errorf("failed to start VMI %s/%s: %v", vmi.Namespace, vmi.Name, err)
	}
	return nil
}

func WaitUntilReady(ctx context.Context, c KubevirtVmisClient, namespace, name string) error {
	log.Printf("waiting for VMI %s/%s to be ready..\n", namespace, name)

	if err := waitForVmiCondition(ctx, c, namespace, name, kvcorev1.VirtualMachineInstanceAgentConnected); err != nil {
		return fmt.Errorf("VMI %s/%s was not ready on time: %v", namespace, name, err)
	}

	return nil
}

func waitForVmiCondition(ctx context.Context, c KubevirtVmisClient, namespace, name string,
	conditionType kvcorev1.VirtualMachineInstanceConditionType) error {
	conditionFn := func(ctx context.Context) (bool, error) {
		updatedVmi, err := c.GetVirtualMachineInstance(namespace, name)
		if err != nil {
			return false, nil
		}
		for _, condition := range updatedVmi.Status.Conditions {
			if condition.Type == conditionType && condition.Status == k8scorev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	}
	const interval = time.Second * 5
	if err := wait.PollImmediateUntilWithContext(ctx, interval, conditionFn); err != nil {
		return fmt.Errorf("failed to wait for VMI %s/%s condition %s: %v", namespace, name, conditionType, err)
	}

	return nil
}
