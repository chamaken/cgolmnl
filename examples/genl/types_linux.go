// +build ignore
package main

/*
#include <linux/genetlink.h>
*/
import "C"

const SizeofGenlmsghdr = C.sizeof_struct_genlmsghdr

type Genlmsghdr C.struct_genlmsghdr
