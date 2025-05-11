#!/bin/bash
spinner=('⠋' '⠙' '⠹' '⠸' '⠼' '⠴' '⠦' '⠧' '⠇' '⠏')
i=0
echo "Starting spinner test..."
for _ in $(seq 1 20); do
  i=$(( (i + 1) % ${#spinner[@]} ))
  printf "\r\033[0;36m%s\033[0m Loading..." "${spinner[$i]}"
  sleep 0.1
done
echo -e "\r\033[0;32m✓\033[0m Done!"
