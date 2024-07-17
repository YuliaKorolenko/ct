#include <iostream>
#include <vector>
#include <map>
#include <cmath>
#include <set>

#define double long double

using namespace std;


enum Type {
    var = 0,
    tnh = 1,
    sgm = 2,
    mul = 3,
    sum = 4,
    had = 5
};

std::map<std::string, Type> typeMap = {{"var", Type::var},
                                       {"tnh", Type::tnh},
                                       {"sgm", Type::sgm},
                                       {"mul", Type::mul},
                                       {"sum", Type::sum},
                                       {"had", Type::had}};


std::map<string, int> varName = {{"w_f", 0},
                                 {"u_f", 1},
                                 {"b_f", 2},
                                 {"w_i", 3},
                                 {"u_i", 4},
                                 {"b_i", 5},
                                 {"w_o", 6},
                                 {"u_o", 7},
                                 {"b_o", 8},
                                 {"w_c", 9},
                                 {"u_c", 10},
                                 {"b_c", 11}};


struct Node {
    int index;
    string name;
    Type type;
    vector<int> previous;
    set<int> next;
    vector<vector<double>> matrix;
    vector<vector<double>> backward;

    Node() {}

    Node(int i, Type type, vector<vector<double>> &matrix, string name) : index(i), type(type), matrix(matrix),
                                                                          name(name) {}

    Node(int i, Type type) : index(i), type(type) {}

    Node(int i, Type type, vector<int> prev) : index(i), type(type), previous(prev) {}

    Node(int i, Type type, vector<int> prev, string name) : index(i), type(type), previous(prev), name(name) {}
};


vector<vector<double>> multiply_el(vector<vector<double>> &a, vector<vector<double>> &b) {
    vector<vector<double>> result(a.size(), vector<double>(a[0].size()));
    for (int i = 0; i < result.size(); ++i) {
        for (int j = 0; j < result[0].size(); ++j) {
            result[i][j] = a[i][j] * b[i][j];
        }
    }
    return result;
}


void countVar(int index, vector<Node> &graph) {
    return;
}

Node getVar(int index, int n, int m, string name) {
    vector<vector<double>> matrix_v(n, vector<double>(m, 0));
    for (int i = 0; i < n; ++i) {
        for (int j = 0; j < m; ++j) {
            cin >> matrix_v[i][j];
        }
    }
    return Node(index, Type::var, matrix_v, name);
}

vector<vector<double>> countSum(vector<Node> &graph, vector<int> &summards) {
    int size_n = graph[summards[0]].matrix.size();
    int size_m = graph[summards[0]].matrix[0].size();
    vector<vector<double>> matrix_sum(size_n, vector<double>(size_m, 0));
    for (int j = 0; j < size_n; ++j) {
        for (int k = 0; k < size_m; ++k) {
            for (int l = 0; l < summards.size(); ++l) {
                matrix_sum[j][k] += graph[summards[l]].matrix[j][k];
            }
        }
    }
    return matrix_sum;
}

Node getSum(int index, vector<Node> &graph, vector<int> summards, string name) {
    for (int j = 0; j < summards.size(); ++j) {
        graph[summards[j]].next.insert(index);
    }
    return Node(index, Type::sum, summards, name);
}

vector<vector<double>> backwardSum(int index, int next_index, vector<Node> &graph) {
    int count = 0;
    for (int i = 0; i < graph[next_index].previous.size(); ++i) {
        if (graph[next_index].previous[i] == index) {
            count++;
        }
    }
    vector<vector<double>> matrix_rlu(graph[index].matrix.size(), vector<double>(graph[index].matrix[0].size(), count));
    return multiply_el(graph[next_index].backward, matrix_rlu);
}


vector<vector<double>> countSgm(vector<Node> &graph, int x) {
    int size_n = graph[x].matrix.size();
    int size_m = graph[x].matrix[0].size();
    vector<vector<double>> matrix_sgm(size_n, vector<double>(size_m, 0));
    for (int j = 0; j < size_n; ++j) {
        for (int k = 0; k < size_m; ++k) {
            matrix_sgm[j][k] = 1.0 / (1 + exp(-graph[x].matrix[j][k]));
        }
    }
    return matrix_sgm;
}


Node getSgm(int index, vector<Node> &graph, int x, string name) {
    graph[x].next.insert(index);

    return Node(index, Type::sgm, {x}, name);
}

vector<vector<double>> backwardSgm(int index, int next_index, vector<Node> &graph) {
    Node current = graph[index];
    Node next = graph[next_index];
    vector<vector<double>> next_val = next.matrix;

    int size_n = current.matrix.size();
    int size_m = current.matrix[0].size();
    vector<vector<double>> back_sgm(size_n, vector<double>(size_m, 0));
    for (int i = 0; i < size_n; ++i) {
        for (int j = 0; j < size_m; ++j) {
            // it is better to use current matrix
            back_sgm[i][j] = next.matrix[i][j] * (1 - next.matrix[i][j]);
        }
    }

    return multiply_el(next.backward, back_sgm);
}

