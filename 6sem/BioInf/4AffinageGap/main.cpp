#include <iostream>
#include <vector>
#include <map>

using namespace std;

int SO = 11;
int SE = 1;

int MIN = INT_MIN + 11 * 10;


std::pair<char, std::pair<int, int>> createVer(char mas, int i, int j) {
    return {mas, {i, j}};
}


vector<std::vector<int>> mp = {{4,  0,  -2, -1, -2, 0,  -2, -1, -1, -1, -1, -2, -1, -1, -1, 1,  0,  0,  -3, -2},
                               {0,  9,  -3, -4, -2, -3, -3, -1, -3, -1, -1, -3, -3, -3, -3, -1, -1, -1, -2, -2},
                               {-2, -3, 6,  2,  -3, -1, -1, -3, -1, -4, -3, 1,  -1, 0,  -2, 0,  -1, -3, -4, -3},
                               {-1, -4, 2,  5,  -3, -2, 0,  -3, 1,  -3, -2, 0,  -1, 2,  0,  0,  -1, -2, -3, -2},
                               {-2, -2, -3, -3, 6,  -3, -1, 0,  -3, 0,  0,  -3, -4, -3, -3, -2, -2, -1, 1,  3},
                               {0,  -3, -1, -2, -3, 6,  -2, -4, -2, -4, -3, 0,  -2, -2, -2, 0,  -2, -3, -2, -3},
                               {-2, -3, -1, 0,  -1, -2, 8,  -3, -1, -3, -2, 1,  -2, 0,  0,  -1, -2, -3, -2, 2},
                               {-1, -1, -3, -3, 0,  -4, -3, 4,  -3, 2,  1,  -3, -3, -3, -3, -2, -1, 3,  -3, -1},
                               {-1, -3, -1, 1,  -3, -2, -1, -3, 5,  -2, -1, 0,  -1, 1,  2,  0,  -1, -2, -3, -2},
                               {-1, -1, -4, -3, 0,  -4, -3, 2,  -2, 4,  2,  -3, -3, -2, -2, -2, -1, 1,  -2, -1},
                               {-1, -1, -3, -2, 0,  -3, -2, 1,  -1, 2,  5,  -2, -2, 0,  -1, -1, -1, 1,  -1, -1},
                               {-2, -3, 1,  0,  -3, 0,  1,  -3, 0,  -3, -2, 6,  -2, 0,  0,  1,  0,  -3, -4, -2},
                               {-1, -3, -1, -1, -4, -2, -2, -3, -1, -3, -2, -2, 7,  -1, -2, -1, -1, -2, -4, -3},
                               {-1, -3, 0,  2,  -3, -2, 0,  -3, 1,  -2, 0,  0,  -1, 5,  1,  0,  -1, -2, -2, -1},
                               {-1, -3, -2, 0,  -3, -2, 0,  -3, 2,  -2, -1, 0,  -2, 1,  5,  -1, -1, -3, -3, -2},
                               {1,  -1, 0,  0,  -2, 0,  -1, -2, 0,  -2, -1, 1,  -1, 0,  -1, 4,  1,  -2, -3, -2},
                               {0,  -1, -1, -1, -2, -2, -2, -1, -1, -1, -1, 0,  -1, -1, -1, 1,  5,  0,  -2, -2},
                               {0,  -1, -3, -2, -1, -3, -3, 3,  -2, 1,  1,  -3, -2, -2, -3, -2, 0,  4,  -3, -1},
                               {-3, -2, -4, -3, 1,  -2, -2, -3, -3, -2, -1, -4, -4, -2, -3, -3, -2, -3, 11, 2},
                               {-2, -2, -3, -2, 3,  -3, 2,  -1, -2, -1, -1, -2, -3, -1, -2, -2, -2, -1, 2,  7}
};

map<char, int> char_to_int = {{'A', 0},
                              {'C', 1},
                              {'D', 2},
                              {'E', 3},
                              {'F', 4},
                              {'G', 5},
                              {'H', 6},
                              {'I', 7},
                              {'K', 8},
                              {'L', 9},
                              {'M', 10},
                              {'N', 11},
                              {'P', 12},
                              {'Q', 13},
                              {'R', 14},
                              {'S', 15},
                              {'T', 16},
                              {'V', 17},
                              {'W', 18},
                              {'Y', 19}};


int is_equal(char A, char B) {
    int newA = char_to_int[A];
    int newB = char_to_int[B];
    int f = mp[newA][newB];
    return f;
}


void print_matrix(std::vector<std::vector<int>> &a) {
    for (int i = 0; i < a.size(); ++i) {
        for (int j = 0; j < a[i].size(); ++j) {
            if (a[i][j] >= 0 && a[i][j] < 10) {
                std::cout << " " << a[i][j] << " ";
            } else if (a[i][j] == MIN) {
                std::cout << " inf ";
            } else {
                std::cout << a[i][j] << " ";
            }

        }
        std::cout << std::endl;
    }
}

