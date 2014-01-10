#include <stdio.h>
#include <errno.h>
#include <libmnl/libmnl.h>
#include "_cgo_export.h"

int attr_parse_wrapper(const struct nlmsghdr *nlh, size_t offset, void *data)
{
	return mnl_attr_parse(nlh, (unsigned int)offset, (mnl_attr_cb_t)GoAttrCb, data);
}

int attr_parse_nested_wrapper(const struct nlattr *nested, void *data)
{
	return mnl_attr_parse_nested(nested, (mnl_attr_cb_t)GoAttrCb, data);
}

int attr_parse_payload_wrapper(const void *payload, size_t payload_len, void *data)
{
	return mnl_attr_parse_payload(payload, (unsigned int)payload_len, (mnl_attr_cb_t)GoAttrCb, data);
}

int
cb_run_wrapper(const void *buf, size_t numbytes, uint32_t seq,
	       uint32_t portid, void *data)
{
	return mnl_cb_run(buf, numbytes, (unsigned int)seq, (unsigned int)portid, (mnl_cb_t)GoCb, data);
}


/*
 * http://stackoverflow.com/questions/1023261/is-there-a-way-to-do-currying-in-c
 * http://gcc.gnu.org/onlinedocs/gcc/Nested-Functions.html#Nested-Functions
 */
// static const mnl_cb_t go_ctlcb_array[NLMSG_MIN_TYPE] = { (mnl_cb_t)GoCtlCb2 };
static mnl_cb_t go_ctlcb_array[NLMSG_MIN_TYPE] = { [0 ... NLMSG_MIN_TYPE - 1] = (mnl_cb_t)GoCtlCb };

int
cb_run2_wrapper(const void *buf, size_t numbytes, uint32_t seq,
		uint32_t portid, void *data, uint16_t *ctl_types, size_t ctl_types_len)
{
	int i;
	mnl_cb_t cb_ctl_array[NLMSG_MIN_TYPE] = { NULL };

	if (ctl_types_len >= NLMSG_MIN_TYPE) {
		errno = EINVAL;
		return MNL_CB_ERROR;
	}

	for (i = 0; i < ctl_types_len; i++) {
		if (ctl_types[i] >= NLMSG_MIN_TYPE) {
			errno = EINVAL;
			return MNL_CB_ERROR;
		}
		cb_ctl_array[ctl_types[i]] = (mnl_cb_t)GoCtlCb;
	}

	return mnl_cb_run2(buf, numbytes, (unsigned int)seq, (unsigned int)portid,
			   (mnl_cb_t)GoCb, data, cb_ctl_array, NLMSG_MIN_TYPE - 1);
}
