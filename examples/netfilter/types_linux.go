// +build ignore
package main

/*
#include <unistd.h>
#include <linux/netfilter/nfnetlink.h>
*/
import "C"

type (
	Size_t		C.size_t
)

// nfct-create-batch
const SizeofNfgenmsg	= C.sizeof_struct_nfgenmsg
type Nfgenmsg		C.struct_nfgenmsg
