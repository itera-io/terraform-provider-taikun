#!/usr/bin/env bash

set -e

remove_optional_id_from_datasources_doc() {
  echo "removing optional id attribute from list-type datasources' documentation"
  target_files=$(ls docs/data-sources/*s.md)
  offending_line="- \*\*id\*\* (String) The ID of this resource."
  modified_files=0
  for target_file in ${target_files}; do
    if ! grep -- "${offending_line}" "${target_file}" &>/dev/null; then
      continue
    fi
    echo "correcting \"${target_file}\""
    modified_files=$((${modified_files}+1))
    sed -i "/${offending_line}/d" "${target_file}"
    optional_header_line=$(grep -n '### Optional' "${target_file}" | cut -d ':' -f 1)
    read_only_header_line=$(grep -n '### Read-Only' "${target_file}" | cut -d ':' -f 1)
    if ! (sed -n "${optional_header_line},${read_only_header_line}p" "${target_file}" | fgrep -- '- **' &>/dev/null); then
      sed -i "${optional_header_line},$((${read_only_header_line}-1))d" "${target_file}"
    fi
  done
  echo "corrected ${modified_files} files"
}

remove_optional_id_from_datasources_doc
