#include <iostream>
#include <map>
#include <set>
#include <stack>
#include <vector>
#include <fstream>

using namespace std;


std::pair<string, string> get_first_second(string s) {
    for (int i = 0; i < s.size(); ++i) {
        if (s[i] == '|') {
            return {s.substr(0, i), s.substr(i + 1, s.size())};
        }
    }
    return {"-", "-"};
}

std::pair<std::pair<string, string>, std::pair<string, string>> split(string s) {
    std::pair<string, string> f_s = get_first_second(s);
    string first = f_s.first;
    string second = f_s.second;
    return {{first.substr(0, first.size() - 1), second.substr(0, second.size() - 1)},
            {first.substr(1, first.size()),     second.substr(1, second.size())}};
}


string createAnswer(vector<pair<string, string>> result, int d, int k) {
    string s1 = "";
    string s2 = "";
    string news1 = "";
    string news2 = "";
    string a = result[result.size() - 1].first;
    int j = 1;
    for (int i = result.size() - 1; i >= 0; i--) {
        if (i % k == 0) {
            s1 += result[i].first;
            s2 += result[i].second;
        }
        if ((result.size() - 1 - i) % k == 0) {
            news1 += result[i].first;
            news2 += result[i].second;
            j = i;
        }
    }

    news1 += result[0].first.substr(result[0].first.size() - j, result[0].first.size());
    news2 += result[0].second.substr(result[0].second.size() - j, result[0].second.size());

    return news1 + news2.substr(news1.size() - (d + k + 1), news2.size());
}


int main() {
    map<std::pair<string, string>, vector<std::pair<string, string>>> graph;
    set<std::pair<string, string>> is_once;


    int k, d;

    std::ifstream in("/Users/yulkorolenko/Desktop/BioInf/3ReconstructionProblem/recprob.txt");
    string line;
    std::getline(in, line);
    int i = 0;
    int line_size = line.size();
    string k_str;
    while (!isspace(line[i])  && line_size > i) {
        k_str += line[i];
        i++;
    }
    k = atoi(k_str.c_str());

    while (isspace(line[i])  && line_size > i) {
        i++;
    }

    string d_str;
    while (!isspace(line[i]) && line_size > i) {
        d_str += line[i];
        i++;
    }
    d = atoi(d_str.c_str());


    if (in.is_open()) {
        while (std::getline(in, line)) {
            string s = line;
            std::pair<std::pair<string, string>, std::pair<string, string>> splited = split(s);
            std::pair<string, string> prefix = splited.first;
            std::pair<string, string> sufix = splited.second;
            if (is_once.find(prefix) == is_once.end()) {
                is_once.insert(prefix);
            } else {
                is_once.erase(prefix);
            }

            if (is_once.find(sufix) == is_once.end()) {
                is_once.insert(sufix);
            } else {
                is_once.erase(sufix);
            }

            if (graph.find(prefix) == graph.end()) {
                graph[prefix] = {sufix};
            } else {
                graph[prefix].push_back(sufix);
            }
        }
    }


    if (is_once.size() != 2) {
        std::cout << "ERROR. More than 2";
    }

    std::pair<string, string> first;
    for (auto v: is_once) {
        if (graph.find(v) != graph.end()) {
            first = v;
            break;
        }
    }

    for (auto v: graph) {
        std::sort(v.second.begin(), v.second.end(), std::greater<>());
    }

    stack<std::pair<string, string>> st;
    st.push(first);
    vector<pair<string, string>> cycle;
    vector<pair<string, string>> res;
    while (!st.empty()) {
        std::pair<string, string> v = st.top();
        if (graph.find(v) == graph.end()) {
            res.push_back(v);
            st.pop();
        } else {
            pair<string, string> i = graph.find(v)->second.back();
            st.push(i);
            graph.find(v)->second.pop_back();
            if (graph.find(v)->second.empty()) {
                graph.erase(v);
            }
        }
    }

    cout << createAnswer(res, d, k - 1);

}
