cgolmnl
========

Go wrapper of libmnl using cgo, under heavy development

sample
------

see examples


requires
--------

libmnl


links
-----

* libmnl: http://netfilter.org/projects/libmnl/


struct
------

nlmsghdr and nlattr has real - Nlmsghdr and Nlattr. mnl_nlmsg_batch and
mnl_socket are opaque, [0]byte as cgo said. these are receiver, see
go_receiver.go



issues
------

### callback ###

in my horrible understanding, Go function which called from C has to be
exported. To follow this, functions which uses callback was implemented in hacky
way. callback is classified into two major - for nlmsghdr, cb_run and
nlattr. both of C functions are wrapped in cb.c

1. call Go function wrapping C function in cb.c. it creates psuedo data param from
   (non exported) requested Go callback function and real data param.
2. wrapped function calls real libmnl function
3. real libmnl function call Go callback, GoAttrCb (attr.go) in case of nlattr,
   GoCb (callback.go) nlmsghdr.
4. Go callback above demultiplex data param into Go function and real data param
5. requested Go callback will be called

cb_run2() is too complicated for me to implement. then I added new one -
cb_run3() which introduce new callback prototype.

    typedef int (*mnl_ctl_cb_t)(const struct nlmsghdr *nlh, uint16_t type, void *data);

As you know though, this is called in case of control type message (NLMSG_...)


### errno ###

I can not find the way of set C's errno. I think it's important for callback and
made callback function type in Go can return err, but ignore, discard in
exported Go callback functions - GoCb and GoCtlCb in callback.go, GoAttrCb in
attr.go. Please tell me how to set C's errno from Go.


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
| mnl_nlmsg_put_header			| NlmsgPutHeader		|				|
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
| mnl_cb_run2				| (could not)			|				|
| (add)					| CbRun3			| ctl dispatcher is mnl_ctl_cb_t|
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
