cgolmnl
=======

Go wrapper of libmnl using cgo, under heavy development


sample
------

see examples


installation
------------

Need running `mktypes.sh` before build.


requires
--------

  * libmnl

  * test reqs (optional): **ginkgo (http://onsi.github.io/ginkgo/)


links
-----

* libmnl: http://netfilter.org/projects/libmnl/


struct
------

nlmsghdr and nlattr has real - Nlmsghdr and Nlattr. mnl_nlmsg_batch and
mnl_socket are opaque as NlmsgBatch and Socket (cgo say [0]byte). there are
receivers, see go_receiver.go



errno
-----

I currently use an incredibly childish C function in set_errno.c

    void SetErrno(int n) { errno = n; }

I can not find the way of tossing up Go callback error, in other words set C's
errno from Go. I am not good at English let me show why I need to do in the code
snippets below

* C library header (lib.h)

```
    typedef int cbf_t(void *data);
    extern int c_func(cbf_t cbfunc, void *data);
```

* wrapper header

```
    #include "lib.h"
    #include "_cgo_export.h"
    extern int wrap(void *data);
```

* C wrapper source

```
    #include "cblib.h"
    int wrap(void *data)
    {
        return c_func((cbf_t)CallFromC, data);
    }
```

* Go

```
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
```

call chain will be:
Go main() -> Go Doit() -> C wrap() -> C c_func() -> Go CallFromC() -> Go cb()

I need to know the way of tossing last Go cb() error up to Go Doit() or
C c_func().



comparison
----------

| original				| cgolmnl			| remarks			|
| ------------------------------------- | ----------------------------- | ----------------------------- |
| mnl_attr_get_type			| Attr.GetType			|				|
| mnl_attr_get_len			| Attr.GetLen			|				|
| mnl_attr_get_payload_len		| Attr.GetPayloadLen		|				|
| mnl_attr_get_payload			| Attr.GetPayload		|				|
| (add)					| Attr.GetPayloadBytes		| returns []byte		|
| mnl_attr_ok				| Attr.Ok			|				|
| mnl_attr_next				| Attr.Next			| 				|
| mnl_attr_type_valid			| Attr.TypeValid		| returns ret, error		|
| mnl_attr_validate			| Attr.Validate			| returns ret, errno		|
| mnl_attr_validate2			| Attr.Validate2		| returns ret, errno		|
| mnl_attr_parse			| Attr.Parse			| returns ret, errno		|
| mnl_attr_parse_nested			| Attr.ParseNested		| returns ret, errno		|
| mnl_attr_parse_payload		| Attr.ParsePayload		| returns ret, errno		|
| mnl_attr_get_u8			| Attr.GetU8			|				|
| mnl_attr_get_u16			| Attr.GetU16			|				|
| mnl_attr_get_u32			| Attr.GetU32			|				|
| mnl_attr_get_u64			| Attr.GetU64			|				|
| mnl_attr_get_str			| Attr.GetStr			|				|
| mnl_attr_put				| Attr.Put			|				|
| (add)					| Attr.PutPtr			|				|
| (add)					| Attr.PutBytes			|				|
| mnl_attr_put_u8			| Attr.PutU8			|				|
| mnl_attr_put_u16			| Attr.PutU16			|				|
| mnl_attr_put_u32			| Attr.PutU32			|				|
| mnl_attr_put_u64			| Attr.PutU64			|				|
| mnl_attr_put_str			| Attr.Putstr			|				|
| mnl_attr_put_strz			| Attr.Putstrz			|				|
| mnl_attr_nest_start			| Attr.NestStart		|				|
| mnl_attr_put_check			| Attr.PutCheck			|				|
| mnl_attr_put_u8_check			| Attr.PutU8Check		|				|
| mnl_attr_put_u16_check		| Attr.PutU16Check		|				|
| mnl_attr_put_u32_check		| Attr.PutU32Check		|				|
| mnl_attr_put_u64_check		| Attr.PutU64Check		|				|
| mnl_attr_put_str_check		| Attr.PutStrCheck		|				|
| mnl_attr_put_strz_check		| Attr.PutStrzCheck		|				|
| mnl_attr_nest_start_check		| Attr.nestStartCheck		|				|
| mnl_attr_nest_end			| Attr.nestEnd			|				|
| mnl_attr_nest_cancel			| Attr.nestCancel		|				|
| ------------------------------------- | ----------------------------- | ----------------------------- |
| mnl_nlmsg_size			| NlmsgSize			|				|
| mnl_nlmsg_get_payload_len		| Nlmsg.GetPayloadLen		|				|
| mnl_nlmsg_put_header			| Nlmsg.PutHeader		|				|
| (add)					| NlmsgBytes			| build new nlmsg from []byte	|
| (add)					| NewNlmsg			| build new nlmsg		|
| mnl_nlmsg_put_extra_header		| Nlmsg.PutExtraHeader		|  				|
| mnl_nlmsg_get_paylod			| Nlmsg.GetPayload		| 				|
| (add)					| Nlmsg.GetPayloadBytes		| returns []byte		|
| mnl_nlmsg_get_payload_offset		| Nlmsg.GetPayloadOffset	| 				|
| (add)					| Nlmsg.GetPayloadOffsetBytes	| returns []byte		|
| mnl_nlmsg_ok				| Nlmsg.Ok			| 				|
| mnl_nlmsg_next			| Nlmsg.Next			|				|
| mnl_nlmsg_get_payload_tail		| Nlmsg.GetPayloadTail		| 				|
| mnl_nlmsg_seq_ok			| Nlmsg.SeqOk			|				|
| mnl_nlmsg_portid_ok			| Nlmsg.PortidOk		| 				|
| mnl_nlmsg_fprintf			| Nlmsg.Fprint			| *os.File, not descriptor	|
| mnl_nlmsg_batch_start			| NewNlmsgBatch			|				|
| mnl_nlmsg_batch_stop			| NlmsgBatch.Stop		| 				|
| mnl_nlmsg_batch_next			| NlmsgBatch.Next		|	 			|
| mnl_nlmsg_batch_reset			| NlmsgBatch.Reset		|	 			|
| mnl_nlmsg_batch_size			| NlmsgBatch.Size		|	 			|
| mnl_nlmsg_batch_head			| NlmsgBatch.Head		|	 			|
| (add)					| NlmsgBatch.HeadBytes		| returns []byte		|
| mnl_nlmsg_batch_current		| NlmsgBatch.Current		|				|
| mnl_nlmsg_batch_is_empty		| NlmsgBatch.IsEmpty		|				|
| ------------------------------------- | ----------------------------- | ----------------------------- |
| mnl_cb_run				| CbRun				| 				|
| mnl_cb_run2				| CbRun2			| changed signature		|
| ------------------------------------- | ----------------------------- | ----------------------------- |
| mnl_socket_get_fd			| Socket.GgetFd			|				|
| mnl_socket_get_portid			| Socket.GetPortid		|				|
| mnl_socket_open			| NewSocket			| 				|
| mnl_socket_open2			| NewSocket2			| 				|
| mnl_socket_fdopen			| NewSocketFd			| 				|
| mnl_socket_bind			| Socket.Bind			|				|
| mnl_socket_sendto			| Socket.Sendto			|				|
| (add)					| Socket.SendNlmsg		|				|
| mnl_socket_recvfrom			| Socket.Recvfrom		|				|
| mnl_socket_close			| Socket.Close			|				|
| mnl_socket_setsockopt			| Socket.Setsockopt		|				|
| mnl_socket_getsockopt			| Socket.Getsockopt		|				|
