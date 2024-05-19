# Specification
This shell script is the simple tool for outputting file which add RDW(record length) into beginning of each line.

```bash:usage
xxd -p -g 0 -c [RecordLength] target.dat | sed s/^/[Hexadecimal record length]0000/g | xxd -p -r > conv_target.dat
```
## Process  
1. Convert the file to hex file and break line per record length.
3. Add RDW('Hexadecimal record length' + '0000') into beginning of each line.
4. Reverse file into binary file.

## About RDW
The first 2 bytes are the record length, and the remaining 2 bytes are reserved for spanned records.
