#include <errno.h>
#include "set_errno.h"

// really?
void SetErrno(int n)
{
	errno = n;
}
