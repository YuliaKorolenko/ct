#include <iostream>
#include <vector>
#include <string>


int findHamming(std::string a, std::string b) {
    int answer = 0;
    for (int i = 0; i < a.size(); ++i) {
        if (a[i] != b[i]) {
            answer++;
        }
    }
    return answer;
}

int patternText(std::string text, std::string pattern, int k) {
    int d_min = INT_MAX;
    for (int i = 0; i < text.size() - k; ++i) {
        d_min = std::min(d_min, findHamming(text.substr(i, k), pattern));
    }
    return d_min;
}

void generated(std::vector<std::string> &generated_patterns, std::string a, int k) {
    if (a.size() == k) {
        generated_patterns.push_back(a);
        return;
    }
    generated(generated_patterns, a + "A", k);
    generated(generated_patterns, a + "T", k);
    generated(generated_patterns, a + "G", k);
    generated(generated_patterns, a + "C", k);
}

int main() {
    int k;
    std::string k_str;
    getline(std::cin, k_str);
    k = stoi(k_str);
    std::vector<std::string> dna;
    std::vector<std::string> generated_patterns(0, "");
    generated(generated_patterns, "", k);
    std::string text;
    while (getline(std::cin, text)) {
        dna.push_back(text);
    }
    std::string pattern_min = "";
    int min = INT_MAX;
    for (int i = 0; i < generated_patterns.size(); ++i) {
        int current_min = 0;
        for (int j = 0; j < dna.size(); ++j) {
            current_min += patternText(dna[j], generated_patterns[i], k);
        }
        if (current_min < min) {
            min = current_min;
            pattern_min = generated_patterns[i];
        }
    }
    std::cout << pattern_min;
}