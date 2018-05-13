import os
import json
from mpl_toolkits.mplot3d import Axes3D
from matplotlib.pyplot import *
import random

# default U25, C25, F25, d10ms, t60s,
# os.system("go run main.go -u wwwKenFMde")

uValues = [5,10,25,50]
cValues = [5,10,25,50,75,100]
fValues = [5,10,25,50]

resulting3DArray = []
bestRun = ''
highestCount = 1

for U in uValues:
  for C in cValues:
    for F in fValues:
      if False:
        commandline = 'go run main.go -u wwwKenFMde -U ' +str(U)+' -C '+str(C)+' -F '+str(F)+' -t 60s > /dev/null'
        print(commandline)
        os.system(commandline)
        sorted_results = []
        with open('results.json') as f:
          sorted_results = json.load(f)
        resultCount = len(sorted_results)
        resulting3DArray.append((U, C, F, resultCount))
        if resultCount > highestCount:
          highestCount = resultCount
          bestRun = commandline
          print('new highest resultCount:'+str(highestCount))
        else:
          print(resultCount)
      else:
        R = random.randint(0, 10000)
        if R > highestCount:
          highestCount = R
        resulting3DArray.append((U, C, F, R))

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
