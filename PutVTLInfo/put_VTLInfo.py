infile = './input.txt'
outfile='./output.csv'
datasetFlag = 0
 
dataName=''
datalimit=''
 
outFile = open(outfile,'w',encoding='utf-8')
 
with open(infile) as file_object:
    for line in file_object:
        if 'GDG BASE' in line:
            dataName=''
            datalimit=''
            target='------'
            idx = line.find(target)
            dataName = line[idx + len(target):]
            datasetFlag=1
        if datasetFlag == 1:
            if 'LIMIT' in line:
                targetst='LIMIT'
                idxst = line.find(targetst)
                targetend='SCRATCH'
                idxend = line.find(targetend)
                datalimit = line[idxst + len(targetst):idxend]
                datalimit = datalimit.replace('-','')
                datasetFlag=0
 
                outFile.write(f'{dataName.strip()},{datalimit.strip()}\n')
outFile.close()