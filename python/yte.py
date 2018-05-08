# Sample Python code for user authorization

import os
import json
import pickle
import matplotlib.pyplot as plt

from googleapiclient.discovery import build
from googleapiclient.errors import HttpError
from google_auth_oauthlib.flow import InstalledAppFlow

# The CLIENT_SECRETS_FILE variable specifies the name of a file that contains
# the OAuth 2.0 information for this application, including its client_id and
# client_secret.
CLIENT_SECRETS_FILE = 'client_secret.json'

# This OAuth 2.0 access scope allows for full read/write access to the
# authenticated user's account and requires requests to use an SSL connection.
SCOPES = ['https://www.googleapis.com/auth/youtube.force-ssl']
API_SERVICE_NAME = 'youtube'
API_VERSION = 'v3'

tar_customURL = 'ROCKETBEANSTV'
tar_chID = ''
nr_videos = 5
include_subscriptions = True

max_pie_chart = 20


def get_authenticated_service():
    flow = InstalledAppFlow.from_client_secrets_file(CLIENT_SECRETS_FILE,
                                                     SCOPES)
    credentials = flow.run_console()
    return build(API_SERVICE_NAME, API_VERSION, credentials=credentials)


def get_chID_by_customURL(service, url):
    results = service.channels().list(part='snippet,contentDetails,statistics',
                                      forUsername=url).execute()
    return results['items'][0]['id']


def get_chName_by_chID(service, inc_id):
    results = service.channels().list(part='snippet,contentDetails,statistics',
                                      id=inc_id).execute()
    try:
        return results['items'][0]['snippet']['title']
    except (KeyError, IndexError):
        return ''


def get_upPlistID_by_chID(service, inc_id):
    results = service.channels().list(part='snippet,contentDetails,statistics',
                                      id=inc_id).execute()
    return results['items'][0]['contentDetails']['relatedPlaylists']['uploads']


def get_vidIDs_by_plistID(service, inc_id, max):
    results = service.playlistItems().list(part='snippet,contentDetails',
                                           playlistId=inc_id,
                                           maxResults=max).execute()
    vidID_list = []
    for x in range(0, len(results['items'])):
        videoID = results['items'][x]['snippet']['resourceId']['videoId']
        vidID_list.append(videoID)
    return vidID_list


def get_comrIDs_by_vidID(service, inc_id):
    results = service.commentThreads().list(part='snippet,replies',
                                            videoId=inc_id,
                                            maxResults=100).execute()
    comr_list = []
    for tl_comment in results['items']:
        # toplevel-comments
        tl_com_snippet = tl_comment['snippet']['topLevelComment']['snippet']
        tl_comrID = tl_com_snippet['authorChannelId']['value']
        if tl_comment not in comr_list:
            comr_list.append(tl_comrID)

        # check for replies
        if tl_comment['snippet']['totalReplyCount'] > 0:
            tl_commentID = tl_comment['snippet']['topLevelComment']['id']
            replies = service.comments().list(part='snippet',
                                              parentId=tl_commentID).execute()
            # add every chID(replyID) to comr_list
            for reply in replies['items']:
                replyID = reply['snippet']['authorChannelId']['value']
                if replyID not in comr_list:
                    comr_list.append(replyID)
    return comr_list


def get_all_comr_by_list_of_vidIDS(vidIDs):
    comrIDs = []
    for videoID in vidIDs:

        # get comrIDs
        tmp_comrIDs = get_comrIDs_by_vidID(service, videoID)

        for comrID in tmp_comrIDs:
            if comrID not in comrIDs:
                comrIDs.append(comrID)
        print('# videoID:'+str(vidIDs.index(videoID)+1)+'/'+str(nr_videos))
    return comrIDs


def get_favPlistID_by_chID(service, inc_id):
    results = service.channels().list(part='snippet,contentDetails,statistics',
                                      id=inc_id).execute()
    try:
        chID_contentDetails = results['items'][0]['contentDetails']
        return chID_contentDetails['relatedPlaylists']['favorites']
    except (KeyError, IndexError):
        return ''


