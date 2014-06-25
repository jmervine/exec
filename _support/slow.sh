#!/bin/bash
function STDERR () {
  cat - 1>&2
}
sleep 5
echo "stdout: foo"
echo "stderr: bar" | STDERR
