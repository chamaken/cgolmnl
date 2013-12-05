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

/**
 * mnl_cb_run3 - callback runqueue for netlink messages
 *
 * dispatch control message by mnl_ctl_cb_t
 * typedef int (*mnl_ctl_cb_t)(const struct nlmsghdr *nlh, uint16_t type, void *data);
 * == GoCtlCb()
 */
int
mnl_cb_run3(const void *buf, size_t numbytes, unsigned int seq,
	    unsigned int portid, mnl_cb_t cb_data, void *data,
	    mnl_ctl_cb_t cb_ctl)
{
	int ret = MNL_CB_OK, len = numbytes;
	const struct nlmsghdr *nlh = buf;

	while (mnl_nlmsg_ok(nlh, len)) {
		/* check message source */
		if (!mnl_nlmsg_portid_ok(nlh, portid)) {
			errno = ESRCH;
			return -1;
		}
		/* perform sequence tracking */
		if (!mnl_nlmsg_seq_ok(nlh, seq)) {
			errno = EPROTO;
			return -1;
		}

		/* dump was interrupted */
		if (nlh->nlmsg_flags & NLM_F_DUMP_INTR) {
			errno = EINTR;
			return -1;
		}

		/* netlink data message handling */
		if (nlh->nlmsg_type >= NLMSG_MIN_TYPE) { 
			if (cb_data){
				ret = cb_data(nlh, data);
				if (ret <= MNL_CB_STOP)
					goto out;
			}
		} if (cb_ctl) {
			ret = cb_ctl(nlh, nlh->nlmsg_type, data);
			if (ret <= MNL_CB_STOP)
				goto out;
		}
		nlh = mnl_nlmsg_next(nlh, &len);
	}
out:
	return ret;
}

int
cb_run3_wrapper(const void *buf, size_t numbytes, uint32_t seq,
		uint32_t portid, void *data)
{
	return mnl_cb_run3(buf, numbytes, (unsigned int)seq,
			   (unsigned int)portid, (mnl_cb_t)GoCb, data,
			   (mnl_ctl_cb_t)GoCtlCb);
}
