#!/bin/bash
xxd -p -g 0 -c 52 target.dat | sed s/^/00340000/g | xxd -p -r > conv_target.dat