char findFirst(std::vector<std::vector<int>> &M, std::vector<std::vector<int>> &I, std::vector<std::vector<int>> &D) {
    int i = M.size() - 1;
    int j = M[0].size() - 1;
    if (M[i][j] > I[i][j] && M[i][j] > D[i][j]) {
        return 'M';
    } else if (D[i][j] > M[i][j] && D[i][j] > I[i][j]) {
        return 'D';
    } else {
        return 'I';
    }
}

int findMax(std::vector<std::vector<int>> &M, std::vector<std::vector<int>> &I, std::vector<std::vector<int>> &D) {
    int i = M.size() - 1;
    int j = M[0].size() - 1;
    return std::max(M[i][j], std::max(I[i][j], D[i][j]));
}

int main() {
    std::string x_string, y_string;
    std::cin >> x_string;
    std::cin >> y_string;


    std::vector<std::vector<int>> M(x_string.size() + 1, std::vector<int>(y_string.size() + 1, MIN));
    M[0][0] = 0;
//  заканчивающегося совпадением или несоответствием символа

    std::vector<std::vector<int>> I(x_string.size() + 1, std::vector<int>(y_string.size() + 1, MIN));
    I[1][0] = -12;
    for (int i = 2; i <= x_string.size(); ++i) {
        I[i][0] = I[i - 1][0] - 1;
    }
//  оценка наилучшего выравнивания x_string[1..i] и y_string[1..j], оканчивающихся пробелом в X.

    std::vector<std::vector<int>> D(x_string.size() + 1, std::vector<int>(y_string.size() + 1, MIN));
    D[0][1] = -12;
    for (int i = 2; i <= x_string.size(); ++i) {
        D[0][i] = D[0][i - 1] - 1;
    }
//  оценка наилучшего выравнивания x_string[1..i] и y_string[1..j], заканчивающихся пробелом в Y.


    map<std::pair<char, std::pair<int, int>>, std::pair<char, std::pair<int, int>>> answer;


    for (int i = 1; i < x_string.size() + 1; ++i) {
        for (int j = 1; j < y_string.size() + 1; ++j) {
            if (i > 1) {
                I[i][j] = std::max(I[i - 1][j] - SE, M[i - 1][j] - SO);
                if (I[i][j] == I[i - 1][j] - SE) {
                    answer[createVer('I', i, j)] = createVer('I', i - 1, j);
                } else {
                    answer[createVer('I', i, j)] = createVer('M', i - 1, j);
                }
            }
            if (j > 1) {
                D[i][j] = std::max(D[i][j - 1] - SE, M[i][j - 1] - SO);
                if (D[i][j] == D[i][j - 1] - SE) {
                    answer[createVer('D', i, j)] = createVer('D', i, j - 1);
                } else {
                    answer[createVer('D', i, j)] = createVer('M', i, j - 1);
                }
            }
            int cur_equal = is_equal(x_string[i - 1], y_string[j - 1]);
            M[i][j] = cur_equal + std::max(std::max(D[i - 1][j - 1], I[i - 1][j - 1]), M[i - 1][j - 1]);
            if (M[i][j] == D[i - 1][j - 1] + cur_equal) {
                answer[createVer('M', i, j)] = createVer('D', i - 1, j - 1);
            } else if (M[i][j] == I[i - 1][j - 1] + cur_equal) {
                answer[createVer('M', i, j)] = createVer('I', i - 1, j - 1);
            } else {
                answer[createVer('M', i, j)] = createVer('M', i - 1, j - 1);
            }
        }
    }

    std::cout << findMax(M, I, D) << std::endl;

    int i = x_string.size();
    int j = y_string.size();
    char start_char = findFirst(M, I, D);
    string answer_x = "";
    string answer_y = "";
    while (i > 0 || j > 0) {
        std::pair<char, std::pair<int, int>> a = answer[createVer(start_char, i, j)];
        char pred_char = a.first;
        int pred_i = a.second.first;
        int pred_j = a.second.second;
        if (pred_i == i - 1 && pred_j == j - 1) {
            answer_x = x_string[i - 1] + answer_x;
            answer_y = y_string[j - 1] + answer_y;
        } else if (pred_i == i - 1) {
            answer_x = x_string[i - 1] + answer_x;
            answer_y = "-" + answer_y;
        } else {
            answer_x = "-" + answer_x;
            answer_y = y_string[j - 1] + answer_y;
        }
        start_char = pred_char;
        i = pred_i;
        j = pred_j;
    }

    std::cout << answer_x << std::endl;
    std::cout << answer_y;
}



