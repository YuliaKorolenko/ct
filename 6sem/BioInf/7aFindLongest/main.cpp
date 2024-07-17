#include <iostream>
#include <vector>

std::string prefix_func(std::string s){
    int n = s.size();
    std::string answer = "";
    std::vector<int> pref_func(n + 1, 0);
    for (int i = 1; i < n ; ++i) {
        int j = pref_func[i - 1];
        while (j > 0 && s[i] != s[j])
            j = pref_func[j - 1];
        if (s[i] == s[j]) {
            if (j + 1 != n && j == answer.size()) {
                answer+=s[j];
            }
            j++;
        }
        pref_func[i] = j;
    }
    return answer;
}

int main() {
    std::string s;
    std::cin >> s;
    int n = s.size();
    int max = 0;
    std::string max_string = "";
    for (int k = n; k > max; --k) {
        std::string answer = prefix_func(s);
        if (max_string.size() < answer.size()) {
            max_string = answer;
            max = max_string.size();
        }
        s = s.substr(1, s.size());
    }

    std::cout << max_string;
}
