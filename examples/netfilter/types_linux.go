// +build ignore
package main

/*
#include <unistd.h>
#include <linux/netfilter/nfnetlink.h>
#include <linux/netfilter/nfnetlink_log.h>
*/
import "C"

type (
	Size_t		C.size_t
)

// nfct-create-batch
const SizeofNfgenmsg	= C.sizeof_struct_nfgenmsg
type Nfgenmsg		C.struct_nfgenmsg

// nf-log
const SizeofNfulnlMsgPacketTimestamp	= C.sizeof_struct_nfulnl_msg_packet_timestamp
type NfulnlMsgPacketTimestamp		C.struct_nfulnl_msg_packet_timestamp
const SizeofNfulnlMsgPacketHw		= C.sizeof_struct_nfulnl_msg_packet_hw
type NfulnlMsgPacketHw			C.struct_nfulnl_msg_packet_hw
const SizeofNfulnlMsgPacketHdr		= C.sizeof_struct_nfulnl_msg_packet_hdr
type NfulnlMsgPacketHdr			C.struct_nfulnl_msg_packet_hdr
const SizeofNfulnlMsgConfigCmd		= C.sizeof_struct_nfulnl_msg_config_cmd
type NfulnlMsgConfigCmd			C.struct_nfulnl_msg_config_cmd
const SizeofNfulnlMsgConfigMode		= C.sizeof_struct_nfulnl_msg_config_mode
type NfulnlMsgConfigMode		C.struct_nfulnl_msg_config_mode
