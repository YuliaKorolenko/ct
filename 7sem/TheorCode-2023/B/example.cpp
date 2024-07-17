#include <iostream>
#include <vector>
#include <fstream>
#include <unordered_set>
#include <random>

std::ifstream in("input.txt");
std::ofstream out("output.txt");
int n, pi, d, k;
std::vector<int> g, logarithm, antiLogarithm;

void clearLeadZeros(std::vector<int> &a) {
    while (!a.empty() && a[a.size() - 1] == 0) {
        a.pop_back();
    }
}

int getHighestBit(int a) {
    int i;
    for (i = 0; a > 0; i++, a >>= 1);
    return i;
}

int mod(int a) {
    while (true) {
        int c = getHighestBit(a) - getHighestBit(pi);
        if (c < 0) {
            return a;
        }
        a ^= pi << c;
    }
}

int times(int a, int b) {
    if (a == 0 || b == 0) {
        return 0;
    }
    return antiLogarithm[(logarithm[a] + logarithm[b]) % n];
}

std::vector<int> times(std::vector<int> const &a, std::vector<int> const &b) {
    std::vector<int> c(a.size() + b.size() - 1, 0);
    for (int i = 0; i < a.size(); i++) {
        for (int j = 0; j < b.size(); j++) {
            c[i + j] ^= times(a[i], b[j]);
        }
    }
    return c;
}

void init() {
    logarithm = std::vector<int>(n + 1);
    antiLogarithm = std::vector<int>(n);

    for (int i = 0, a = 1; i < n; i++, a = mod(a << 1)) {
        std::cerr << i << "!" << a << std::endl;
        logarithm[a] = i;
        antiLogarithm[i] = a;
    }
    std::cout << " here ";

    for (int i = 0; i < logarithm.size(); i++){
        std::cout << logarithm[i] << " ";
    }

    std::cout << std::endl;

    for (int i = 0; i < antiLogarithm.size(); i++){
        std::cout << antiLogarithm[i] << " ";
    }

    std::unordered_set<int> used;
    g = std::vector<int>(1, 1);
    for (int i = 1, a = 2; i < d; i++, a = mod(a << 1)) {
        for (int b = a; used.find(b) == used.end(); b = times(b, b)) {
            used.insert(b);
            std::vector c(2, 1);
            c[0] = b;
            g = times(g, c);
        }
    }

    clearLeadZeros(g);
    k = n + 1 - g.size();
}

void print(std::vector<int> const &a) {
    for (int i: a) {
        out << i << " ";
    }
    out << std::endl;
}

std::vector<int> encode(std::vector<int> const &c) {
    std::vector<int> y(n, 0);
    for (int i = 0; i < c.size(); i++) {
        y[i + n - k] = c[i];
    }
    std::vector<int> y2 = y;
    clearLeadZeros(y2);
    while (y2.size() >= g.size()) {
        for (int i = 0; i < g.size(); i++) {
            y2[i + y2.size() - g.size()] ^= g[i];
        }
        clearLeadZeros(y2);
    }
    for (int i = 0; i < n - k; i++) {
        y[i] = y2[i];
    }
    return y;
}

int eval(std::vector<int> const &a, int x) {
    int y = 0;
    for (int i = 0; i < a.size(); i++) {
        y ^= times(a[i], antiLogarithm[x * i % n]);
    }
    return y;
}

std::vector<int> decode(std::vector<int> y) {
    std::vector<int> S(d - 1);
    for (int i = 1; i < d; i++) {
        S[i - 1] = eval(y, i);
    }
    std::vector<int> Lambda(1, 1), B(1, 1);
    int m = 0, l = 0;
    for (int r = 1; r < d; r++) {
        int delta = 0;
        for (int i = 0; i < Lambda.size() && i <= l; i++) {
            delta ^= times(Lambda[i], S[r - i - 1]);
        }
        if (delta != 0) {
            std::vector<int> T(std::max(r - m + B.size(), Lambda.size()), 0);
            for (int i = 0; i < B.size(); i++) {
                T[i + r - m] = times(B[i], delta);
            }
            for (int i = 0; i < Lambda.size(); i++) {
                T[i] ^= Lambda[i];
            }
            clearLeadZeros(T);
            if (2 * l < r) {
                B = Lambda;
                delta = delta == 1 ? 1 : antiLogarithm[n - logarithm[delta]];
                for (int &i: B) {
                    i = times(i, delta);
                }
                l = r - l;
                m = r;
            }
            Lambda = T;
        }
    }
    if (eval(Lambda, 0) == 0) {
        y[0] ^= 1;
    }
    for (int i = 1; i < n; i++) {
        if (eval(Lambda, i) == 0) {
            y[n - i] ^= 1;
        }
    }
    return y;
}

void simulate(double p, int maxIterations, int maxErrors) {
    int e = 0;
    std::mt19937 gen(3);
    for (int i = 0; i < maxIterations; i++) {
        std::vector<int> c(k);
        for (int j = 0; j < k; j++) {
            c[j] = gen() % 2;
        }
        std::vector<int> y = encode(c);
        std::vector<int> badY = y;
        for (int j = 0; j < n; j++) {
            if (p > (double) gen() / std::mt19937::max()) {
                badY[j] ^= 1;
            }
        }
        badY = decode(badY);
        for (int j = 0; j < n; j++) {
            if (y[j] != badY[j]) {
                e++;
                break;
            }
        }
        if (e == maxErrors) {
            out << (double) e / (i + 1) << std::endl;
            return;
        }
    }
    out << (double) e / maxIterations << std::endl;
}

int main() {
    in >> n >> pi >> d;
    init();
    out << k << std::endl;
    print(g);
    
    std::string command;
    while (in >> command) {
        if (command == "Encode") {
            std::vector<int> c(k);
            for (int i = 0; i < k; i++) {
                in >> c[i];
            }
            print(encode(c));
        }
        if (command == "Decode") {
            std::vector<int> y(n);
            for (int i = 0; i < n; i++) {
                in >> y[i];
            }
            print(decode(y));
        }
        if (command == "Simulate") {
            double p;
            int maxIterations, maxErrors;
            in >> p >> maxIterations >> maxErrors;
            simulate(p, maxIterations, maxErrors);
        }
    }
    return 0;
}
