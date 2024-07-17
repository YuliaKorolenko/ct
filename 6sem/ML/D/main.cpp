#include <iostream>
#include <vector>
#include <set>
#include <cmath>

using namespace std;

struct Node {
    int index;
    vector<int> previous;
    set<int> next;
    vector<vector<double>> matrix;
    double alpha;
    vector<vector<double>> backward;

    Node() {}

    Node(int i, vector<vector<double>> &matrix) : index(i), matrix(matrix) {}

    Node(int i, vector<int> prev) : index(i), previous(prev) {}

    Node(int i, vector<int> prev, double alpha) : index(i), previous(prev), alpha(alpha) {}
};

struct LSTMNode {
    Node o;
    Node c;
    Node h;
    Node i;
    Node x;
    Node f;
};

struct Factors {
    vector<vector<double>> w_;
    vector<vector<double>> u_;
    vector<vector<double>> b_;
    vector<vector<double>> w_backward;
    vector<vector<double>> u_backward;
    vector<vector<double>> b_backward;

    Factors(int n) {
        w_ = vector<vector<double>>(n, vector<double>(n, 0));
        w_backward = vector<vector<double>>(n, vector<double>(n, 0));
        u_ = vector<vector<double>>(n, vector<double>(n, 0));
        u_backward = vector<vector<double>>(n, vector<double>(n, 0));
        b_ = vector<vector<double>>(n, vector<double>(1, 0));
        b_backward = vector<vector<double>>(n, vector<double>(1, 0));
    }
};

Factors read_factor(int n) {
    Factors factor = Factors(n);
    for (int i = 0; i < n; ++i) {
        for (int j = 0; j < n; ++j) {
            cin >> factor.w_[i][j];
        }
    }

    for (int i = 0; i < n; ++i) {
        for (int j = 0; j < n; ++j) {
            cin >> factor.u_[i][j];
        }
    }

    for (int i = 0; i < n; ++i) {
        cin >> factor.b_[i][0];
    }
    return factor;
}

LSTMNode readFirstLstmNode(int n) {
    LSTMNode lstmNode;

    lstmNode.h.matrix = vector<vector<double>>(n, vector<double>(1, 0));
    lstmNode.c.matrix = vector<vector<double>>(n, vector<double>(1, 0));

    for (int i = 0; i < n; ++i) {
        cin >> lstmNode.h.matrix[i][0];
    }

    for (int i = 0; i < n; ++i) {
        cin >> lstmNode.c.matrix[i][0];
    }

    return lstmNode;
}

vector<vector<double>> readVector(int n, int m) {
    vector<vector<double>> x = vector<vector<double>>(n, vector<double>(m, 0));
    for (int i = 0; i < n; ++i) {
        for (int j = 0; j < m; ++j) {
            cin >> x[i][j];
        }
    }
    return x;
}


void readX(vector<LSTMNode> &graph, int n, int m) {
    for (int i = 1; i <= m; ++i) {
        graph[i].x.matrix = readVector(n, 1);
    }
}

void readO(vector<LSTMNode> &graph, int n, int m) {
    for (int i = 1; i <= m; ++i) {
        graph[i].o.matrix = readVector(n, 1);
    }
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

vector<vector<double>> sumMatrix(vector<vector<double>> a, vector<vector<double>> b, vector<vector<double>> c) {
    vector<vector<double>> matrix_sum(a.size(), vector<double>(a[0].size(), 0));
    for (int j = 0; j < a.size(); ++j) {
        for (int k = 0; k < a[0].size(); ++k) {
            matrix_sum[j][k] = a[j][k] + b[j][k] + c[j][k];
        }
    }
    return matrix_sum;
}

vector<vector<double>> sum2Matrix(vector<vector<double>> a, vector<vector<double>> b) {
    vector<vector<double>> matrix_sum(a.size(), vector<double>(a[0].size(), 0));
    for (int j = 0; j < a.size(); ++j) {
        for (int k = 0; k < a[0].size(); ++k) {
            matrix_sum[j][k] = a[j][k] + b[j][k];
        }
    }
    return matrix_sum;
}


vector<vector<double>> hadMatrix(vector<vector<double>> a, vector<vector<double>> b) {
    vector<vector<double>> matrix_had(a.size(), vector<double>(a[0].size(), 1));
    for (int i = 0; i < a.size(); ++i) {
        for (int j = 0; j < a[0].size(); ++j) {
            matrix_had[i][j] = a[i][j] * b[i][j];
        }
    }
    return matrix_had;
}


vector<vector<double>> backwardHad(vector<vector<double>> a, vector<vector<double>> b) {
    return a;
}


vector<vector<double>> tnhMatrix(vector<vector<double>> a) {
    vector<vector<double>> matrix_tnh(a.size(), vector<double>(a[0].size(), 0));
    for (int i = 0; i < a.size(); ++i) {
        for (int j = 0; j < a[0].size(); ++j) {
            matrix_tnh[i][j] = tanh(a[i][j]);
        }
    }
    return matrix_tnh;
}

vector<vector<double>> sigma(vector<vector<double>> a) {
    vector<vector<double>> matrix_sum(a.size(), vector<double>(a[0].size(), 0));
    for (int j = 0; j < a.size(); ++j) {
        for (int k = 0; k < a[0].size(); ++k) {
            matrix_sum[j][k] = 1.0 / (1 + exp(-a[j][k]));
        }
    }
    return matrix_sum;
}

vector<vector<double>> countVector(vector<LSTMNode> &graph, int index, Factors &factor) {
    return sumMatrix(mulMatrix(factor.w_, graph[index].x.matrix),
                     mulMatrix(factor.u_, graph[index - 1].h.matrix),
                     factor.b_);


}

int main() {
    int n, m;
    cin >> n;

    Factors factor_f = read_factor(n);
    Factors factor_i = read_factor(n);
    Factors factor_o = read_factor(n);
    Factors factor_c = read_factor(n);

    cin >> m;

    vector<LSTMNode> graph(m + 1, LSTMNode());

    graph[0] = readFirstLstmNode(n);

    readX(graph, n, m);

    graph[m].h.backward = readVector(n, 1);
    graph[m].c.backward = readVector(n, 1);

    readO(graph, n, m);

    for (int i = 1; i < graph.size(); ++i) {
        graph[i].f.matrix = sigma(countVector(graph, i, factor_f));
        graph[i].i.matrix = sigma(countVector(graph, i, factor_i));
        graph[i].o.matrix = sigma(countVector(graph, i, factor_o));
        graph[i].c.matrix = sum2Matrix(hadMatrix(graph[i].f.matrix, graph[i - 1].c.matrix),
                                       hadMatrix(tnhMatrix(countVector(graph, i, factor_c)), graph[i].i.matrix));
        graph[i].h.matrix = hadMatrix(graph[i].o.matrix, graph[i].c.matrix);
    }


    for (int i = m; i > 0; --i) {
        graph[i].
    }


//    вывод

    for (int i = 1; i < graph.size(); ++i) {
        for (int j = 0; j < m; ++j) {
            cout << graph[i].o.matrix[j][0] << " ";
        }
    }

    for (int j = 0; j < m; ++j) {
        cout << graph[m].h.matrix[j][0] << " ";
    }

    for (int j = 0; j < m; ++j) {
        cout << graph[m].c.matrix[j][0] << " ";
    }


}

