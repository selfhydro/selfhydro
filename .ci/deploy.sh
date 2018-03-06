#!/bin/sh

set -e

scp selfhydro pi@water.local:/selfhydro/

ssh pi@water.local 'nohup sudo ./selfhydro/selfhydro &'
