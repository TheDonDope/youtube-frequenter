import json
import matplotlib.pyplot as plt

def draw_piechart(inc_list):
    draw_list = []
    labels = ()
    sizes = []
    max = 15
    if len(inc_list) > max:
        draw_list = inc_list[:max]
    else:
        draw_list = inc_list
    for elem in draw_list:
        labels = labels + (elem['key']+' '+str(elem['value']),)
        sizes.append(elem['value'])
    _, ax1 = plt.subplots()
    ax1.pie(sizes, labels=labels, autopct='%1.1f%%',
            shadow=False, startangle=90)
    ax1.axis('equal')
    plt.show()


if __name__ == '__main__':
    sorted_results = []
    with open('results.json') as f:
      sorted_results = json.load(f)
    draw_piechart(sorted_results)
