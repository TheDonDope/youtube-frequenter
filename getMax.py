import os
import json
from mpl_toolkits.mplot3d import Axes3D
from matplotlib.pyplot import *
import random

useRandomResults = False

# default U25, C25, F25, d10ms, t60s,
# os.system("go run main.go -u wwwKenFMde")
customURL = 'wwwKenFMde'
filePath = 'results.json'
# filePath = 'results/custom-url-' + customURL+'/custom-url-results.json'
uValues = [25,50]
cValues = [10,50]
fValues = [25,50]

resulting3DArray = []
bestRun = ''
highestCount = 1

for U in uValues:
  for C in cValues:
    for F in fValues:
      if useRandomResults:
        R = random.randint(0, 10000)
        if R > highestCount:
          highestCount = R
        resulting3DArray.append((U, C, F, R))
      else:
        commandline = 'go run main.go -u wwwKenFMde -U ' +str(U)+' -C '+str(C)+' -F '+str(F)+' -t 60s -o results> /dev/null'
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
          print('new highest resultCount:'+str(highestCount))
        else:
          print(resultCount)


print('bestRun:' + bestRun)
fig = matplotlib.pyplot.figure()
ax = fig.add_subplot(111, projection='3d')
for u, c, f, r in resulting3DArray:
  if r > 0:
    ax.scatter(u, c, f, color='g', alpha=1. * r / highestCount)
ax.set_xlabel('U')
ax.set_ylabel('C')
ax.set_zlabel('F')
show()
