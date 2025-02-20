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

package main

import (
	"log"
	"os"

	"github.com/kiagnose/kiagnose/kiagnose/environment"

	"github.com/kiagnose/kiagnose/checkups/kubevirt-vm-latency/vmlatency"
)

func main() {
	const errMessagePrefix = "Kubevirt VM latency checkup failed"
	env := environment.EnvToMap(os.Environ())

	workingNamespace, err := environment.ReadNamespaceFile()
	if err != nil {
		log.Fatalf("%s: %v\n", errMessagePrefix, err)
	}

	if err := vmlatency.Run(env, workingNamespace); err != nil {
		log.Fatalf("%s: %v\n", errMessagePrefix, err)
	}
}
