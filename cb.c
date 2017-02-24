#include <stdio.h>
#include <errno.h>
#include <libmnl/libmnl.h>
#include "_cgo_export.h"

/*
 * These are C function wrapper to call from Go. see
 *   https://groups.google.com/forum/#!topic/golang-nuts/PRcvOJqItow
 *
 * Unfortunately you can't pass a Go func to C code and have the C code
 * call it.  The best you can do is pass a Go func to C code and have the C
 * code turn around and pass the Go func back to a Go function that then
 * calls the func.
 */

// int attr_parse_wrapper(const struct nlmsghdr *nlh, size_t offset, void *data)
int attr_parse_wrapper(const struct nlmsghdr *nlh, size_t offset, uintptr_t data)
{
	return mnl_attr_parse(nlh, (unsigned int)offset, (mnl_attr_cb_t)GoAttrCb, (void *)data);
}

int attr_parse_nested_wrapper(const struct nlattr *nested, uintptr_t data)
{
	return mnl_attr_parse_nested(nested, (mnl_attr_cb_t)GoAttrCb, (void *)data);
}

int attr_parse_payload_wrapper(const void *payload, size_t payload_len, uintptr_t data)
{
	return mnl_attr_parse_payload(payload, (unsigned int)payload_len, (mnl_attr_cb_t)GoAttrCb, (void *)data);
}

int
cb_run_wrapper(const void *buf, size_t numbytes, uint32_t seq,
	       uint32_t portid, uintptr_t data)
{
	return mnl_cb_run(buf, numbytes, (unsigned int)seq, (unsigned int)portid, (mnl_cb_t)GoCb, (void *)data);
}


/*
 * http://stackoverflow.com/questions/1023261/is-there-a-way-to-do-currying-in-c
 * http://gcc.gnu.org/onlinedocs/gcc/Nested-Functions.html#Nested-Functions
 */

int
cb_run2_wrapper(const void *buf, size_t numbytes, uint32_t seq,
		uint32_t portid, uintptr_t data, uint16_t *ctl_types, size_t ctl_types_len)
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
			   (mnl_cb_t)GoCb, (void *)data, cb_ctl_array, NLMSG_MIN_TYPE - 1);
}
