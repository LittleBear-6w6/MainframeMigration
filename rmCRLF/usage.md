# Specification
This shell script is the simple tool for outputting file which add RDW(record length) into beginning of each line.

```bash:usage
xxd -p -g 0 -c [RecordLength]  target.dat | sed s/0d0a$//g | xxd -p -r > conv_target.dat
```
## Process  
1. Convert the file to hex file and break line per record length.
3. Remove CRLF(0d0a) on each line
4. Reverse file into binary file.
