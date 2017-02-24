// Go wrapper of libmnl using cgo
//
// ---- Citing the original libmnl
//
// libmnl is a minimalistic user-space library oriented to Netlink developers.
// There are a lot of common tasks in parsing, validating, constructing of
// both the Netlink header and TLVs that are repetitive and easy to get wrong.
// This library aims to provide simple helpers that allows you to avoid
// re-inventing the wheel in common Netlink tasks.
//
//     "Simplify, simplify" -- Henry David Thoureau. Walden (1854)
//
// The acronym libmnl stands for LIBrary Minimalistic NetLink.
//
// libmnl homepage is:
//      http://www.netfilter.org/projects/libmnl/
//
// Main Features
// - Small: the shared library requires around 30KB for an x86-based computer.
// - Simple: this library avoids complex abstractions that tend to hide Netlink
//   details. It avoids elaborated object-oriented infrastructure and complex
//   callback-based workflow.
// - Easy to use: the library simplifies the work for Netlink-wise developers.
//   It provides functions to make socket handling, message building,
//   validating, parsing and sequence tracking, easier.
// - Easy to re-use: you can use this library to build your own abstraction
//   layer upon this library, if you want to provide another library that
//   hides Netlink details to your users.
// - Decoupling: the interdependency of the main bricks that compose this
//   library is reduced, i.e. the library provides many helpers, but the
//   programmer is not forced to use them.
//
// Licensing terms
//   This library is released under the LGPLv2.1 or any later (at your option).
//
// Dependencies
//   You have to install the Linux kernel headers that you want to use to develop
//   your application. Moreover, this library requires that you have some basics
//   on Netlink.
//
// Git Tree
//   The current development version of libmnl can be accessed at:
//   http://git.netfilter.org/cgi-bin/gitweb.cgi?p=libmnl.git;a=summary
//
// Using libmnl
//   You can access several example files under examples/ in the libmnl source
//   code tree.
//
package cgolmnl

// XXX: can not do coverage test: https://code.google.com/p/go/issues/detail?id=6333

// Netlink message:
//
//	|<----------------- 4 bytes ------------------->|
//	|<----- 2 bytes ------>|<------- 2 bytes ------>|
//	|-----------------------------------------------|
//	|      Message length (including header)        |
//	|-----------------------------------------------|
//	|     Message type     |     Message flags      |
//	|-----------------------------------------------|
//	|           Message sequence number             |
//	|-----------------------------------------------|
//	|                 Netlink PortID                |
//	|-----------------------------------------------|
//	|                                               |
//	.                   Payload                     .
//	|_______________________________________________|
//
// There is usually an extra header after the the Netlink header (at the
// beginning of the payload). This extra header is specific of the Netlink
// subsystem. After this extra header, it comes the sequence of attributes
// that are expressed in Type-Length-Value (TLV) format.
type Nlmsg struct {
	*Nlmsghdr
	buf []byte
}

// Netlink Type-Length-Value (TLV) attribute:
//
//	|<-- 2 bytes -->|<-- 2 bytes -->|<-- variable -->|
//	-------------------------------------------------
//	|     length    |      type     |      value     |
//	-------------------------------------------------
//	|<--------- header ------------>|<-- payload --->|
//
// The payload of the Netlink message contains sequences of attributes that are
// expressed in TLV format.
type Nlattr nlattr
