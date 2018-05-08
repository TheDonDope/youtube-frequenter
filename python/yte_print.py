import pickle
import matplotlib.pyplot as plt

file_name = 'Martin Sellner'


def draw_piechart(inc_list):
    draw_list = []
    labels = ()
    sizes = []
    max = 15
    if len(inc_list) > max:
        for i in range(len(inc_list)-1, len(inc_list)-max-1, -1):
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
    sorted_list = []
    with open('dumps/'+file_name+'.pkl', 'rb') as input:
        sorted_list = pickle.load(input)
    draw_piechart(sorted_list)
