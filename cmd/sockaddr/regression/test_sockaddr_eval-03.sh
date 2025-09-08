#!/bin/sh --
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e
exec 2>&1
exec ../sockaddr eval '. | include "name" "lo0" | include "type" "IPv6" | sort "address" | join "address" " "'
