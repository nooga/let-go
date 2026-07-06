//go:build linux && !tinygo

/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import "github.com/nooga/let-go/pkg/vm"

var (
	waitResultMapping  *vm.StructMapping
	unameResultMapping *vm.StructMapping
	spawnResultMapping *vm.StructMapping
)

func initSyscallTypeMappings() {
	waitResultMapping = vm.RegisterStruct[WaitResult]("syscall/WaitResult")
	unameResultMapping = vm.RegisterStruct[UnameResult]("syscall/UnameResult")
	spawnResultMapping = vm.RegisterStruct[SpawnResult]("syscall/SpawnResult")
}
