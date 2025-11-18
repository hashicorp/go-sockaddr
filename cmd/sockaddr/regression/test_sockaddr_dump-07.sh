#!/bin/sh --
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: MPL-2.0


set -e
exec 2>&1
exec ../sockaddr dump -H '[2001:db8::7]:22'
