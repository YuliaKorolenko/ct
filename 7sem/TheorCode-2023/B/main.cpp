#include <iostream>
#include <vector>
#include <set>
#include <fstream>
#include <chrono>
#include <random>
#include <algorithm>

using namespace std;

int n, d, gpol;
std::vector<int> lg;
std::vector<int> alog;
const int MAX_LOG = 30;
int power_pol;
int k;
vector<int> g;

int mod(int pol){
    for (int i = MAX_LOG; i >= power_pol; i--){
        if (pol & (1 << i)){
            pol = pol ^ (gpol << (i - power_pol));
        }
    }
    return pol;
}


struct NoiseGenerator{
    NoiseGenerator() : rd(), gen(2008/*rd()*/) {}
    bool gen_error(double prob) {
        long seed = chrono::system_clock::now().time_since_epoch().count();
        static default_random_engine generator(seed);
        uniform_real_distribution<double> realDistribution(0.0, 1.0);
        return realDistribution(generator) < prob;
    }
    std::random_device rd;
    std::mt19937 gen;
};


void calc_power(){
    for (int i = MAX_LOG; i > 0; i--){
        if (gpol & (1 << i)){
            power_pol = i;
            break;
        }
    }
}

std::vector<int> get_random_vector(int n){
    std::vector<int> array(n);  
    for (int i = 0; i < n; i++) {
        array[i] = rand() % 2;
    }
  
    return array;
}

void pre_calc(){
    // log i - многочлен log[i] - степень
    calc_power();
    lg.resize(n + 1);
    alog.resize(n);
    int pol = 1;
    for (int i = 1; i <= n; i++){
        lg[pol] = i - 1;
        alog[i - 1] = pol;
        
        pol = mod(pol * 2);
    }
    
}

int add(int a, int b){
    return a ^ b;
}

int mul(int a, int b){
    if (a == 0 || b == 0){
        return 0;
    }
    return alog[(lg[a] + lg[b]) % n];
}

vector<int> mul_vectors(vector<int> a, vector<int> b){
    vector<int> result(a.size() + b.size() - 1);
    for (int i = 0; i < a.size(); i++) {
        for (int j = 0; j < b.size(); j++) {
            result[i + j] = result[i + j] ^ (a[i] * b[j]);
        }
    }
    
    return result;
}

vector<int> encode(vector<int> & v){
    int r = (int) g.size() - 1;

    vector<int> encodedv(r, 0);
    for (int i = 0; i < v.size(); i++){
        encodedv.push_back(v[i]);
    }
    
    vector<int> new_enc = encodedv;
    for (size_t i = new_enc.size() - 1; i >= r; i--){
        if (new_enc[i] == 1){
            for (int j = 0; j < g.size(); j++){
                new_enc[i + j - g.size() + 1] ^= g[j];    
            }
        }
    }
    
    for (int i = 0; i < r; i++){
        encodedv[i] = new_enc[i];
    }
    
    return encodedv;
}

int get_sindrome(vector<int> & v, vector<int> & syndrome){
    int errors_count = 0;
    for (int i = 1; i < d; i++) {
        for (int j = 0; j < v.size(); j++) {
            if (v[j]){
                syndrome[i - 1] ^= alog[(j * i) % n];
            }
        }
        if (syndrome[i - 1] != 0){
            errors_count++;
        }
    }
    return errors_count;
}

void decode(vector<int> & v){
    vector<int> syndrome(d - 1, 0);
    int errors_count = get_sindrome(v, syndrome);

    if (errors_count == 0){
        return;
    }
    
    int L = 0;
    vector<int> LAMBDA = {1};
    vector<int> B = {1};
    int m = 0;
    for (int r = 1; r < d; r++){
        int resid = 0;
        for (int j = 0; j <= L; j++){
            resid ^= mul(syndrome[r - j - 1], LAMBDA[j]);
        }

        if (resid){
            vector<int> T(max(LAMBDA.size(), r - m + B.size()), 0);
            for (int i = 0; i < LAMBDA.size(); i++){
                T[i] = LAMBDA[i];
            }    
            for (int i = 0; i < B.size(); i++){
                T[r - m + i] ^= mul(resid, B[i]);
            }
            if (2 * L <= r - 1){
                B = LAMBDA;
                int inv = alog[(n - lg[resid]) % n];
                for (int j = 0; j < B.size(); j++){
                    B[j] = mul(B[j], inv);
                }
                L = r - L;
                m = r;
            }
            LAMBDA = T;
        }
    }
    if (LAMBDA.size() - 1 != L){
        return;
    }

    int err_count = 0;
    for (int i = 0; i < n; i++){
        int res = LAMBDA[0];
        for (int j = 1; j < LAMBDA.size(); j++){
            res ^= mul(alog[(j * i) % n], LAMBDA[j]);
        }
        if (res == 0){
            int index = (n - i) % n;
            v[index] = !v[index];
            err_count++;
            if (err_count == L ){
                break;
            }
        }
    }
    
}

