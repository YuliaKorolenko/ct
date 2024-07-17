eps = 0.01
def checker():
    f1 = open("kostyaanswer.txt", "r")
    f2 = open("myanswer.txt", "r")
    lines1 = f1.readlines()
    lines2 = f2.readlines()
    if (len(lines1) != len(lines2)):
        print("WA different number of lines")
        return
    for i in range(len(lines1)):
        arr1 = list(map(float, lines1[i].split()))
        arr2 = list(map(float, lines2[i].split()))
        if (len(arr1) != len(arr2)):
            print("WA different number of floats on line", i + 1)
            return
        for j in range(len(arr1)):
            if (abs(arr1[j] - arr2[j]) > eps):
                print("WA numbers", j + 1, "on line", i + 1, "are different")
                print(arr1[j], "vs", arr2[j])
                return
    print("OK")


checker()