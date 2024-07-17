#include <iostream>
#include <fstream>
#include <vector>

using namespace std;


std::ifstream in("/Users/yulkorolenko/Desktop/BioInf/6AdditivePhylogeny/phylog.txt");

struct LimbInf {
    int limb_size;
    int k;
    int j;
    int length;
    int cur_i;
};


struct DfsResult {
    int start;
    int end;
    int lenByStart;
    int lenByEnd;
};


LimbInf findLimbInf(vector<vector<int>> &matrix) {
    int i = matrix.size() - 1;
    LimbInf inf{};
    inf.cur_i = i;
    inf.limb_size = INT16_MAX;
    for (int j = 0; j < matrix.size(); ++j) {
        for (int k = 0; k < matrix.size(); ++k) {
            if (i == j || j == k || i == k) {
                continue;
            }
            int cur_limb = (matrix[j][i] + matrix[i][k] - matrix[j][k]) / 2;
            if (inf.limb_size > cur_limb) {
                inf.limb_size = cur_limb;
                inf.j = j;
                inf.k = k;
                inf.length = matrix[i][j] - inf.limb_size;
            }
        }
    }

    for (int j = 0; j < matrix.size(); ++j) {
        matrix[i][j] -= inf.limb_size;
        matrix[j][i] -= inf.limb_size;
    }

    return inf;
}

vector<vector<int>> getMatrix() {
    string line;
    std::getline(in, line);
    int n = stoi(line);
    vector<vector<int>> answer(n, vector<int>());

    for (int i = 0; i < n; ++i) {
        std::getline(in, line);
        for (int j = 0; j < line.size(); ++j) {
            int size = line.size();
            if (!isspace(line[j])) {
                string number;
                while (!isspace(line[j]) && j < size) {
                    number += line[j];
                    j++;
                }
                answer[i].push_back(stoi(number));
            }
        }
    }

    return answer;
}


int findPath(int cur, int end, vector<vector<pair<int, int>>> &tree, vector<bool> &used, vector<int> &answer) {
    used[cur] = true;
    if (cur == end) {
        answer.push_back(cur);
        return true;
    }
    for (int i = 0; i < tree[cur].size(); ++i) {
        if (used[tree[cur][i].first]) {
            continue;
        }
        bool isTrue = findPath(tree[cur][i].first, end, tree, used, answer);
        if (isTrue) {
            answer.push_back(cur);
            return true;
        }
    }
    return false;
}

int findWeight(vector<pair<int, int>> &pairs, int j) {
    for (auto value: pairs) {
        if (value.first == j) {
            return value.second;
        }
    }
    return -1;
}

void deleteEdge(int i, int j, vector<vector<pair<int, int>>> &tree) {
    for (int k = 0; k < tree[i].size(); ++k) {
        if (tree[i][k].first == j) {
            tree[i].erase(tree[i].cbegin() + k);
        }
    }

    for (int k = 0; k < tree[j].size(); ++k) {
        if (tree[j][k].first == i) {
            tree[j].erase(tree[j].cbegin() + k);
        }
    }
}

DfsResult findEdge(vector<vector<pair<int, int>>> &tree, vector<int> &path, int length) {
    DfsResult res{};
    int cur_length = 0;
    for (int i = 0; i < path.size(); ++i) {
        int weight = findWeight(tree[path[i]], path[i + 1]);
        if (cur_length < length && length < cur_length + weight) {
            res.start = path[i];
            res.end = path[i + 1];
            res.lenByEnd = cur_length + weight - length;
            res.lenByStart = weight - res.lenByEnd;
            if (res.lenByEnd != 0 && res.lenByStart != 0) {
                deleteEdge(path[i], path[i + 1], tree);
            }
            return res;
        }
        cur_length += weight;
    }
    return res;
}


vector<vector<pair<int, int>>> phylogenyRec(vector<vector<int>> matrix, int n) {
    if (matrix.size() == 2) {
        vector<vector<pair<int, int>>> tree(n, vector<pair<int, int>>());
        tree[0].emplace_back(1, matrix[0][1]);
        tree[1].emplace_back(0, matrix[0][1]);
        return tree;
    }

    LimbInf inf = findLimbInf(matrix);

    for (int i = 0; i < matrix.size(); ++i) {
        matrix[i].pop_back();
    }
    matrix.pop_back();

    vector<vector<pair<int, int>>> tree = phylogenyRec(matrix, n);

    vector<bool> used(tree.size(), false);
    vector<int> answer;
    findPath(inf.j, inf.k, tree, used, answer);

    reverse(answer.begin(), answer.end());

    DfsResult res = findEdge(tree, answer, inf.length);

    if (res.lenByEnd == 0) {
        tree[inf.cur_i].emplace_back(res.end, inf.limb_size);
        tree[res.end].emplace_back(inf.cur_i, inf.limb_size);
    } else if (res.lenByEnd == 0) {
        tree[inf.cur_i].emplace_back(res.start, inf.limb_size);
        tree[res.start].emplace_back(inf.cur_i, inf.limb_size);
    } else {
        int newVertex = tree.size();
        vector<pair<int, int>> newVector;
        tree.push_back(newVector);
        tree[newVertex].emplace_back(res.start, res.lenByStart);
        tree[res.start].emplace_back(newVertex, res.lenByStart);

        tree[newVertex].emplace_back(res.end, res.lenByEnd);
        tree[res.end].emplace_back(newVertex, res.lenByEnd);

        tree[inf.cur_i].emplace_back(newVertex, inf.limb_size);
        tree[newVertex].emplace_back(inf.cur_i, inf.limb_size);
    }

    return tree;
}

int main() {
    vector<vector<int>> matrix = getMatrix();

    vector<vector<pair<int, int>>> tree = phylogenyRec(matrix, matrix.size());

    cout << "Graph" << endl;
    for (int i = 0; i < tree.size(); ++i) {
        for (int j = 0; j < tree[i].size(); ++j) {
            cout << i << "->" << tree[i][j].first << ":" << tree[i][j].second << endl;
        }
    }
}
