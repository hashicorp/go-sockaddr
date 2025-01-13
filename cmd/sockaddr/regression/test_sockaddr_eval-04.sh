#!/bin/sh --
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e
exec 2>&1
cat <<'EOF' | exec ../sockaddr eval -
{{GetAllInterfaces | include "name" "lo0" | printf "%v"}}
EOF