long double simulate(long double noise_level, int num_of_iterations, int max_errors){
    int errors_numbers = 0;
    int all_mistake_count = 0;
    for (int i = 0; i < num_of_iterations; i++){
        std::vector<int> random_vector = get_random_vector(k);
        std::vector<int> encoded_vect = encode(random_vector);
  
        NoiseGenerator noise_gen;
        std::vector<int> noise_vect = encoded_vect;
        for (int i = 0; i < encoded_vect.size(); i++){
            if (noise_gen.gen_error(noise_level)){
                noise_vect[i] = 1 - noise_vect[i];
            }
        }
        
        decode(noise_vect);

        int is_mistake = 0;
        for (int i = 0; i < noise_vect.size(); i++){
            if (encoded_vect[i] != noise_vect[i]){
                is_mistake = 1;
                break;
            }
        }
        all_mistake_count += is_mistake;
        if (all_mistake_count >= max_errors){
            return (long double)all_mistake_count / (i + 1);
        }    
    }
    return (long double)all_mistake_count / num_of_iterations;
}

void get_cyclotomical(){
    // 𝑔(𝑥) = НОК {𝑀𝑖(𝑥), 𝑖 = 𝑏, 𝑏 + 1, . . . , 𝑏 + 𝑑 − 2} .
    // b == 1 
    int used_all = 1 + d - 2;
    vector<set<int>> cycl_classes;
    set<int> used;
    // иду по номеру класса, пока есть не used
    for (int i = 0; i < used_all; i++){
        if (used.find(i) != used.end()) {
            continue;
        }

        set<int> cur_class;
        int cur_el = i;

        while (cur_class.find(cur_el) == cur_class.end()){
            cur_class.insert(cur_el);   
            used.insert(cur_el);
            cur_el = (cur_el * 2) % n;
        }
        
        cycl_classes.push_back(cur_class);
    }

    g.push_back(1);
    for (int ij = 1; ij < cycl_classes.size(); ij++){
        vector<int> cur_m_i = {1};
        for (auto value : cycl_classes[ij]){
            vector<int> m_for_r = cur_m_i;
            for (int j = 0; j < m_for_r.size(); j++){
                m_for_r[j] = mul(alog[value], m_for_r[j]);
            }
            cur_m_i.push_back(0);
            for (int j = 0; j < m_for_r.size(); j++){
                cur_m_i[j+1] ^= m_for_r[j];
            }
            
        }

        reverse(cur_m_i.begin(), cur_m_i.end());
        g = mul_vectors(g, cur_m_i);    
    } 
}

int main(){
    std::ifstream file("input.txt");
    freopen("output.txt", "w", stdout);
    file >> n >> gpol >> d;

    pre_calc();
    
    get_cyclotomical();

    k = n + 1 - g.size();

    cout << k << endl;

    for (int i = 0; i < g.size(); i++){
        cout << g[i] << " ";
    }
    cout << endl;

    std::string line;
    while (std::getline(file, line)) {
        std::string type;
        file >> type;
        int vec_size;
        if (type == "Encode"){
            vec_size = k;
        } else if (type == "Decode"){
            vec_size = n;
        } else if (type == "Simulate"){
            vec_size = 0;
        }

        std::vector<int> vector_for_type(vec_size);
  
        for (int i = 0; i < vec_size; i++){
            file >> vector_for_type[i];
        }
        if (type == "Encode"){
            std::vector<int> encode_result = encode(vector_for_type);
            for (int i = 0; i < encode_result.size(); i++){
                cout << encode_result[i] << " ";
            }
            cout << endl;
        } else if (type == "Decode"){
            decode(vector_for_type);
            for (int i = 0; i < vector_for_type.size(); i++){
                std::cout << vector_for_type[i] << " ";
            }
            cout << endl;
        } else if (type == "Simulate"){
            long double noise_level;
            int num_of_iterations, max_errors;
            file >> noise_level >> num_of_iterations >> max_errors;
            cout << simulate(noise_level, num_of_iterations, max_errors) << endl;
        }
    }  
}