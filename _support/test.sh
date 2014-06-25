#!/bin/bash
function STDERR () {
  cat - 1>&2
}
echo "stdout: foo"
echo "stderr: bar" | STDERR
