#!/bin/sh
##
## Author(s):
##  - Cedric GESTES <gestes@aldebaran-robotics.com>
##
## Copyright (C) 2014 Aldebaran Robotics.com
##

qipath=~/src/qi
builddir=build-sys-linux-x86_64

export LD_LIBRARY_PATH="${qipath}/sdk/libqimessaging/${builddir}/sdk/lib"
export GOPATH="${qipath}/sdk/libqimessaging/bindings/go/"
export CGO_LDFLAGS="-L${qipath}/sdk/libqimessaging/${builddir}/sdk/lib/"
export CGO_CFLAGS="-I${qipath}/sdk/libqi -I${qipath}/sdk/libqimessaging/c/"
