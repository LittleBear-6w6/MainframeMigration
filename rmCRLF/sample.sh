#!/bin/bash
xxd -p -g 0 -c 999 target.dat | sed s/0d0a$//g | xxd -p -r > conv_target.dat
