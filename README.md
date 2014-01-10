cgolmnl
========

Go wrapper of libmnl using cgo, under heavy development

sample
------

see examples


installation
------------

To enable mmap, specify ``nlmmap'' tag. go build -tags nlmmap


requires
--------

libmnl

  * test reqs (optional): **ginkgo (http://onsi.github.io/ginkgo/)


links
-----

* libmnl: http://netfilter.org/projects/libmnl/


struct
------

nlmsghdr and nlattr has real - Nlmsghdr and Nlattr. mnl_nlmsg_batch and
mnl_socket are opaque, [0]byte as cgo said. there are receivers, see
go_receiver.go



issues
------

### callback ###

in my horrible understanding, Go function which called from C has to be
exported. To follow this, functions which uses callback was implemented in a
hacky way. callback is classified into two major - for nlmsghdr, cb_run and
nlattr. both of C functions are wrapped in cb.c

1. call Go function wrapping C function in cb.c. it creates pseudo data param from
   (non exported, real) Go callback function and real data param.
2. wrapping C function calls real libmnl function
3. real libmnl function call Go callback, GoAttrCb (attr.go) in case of nlattr,
   GoCb (callback.go) nlmsghdr.
4. Go callback above demultiplex data param into Go function and real data param
5. call real Go callback

To wrap cb_run2() I added new one -

    typedef int (*mnl_ctl_cb_t)(const struct nlmsghdr *nlh, uint16_t type, void *data);

As you know though, this is called in case of control type message (NLMSG_...)
and function signature has changed.

    func CbRun2(buf []byte, seq, portid uint32, cb_data MnlCb, data interface{},
                cb_ctl MnlCtlCb, ctltypes []uint16) (int, error) {

The last parameters is list of controll message type, e.g NLMSG_NOOP, NLMSG_ERROR...
cb_ctl will be called in case of being specified in the list. See callbacl_test.go
CbRun2 context for example.
    

### errno ###

I currently use an incredibly childish C function in set_errno.c

    void SetErrno(int n) { errno = n; }

I can not find the way of tossing up Go callback error, in other words set C's
errno from Go. I am not good at English let me show why I need to do in the code
snippets below

* C library header (lib.h)

    typedef int cbf_t(void *data);
    extern int c_func(cbf_t cbfunc, void *data);

* C wrapper header

    #include "lib.h"
    #include "_cgo_export.h"
    extern int wrap(void *data);

* C wrapper source

    #include "cblib.h"
    int wrap(void *data)
    {
        return c_func((cbf_t)CallFromC, data);
    }

* Go
    /*
    #include "cbwrap.h"
    */
    import "C"
    import "unsafe"

    type Cb_t func(interface {}) (int, error)

    func Doit(cbfunc Cb_t, data interface{}) (int, error) {
        // multiplexing
        pseudo_data := [2]unsafe.Pointer{unsafe.Pointer(&cbfunc), unsafe.Pointer(&data)}
        return C.wrap(unsafe.Pointer(&pseudo_data))
    }

    //export CallFromC
    func CallFromC(pseudo_data interface) C.int {
        // demultiplexing
        args := *(*[2]unsafe.Porinter)(pseudo_data)
        cbfunc := *(*Cb_t)(args[0])
        real_data := *(*interface{})(args[1])
        ret, err := cbfunc(real_data)
	// set C errno here
    }
     
    func cb(data interface{}) (int, error) {
        i := data.(int)
        if i < 0:
	    return -1, syscall.Errno(-i)
        else:
	    return i, syscall.Errno(0)
    }

    func main() {
        Doit(cb, 7)
    }

call chain will be:
Go main() -> Go Doit() -> C wrap() -> C c_func() -> Go CallFromC() -> Go cb()

I need to know the way of tossing last Go cb() error up to Go Doit() or
C c_func().



comparison
----------

| original				| cgolmnl			| remarks			|
| ------------------------------------- | ----------------------------- | ----------------------------- |
| mnl_attr_get_type			| AttrGetType			|				|
| mnl_attr_get_len			| AttrGetLen			|				|
| mnl_attr_get_payload_len		| AttrGetPayloadLen		|				|
| mnl_attr_get_payload			| AttrGetPayload		|				|
| (add)					| AttrGetPayloadBytes		| returns []byte		|
| mnl_attr_ok				| AttrOk			|				|
| mnl_attr_next				| AttrNext			| 				|
| mnl_attr_type_valid			| AttrTypeValid			| returns ret, error		|
| mnl_attr_validate			| AttrValidate			| returns ret, errno		|
| mnl_attr_validate2			| AttrValidate2			| returns ret, errno		|
| mnl_attr_parse			| AttrParse			| returns ret, errno		|
| mnl_attr_parse_nested			| AttrParseNested		| returns ret, errno		|
| mnl_attr_parse_payload		| AttrParsePayload		| returns ret, errno		|
| mnl_attr_get_u8			| AttrGetU8			|				|
| mnl_attr_get_u16			| AttrGetU16			|				|
| mnl_attr_get_u32			| AttrGetU32			|				|
| mnl_attr_get_u64			| AttrGetU64			|				|
| mnl_attr_get_str			| AttrGetStr			|				|
| mnl_attr_put				| AttrPut			|				|
| (add)					| AttrPutData			|				|
| (add)					| AttrPutBytes			|				|
| mnl_attr_put_u8			| AttrPutU8			|				|
| mnl_attr_put_u16			| AttrPutU16			|				|
| mnl_attr_put_u32			| AttrPutU32			|				|
| mnl_attr_put_u64			| AttrPutU64			|				|
| mnl_attr_put_str			| AttrPutstr			|				|
| mnl_attr_put_strz			| AttrPutstrz			|				|
| mnl_attr_nest_start			| AttrNestStart			|				|
| mnl_attr_put_check			| AttrPutCheck			|				|
| mnl_attr_put_u8_check			| AttrPutU8Check		|				|
| mnl_attr_put_u16_check		| AttrPutU16Check		|				|
| mnl_attr_put_u32_check		| AttrPutU32Check		|				|
| mnl_attr_put_u64_check		| AttrPutU64Check		|				|
| mnl_attr_put_str_check		| AttrPutStrCheck		|				|
| mnl_attr_put_strz_check		| AttrPutStrzCheck		|				|
| mnl_attr_nest_start_check		| AttrnestStartCheck		|				|
| mnl_attr_nest_end			| AttrnestEnd			|				|
| mnl_attr_nest_cancel			| AttrnestCancel		|				|
| ------------------------------------- | ----------------------------- | ----------------------------- |
| mnl_nlmsg_size			| NlmsgSize			|				|
| mnl_nlmsg_get_payload_len		| NlmsgGetPayloadLen		|				|
| mnl_nlmsg_put_header			| NlmsgPutHeader		| require unsafe.Pointer	|
| (add)					| NlmsgPutHeaderBytes		| wrap above, require []byte	|
| mnl_nlmsg_put_extra_header		| NlmsgPutExtraHeader		|  				|
| mnl_nlmsg_get_paylod			| NlmsgGetPayload		| 				|
| (add)					| NlmsgGetPayloadBytes		| returns []byte		|
| mnl_nlmsg_get_payload_offset		| NlmsgGetPayloadOffset		| 				|
| (add)					| NlmsgGetPayloadOffsetBytes	| returns []byte		|
| mnl_nlmsg_ok				| NlmsgOk			| 				|
| mnl_nlmsg_next			| NlmsgNext			|				|
| mnl_nlmsg_get_payload_tail		| NlmsgGetPayloadTail		| 				|
| mnl_nlmsg_seq_ok			| NlmsgSeqOk			|				|
| mnl_nlmsg_portid_ok			| NlmsgPortidOk			| 				|
| mnl_nlmsg_fprintf			| NlmsgFprint			| *os.File, not descriptor	|
| mnl_nlmsg_batch_start			| NlmsgBatchStart		|				|
| mnl_nlmsg_batch_stop			| NlmsgBatchStop		| 				|
| mnl_nlmsg_batch_next			| NlmsgBatchNext		|	 			|
| mnl_nlmsg_batch_reset			| NlmsgBatchReset		|	 			|
| mnl_nlmsg_batch_size			| NlmsgBatchSize		|	 			|
| mnl_nlmsg_batch_head			| NlmsgBatchHead		|	 			|
| (add)					| NlmsgBatchHeadBytes		| returns []byte		|
| mnl_nlmsg_batch_current		| NlmsgBatchCurrent		|				|
| mnl_nlmsg_batch_is_empty		| NlmsgBatchIsEmpty		|				|
| ------------------------------------- | ----------------------------- | ----------------------------- |
| mnl_cb_run				| CbRun				| 				|
| mnl_cb_run2				| CbRun2			| changed signature		|
| ------------------------------------- | ----------------------------- | ----------------------------- |
| mnl_socket_get_fd			| SocketGgetFd			|				|
| mnl_socket_get_portid			| SocketGetPortid		|				|
| mnl_socket_open			| SocketOpen			| 				|
| mnl_socket_bind			| SocketBind			|				|
| mnl_socket_sendto			| SocketSendto			|				|
| (add)					| SocketSendNlmsg		|				|
| mnl_socket_recvfrom			| SocketRecvfrom		|				|
| mnl_socket_close			| SocketClose			|				|
| mnl_socket_setsockopt			| SocketSetsockopt		|				|
| mnl_socket_getsockopt			| SocketGetsockopt		|				|