vector<vector<double>> mulMatrix(vector<vector<double>> a, vector<vector<double>> b) {
    vector<vector<double>> matrix_mul(a.size(), vector<double>(b[0].size(), 0));
    for (int first_ = 0; first_ < a.size(); ++first_) {
        for (int second_ = 0; second_ < b[0].size(); ++second_) {
            for (int i = 0; i < a[0].size(); ++i) {
                matrix_mul[first_][second_] += a[first_][i] * b[i][second_];
            }
        }
    }
    return matrix_mul;
}

vector<vector<double>> countMul(vector<Node> &graph, int a, int b) {
    int size_a_n = graph[a].matrix.size();
    int size_b_m = graph[b].matrix[0].size();

    vector<vector<double>> matrix_mul(size_a_n, vector<double>(size_b_m, 0));
    matrix_mul = mulMatrix(graph[a].matrix, graph[b].matrix);
    return matrix_mul;
}

Node getMul(int index, vector<Node> &graph, int a, int b) {
    graph[a].next.insert(index);
    graph[b].next.insert(index);
    return Node(index, Type::mul, {a, b});
}

vector<vector<double>> transponate(vector<vector<double>> &a) {
    vector<vector<double>> transponated(a[0].size(), vector<double>(a.size()));
    for (int i = 0; i < a.size(); ++i) {
        for (int j = 0; j < a[0].size(); ++j) {
            transponated[j][i] = a[i][j];
        }
    }
    return transponated;
}

vector<vector<double>> sumMatrix(vector<vector<double>> a, vector<vector<double>> b) {
    vector<vector<double>> result(a.size(), vector<double>(a[0].size(), 0));
    for (int i = 0; i < a.size(); ++i) {
        for (int j = 0; j < a[0].size(); ++j) {
            result[i][j] = a[i][j] + b[i][j];
        }
    }
    return result;
}


vector<vector<double>> backwardMul(int index, int next_index, vector<Node> &graph) {
    Node next = graph[next_index];
    vector<vector<double>> dnext_dcur;
    int is_left = 1;
    if (next.previous[0] != index) {
        dnext_dcur = graph[next.previous[0]].matrix;
        is_left = 0;
    } else if (next.previous[1] != index) {
        dnext_dcur = graph[next.previous[1]].matrix;
    } else {
        is_left = 2;
        dnext_dcur = graph[next.previous[1]].matrix;
    }
    vector<vector<double>> tran_dnext_dcur = transponate(dnext_dcur);

    if (is_left == 1) {
        return mulMatrix(next.backward, tran_dnext_dcur);
    } else if (is_left == 0) {
        return mulMatrix(tran_dnext_dcur, next.backward);
    } else {
        return sumMatrix(mulMatrix(next.backward, tran_dnext_dcur), mulMatrix(tran_dnext_dcur, next.backward));
    }
}

vector<vector<double>> countTnh(vector<Node> &graph, int x) {
    int size_n = graph[x].matrix.size();
    int size_m = graph[x].matrix[0].size();

    vector<vector<double>> matrix_tnh(size_n, vector<double>(size_m, 0));
    for (int i = 0; i < size_n; ++i) {
        for (int j = 0; j < size_m; ++j) {
            matrix_tnh[i][j] = tanh(graph[x].matrix[i][j]);
        }
    }
    return matrix_tnh;
}

Node getTnh(int index, vector<Node> &graph, int x) {
    graph[x].next.insert(index);
    return Node(index, Type::tnh, {x});
}


vector<vector<double>> backwardTnh(int index, int next_index, vector<Node> &graph) {
    Node next = graph[next_index];
    Node current = graph[index];
    vector<vector<double>> backward_tnh(current.matrix.size(), vector<double>(current.matrix[0].size(), 0));
    for (int i = 0; i < next.matrix.size(); ++i) {
        for (int j = 0; j < next.matrix[0].size(); ++j) {
            double tmp = cosh(current.matrix[i][j]);
            backward_tnh[i][j] = 1.0 / (tmp * tmp);
        }
    }
    return multiply_el(next.backward, backward_tnh);
}


vector<vector<double>> countHad(vector<Node> &graph, vector<int> &prev) {
    int size_n = graph[prev[0]].matrix.size();
    int size_m = graph[prev[0]].matrix[0].size();
    vector<vector<double>> matrix_had(size_n, vector<double>(size_m, 1));
    for (int i = 0; i < size_n; ++i) {
        for (int j = 0; j < size_m; ++j) {
            for (int k = 0; k < prev.size(); ++k) {
                matrix_had[i][j] *= graph[prev[k]].matrix[i][j];
            }
        }
    }
    return matrix_had;
}

