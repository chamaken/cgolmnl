#ifndef __CGOLMNL_CB_H__
#define __CGOLMNL_CB_H__

extern int attr_parse_wrapper(const struct nlmsghdr *nlh, size_t offset, uintptr_t data);
extern int attr_parse_nested_wrapper(const struct nlattr *attr, uintptr_t data);
extern int attr_parse_payload_wrapper(const void *payload, size_t payload_len, uintptr_t data);

typedef int (*mnl_ctl_cb_t)(const struct nlmsghdr *nlh, uint16_t type, void *data);
extern int cb_run2_wrapper(const void *buf, size_t numbytes, uint32_t seq,
			   uint32_t portid, uintptr_t data, uint16_t *ctl_types, size_t ctl_types_len);
extern int cb_run_wrapper(const void *buf, size_t numbytes, uint32_t seq,
			  uint32_t portid, uintptr_t data);
#endif
