import os
import json
from mpl_toolkits.mplot3d import Axes3D
from matplotlib.pyplot import *
import random

useRandomResults = False

# default U25, C25, F25, d10ms, t60s,
channelID = 'UCCj956IF62FbT7Gouszaj9w'
filePath = 'output/channel-id-' + channelID + \
    '/channel-id-' + channelID + '-results.json'

# max values: u=50, c=100, f=50
uValues = [5, 12, 25, 38, 50]
cValues = [5, 12, 25, 50, 100]
fValues = [5, 12, 25, 38, 50]

# in seconds
timeout = 5

numOfTotalCalc = len(uValues) * len(cValues) * len(fValues)
timeoutPerCalculation = timeout
numOfCurrentCalc = 1
resulting3DArray = []
bestRun = ''
highestCount = 1

for U in uValues:
    for C in cValues:
        for F in fValues:
            if useRandomResults:
                R = random.randint(0, 10000)
                if R > highestCount:
                    bestrun = 'U:'+str(U)+'C:'+str(C)+'F:'+str(F)
                    highestCount = R
                resulting3DArray.append((U, C, F, R))
            else:
                # start youtube-frequenter
                binaryPath = './cmd/youtube-frequenter/main.go'
                arguments = '-c '+channelID+' ' +\
                            '-U '+str(U)+' ' +\
                            '-C '+str(C)+' ' +\
                            '-F '+str(F)+' ' +\
                            '-t '+str(timeoutPerCalculation)+'s'
                commandline = 'go run '+binaryPath+' '+arguments+' > /dev/null'
                os.system(commandline)

                # and get the results
                sorted_results = []
                with open(filePath) as f:
                    sorted_results = json.load(f)
                resultCount = len(sorted_results)
                resulting3DArray.append((U, C, F, resultCount, sorted_results))

                # and get the best parameter configuration
                if resultCount > highestCount:
                    highestCount = resultCount
                    bestRun = commandline

                # and print what happened since whole cycle
                countPrefix = '['+str(numOfCurrentCalc) + \
                    '/'+str(numOfTotalCalc)+']'
                print(countPrefix+' '+commandline+'=> #channels:' + str(resultCount) +
                      '(highest:'+str(highestCount)+')')
            numOfCurrentCalc += 1


print('bestRun:' + bestRun)
fig = matplotlib.pyplot.figure()
ax = fig.add_subplot(111, projection='3d')
for u, c, f, r, sr in resulting3DArray:
    if r > 0:
        ax.scatter(u, c, f, s=(1.*r)/highestCount*800, color=[(1.*u/50, 1.*c/100, 1.*f/50)],
                   alpha=(1.*r) / highestCount)
ax.set_xlabel('max uploads')
ax.set_ylabel('max comments')
ax.set_zlabel('max favorite')
show()
