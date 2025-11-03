#!/bin/sh --
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: MPL-2.0


set -e
exec 2>&1
exec ../sockaddr eval 'GetAllInterfaces | include "name" "lo0" | printf "%v"' 'GetAllInterfaces | include "name" "lo0" | printf "%v"'
