#!/bin/sh
##
## Author(s):
##  -  <gestes@aldebaran-robotics.com>
##
## Copyright (C) 2011 Aldebaran Robotics
##

: ${BUILD:=build-sys-linux-x86_64}

if ! [ -d "../../${BUILD}" ] ; then
  echo "qi-messaging build folder not found..."
  echo "Please set BUILD to a valid value, current is value is '$BUILD'"
  echo "and '../../$BUILD' does not exists"
  exit 2
fi

QIMESSAGING_BUILD="../../${BUILD}" make
QIMESSAGING_BUILD="../../${BUILD}" make -C test
VERBOSE=5 LD_LIBRARY_PATH="../../${BUILD}/sdk/lib" ./test/test