Node getHad(int index, vector<Node> &graph, vector<int> prev, string name) {
    for (int i = 0; i < prev.size(); ++i) {
        graph[prev[i]].next.insert(index);
    }

    return Node(index, Type::had, prev, name);
}

vector<vector<double>> backwardHad(int index, int next_index, vector<Node> &graph) {
    Node next = graph[next_index];
    Node current = graph[index];

    int count = 0;
    for (int i = 0; i < next.previous.size(); ++i) {
        if (next.previous[i] == index) {
            count++;
        }
    }

    vector<vector<double>> back_had(current.matrix.size(), vector<double>(current.matrix[0].size(), count));


    for (int i = 0; i < next.previous.size(); ++i) {
        if (next.previous[i] != index) {
            back_had = multiply_el(graph[next.previous[i]].matrix, back_had);
        }
        if (next.previous[i] == index && count > 1) {
            back_had = multiply_el(graph[next.previous[i]].matrix, back_had);
            count--;
        }
    }
    return multiply_el(next.backward, back_had);
}

void insertBackward(int index, vector<Node> &graph) {
    int size_n = graph[index].matrix.size();
    int size_m = graph[index].matrix[0].size();
    vector<vector<double>> matrix_back(size_n, vector<double>(size_m));
    for (int i = 0; i < size_n; ++i) {
        for (int j = 0; j < size_m; ++j) {
            cin >> matrix_back[i][j];
        }
    }
    graph[index].backward = matrix_back;
}

void add_to_back(int index, vector<vector<double>> back_2, vector<Node> &graph) {
    for (int i = 0; i < graph[index].backward.size(); ++i) {
        for (int j = 0; j < graph[index].backward[0].size(); ++j) {
            graph[index].backward[i][j] += back_2[i][j];
        }
    }
}

