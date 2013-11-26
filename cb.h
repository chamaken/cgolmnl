#ifndef __CGOMNL_CB_H__
#define __CGOMNL_CB_H__
extern int attr_parse_wrapper(const struct nlmsghdr *nlh, size_t offset, void *data);
extern int attr_parse_nested_wrapper(const struct nlattr *attr, void *data);
extern int attr_parse_payload_wrapper(const void *payload, size_t payload_len, void *data);

typedef int (*mnl_ctl_cb_t)(const struct nlmsghdr *nlh, uint16_t type, void *data);
extern int _mnl_cb_run3(const void *buf, size_t numbytes, unsigned int seq,
		       unsigned int portid, mnl_cb_t cb_data, void *data,
		       mnl_ctl_cb_t cb_ctl, unsigned int cb_ctl_max);
extern int mnl_cb_run3(const void *buf, size_t numbytes, unsigned int seq,
		       unsigned int portid, mnl_cb_t cb_data, void *data,
		       mnl_ctl_cb_t cb_ctl);
extern int _cb_run3_wrapper(const void *buf, size_t numbytes, uint32_t seq,
			   uint32_t portid, void *data, unsigned int cb_ctl_max);
extern int cb_run3_wrapper(const void *buf, size_t numbytes, uint32_t seq,
			   uint32_t portid, void *data);
extern int cb_run_wrapper(const void *buf, size_t numbytes, uint32_t seq,
			  uint32_t portid, void *data);
#endif
