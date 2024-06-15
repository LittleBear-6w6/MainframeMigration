infile = './input.txt'
outfile='./output.txt'
datasetFlag = 0

outFile = open(outfile,'w',encoding='utf-8')

with open(infile) as file_obuject:
    for line in file_obuject:
        if 'GDG BASE' in line:
            target='------'
            idx = line.find(target)
            dataName = line[idx + len(target):]
            datasetFlag=1
        
        if datasetFlag == 1:
            if 'LIMIT-----------------' in line:
                targetst='-----------------'
                idxst = line.find(targetst)
                targetend='SCRATCH'
                idxend = line.find(targetend)
                datalimit = line[idxst + len(targetst):idxend]
                datasetFlag=0

                outFile.write(f'{dataName.strip()},{datalimit.strip()}\n')
outFile.close()