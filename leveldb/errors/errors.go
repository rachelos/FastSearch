// Copyright (c) 2014, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package errors provides common error types used throughout leveldb.
package errors

import (
	"errors"
	"fmt"

	"gitee.com/rachel_os/fastsearch/leveldb/storage"
	"gitee.com/rachel_os/fastsearch/leveldb/util"
)

// Common errors.
var (
	ErrNotFound    = New("leveldb: not found")
	ErrReleased    = util.ErrReleased
	ErrHasReleaser = util.ErrHasReleaser
)

// New returns an error that formats as the given text.
func New(text string) error {
	return errors.New(text)
}

// ErrCorrupted is the type that wraps errors that indicate corruption in
// the database.
type ErrCorrupted struct {
	Fd  storage.FileDesc
	Err error
}

func (e *ErrCorrupted) Error() string {
	if !e.Fd.Zero() {
		return fmt.Sprintf("%v [file=%v]", e.Err, e.Fd)
	}
	return e.Err.Error()
}

// NewErrCorrupted creates new ErrCorrupted error.
func NewErrCorrupted(fd storage.FileDesc, err error) error {
	return &ErrCorrupted{fd, err}
}

// IsCorrupted returns a boolean indicating whether the error is indicating
// a corruption.
func IsCorrupted(err error) bool {
	switch err.(type) {
	case *ErrCorrupted:
		return true
	case *storage.ErrCorrupted:
		return true
	}
	return false
}

// ErrMissingFiles is the type that indicating a corruption due to missing
// files. ErrMissingFiles always wrapped with ErrCorrupted.
type ErrMissingFiles struct {
	Fds []storage.FileDesc
}

func (e *ErrMissingFiles) Error() string { return "file missing" }

// SetFd sets 'file info' of the given error with the given file.
// Currently only ErrCorrupted is supported, otherwise will do nothing.
func SetFd(err error, fd storage.FileDesc) error {
	switch x := err.(type) {
	case *ErrCorrupted:
		x.Fd = fd
		return x
	}
	return err
}