int main() {

    cout.precision(30);
    int n, m;
    cin >> n;

    vector<Node> graph;


    for (int i = 0; i < varName.size(); ++i) {
        if ((i + 1) % 3 == 0) {
            graph.push_back(getVar(i, n, 1, "b_"));
        } else {
            graph.push_back(getVar(i, n, n, "w|c"));
        }
    }

    cin >> m;

    int startX = 11 + 1;

    Node h0 = getVar(startX + m, n, 1, "h0"); // h0
    Node c0 = getVar(startX + m + 1, n, 1, "c0"); //c0

    for (int i = startX; i < startX + m; ++i) {
        graph.push_back(getVar(i, n, 1, "x_"));
    }

    graph.push_back(c0);
    graph.push_back(h0);

    vector<vector<double>> hm = getVar(0, n, 1, "hm").matrix;
    vector<vector<double>> cm = getVar(0, n, 1, "cm").matrix;

    int start = startX + m + 2;
    int counter = start;

    for (int i = 0; i < m; ++i) {
        // f_t
        graph.push_back(getMul(counter + 0, graph, varName["w_f"], startX + i)); // w_f * x_t
        graph.push_back(getMul(counter + 1, graph, varName["u_f"], counter - 1)); // u_f * h_t-1
        graph.push_back(
                getSum(counter + 2, graph, {counter, counter + 1, varName["b_f"]},
                       "w_f * x_t + u_f * h_t-1 + b_f")); // w_f * x_t + u_f * h_t-1 + b_f
        graph.push_back(getSgm(counter + 3, graph, counter + 2, "f"));   // sigma(.....) f_t

        // i_t
        graph.push_back(getMul(counter + 4, graph, varName["w_i"], startX + i)); // w_i * x_t
        graph.push_back(getMul(counter + 5, graph, varName["u_i"], counter - 1)); // u_i * h_t-1
        graph.push_back(getSum(counter + 6, graph,
                               {counter + 4, counter + 5, varName["b_i"]},
                               "w_i * x_t + u_i * h_t-1 + b_i")); // w_i * x_t + u_i * h_t-1 + b_i
        graph.push_back(getSgm(counter + 7, graph, counter + 6, "i"));   // sigma(.....)

        // ot
        graph.push_back(getMul(counter + 8, graph, varName["w_o"], startX + i)); // w_o * x_t
        graph.push_back(getMul(counter + 9, graph, varName["u_o"], counter - 1)); // u_o * h_t-1
        graph.push_back(getSum(counter + 10, graph,
                               {counter + 8, counter + 9, varName["b_o"]},
                               "w_o * x_t + u_o * h_t-1 + b_o")); // w_o * x_t + u_o * h_t-1 + b_i
        graph.push_back(getSgm(counter + 11, graph, counter + 10, "o"));   // sigma(.....)  o_t


        // ct
        graph.push_back(getMul(counter + 12, graph, varName["w_c"], startX + i)); // w_o * x_t
        graph.push_back(getMul(counter + 13, graph, varName["u_c"], counter - 1)); // u_o * h_t-1
        graph.push_back(getSum(counter + 14, graph,
                               {counter + 12, counter + 13, varName["b_c"]},
                               "w_c * x_t + u_c * h_t-1 + b_c")); // w_c * x_t + u_c * h_t-1 + b_c
        graph.push_back(getTnh(counter + 15, graph, counter + 14));   // tnh(.....)
        graph.push_back(getHad(counter + 16, graph, {counter + 7, counter + 15}, "i_t ** tnh(.....)"));   // i_t * prev

        graph.push_back(getHad(counter + 17, graph, {(counter + 3), counter - 2}, "f_t∘c_t−1"));   // f_t∘c_t−1
        graph.push_back(getSum(counter + 18, graph, {counter + 16, counter + 17}, "c"));   // c_t

        // ht
        graph.push_back(getHad(counter + 19, graph, {(counter + 11), counter + 18}, "h")); // o_t * c_t
        counter += 20;
    }

    graph[counter - 1].backward = hm;
    graph[counter - 2].backward = cm;


    for (int i = 0; i < graph.size(); ++i) {
        switch (graph[i].type) {
            case Type::var: {
                countVar(i, graph);
                break;
            }
            case Type::sum: {
                graph[i].matrix = countSum(graph, graph[i].previous);
                break;
            }
            case Type::tnh:
                graph[i].matrix = countTnh(graph, graph[i].previous[0]);
                break;
            case Type::sgm:
                graph[i].matrix = countSgm(graph, graph[i].previous[0]);
                break;
            case Type::mul:
                graph[i].matrix = countMul(graph, graph[i].previous[0], graph[i].previous[1]);
                break;
            case Type::had:
                graph[i].matrix = countHad(graph, graph[i].previous);
                break;
        }
    }

    for (int i = m; i > 0; --i) {
        insertBackward(counter - 20 * (m - i) - 9, graph);
    }

    for (int i = (int) graph.size() - 1; i >= 0; --i) {
        if (graph[i].backward.empty()) {
            graph[i].backward = vector<vector<double>>(graph[i].matrix.size(),
                                                       vector<double>(graph[i].matrix[0].size(), 0));
        }


        for (int value: graph[i].next) {
            switch (graph[value].type) {
                case Type::var: {
                    break;
                }
                case Type::sum: {
                    add_to_back(i, backwardSum(i, value, graph), graph);
                    break;
                }
                case Type::tnh:
                    add_to_back(i, backwardTnh(i, value, graph), graph);
                    break;
                case Type::sgm: {
                    add_to_back(i, backwardSgm(i, value, graph), graph);
                    break;
                }
                case Type::mul: {
                    add_to_back(i, backwardMul(i, value, graph), graph);
                    break;
                }
                case Type::had: {
                    add_to_back(i, backwardHad(i, value, graph), graph);
                    break;
                }
            }
        }
    }
    cerr << "-----------------------" << endl;
    //cerr << sizeof(double) << " " << sizeof(double) << endl;

    for (int i = 1; i <= m; ++i) {
        int index_o = counter - 20 * (m - i) - 9;
        for (int j = 0; j < n; ++j) {
            cout << graph[index_o].matrix[j][0] << " ";
        }
        cout << endl;
    }

    for (int i = 0; i < n; ++i) {
        cout << graph[counter - 1].matrix[i][0] << " ";
    }

    cout << endl;

    for (int i = 0; i < n; ++i) {
        cout << graph[counter - 2].matrix[i][0] << " ";
    }

    cout << endl;

    for (int i = startX + m - 1; i >= startX; --i) {
        for (int j = 0; j < n; ++j) {
            cout << graph[i].backward[j][0] << " ";
        }
        cout << endl;
    }

    for (int i = 0; i < n; ++i) {
        cout << graph[start - 1].backward[i][0] << " ";
    }

    cout << endl;

    for (int i = 0; i < n; ++i) {
        cout << graph[start - 2].backward[i][0] << " ";
    }

    cout << endl;

    for (int i = 0; i < varName.size(); ++i) {
        if ((i + 1) % 3 == 0) {
            for (int j = 0; j < n; ++j) {
                cout << graph[i].backward[j][0] << " ";
            }
            cout << endl;
        } else {
            for (int j = 0; j < n; ++j) {
                for (int k = 0; k < n; ++k) {
                    cout << graph[i].backward[j][k] << " ";
                }
                cout << endl;
            }
        }
    }


}
