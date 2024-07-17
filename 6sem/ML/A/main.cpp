#include <iostream>
#include <cmath>
#include <vector>


using namespace std;


vector<int> getBitMask(int c, int n) {
    vector<int> bitMask;
    for (int i = 0; i < n; ++i) {
        bitMask.push_back((c >> i) & 1);
    }
    return bitMask;
}

int main() {
    int n;
    cin >> n;
    vector<int> zeros;
    vector<int> ones;

    for (int i = 0; i < pow(2, n); ++i) {
        int a;
        cin >> a;
        if (a) {
            ones.push_back(i);
        } else {
            zeros.push_back(i);
        }
    }
    bool isZeros = false;
    vector<int> cur;

    if (zeros.size() > ones.size()) {
        cur = ones;
    } else {
        cur = zeros;
        isZeros = true;
    }

    if (cur.size() == 0) {
        cout << 1 << endl;
        cout << 1 << endl;
        for (int i = 0; i < n; ++i) {
            cout << 0 << " ";
        }
        if (isZeros) {
            cout << 1 << endl;
        } else {
            cout << -1 << endl;
        }
        return 0;
    }

    vector<vector<double>> matrix(cur.size(), vector<double>(n));
    vector<double> shift(cur.size(), -n + 0.5);

    for (int i = 0; i < cur.size(); ++i) {
        vector<int> bitMask = getBitMask(cur[i], n);
        for (int j = 0; j < n; ++j) {
            if (bitMask[j]) {
                matrix[i][j] = 1;
            } else {
                matrix[i][j] = -1;
                shift[i] += 1;
            }
        }
    }

    vector<double> vector2;
    double shift2 = 0.5;

    if (!isZeros) {
        vector2 = vector<double>(cur.size(), 1);
        shift2 = -1 * shift2;
    } else {
        vector2 = vector<double>(cur.size(), -1);
    }

    cout << 2 << endl;
    cout << cur.size() << " 1" << endl;
    for (int i = 0; i < matrix.size(); ++i) {
        for (int j = 0; j < matrix[i].size(); ++j) {
            cout << matrix[i][j] << " ";
        }
        cout << shift[i] << endl;
    }

    for (int i = 0; i < cur.size(); ++i) {
        cout << vector2[i] << " ";
    }
    cout << shift2;


}
