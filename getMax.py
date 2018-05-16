import os
import json
import matplotlib.pyplot as plt
import random
from mpl_toolkits.mplot3d import Axes3D
import math
from printPieChart import draw_piechart

useRandomResults = False

# default U25, C25, F25, d10ms, t60s,
customURL = 'wwwKenFMde'
channelID = ''
if channelID != '':
    filePath = 'output/channel-id-' + channelID + \
               '/channel-id-' + channelID + '-results.json'
else:
    filePath = 'output/custom-url-' + customURL + \
               '/custom-url-' + customURL + '-results.json'


# timeout which terminates each youtube-frequenter process
# in seconds
timeout = 10

# max values: u=50, c=100, f=50
uMax = 10
cMax = 25
fMax = 25

uValues = []
cValues = []
fValues = []

for u in range(1, uMax, int(uMax/3.)):
    uValues.append(u)

for c in range(1, cMax, int(cMax/3.)):
    cValues.append(c)

for f in range(1, fMax, int(fMax/3.)):
    fValues.append(f)

dataPointSize = 500

numOfTotalCalc = len(uValues) * len(cValues) * len(fValues)
timeoutPerCalculation = timeout
numOfCurrentCalc = 1
resulting3DArray = []
bestRun = ''
highestCount = 1
bestSorted_results = []

# search through valid parameter configurations
for U in uValues:
    for C in cValues:
        for F in fValues:
            resultCount = 0
            commandline = ''
            if useRandomResults:
                # fake commandline
                commandline = 'U:' + str(U) + 'C:' + str(C) + 'F:' + str(F)

                # get results ''for commandline''
                resultCount = random.randint(0, 10000)

                # save best random configuration
                if resultCount > highestCount:
                    bestrun = commandline
                    highestCount = resultCount

                    bestSorted_results = [{'key': 'channel1', 'value': 200},
                                          {'key': 'channel2', 'value': 50}]
                resulting3DArray.append((U, C, F, resultCount, ['a']))
            else:
                # start youtube-frequenter
                binaryPath = './cmd/youtube-frequenter/main.go'
                target = ''
                if channelID != '':
                    target = '-c ' + channelID
                else:
                    target = '-u ' + customURL

                arguments = target + ' ' +\
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

                # and get the best parameter configuration
                if resultCount > highestCount:
                    highestCount = resultCount
                    bestRun = commandline
                    bestSorted_results = sorted_results

                resulting3DArray.append((U, C, F, resultCount, sorted_results))
            # and print what happened since whole cycle
            countPrefix = '['+str(numOfCurrentCalc) + \
                '/'+str(numOfTotalCalc)+']'

            print(countPrefix+' '+commandline+'=> #channels:' + str(resultCount) +
                  '(highest:'+str(highestCount)+')')
            numOfCurrentCalc += 1

if __name__ == '__main__':

    print('bestRun:' + bestRun)
    fig = plt.figure()
    ax = fig.add_subplot(211, projection='3d')

    # add every datapoint to 3d-plot
    for u, c, f, r, sr in resulting3DArray:
        if r > 0:
            ax.scatter(u, c, f, s=(1.*r)/highestCount*dataPointSize, color=[(1.*u/50, 1.*c/100, 1.*f/50)],
                       alpha=(1. * r) / highestCount)

    # set axis names
    ax.set_xlabel('max uploads')
    ax.set_ylabel('max comments')
    ax.set_zlabel('max favorite')

    # append piechart with pos=212(?)
    draw_piechart(fig, bestSorted_results)

    # view graphic
    plt.show()
