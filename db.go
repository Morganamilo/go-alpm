// db.go - Functions for database handling.
//
// Copyright (c) 2013 The go-alpm Authors
//
// MIT Licensed. See LICENSE for details.

package alpm

/*
#include <alpm.h>
*/
import "C"

import (
	"unsafe"
)

// DB structure representing a alpm database.
type DB struct {
	ptr    *C.alpm_db_t
	handle Handle
}

// RegisterSyncDB Loads a sync database with given name and signature check level.
func (h *Handle) RegisterSyncDB(dbname string, siglevel SigLevel) (*DB, error) {
	cName := C.CString(dbname)
	defer C.free(unsafe.Pointer(cName))

	db := C.alpm_register_syncdb(h.ptr, cName, C.int(siglevel))
	if db == nil {
		return nil, h.LastError()
	}
	return &DB{db, *h}, nil
}

func (db *DB) Unregister() error {
	ok := C.alpm_db_unregister(db.ptr)
	if ok != 0 {
		return db.handle.LastError()
	}

	return nil
}

func (h *Handle) UnregisterAllSyncDBs() error {
	ok := C.alpm_unregister_all_syncdbs(h.ptr)
	if ok != 0 {
		return h.LastError()
	}

	return nil
}

// Name returns name of the db
func (db *DB) Name() string {
	return C.GoString(C.alpm_db_get_name(db.ptr))
}

// Servers returns host server URL.
func (db *DB) Servers() StringList {
	ptr := C.alpm_db_get_servers(db.ptr)
	return makeStringList(ptr)
}

// SetServers sets server list to use.
func (db *DB) SetServers(servers StringList) {
	C.alpm_db_set_servers(db.ptr, (*C.alpm_list_t)(unsafe.Pointer(servers.List)))
}

// Doesn'twork because this function strdups the server
// AddServers adds a string to the server list.
/*func (db *DB) AddServer(server string) {
	Csrv := C.CString(server)
	//defer C.free(unsafe.Pointer(Csrv))
	C.alpm_db_add_server(db.ptr, Csrv)
}*/

// SetUsage sets the Usage of the database
func (db *DB) SetUsage(usage Usage) {
	C.alpm_db_set_usage(db.ptr, C.int(usage))
}

// Name searches a package in db.
func (db *DB) Pkg(name string) (*Package, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ptr := C.alpm_db_get_pkg(db.ptr, cName)
	if ptr == nil {
		return nil, db.handle.LastError()
	}
	return &Package{ptr, db.handle}, nil
}

// PkgCache returns the list of packages of the database
func (db *DB) PkgCache() PackageList {
	pkgcache := C.alpm_db_get_pkgcache(db.ptr)
	return makePackageList(pkgcache, db.handle)
}

func (db *DB) Search(targets StringList) PackageList {
	pkglist := C.alpm_db_search(db.ptr, (*C.alpm_list_t)(unsafe.Pointer(targets.List)))
	return makePackageList(pkglist, db.handle)
}
