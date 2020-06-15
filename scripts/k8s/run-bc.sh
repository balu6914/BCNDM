#!/bin/bash
. env.sh
echo "starting bc"
ssh $BCHOSTUSER@$BCHOSTIP << EOF
  cd datapace
  make runbcdev
EOF
echo "bc started"
