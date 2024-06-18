infile = './input.txt'
outfile='./VTLInfo_daset_volser.csv'
datasetFlag = 0
 
dataName=''
datalimit=''
 
outFile = open(outfile,'w',encoding='utf-8')
 
with open(infile) as file_object:
    for line in file_object:
        if '0   NONVSAM ---- ' in line:
            dataName=''
            datalimit=''
            target='NONVSAM'
            idx = line.find(target)
            dataName = line[idx + len(target):]
            dataName = dataName.replace('-','')
            datasetFlag = 1
        if datasetFlag == 1:
            if 'VOLSER' in line:
                targetst='VOLSER'
                idxst = line.find(targetst)
                targetend='DEVTYPE'
                idxend = line.find(targetend)
                datalimit = line[idxst + len(targetst):idxend]
                datalimit = datalimit.replace('-','')
                datasetFlag=0
 
                outFile.write(f'{dataName.strip()},{datalimit.strip()}\n')
outFile.close()