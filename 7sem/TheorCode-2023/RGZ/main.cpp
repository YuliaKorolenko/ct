#include <iostream>
#include <vector>

#include <chrono>
#include <random>
#include <algorithm>
#include <set>

using namespace std;

#define int long long

std::vector<int> lg;
std::vector<int> alog;

vector<int> g;
vector<set<int>> cycl_classes;
vector<set<int>> cycl_classes_2047;
std::set<int> quadratic_ded;

// 2 ^ 11 - 1
int size_n = 2047;
int power_pol;
const int MAX_LOG = 30;
// primitive_poly = x1^11 + x1^2 + 1
int gpol = 1 + 4 + 2048;
int n = 89;
int alpha = 23;

int mod(int pol){
    for (int i = MAX_LOG; i >= power_pol; i--){
        if (pol & (1 << i)){
            pol = pol ^ (gpol << (i - power_pol));
        }
    }
    return pol;
}

void calc_power(){
    for (int i = MAX_LOG; i > 0; i--){
        if (gpol & (1 << i)){
            power_pol = i;
            break;
        }
    }
}

void pre_calc(){
    // log i - многочлен log[i] - степень
    calc_power();
    lg.resize(size_n + 1);
    alog.resize(size_n);
    int pol = 1;
    for (int i = 1; i <= size_n; i++){
        lg[pol] = i - 1;
        alog[i - 1] = pol;
        
        pol = mod(pol * 2);
    }
    
}

void print_vec(vector<int> vec){
    for (int i = 0; i < vec.size(); i++){
        cout << vec[i] << " ";
    }
    cout << endl;
}

void create_quad_ded(){
    for (int i = 1; i < n; i++){
        quadratic_ded.insert(i*i % n);
    }
}

void get_cyclotomical(){
    // 𝑔(𝑥) = НОК {𝑀𝑖(𝑥), 𝑖 = 𝑏, 𝑏 + 1, . . . , 𝑏 + 𝑑 − 2} .
    // b == 1 
    // иду по номеру класса, пока есть не used
    set<int> used;
    for (int i = 0; i < n; i++){
        if (quadratic_ded.find(i) == quadratic_ded.end() || used.find(i) != used.end()){
            continue;
        }
        set<int> cur_class;
        set<int> cur_class_2047;
        int cur_el = i;

        while (cur_class.find(cur_el) == cur_class.end()){
            cur_class.insert(cur_el); 
            cur_class_2047.insert(alog[cur_el * alpha]);
            used.insert(cur_el);  
            cur_el = (cur_el * 2) % n;
        }
        
        cycl_classes.push_back(cur_class);
        cycl_classes_2047.push_back(cur_class_2047);
    }
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

int mul(int a, int b){
    if (a == 0 || b == 0){
        return 0;
    }
    return alog[(lg[a] + lg[b]) % size_n];
}

void get_g(){
    g.push_back(1);
    for (int ij = 0; ij < cycl_classes_2047.size(); ij++){
        vector<int> cur_m_i = {1};
        for (auto value : cycl_classes_2047[ij]){
            vector<int> m_for_r = cur_m_i;
            for (int j = 0; j < m_for_r.size(); j++){
                m_for_r[j] = mul(value, m_for_r[j]);
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

void print_quadratic(){
    cout << " Quadratic ded: ";
    cout << endl;
    for (auto value : quadratic_ded){
        cout << value << " ";
    }
    cout << endl;
}

void print_cyclonomical(vector<set<int>> cycl){
    cout << " Print cyclomocal: ";
    cout << endl;

    for (int i = 0; i < cycl.size(); i++){
        for (auto value : cycl[i]){
            cout << value << " ";
        }
        cout << endl;
    }
    cout << endl;
}

void print_g(){
    int cur_ex = g.size() - 1;
    for (int i = 0; i < g.size(); i++){
        if (g[i] && cur_ex != 0){
            cout << "x^" << cur_ex << " + ";
        }
        if (cur_ex == 0){
            cout << "1";
        }
        cur_ex--;
    }
    
}

vector<vector<int>> create_matrix(){
    vector<vector<int>> matrix(g.size(), vector<int>(n, 0));
    for (int j = 0; j < g.size(); j++){
        for (int i = 0; i < g.size(); i++){
                matrix[j][j+i] = g[i];   
        }
    }
    return matrix;
}


vector<vector<int>> create_matrix_with_parity(){
    int parity = 0;
    for (int i = 0; i < g.size(); i++){
        parity += g[i];
    }
    
    parity = parity % 2;
    vector<vector<int>> matrix(g.size(), vector<int>(n + 1, 0));
    for (int j = 0; j < g.size(); j++){
        matrix[j][0] = parity;
        for (int i = 0; i < g.size(); i++){
                matrix[j][j+i+1] = g[i];   
        }
    }
    return matrix;
}

void print_matrix(vector<vector<int>> matrix){
    for (int i = 0; i < matrix.size(); i++){
        for (int j = 0; j < matrix[0].size(); j++){
            cout << matrix[i][j];
        }
        cout << endl;
    }
}

signed main(){
    pre_calc();

    create_quad_ded();
    print_quadratic();

    get_cyclotomical();
    print_cyclonomical(cycl_classes);
    print_cyclonomical(cycl_classes_2047);

    get_g();
    cout << " G in 0-1 " << endl;
    print_vec(g);
    cout << endl;

    cout << " G " << endl;
    print_g();
    cout << endl;

    vector<vector<int>> matrix = create_matrix();
    cout << " Matrix with size rows: " << matrix.size() << " columns: " << matrix[0].size() << endl;
    print_matrix(matrix);
    
    vector<vector<int>> matrix_parity = create_matrix_with_parity();
    cout << " Matrix with size rows: " << matrix_parity.size() << " columns: " << matrix_parity[0].size() << endl;
    print_matrix(matrix_parity);
}