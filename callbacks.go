// callbacks.go - Handles libalpm callbacks.
//
// Copyright (c) 2013 The go-alpm Authors
//
// MIT Licensed. See LICENSE for details.

package alpm

/*
#include <stdint.h>
#include <alpm.h>
void logCallback(alpm_loglevel_t level, char *cstring);
void go_alpm_log_cb(alpm_loglevel_t level, const char *fmt, va_list arg);
void go_alpm_set_logging(alpm_handle_t *handle);
void go_alpm_set_question(alpm_handle_t *handle);
*/
import "C"

import (
	"unsafe"
)

type logCallbackSig func(LogLevel, string)
type questionCallbackSig func(*QuestionAny)

var DefaultLogLevel = LogWarning

func DefaultLogCallback(lvl LogLevel, s string) {
	if lvl <= DefaultLogLevel {
		print("go-alpm: ", s)
	}
}

var globalLogCallback logCallbackSig
var globalQuestionCallback questionCallbackSig

//export logCallback
func logCallback(level C.alpm_loglevel_t, cstring *C.char) {
	globalLogCallback(LogLevel(level), C.GoString(cstring))
}

//export questionCallback
func questionCallback(question *C.alpm_question_t) {
	q := (*QuestionAny)(unsafe.Pointer(question))
	globalQuestionCallback(q)
}

func (h *Handle) SetLogCallback(cb logCallbackSig) {
	globalLogCallback = cb
	C.go_alpm_set_logging(h.ptr)
}

func (h *Handle) SetQuestionCallback(cb questionCallbackSig) {
	globalQuestionCallback = cb
	C.go_alpm_set_question(h.ptr)
}
