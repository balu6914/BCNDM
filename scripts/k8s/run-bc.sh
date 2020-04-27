#!/bin/bash
. env.sh
echo "starting bc"
ssh ubuntu@$BCHOSTIP << EOF
  cd datapace
  make runbcdev
EOF
echo "bc started"
