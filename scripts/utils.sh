#!/bin/bash
set -e

transform_action() {
    local action="$1"
    case "$action" in
        *ing) echo "${action%ing}ed" ;;
        *) echo "$action" ;;
    esac
}

calculate_elapsed_time() {
    local start="$1"
    local end="$2"
    local elapsed=$(echo "$end - $start" | bc)
    if (( $(echo "$elapsed >= 60" | bc -l) )); then
        local minutes=$(echo "$elapsed / 60" | bc)
        local seconds=$(echo "$elapsed % 60" | bc)
        printf "%.2fs (%dm %.2fs)" "$elapsed" "$minutes" "$seconds"
    else
        printf "%.2fs" "$elapsed"
    fi
}

with_progress_bar() {
    local msg="$1"
    local cmd="$2"
    local action="${3:-processing}"
    local verbose="$4"
    local spinner=(⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏)
    local i=0
    local start_time=$(date +%s.%N)
    local tmp_output=$(mktemp)
    trap "rm -f '$tmp_output'" EXIT

    bash -c "$cmd" >"$tmp_output" 2>&1 &
    local pid=$!

    echo -ne "\033[0;36m${spinner[0]}\033[0m $action $msg"

    if [ "$verbose" == "--verbose" ]; then
        tail -f "$tmp_output" &
        local tail_pid=$!
    fi

    while kill -0 $pid 2>/dev/null; do
        printf "\r\033[0;36m${spinner[$i]}\033[0m $action $msg..."
        sleep 0.2
        i=$(( (i + 1) % ${#spinner[@]} ))
    done

    wait $pid
    local exit_code=$?
    local end_time=$(date +%s.%N)
    local elapsed_formatted=$(calculate_elapsed_time "$start_time" "$end_time")
    local transformed_action=$(transform_action "$action")

    if [ "$verbose" == "--verbose" ]; then
        kill $tail_pid
    fi

    if [ $exit_code -eq 0 ]; then
        echo -e "\r\033[0;32m✓\033[0m Successfully $transformed_action $msg in $elapsed_formatted"
    else
        echo -e "\r\033[0;31m✗\033[0m Failed to $action $msg (after $elapsed_formatted)"
        cat "$tmp_output"
    fi
    rm -f "$tmp_output"
}
