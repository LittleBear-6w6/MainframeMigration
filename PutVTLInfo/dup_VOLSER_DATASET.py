import csv
 
infile = './sorted_VTLInfo_daset_volser.csv'
outfile='./VTLInfo_dup_daset_volser.csv'
 
iniflag = 0
bvolser = ''
str = ''
 
outFile = open(outfile,'w',encoding='utf-8')
 
with open(infile) as file_object:
    csvreader = csv.reader(file_object)
    for row in csvreader:
        dataName,volser = row
        if (iniflag == 0):
            str = volser.strip()
            iniflag = 1
        else:
            if(bvolser != volser):
                outFile.write(f'{str.strip()}\n')
                str = ''
                str = volser.strip()
                bvolser = volser
        str += ',' + dataName.strip()
outFile.close()