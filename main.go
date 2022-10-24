package main

// #cgo CFLAGS: -I${SRCDIR}/target
// #cgo LDFLAGS: -L${SRCDIR}/target -lexample
//
// #include <stdlib.h>
// #include <libexample.h>
// #include <stdio.h>
// static void* pass(graal_isolatethread_t* thread, void* in, void* inSize, void* outSize) {
//   return pass_bytes(thread, malloc, in, (int*)inSize, (int*)outSize);
// }
// typedef void* void_pointer;
import "C"
import (
	"fmt"
	"unsafe"
)

type javaCgo struct {
	isolate *C.graal_isolate_t
}

type JavaCgo interface {
	TestPassBytes()
}

func New() (JavaCgo, error) {
	var isolate *C.graal_isolate_t
	var thread *C.graal_isolatethread_t

	param := &C.graal_create_isolate_params_t{
		reserved_address_space_size: 1024 * 1024 * 500,
	}

	if C.graal_create_isolate(param, &isolate, &thread) != 0 {
		return nil, fmt.Errorf("failed to initialize")
	}

	return &javaCgo{
		isolate: isolate,
	}, nil
}

func (j *javaCgo) attachThread() (*C.graal_isolatethread_t, error) {
	thread := C.graal_get_current_thread(j.isolate)
	if thread != nil {
		return thread, nil
	}

	var newThread *C.graal_isolatethread_t
	if C.graal_attach_thread(j.isolate, &newThread) != 0 {
		return nil, fmt.Errorf("failed to attach thread")
	}

	return newThread, nil
}

func (j *javaCgo) TestPassBytes() {
	thread, err := j.attachThread()
	if err != nil {
		panic(err.Error())
	}

	var inPointer unsafe.Pointer
	inPointer = C.CBytes([]byte{5, 0, 7})
	var outSize C.int
	inSize := 3
	res := C.pass(thread, inPointer, unsafe.Pointer(C.void_pointer(&inSize)), unsafe.Pointer(C.void_pointer(&outSize)))
	defer C.free(unsafe.Pointer(res))
	goBytes := C.GoBytes(res, outSize)
	fmt.Printf("bytes %v %d\n", goBytes[:int(outSize)], int(outSize))
}

func main() {
	javaCgo, err := New()
	if err != nil {
		println(err)
		return
	}

	javaCgo.TestPassBytes()
}
