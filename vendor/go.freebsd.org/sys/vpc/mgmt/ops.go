// Go interface to VPC Management instance operations.
//
// SPDX-License-Identifier: BSD-2-Clause-FreeBSD
//
// Copyright (C) 2018 Sean Chittenden <seanc@joyent.com>
// Copyright (c) 2018 Joyent, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

package mgmt

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"go.freebsd.org/sys/vpc"
)

// _MgmtCmd is the encoded type of operations that can be performed on a VPC
// Management handle.
type _MgmtCmd vpc.Cmd

// Ops that can be encoded into a vpc.Cmd
const (
	_OpInvalid         = vpc.Op(0)
	_OpCountType       = vpc.Op(1)
	_OpObjHeaderGetAll = vpc.Op(2)

	_CountTypeCmd       _MgmtCmd = _MgmtCmd(vpc.InBit|vpc.OutBit|(vpc.Cmd(vpc.ObjTypeMgmt)<<16)) | _MgmtCmd(_OpCountType)
	_ObjHeaderGetAllCmd _MgmtCmd = _MgmtCmd(vpc.InBit|vpc.OutBit|(vpc.Cmd(vpc.ObjTypeMgmt)<<16)) | _MgmtCmd(_OpObjHeaderGetAll)
)

// CountType obtains a count of VPC objects.
func (m *Mgmt) CountType(objType vpc.ObjType) (uint32, error) {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	// vpc_ctl(2): Input is a uint16 representing a type and the output is a
	// uint32.

	in := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(in, uint64(objType))
	if n < 2 {
		in = in[:n]
	} else {
		panic(fmt.Sprintf("invariant: ObjType size too big for kernel interface input (want/got: 2/%d", n))
	}

	out := make([]byte, binary.MaxVarintLen64)
	if err := vpc.Ctl(m.h, vpc.Cmd(_CountTypeCmd), in, &out); err != nil {
		return 0, errors.Wrapf(err, "unable to get count of VPC %s objects", objType)
	}

	count, n := binary.Uvarint(out)
	if n > 0 && n <= 4 {
		return uint32(count), nil
	}

	panic(fmt.Sprintf("invariant: ObjType size too big for kernel interface output (want/got: 4/%d", n))
}

// Close closes the VPC Handle descriptor.  Created VPC Switches will not be
// destroyed when the VPCSW is closed if the VPC Switch has been Committed.
func (m *Mgmt) Close() error {
	// Note for future reviewers: New/Free would have been more symmetric but I
	// used io.Closer's interface so that this handle could be managed in the same
	// way as any other io descriptor.

	if m.h.FD() <= 0 {
		return nil
	}

	if err := m.h.Close(); err != nil {
		return errors.Wrap(err, "unable to close VPC Management handle")
	}

	return nil
}