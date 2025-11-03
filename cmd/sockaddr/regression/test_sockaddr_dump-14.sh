#!/bin/sh --
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: MPL-2.0


set -e
exec 2>&1
../sockaddr dump -4 '::c0a8:1'
../sockaddr dump -6 '::c0a8:1'
../sockaddr dump -4 '::c0a8:1/112'
../sockaddr dump -6 '::c0a8:1/112'
../sockaddr dump -4 '0:0:0:0:0:ffff:c0a8:1/112'
../sockaddr dump -6 '0:0:0:0:0:ffff:c0a8:1/112'
../sockaddr dump -4 '[0:0:0:0:0:ffff:c0a8:1/112]'
../sockaddr dump -6 '[0:0:0:0:0:ffff:c0a8:1/112]'
