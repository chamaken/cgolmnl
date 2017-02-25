#include <sys/types.h>
#include <stdbool.h>

struct mnl_nlmsg_batch {
	/* the buffer that is used to store the batch. */
	void *buf;
	size_t limit;
	size_t buflen;
	/* the current netlink message in the batch. */
	void *cur;
	bool overflow;
};

size_t mnl_nlmsg_batch_limit(struct mnl_nlmsg_batch *b)
{
	return b->limit;
}
