import os
import json
from mpl_toolkits.mplot3d import Axes3D
from matplotlib.pyplot import *
import random

useRandomResults = False

# default U25, C25, F25, d10ms, t60s,
# os.system("go run main.go -u wwwKenFMde")
channelID = 'UC_dZp8bZipnjntBGLVHm6rw'
filePath = 'output/channel-id-' + channelID + \
    '/channel-id-' + channelID+'-results.json'
uValues = [10,15,25]
cValues = [10,25,50]
fValues = [10,15,25]
timeout = 30 # in seconds

numOfTotalCalculations = len(uValues) * len(cValues) * len(fValues)
timeoutPerCalculation = timeout
numOfCurrentCalculation = 1
resulting3DArray = []
bestRun = ''
highestCount = 1

for U in uValues:
  for C in cValues:
    for F in fValues:
      print('['+str(numOfCurrentCalculation)+'/'+str(numOfTotalCalculations)+']')
      if useRandomResults:
        R = random.randint(0, 10000)
        if R > highestCount:
          highestCount = R
        resulting3DArray.append((U, C, F, R))
      else:
        commandline = 'go run main.go -c '+channelID+' -U ' + \
            str(U)+' -C '+str(C)+' -F '+str(F)+' -t ' + \
            str(timeoutPerCalculation)+'s > /dev/null'
        print(commandline)
        os.system(commandline)
        sorted_results = []
        with open(filePath) as f:
          sorted_results = json.load(f)
        resultCount = len(sorted_results)
        resulting3DArray.append((U, C, F, resultCount))
        if resultCount > highestCount:
          highestCount = resultCount
          bestRun = commandline
        print("resultcount:"+str(resultCount)+' ,highest resultcount:'+str(highestCount))
      numOfCurrentCalculation += 1



print('bestRun:' + bestRun)
fig = matplotlib.pyplot.figure()
ax = fig.add_subplot(111, projection='3d')
for u, c, f, r in resulting3DArray:
  if r > 0:
    ax.scatter(u, c, f, s=(1.*r)/highestCount*1000, color='g',
               alpha=(1.*r)/ highestCount)
ax.set_xlabel('max uploads')
ax.set_ylabel('max comments')
ax.set_zlabel('max favorite')
show()
