#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source_dir="$(cd "${script_dir}/.." && pwd)"
cmdref_dir="${script_dir}/cmdref"

generators=(
    "bugtool/cilium-bugtool cmdref"
    "cilium-cli/cilium cmdref"
    "cilium-dbg/cilium-dbg cmdref"
    "cilium-health/cilium-health cmdref"
    "clustermesh-apiserver/clustermesh-apiserver cmdref"
    "daemon/cilium-agent cmdref"
    "operator/cilium-operator cmdref"
    "operator/cilium-operator-alibabacloud cmdref"
    "operator/cilium-operator-aws cmdref"
    "operator/cilium-operator-azure cmdref"
    "operator/cilium-operator-generic cmdref"
)

for g in "${generators[@]}" ; do
    ${source_dir}/${g} "${cmdref_dir}"
done
