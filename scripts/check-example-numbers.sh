#!/usr/bin/env bash
#
# check-example-numbers.sh — validate the examples/ directory numbering.
#
# Rules:
#   1. Every example dir matches `NN[_suffix]` or letter sub-example `NNx`/`NNx_suffix`
#      (e.g. 97a_option_quote). Non-numeric dirs (pkg/, graphify-out/) are ignored.
#   2. The main numeric sequence 0..N must be contiguous — no gaps. Letter sub-examples
#      (97a-97e) are excluded from the sequence; connection-group variants (00_connect,
#      00_connect_ha, 00_rsa_connect) count as a single "00".
#   3. No unexpected duplicate numeric prefixes (only 00-group and letter sub-examples
#      may share a prefix).
#   4. The header of docs/EXAMPLES.md must report the count of standalone examples.
#
# Usage: scripts/check-example-numbers.sh
# Exit 0 if all checks pass, nonzero otherwise.

set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
EXAMPLES_DIR="$ROOT/examples"
DOC="$ROOT/docs/EXAMPLES.md"

FAILED=0
fail() {
	echo "  FAIL: $*" >&2
	FAILED=1
}

declare -a dirs=()
while IFS= read -r -d '' d; do
	dirs+=("$(basename "$d")")
done < <(find "$EXAMPLES_DIR" -mindepth 1 -maxdepth 1 -type d -print0 | sort -z)

# --- Rule 1: name shape ----------------------------------------------------
declare -a numeric=()          # valid dir names
declare -a main_nums=()        # each main example contributes its numeric prefix once
declare -a letter_seen=()      # prefixes that have a letter sub-example
declare -a dup_prefixes=()     # prefixes seen more than once (allowed for 00 + letter)

for name in "${dirs[@]}"; do
	case "$name" in
		pkg | graphify-out) continue ;; # non-example dirs
	esac
	if [[ ! "$name" =~ ^([0-9]+)([a-z])?(_[A-Za-z0-9_]*)?$ ]]; then
		fail "$name: name must look like NN, NN_name, NNx or NNx_name"
		continue
	fi
	numeric+=("$name")
done

# --- Rule 2 + 3: collect unique main prefixes, detect dup prefixes ---------
declare -A counts=()
for name in "${numeric[@]}"; do
	# strip letter sub-example suffix (e.g. 97a_option_quote -> 97a)
	num="${name%%_*}"
	# num is like "97a" or "97" or "00"
	base="${num:0:${#num}-1}"
	last="${num: -1}"
	if [[ "$last" =~ [a-z] ]] && [[ "$base" =~ ^[0-9]+$ ]]; then
		# letter sub-example
		letter_seen+=("$base")
		counts["$num"]=$(( ${counts["$num"]:-0} + 1 ))
		continue
	fi
	main_nums+=("$num")
	counts["$num"]=$(( ${counts["$num"]:-0} + 1 ))
	if [ "${counts["$num"]}" -ge 2 ]; then
		dup_prefixes+=("$num")
	fi
done

# Duplicate prefixes allowed: 00 connection variants, or any prefix that also
# has a letter sub-example (e.g. 97 + 97a-97e).
allow_dup() {
	local p="$1"
	if [ "$p" == "00" ]; then return 0; fi
	for l in "${letter_seen[@]:-}"; do
		if [ "$l" == "$p" ]; then return 0; fi
	done
	return 1
}
for d in "${dup_prefixes[@]:-}"; do
	if ! allow_dup "$d"; then
		fail "duplicate numeric prefix '$d'"
	fi
done

letter_ok=1
for l in "${letter_seen[@]:-}"; do
	if [ "${counts["$l"]:-0}" -eq 0 ]; then
		fail "letter sub-example prefix '$l' has no base example '$l'"
	fi
done

# --- Rule 2: contiguous main sequence 0..N --------------------------------
if [ "${#main_nums[@]}" -gt 0 ]; then
	min="${main_nums[0]}"
	max="${main_nums[0]}"
	declare -A seen=()
	for n in "${main_nums[@]}"; do
		seen["$n"]=1
		if ((10#$n < 10#$min)); then min="$n"; fi
		if ((10#$n > 10#$max)); then max="$n"; fi
	done
	missing=0
	for ((i = 10#$min; i <= 10#$max; i++)); do
		# pad to 2+ digits to match dir naming (00, 65, 110)
		printf -v padded "%02d" "$i"
		if [ -z "${seen["$padded"]:-}" ] && [ -z "${seen["$i"]:-}" ]; then
			missing=$((missing + 1))
			if [ "$missing" -le 10 ]; then
				fail "missing example $padded in sequence $min..$max"
			fi
		fi
	done
fi

# --- Rule 4: EXAMPLES.md header count --------------------------------------
# The doc lists each example dir as a row (00_connect, 00_connect_ha and
# 00_rsa_connect are separate rows). Letter sub-examples (97a-97e) are NOT
# separate rows. So standalone count = number of numeric dirs minus letter subs.
letter_dirs=0
for name in "${numeric[@]}"; do
	num="${name%%_*}"
	if [[ "$num" =~ ^[0-9]+[a-z]$ ]]; then
		letter_dirs=$((letter_dirs + 1))
	fi
done
standalone=$(( ${#numeric[@]} - letter_dirs ))

header="$(grep -m1 -oE 'All [0-9]+ standalone examples' "$DOC" 2>/dev/null | grep -oE '[0-9]+')"
if [ -z "$header" ]; then
	fail "could not parse the example count from $DOC"
elif [ "$header" != "$standalone" ]; then
	fail "docs/EXAMPLES.md says $header standalone examples, but examples/ has $standalone"
fi

echo "check-example-numbers: ${#numeric[@]} numeric dirs / ${standalone} standalone examples"
if [ "$FAILED" -eq 1 ]; then
	echo "check-example-numbers: FAILED"
	exit 1
fi
echo "check-example-numbers: OK"
exit 0