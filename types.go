// types.go - libalpm types.
//
// Copyright (c) 2013 The go-alpm Authors
//
// MIT Licensed. See LICENSE for details.

package alpm

// #cgo CFLAGS: -D_FILE_OFFSET_BITS=64
// #include <alpm.h>
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Depend C.alpm_depend_t

func (dep *Depend) Name() string {
	return C.GoString(dep.name)
}

func (dep *Depend) Version() string {
	return C.GoString(dep.version)
}

func (dep *Depend) Description() string {
	return C.GoString(dep.desc)
}

func (dep *Depend) NameHash() uint {
	return uint(dep.name_hash)
}

func (dep *Depend) Mod() DepMod {
	return DepMod(dep.mod)
}

func (dep *Depend) String() string {
	str := C.alpm_dep_compute_string((*C.alpm_depend_t)(dep))
	return C.GoString(str)
}

func (dep *Depend) Free() {
	C.alpm_dep_free((*C.alpm_depend_t)(dep))
}

// File provides a description of package files.
type File struct {
	Name string
	Size int64
	Mode uint32
}

func convertFilelist(files *C.alpm_filelist_t) []File {
	size := int(files.count)
	items := make([]File, size)

	rawItems := reflect.SliceHeader{
		Len:  size,
		Cap:  size,
		Data: uintptr(unsafe.Pointer(files.files))}

	cFiles := *(*[]C.alpm_file_t)(unsafe.Pointer(&rawItems))

	for i := 0; i < size; i++ {
		items[i] = File{
			Name: C.GoString(cFiles[i].name),
			Size: int64(cFiles[i].size),
			Mode: uint32(cFiles[i].mode)}
	}
	return items
}

// Internal alpm list structure.
type list struct {
	Data unsafe.Pointer
	Prev *list
	Next *list
}

// Iterates a function on a list and stop on error.
func (l *list) forEach(f func(unsafe.Pointer) error) error {
	for ; l != nil; l = l.Next {
		err := f(l.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

type StringList struct {
	*list
}

func (l StringList) ForEach(f func(string) error) error {
	return l.forEach(func(p unsafe.Pointer) error {
		return f(C.GoString((*C.char)(p)))
	})
}

func (l StringList) Slice() []string {
	slice := []string{}
	l.ForEach(func(s string) error {
		slice = append(slice, s)
		return nil
	})
	return slice
}

type BackupFile struct {
	Name string
	Hash string
}

type BackupList struct {
	*list
}

func (l BackupList) ForEach(f func(BackupFile) error) error {
	return l.forEach(func(p unsafe.Pointer) error {
		bf := (*C.alpm_backup_t)(p)
		return f(BackupFile{
			Name: C.GoString(bf.name),
			Hash: C.GoString(bf.hash),
		})
	})
}

func (l BackupList) Slice() (slice []BackupFile) {
	l.ForEach(func(f BackupFile) error {
		slice = append(slice, f)
		return nil
	})
	return
}

type QuestionAny struct {
	ptr *C.alpm_question_any_t
}

func (question QuestionAny) SetAnswer(answer bool) {
	if answer {
		question.ptr.answer = 1
	} else {
		question.ptr.answer = 0
	}
}

type QuestionInstallIgnorepkg struct {
	ptr *C.alpm_question_install_ignorepkg_t
}

func (question QuestionAny) Type() QuestionType {
	return QuestionType(question.ptr._type)
}

func (question QuestionAny) Answer() bool {
	return question.ptr.answer == 1
}

func (question QuestionAny) QuestionInstallIgnorepkg() (QuestionInstallIgnorepkg, error) {
	if question.Type() == QuestionTypeInstallIgnorepkg {
		return *(*QuestionInstallIgnorepkg)(unsafe.Pointer(&question)), nil
	}

	return QuestionInstallIgnorepkg{}, fmt.Errorf("Can not convert to QuestionInstallIgnorepkg")
}

func (question QuestionAny) QuestionSelectProvider() (QuestionSelectProvider, error) {
	if question.Type() == QuestionTypeSelectProvider {
		return *(*QuestionSelectProvider)(unsafe.Pointer(&question)), nil
	}

	return QuestionSelectProvider{}, fmt.Errorf("Can not convert to QuestionInstallIgnorepkg")
}

func (question QuestionAny) QuestionReplace() (QuestionReplace, error) {
	if question.Type() == QuestionTypeReplacePkg {
		return *(*QuestionReplace)(unsafe.Pointer(&question)), nil
	}

	return QuestionReplace{}, fmt.Errorf("Can not convert to QuestionReplace")
}

func (question QuestionInstallIgnorepkg) SetInstall(install bool) {
	if install {
		question.ptr.install = 1
	} else {
		question.ptr.install = 0
	}
}

func (question QuestionInstallIgnorepkg) Type() QuestionType {
	return QuestionType(question.ptr._type)
}

func (question QuestionInstallIgnorepkg) Install() bool {
	return question.ptr.install == 1
}

func (question QuestionInstallIgnorepkg) Pkg(h *Handle) Package {
	return Package{
		question.ptr.pkg,
		*h,
	}
}

type QuestionReplace struct {
	ptr *C.alpm_question_replace_t
}

func (question QuestionReplace) Type() QuestionType {
	return QuestionType(question.ptr._type)
}

func (question QuestionReplace) SetReplace(replace bool) {
	if replace {
		question.ptr.replace = 1
	} else {
		question.ptr.replace = 0
	}
}

func (question QuestionReplace) Replace() bool {
	return question.ptr.replace == 1
}

func (question QuestionReplace) NewPkg(h *Handle) *Package {
	return &Package{
		question.ptr.newpkg,
		*h,
	}
}

func (question QuestionReplace) OldPkg(h *Handle) *Package {
	return &Package{
		question.ptr.oldpkg,
		*h,
	}
}

func (question QuestionReplace) newDB(h *Handle) *DB {
	return &DB{
		question.ptr.newdb,
		*h,
	}
}

type QuestionSelectProvider struct {
	ptr *C.alpm_question_select_provider_t
}

func (question QuestionSelectProvider) Type() QuestionType {
	return QuestionType(question.ptr._type)
}

func (question QuestionSelectProvider) SetUseIndex(index int) {
	question.ptr.use_index = C.int(index)
}

func (question QuestionSelectProvider) UseIndex() int {
	return int(question.ptr.use_index)
}

func (question QuestionSelectProvider) Providers(h *Handle) PackageList {
	return PackageList{
		(*list)(unsafe.Pointer(question.ptr.providers)),
		*h,
	}
}

func (question QuestionSelectProvider) Dep() *Depend {
	return (*Depend)(question.ptr.depend)
}
