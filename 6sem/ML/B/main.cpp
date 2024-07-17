#include <iostream>
#include <vector>
#include <map>
#include <cmath>
#include <set>
#include <iomanip>

#define double long double

using namespace std;


enum Type {
    var = 0,
    tnh = 1,
    rlu = 2,
    mul = 3,
    sum = 4,
    had = 5
};

std::map<std::string, Type> typeMap = {{"var", Type::var},
                                       {"tnh", Type::tnh},
                                       {"rlu", Type::rlu},
                                       {"mul", Type::mul},
                                       {"sum", Type::sum},
                                       {"had", Type::had}};


struct Node {
    int index;
    Type type;
    vector<int> previous;
    set<int> next;
    vector<vector<double>> matrix;
    double alpha;
    vector<vector<double>> backward;

    Node() {}

    Node(int i, Type type, vector<vector<double>> &matrix) : index(i), type(type), matrix(matrix) {}

    Node(int i, Type type) : index(i), type(type) {}

    Node(int i, Type type, vector<int> prev) : index(i), type(type), previous(prev) {}

    Node(int i, Type type, vector<int> prev, double alpha) : index(i), type(type), previous(prev), alpha(alpha) {}
};


vector<vector<double>> multiply_el(vector<vector<double>> &a, vector<vector<double>> &b) {
    vector<vector<double>> result(a.size(), vector<double>(a[0].size()));
    for (int i = 0; i < result.size(); ++i) {
        for (int j = 0; j < result[0].size(); ++j) {
            result[i][j] = a[i][j] * b[i][j];
        }
    }
    return result;

    switch (a) {
        case 1:
            a++;
        case 2:
            a++;
        case 3:
            a++;
    }
}


void countVar(int index, vector<Node> &graph) {
    for (int j = 0; j < graph[index].matrix.size(); ++j) {
        for (int l = 0; l < graph[index].matrix[0].size(); ++l) {
            cin >> graph[index].matrix[j][l];
        }
    }
}

Node getVar(int index) {
    int a, b;
    cin >> a >> b;
    vector<vector<double>> matrix_v(a, vector<double>(b, 0));
    return Node(index, Type::var, matrix_v);
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

Node getSum(int index, vector<Node> &graph) {
    int count;
    cin >> count;
    vector<int> summards;
    for (int j = 0; j < count; ++j) {
        int a;
        cin >> a;
        a--;
        graph[a].next.insert(index);
        summards.push_back(a);
    }
    return Node(index, Type::sum, summards);
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


vector<vector<double>> countRlu(vector<Node> &graph, int x, double alpha) {
    int size_n = graph[x].matrix.size();
    int size_m = graph[x].matrix[0].size();
    vector<vector<double>> matrix_rlu(size_n, vector<double>(size_m, 0));
    for (int j = 0; j < size_n; ++j) {
        for (int k = 0; k < size_m; ++k) {
            if (graph[x].matrix[j][k] < 0) {
                matrix_rlu[j][k] = graph[x].matrix[j][k] / alpha;
            } else {
                matrix_rlu[j][k] = graph[x].matrix[j][k];
            }
        }
    }
    return matrix_rlu;
}


Node getRlu(int index, vector<Node> &graph) {
    int x;
    double alpha;
    cin >> alpha >> x;
    x--;
    graph[x].next.insert(index);

    return Node(index, Type::rlu, {x}, alpha);
}

vector<vector<double>> backwardRlu(int index, int next_index, vector<Node> &graph) {
    Node current = graph[index];
    Node next = graph[next_index];
    vector<vector<double>> next_val = next.matrix;

    int size_n = current.matrix.size();
    int size_m = current.matrix[0].size();
    vector<vector<double>> back_rlu(size_n, vector<double>(size_m, 0));
    for (int i = 0; i < size_n; ++i) {
        for (int j = 0; j < size_m; ++j) {
            if (next_val[i][j] < 0) { // it is better to use current matrix
                back_rlu[i][j] = 1.0 / next.alpha;
            } else {
                back_rlu[i][j] = 1;
            }
        }
    }

    return multiply_el(next.backward, back_rlu);
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

Node getMul(int index, vector<Node> &graph) {
    int a, b;
    cin >> a >> b;
    a--;
    b--;
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

Node getTnh(int index, vector<Node> &graph) {
    int x;
    cin >> x;
    x--;
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

Node getHad(int index, vector<Node> &graph) {
    int count;
    cin >> count;
    vector<int> prev;
    for (int i = 0; i < count; ++i) {
        int a;
        cin >> a;
        a--;
        graph[a].next.insert(index);
        prev.push_back(a);
    }

    return Node(index, Type::had, prev);
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
    int n, m, k;
    cin >> n >> m >> k;

    vector<Node> graph;
    for (int i = 0; i < n; ++i) {
        string cur_type;
        cin >> cur_type;
        switch (typeMap[cur_type]) {
            case Type::var: {
                graph.push_back(getVar(i));
                break;
            }
            case Type::sum: {
                graph.push_back(getSum(i, graph));
                break;
            }
            case Type::tnh:
                graph.push_back(getTnh(i, graph));
                break;
            case Type::rlu:
                graph.push_back(getRlu(i, graph));
                break;
            case Type::mul:
                graph.push_back(getMul(i, graph));
                break;
            case Type::had:
                graph.push_back(getHad(i, graph));
                break;
        }

    }

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
            case Type::rlu:
                graph[i].matrix = countRlu(graph, graph[i].previous[0], graph[i].alpha);
                break;
            case Type::mul:
                graph[i].matrix = countMul(graph, graph[i].previous[0], graph[i].previous[1]);
                break;
            case Type::had:
                graph[i].matrix = countHad(graph, graph[i].previous);
                break;
        }
    }

    for (int i = graph.size() - k; i < graph.size(); ++i) {
        insertBackward(i, graph);
    }

    for (int i = (int) graph.size() - 1 - k; i >= 0; --i) {
        graph[i].backward = vector<vector<double>>(graph[i].matrix.size(),
                                                   vector<double>(graph[i].matrix[0].size(), 0));
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
                case Type::rlu: {
                    add_to_back(i, backwardRlu(i, value, graph), graph);
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
//    cerr << "-----------------------" << endl;
    //cerr << sizeof(double) << " " << sizeof(double) << endl;

    for (int i = graph.size() - k; i < graph.size(); ++i) {
        for (int j = 0; j < graph[i].matrix.size(); ++j) {
            for (int l = 0; l < graph[i].matrix[0].size(); ++l) {
                cout << graph[i].matrix[j][l] << " ";
            }
            cout << endl;
        }
    }


    for (int i = 0; i < m; ++i) {
        for (int j = 0; j < graph[i].backward.size(); ++j) {
            for (int l = 0; l < graph[i].backward[0].size(); ++l) {
                cout << graph[i].backward[j][l] << " ";
            }
            cout << endl;
        }
    }


}
