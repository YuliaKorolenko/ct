#include <iostream>
#include <fstream>
#include <map>
#include <cctype>
#include <vector>

using namespace std;

//Модель HMM можно описать как: HMM= <S, I, T, E>
//S={S_1,...,S_M} – конечное множество состояний;

//I=(P(q_0=S_i)=\pi_i) – вектор вероятностей, того что система находится в состоянии i в момент времени 0;

//T=|P(q_t=S_j|q_{t-1}=S_i)=a_{ij}| – матрица вероятностей перехода из состояния i в состояние j для \forall t\in[1,T_m];

//E=(E_1,...,E_M), E_i=f(o_t|q_t=S_i)– вектор распределений наблюдаемых случайных величин, соответствующих каждому из состояний,
//заданных как функции плотности или распределения, определенные на O (совокупности наблюдаемых значений для всех состояний).

std::ifstream in("/Users/yulkorolenko/Desktop/BioInf/8Viterbi/viter.txt");

map<char, int> getIntMap() {
    string alphabet;
    map<char, int> alph_map;
    std::getline(in, alphabet);
    int count = 0;
    for (int i = 0; i < alphabet.size(); ++i) {
        if (!isspace(alphabet[i])) {
            alph_map[alphabet[i]] = count;
            count++;
        }
    }
    return alph_map;
}

vector<vector<double>> getTable(int n, int m) {
    string line;
    std::getline(in, line);
    vector<vector<double>> table(n, vector<double>(m));
    int i = 0;
    while (std::getline(in, line) && i < n) {
        int k = 1;
        int j = 0;
        string cur_double = "";
        while (k < line.size()) {
            if (isspace(line[k - 1]) && !isspace(line[k])) {
                while (!isspace(line[k]) && k < line.size()) {
                    cur_double += line[k];
                    k++;
                }
                table[i][j] = ::strtod(cur_double.c_str(), 0);
                j++;
                cur_double = "";
            }
            k++;
        }
        i++;
    }
    return table;
}

map<int, char> getInverseStates(const map<char, int>& states) {
    map<int, char> states_st_int;
    for (auto state: states) {
        states_st_int[state.second] = state.first;
    }
    return states_st_int;
}

int main() {
    string line;
    string emit_string;
    std::getline(in, emit_string);
    std::getline(in, line);
    map<char, int> alphabet = getIntMap();
    std::getline(in, line);
    map<char, int> states = getIntMap();
    std::getline(in, line);
    vector<vector<double>> trans = getTable(states.size(), states.size());

    vector<vector<double>> emis = getTable(states.size(), alphabet.size());

    map<int, char> states_st_int = getInverseStates(states);


    vector<vector<double>> viterbi(states.size(), vector<double>(emit_string.size()));
    vector<vector<double>> save_answ(states.size(), vector<double>(emit_string.size()));

    for (int k = 0; k < states.size(); ++k) {
        viterbi[k][0] = emis[k][alphabet[emit_string[0]]];
    }

    for (int i = 0; i < emit_string.size(); ++i) {
        for (int k = 0; k < states.size(); ++k) {

            double max_V_T = viterbi[0][i] * trans[0][k];
            int l = 0;
            for (int j = 0; j < states.size(); ++j) {
                if (max_V_T < viterbi[j][i] * trans[j][k]) {
                    max_V_T = viterbi[j][i] * trans[j][k];
                    l = j;
                }
            }

            viterbi[k][i + 1] = emis[k][alphabet[emit_string[i + 1]]] * max_V_T;
            save_answ[k][i + 1] = l;
        }
    }

    string answer;
    int l = 0;

    for (int i = 0; i < states.size(); ++i) {
        if (viterbi[l].back() < viterbi[i].back()){
            l = i;
        }
    }

    for (int i = emit_string.size() - 1; i >= 0 ; --i) {
        answer = states_st_int[l] + answer;
        l = save_answ[l][i];
    }

    cout << answer;

}
