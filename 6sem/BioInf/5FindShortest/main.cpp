#include <iostream>
#include <vector>
#include <fstream>
#include <set>


using namespace std;

struct Node {
    int number;
    int value;
    int startCycle;
    int endCycle;
};

std::ifstream in("/Users/yulkorolenko/Desktop/BioInf/5FindShortest/shortest.txt");

vector<int> getCycle(string line) {
    vector<int> cycle;
    line = line.substr(1, line.size() - 2);
    for (int i = 0; i < line.size(); ++i) {
        string cur_int;
        if (line[i] == '-' || line[i] == '+') {
            while (!isspace(line[i])) {
                cur_int += line[i];
                i++;
            }
            cycle.push_back(atoi(cur_int.c_str()));
        }
    }

    return cycle;
}

set<int> getSetVertex(vector<int> &start, vector<int> &end) {
    set<int> answer;
    for (int value: start) {
        if (answer.find(abs(value)) == answer.end()) {
            answer.insert(abs(value));
        }
    }
    for (int value: end) {
        if (answer.find(abs(value)) == answer.end()) {
            answer.insert(abs(value));
        }
    }
    return answer;
}


Node createNode(int value, int number) {
    Node n{};
    n.value = value;
    n.number = number;
    return n;
}

int numberToIndex(int a, bool is_exit) {
    if (a < 0) {
//        0 <- 1
        if (is_exit) {
            return (abs(a) - 1) * 2;
        } else {
            return (abs(a) - 1) * 2 + 1;
        }
    } else {
//        0 -> 1
        if (is_exit) {
            return (a - 1) * 2 + 1;
        } else {
            return (a - 1) * 2;
        }
    }
}


vector<vector<int>> findCycles(vector<Node> &nodes) {
    vector<bool> visited(nodes.size(), false);
    vector<vector<int>> cycles;
    for (int i = 0; i < nodes.size(); ++i) {
        if (!visited[i]) {
            int curVertex = i;
            bool is_first = true;
            vector<int> cur_cycle;
            while (!visited[curVertex]) {
                visited[curVertex] = true;
                cur_cycle.push_back(curVertex);
                if (is_first) {
                    curVertex = nodes[curVertex].endCycle;
                    is_first = false;
                } else {
                    curVertex = nodes[curVertex].startCycle;
                    is_first = true;
                }
            }
            if (cur_cycle.size() > 2) {
                cycles.push_back(cur_cycle);
            }
        }
    }
    return cycles;

}


void twoBreakGraph(vector<Node> &nodes, vector<int> cycle) {
    nodes[cycle[0]].startCycle = cycle[1];
    nodes[cycle[1]].startCycle = cycle[0];

    nodes[cycle[cycle.size() - 1]].startCycle = cycle[2];
    nodes[cycle[2]].startCycle = cycle[cycle.size() - 1];
}


vector<Node> prepareAllVertex(vector<int> &start, vector<int> &end) {
    set<int> set_vertex = getSetVertex(start, end);
    vector<Node> nodes(set_vertex.size() * 2);

    int i = 0;
    for (int value: set_vertex) {
        nodes[i] = createNode(value, i);
        nodes[i + 1] = createNode(value, i + 1);
        i = i + 2;
    }

    for (int j = 1; j < start.size(); ++j) {
        int prev_value = numberToIndex(start[j - 1], true);
        int cur_value = numberToIndex(start[j], false);
        nodes[prev_value].startCycle = cur_value;
        nodes[cur_value].startCycle = prev_value;
    }

    int prev = numberToIndex(start[start.size() - 1], true);
    int cur = numberToIndex(start[0], false);
    nodes[prev].startCycle = cur;
    nodes[cur].startCycle = prev;


    for (int j = 1; j < end.size(); ++j) {
        int prev_value = numberToIndex(end[j - 1], true);
        int cur_value = numberToIndex(end[j], false);
        nodes[prev_value].endCycle = cur_value;
        nodes[cur_value].endCycle = prev_value;
    }

    prev = numberToIndex(end[end.size() - 1], true);
    cur = numberToIndex(end[0], false);
    nodes[prev].endCycle = cur;
    nodes[cur].endCycle = prev;

    return nodes;

}


void printCurrentPosition(vector<Node> &nodes) {
    string answer = "";
    vector<bool> visited(nodes.size(), false);
    for (int i = 0; i < visited.size(); ++i) {
        if (!visited[i]) {
            int curValue = i;
            string cur_cycle = "";
            while (!visited[curValue]) {
                visited[curValue] = true;
                if (!cur_cycle.empty()) {
                    cur_cycle += " ";
                }
                if (curValue % 2 == 0) {
                    cur_cycle += "+";
                    cur_cycle += to_string((curValue / 2) + 1);
                    curValue++;
                    visited[curValue] = true;
                } else {
                    cur_cycle += "-";
                    cur_cycle += to_string((curValue / 2) + 1);
                    curValue--;
                    visited[curValue] = true;
                }
                curValue = nodes[curValue].startCycle;
            }
            answer += '(' + cur_cycle + ')';
        }
    }
    cout << answer << endl;
}

int main() {
    string line;
    std::getline(in, line);
    vector<int> start = getCycle(line);
    std::getline(in, line);
    vector<int> end = getCycle(line);

    vector<Node> nodes = prepareAllVertex(start, end);

    printCurrentPosition(nodes);
    while (true) {
        vector<vector<int>> cycles = findCycles(nodes);
        if (cycles.size() == 0) {
            break;
        }
        for (int i = 0; i < cycles.size(); ++i) {
            twoBreakGraph(nodes, cycles[i]);
            printCurrentPosition(nodes);
        }
    }

}
