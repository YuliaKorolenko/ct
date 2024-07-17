#include <iostream>
#include <vector>
#include <string>
#include <fstream>
#include <map>
#include <random>
#include <chrono>

using namespace std;

std::mt19937 rnd(chrono::steady_clock::now().time_since_epoch().count());


std::ifstream in("/Users/yulkorolenko/Desktop/BioInf/2RandomizedMotifSearch/randMot.txt");

map<char, int> char_to_int = {{'A', 0},
                              {'C', 1},
                              {'G', 2},
                              {'T', 3}};


vector<string> getRandomMotif(int k, vector<string> &DNA) {
    int n = DNA[0].size();
    vector<string> motif;

    for (int i = 0; i < DNA.size(); ++i) {
        int first = rnd() % (n - k + 1);
        motif.push_back(DNA[i].substr(first, k));
    }
    return motif;
}

vector<vector<int>> getCount(vector<string> &motif) {
    int mot_size = motif[0].size();
    vector<vector<int>> count(char_to_int.size(), vector<int>(mot_size, 0));
    for (int i = 0; i < mot_size; ++i) {
        for (int j = 0; j < motif.size(); ++j) {
            count[char_to_int[motif[j][i]]][i] += 1;
        }
    }
    return count;
}

vector<string> getMotifsByProfile(vector<vector<double>> &profile, vector<string> &DNA, int k) {
    vector<string> motifs(DNA.size());
    vector<double> scores(DNA.size(), INT16_MIN);
    for (int i = 0; i < DNA.size(); ++i) {
        for (int j = 0; j < DNA[0].size() - k; ++j) {
            string str = DNA[i].substr(j, k);
            double cur_score = 1;
            for (int l = 0; l < k; ++l) {
                cur_score *= profile[char_to_int[str[l]]][l];
            }
            if (cur_score > scores[i]) {
                motifs[i] = str;
                scores[i] = cur_score;
            }
        }
    }
    return motifs;
}

vector<vector<double>> getProfile(vector<string> &motif) {
    int mot_size = motif[0].size();
    vector<vector<int>> count = getCount(motif);
    vector<vector<double>> profile(count.size(), vector<double>(mot_size));
    for (int i = 0; i < count.size(); ++i) {
        for (int j = 0; j < mot_size; ++j) {
            profile[i][j] = (double) (1 + count[i][j]) / (mot_size + 4);
        }
    }
    return profile;
}


int score(vector<string> &motifs) {
    vector<vector<int>> count = getCount(motifs);
    int score_int = 0;
    for (int j = 0; j < count[0].size(); ++j) {
        int max = 0;
        for (int i = 0; i < count.size(); ++i) {
            if (count[i][j] > max) {
                max = count[i][j];
            }
        }
        score_int += motifs.size() - max;
    }
    return score_int;
}

vector<string> randomMotifSearch(int k, vector<string> &DNA) {
    vector<string> bestMotifs = getRandomMotif(k, DNA);
    while (true) {
        vector<vector<double>> profile = getProfile(bestMotifs);
        vector<string> motifs = getMotifsByProfile(profile, DNA, k);
        if (score(motifs) < score(bestMotifs)) {
            bestMotifs = motifs;
        } else {
            return bestMotifs;
        }
    }
}

int main() {
    int k, n;

    string line;
    std::getline(in, line);
    int line_size = line.size();
    int i = 0;
    string k_str;
    while (!isspace(line[i])  && line_size > i) {
        k_str += line[i];
        i++;
    }
    k = atoi(k_str.c_str());

    while (isspace(line[i])  && line_size > i) {
        i++;
    }

    string n_str;
    while (!isspace(line[i]) && line_size > i) {
        n_str += line[i];
        i++;
    }
    n = atoi(n_str.c_str());

    vector<string> DNA;
    for (int i = 0; i < n; ++i) {
        std::getline(in, line);
        DNA.push_back(line);
    }

    vector<string> BestMotifs;
    for (int j = 0; j < DNA.size(); ++j) {
        BestMotifs.push_back(DNA[j].substr(0, k));
    }

    i = 0;
    while (i < 1000) {
        vector<string> motifs = randomMotifSearch(k, DNA);
        if (score(motifs) < score(BestMotifs)) {
            BestMotifs = motifs;
            cout << score(BestMotifs) << endl;
        }
        i++;

    }

    for (int i = 0; i < BestMotifs.size(); ++i) {
        cout << BestMotifs[i] << endl;
    }
}