def get_subIDs_by_chID(service, inc_id):
    subbed_chIDs = []
    try:
        results = service.subscriptions().list(part='snippet,contentDetails',
                                               channelId=inc_id).execute()
        for item in results['items']:
            subbed_chIDs.append(item['snippet']['resourceId']['channelId'])
    except (KeyError, IndexError, HttpError):
        pass
    return subbed_chIDs


def get_chID_by_vidID(service, inc_id):
    results = service.videos().list(part='snippet, contentDetails, statistics',
                                    id=inc_id).execute()
    try:
        return results['items'][0]['snippet']['channelId']
    except (KeyError, IndexError):
        return ''


def print_dict(dicto):
    print('######## new dict #########')
    print(json.dumps(dicto, indent=4))


def draw_piechart(inc_list):
    draw_list = []
    labels = ()
    sizes = []
    if len(inc_list) > max_pie_chart:
        for i in range(len(inc_list)-1, len(inc_list)-max_pie_chart-1, -1):
            draw_list.append(inc_list[i])
    else:
        draw_list = inc_list
    for elem in draw_list:
        labels = labels + (elem[0]+' '+str(elem[1]),)
        sizes.append(elem[1])
    fig1, ax1 = plt.subplots()
    ax1.pie(sizes, labels=labels, autopct='%1.1f%%',
            shadow=False, startangle=90)
    ax1.axis('equal')
    plt.show()


if __name__ == '__main__':
    # When running locally, disable OAuthlib's HTTPs verification. When
    # running in production *do not* leave this option enabled.
    os.environ['OAUTHLIB_INSECURE_TRANSPORT'] = '1'
    service = get_authenticated_service()

    rel_chDict = {}
    sorted_list = []

    # get all videos
    if tar_chID == '':
        tar_chID = get_chID_by_customURL(service, tar_customURL)
    tar_chName = get_chName_by_chID(service, tar_chID)
    tar_upPlistID = get_upPlistID_by_chID(service, tar_chID)
    tar_vidIDs = get_vidIDs_by_plistID(service, tar_upPlistID, nr_videos)

    # get all commentators of all videos
    comrIDs = get_all_comr_by_list_of_vidIDS(tar_vidIDs)

    # cycle through commentators and get wanted data
    for comrID in comrIDs:
        print('# comrID:' + str(comrIDs.index(comrID)+1) +
              '/' + str(len(comrIDs)))
        # get favoritePlaylists
        rel_chNames = []
        favPlistID = get_favPlistID_by_chID(service, comrID)
        # check whether fav playlist exists
        if favPlistID != '':
            rel_vidIDs = get_vidIDs_by_plistID(service, favPlistID, 50)
            for rel_vidID in rel_vidIDs:
                rel_chID = get_chID_by_vidID(service, rel_vidID)
                rel_chName = get_chName_by_chID(service, rel_chID)

                if (rel_chID not in rel_chNames or
                   rel_chName not in rel_chNames):
                    if rel_chName != '':
                        rel_chNames.append(rel_chName)
                    else:
                        rel_chNames.append(rel_chID)

        # get supscriptions for commentator
        if include_subscriptions:
            subbed_chIDs = get_subIDs_by_chID(service, comrID)
            for subbed_chID in subbed_chIDs:
                rel_chName = get_chName_by_chID(service, subbed_chID)
                if (subbed_chID not in rel_chNames or
                   rel_chName not in rel_chNames):
                    if rel_chName != '':
                        rel_chNames.append(rel_chName)
                    else:
                        rel_chNames.append(subbed_chID)

        # add related channels to dict
        for rel_chName in rel_chNames:
            if rel_chName == '' or rel_chName == tar_chName:
                continue
            if rel_chName in rel_chDict:
                rel_chDict[rel_chName] += 1
            else:
                rel_chDict[rel_chName] = 1

        # print sorted dict
        sorted_list = sorted(rel_chDict.items(), key=lambda x: x[1])
        print(sorted_list)
        with open('dumps/'+tar_chName+'.pkl', 'wb') as output:
            pickle.dump(sorted_list, output, pickle.HIGHEST_PROTOCOL)

    draw_piechart(sorted_list)
