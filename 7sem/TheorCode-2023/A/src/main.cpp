#include <iostream>
#include <fstream>
#include <vector>
#include <algorithm>
#include <map>
#include <cmath>
#include <cstdint>
#include <cstdlib>
#include <ctime>
#include <chrono>
#include <random>
  
class FileException : public std::exception {
public:
    const char* what() const throw() {
        return "Ошибка открытия файла.";
    }
};
  
class NonZeroRange {
public:
    int left; 
    int right;
  
    NonZeroRange(int left, int right) : left(left), right(right) {}
};
  
// create array with information about first non zeros elements
std::vector<NonZeroRange> non_zero;
std::vector<int> is_alive;
std::vector<int> is_died;
std::vector<int> active_children;
std::vector<int> parents;
std::vector<int> ans_for_parents;
std::vector<long double> weights;
std::vector<int> decode_result;
  
class VertexGrid {
public:
    int index;
    std::vector<int> m_i; 
    int number_of_layer;
    std::vector<std::pair<int, int>> children;   
  
    VertexGrid(int i, std::vector<int> m_i, int number_of_layer, std::vector<std::pair<int, int>> children) 
        : index(i), m_i(m_i), number_of_layer(number_of_layer), children(children) {}  
};
  
  
struct NoiseGenerator{
    NoiseGenerator() : rd(), gen(2008/*rd()*/) {}
    long double gen_gauss(long double mean, long double sigma) {
        long seed = std::chrono::system_clock::now().time_since_epoch().count();
        static std::default_random_engine generator(seed);
        std::normal_distribution<double> normalDistribution(mean, sigma);
        return normalDistribution(generator);
    }
  
    std::vector<long double> add_noise(std::vector<int> & encoded, long double noise_level, int n, int k) {
        std::vector<long double> new_vect(encoded.size());
        long double decibel = pow(10, -noise_level / 10);
        long double deviation = (0.5 * n / k) * decibel;
        long double sigma = sqrt(deviation);
  
        for (int i = 0; i < encoded.size(); i++) {
            new_vect[i] = ((encoded[i] == 1) ? -1.0 : 1.0) + this->gen_gauss(0.0, sigma);
        }
        return new_vect;
    } 
  
    std::random_device rd;
    std::mt19937 gen;
};
  
std::vector<VertexGrid> grid;
  
std::vector<std::vector<int> > read_matrix(std::ifstream & file){
    if (!file) {
        throw FileException();
        std::cerr << "Ошибка открытия файла." << std::endl;
    }
  
    int rows, columns;
    file >> columns >> rows;
  
    std::vector<std::vector<int> > matrix(rows, std::vector<int>(columns));
  
    for (int i = 0; i < rows; ++i) {
        for (int j = 0; j < columns; ++j) {
            file >> matrix[i][j];
        }
    }
    return matrix;
}
  
std::vector<long double> get_random_vector(int n){
    std::vector<long double> array(n);  
    for (int i = 0; i < n; i++) {
        array[i] = rand() % 2;
    }
  
    return array;
}
  
void displayMatrix(std::vector<std::vector<int> > matrix){
    std::cout << "Матрица:" << std::endl;
    for (int i = 0; i < matrix.size(); ++i) {
        for (int j = 0; j < matrix[i].size(); ++j) {
            std::cout << matrix[i][j] << " ";
        }
        std::cout << std::endl;
    }
}
  
std::vector<int> add_vectors(const std::vector<int>& vec1, const std::vector<int>& vec2) {
    std::vector<int> result;
      
    for (size_t i = 0; i < vec1.size(); i++) {
        result.push_back((vec1[i] + vec2[i]) % 2);
    }
    return result;
}
  
NonZeroRange calc_span(const std::vector<int>& vec){
    int j1 = 0;
    while (j1 < vec.size() && vec[j1] == 0){
        j1++;
    }
    int j2 = vec.size() - 1;
    while (j2 > -1 && vec[j2] == 0){
        j2--;
    }
    return NonZeroRange(j1, j2);
}
  
std::vector<std::vector<int> > to_minimal_span(std::vector<std::vector<int> > matrix){
    int rows_count = matrix.size();
    int columns_count = matrix[0].size();
  
    // for sort
    std::vector<std::vector<int> > sort_result(columns_count);
  
    for (int i = 0; i < matrix.size(); ++i){
        non_zero.push_back(calc_span(matrix[i]));
    }
  
    bool is_find = true;
    while (is_find){
        is_find = false;
        for (int i = 0; i < rows_count; i++){
            for (int j = i + 1; j < rows_count; j++){
                if (i == j){
                    continue;
                }
                int del = -1;
  
                if (non_zero[i].left == non_zero[j].left || non_zero[i].right == non_zero[j].right){
                    is_find = true;
                    std::vector<int> result_vect = add_vectors(matrix[i], matrix[j]);
                    NonZeroRange range_for_result = calc_span(result_vect);
                    if (non_zero[i].right < non_zero[j].right || non_zero[i].left < non_zero[j].left){
                        del = j;
                    } else {
                        del = i;
                    }
                    matrix[del] = result_vect;
                    non_zero[del] = range_for_result;
                }
            }
        }
    }
  
    std::vector<int> numbers(rows_count);
    for (int i = 0; i < rows_count; ++i) {
        numbers[i] = i;
    }
  
    std::sort(numbers.begin(), numbers.end(), [&](int i1, int i2){
        return non_zero[i1].left < non_zero[i2].left;
    });
  
    for (int i = 0; i < numbers.size(); i++)
    {
        int j = numbers[i];
        if (numbers[i] < i){
            continue;
        }
        std::swap(matrix[i], matrix[j]);
        std::swap(non_zero[i], non_zero[j]);
    }
      
    return matrix;
}
  
  
void print_vertex(VertexGrid vertex){
    std::cout << std::endl << vertex.index << std::endl;
    std::cout << "m: (";
    for (int i = 0; i < vertex.m_i.size(); i++){
        std::cout << vertex.m_i[i] << " ";
    }
    std::cout << ")" << std::endl;
  
      
    for (int i = 0; i < vertex.children.size(); i++){
        std::cout << "~~ chldren: " << vertex.children[i].first << " weight: " << vertex.children[i].second << " ~~  ";
    }
      
    std::cout << std::endl;
  
}
  
std::vector<int> encode(std::vector<long double> vect, std::vector<std::vector<int> >& matrix){
    int rows_count = matrix.size();
    int columns_count = matrix[0].size();
    std::vector<int> encode_vector(columns_count, 0);
    for (int i = 0; i < columns_count; i++){
        for (int j = 0; j < rows_count; j++){
            encode_vector[i] += matrix[j][i] * vect[j];
        }
        encode_vector[i] = encode_vector[i] % 2;
    }
    return encode_vector;
}
  
void decode(std::vector<long double> & vect, std::vector<VertexGrid> & grid, std::vector<int> & answer){
    for (int i = 0; i < weights.size(); i++){
        weights[i] = INT64_MAX;    
    }
    weights[0] = 0;
    ans_for_parents[0] = -1;
    int layers_count = vect.size();
    int i = 0;
    while (i < grid.size()){
        for (int j = 0; j < grid[i].children.size(); j++){
            std::pair<int, int> child = grid[i].children[j];
            long double cur_weight = (child.second ? 1 : -1) * vect[grid[child.first].number_of_layer];
            if (cur_weight + weights[i] < weights[child.first]){
                weights[child.first] = cur_weight + weights[i];
                parents[child.first] = i;
                ans_for_parents[child.first] = child.second;
            }
        }  
        i++;
    }
 
    int pt = 0;
  
    int j = grid.size() - 1;
    while (j > 0){
        answer[pt++] = ans_for_parents[j];
        j = parents[j];
    }
  
  
    std::reverse(answer.begin(), answer.end());
}
  
  
long double simulate(int noise_level, int num_of_iteration, int max_errors, std::vector<std::vector<int> >& matrix){
    int j = 0;
    int all_mistake_count = 0;
    long double frequency = 0;
    int rows_count = matrix.size();
    int columns_count = matrix[0].size();
    while (j < num_of_iteration){
        j++;
        std::vector<long double> random_vector = get_random_vector(columns_count);
        std::vector<int> encoded_vect = encode(random_vector, matrix);
  
        NoiseGenerator ng = NoiseGenerator();
        std::vector<long double> noise_vect = ng.add_noise(encoded_vect, noise_level, columns_count, rows_count);
        decode(noise_vect, grid, decode_result);
  
        int is_mistake = 0;
        for (int i = 0; i < columns_count; i++){
            if (encoded_vect[i] != decode_result[i]){
                is_mistake = 1;
                break;
            }
        }
        all_mistake_count += is_mistake;
        if (all_mistake_count >= max_errors){
            return (long double)all_mistake_count / j;
        }    
    }
      
    return (long double)all_mistake_count / j;
}
  
void print_grid(std::vector<VertexGrid> grid, int rows_count){
    std::cout << std::endl << "GRID";
    int j = 0;
    for (int i = -1; i < rows_count; i++){
        while (j < grid.size() && grid[j].number_of_layer == i){
            print_vertex(grid[j]);
            j += 1;
        }
        std::cout << "-----------------------------" << std::endl;
    }
}
  
void add_one_to_vector(std::vector<int> & cur_m){
    int j = 0;
    int prev = 1;
    while (j < cur_m.size()){
        if (cur_m[j] == 1){
            cur_m[j] = 0;
        } else if (cur_m[j] == 0){
            cur_m[j] = 1;
            break;
        }
        j += 1;
    }     
}
  
std::vector<int> is_similar_vertex(std::vector<int> & parent, std::vector<int> & children){
    std::vector<int> result_sim;
  
    for (int i = 0; i < parent.size(); i++){
        if (parent[i] != -1){
            if (parent[i] != children[i]){
                if (children[i] != -1){
                    return result_sim;
                } else {
                    result_sim.push_back(parent[i]);
                    continue;
                }
            }
        }
        result_sim.push_back(children[i]);
      
    } 
    return result_sim;
}
  
int multiply_two_vectors(std::vector<int> & first, std::vector<std::vector<int> >  & span, int layer){
    int res = 0;
    for (int i = 0; i < first.size(); i++){
        res += ((first[i] == -1) ? 0 : first[i]) * ((span[i][layer] == -1) ? 0 : span[i][layer]);
    }
    return res % 2;
}
  
int binaryToDecimal(const std::vector<int>& binary_vector) {
    int decimal = 0;
    int power = 1;
  
    for (int i = 0; i < binary_vector.size(); ++i) {
        if (binary_vector[i] == -1){
            continue;
        }
        decimal += binary_vector[i] * power;
        power *= 2;
    }
  
    return decimal;
}
  
void create_next_layer(int layer, std::vector<VertexGrid> & grid, std::vector<std::vector<int> > &span_matrix){
    int active_count = 0;
    int index = grid.size() - 1;
    std::vector<int> cur_m(non_zero.size(), -1);
    for (int i = 0; i < non_zero.size(); i++){
        if (non_zero[i].left <= layer && non_zero[i].right > layer){
            active_count += 1;
            cur_m[i] = 0;
        }
    }
  
    std::vector<std::pair<int, int>> start_vector(0);
  
    int children_number = std::pow(2, active_count);
  
    active_children.push_back(children_number);
  
    for (int i = 0; i < children_number; i++){
        VertexGrid vert = VertexGrid(index + 1 + i, cur_m, layer, start_vector);
        grid.push_back(vert);
        add_one_to_vector(cur_m);
    }
  
    // проходимся по элементам предыдущего слоя
    int j = index;
    while (j > -1 && grid[j].number_of_layer == layer - 1){
        std::vector<int> sons_vectors;
        int n_l = grid[j].number_of_layer + 1;
        if (is_alive[n_l] == -1){
            if (is_died[n_l] == -1){
                int son = active_children[n_l] + j;
                sons_vectors.push_back(son);
             } else {
                std::vector<int> sons_vec = grid[j].m_i;
                sons_vec[is_died[n_l]] = -1;
                int pos = binaryToDecimal(sons_vec);
                sons_vectors.push_back(index + 1 + pos);
             }
        } else {
            std::vector<int> sons_vec = grid[j].m_i;
            if (is_died[n_l] != -1){
                sons_vec[is_died[n_l]] = -1;
            }
            for (int bt = 0; bt < 2; bt++){
                sons_vec[is_alive[n_l]] = bt;
                int pos = binaryToDecimal(sons_vec);
                sons_vectors.push_back(index + 1 + pos);
            }       
        }
  
        for (int i = 0; i < sons_vectors.size(); i++){
            std::vector<int> ver = is_similar_vertex(grid[j].m_i, grid[sons_vectors[i]].m_i);
            grid[j].children.push_back(std::make_pair(sons_vectors[i], multiply_two_vectors(ver, span_matrix, layer))); 
        }
          
        j -= 1;
    }
}
  
std::vector<VertexGrid> create_grid(std::vector<std::vector<int> > &span_matrix){
    int rows_count = span_matrix.size();
    int columns_count = span_matrix[0].size();
  
    std::vector<VertexGrid> grid;
  
    int start_active = -1;
    int next_active = -1;
    int layer = -1;
  
    std::vector<int> cur_m(rows_count, -1);
  
    int j = 0;
  
    std::vector<std::pair<int, int>> cur_map;
  
    VertexGrid firt_vertex = VertexGrid(0, cur_m, layer, cur_map);
    grid.push_back(firt_vertex);
    active_children.push_back(1);
  
 
    while (layer < columns_count - 1){
        layer++;
        create_next_layer(layer, grid, span_matrix);
    }
  
    return grid;
}
  
  
int main(int argc, char *argv[])
{   
    srand(time(0));
    std::ifstream file("input.txt");
    freopen("output.txt", "w", stdout);
    std::vector<std::vector<int> > matrix = read_matrix(file);
  
    std::vector<std::vector<int> > span_matrix = to_minimal_span(matrix);
  
    is_alive = std::vector<int>(matrix[0].size(), -1);
    is_died = std::vector<int>(matrix[0].size(), -1);
    for (int i = 0; i < non_zero.size(); i++){
        is_alive[non_zero[i].left] = i;
        is_died[non_zero[i].right] = i;
    }
      
  
    grid = create_grid(span_matrix);
 
    parents.resize(grid.size());
    ans_for_parents.resize(grid.size());
    weights.resize(grid.size());
    decode_result.resize(matrix[0].size());
  
    std::string line;
  
    for (int i = 0; i < active_children.size(); i++){
        std::cout << active_children[i] << " ";
    }
    std::cout << std::endl;
      
    while (std::getline(file, line)) {
        std::string type;
  
        file >> type;
        int n;
        if (type == "Encode"){
            n = matrix.size();
        } else if (type == "Decode"){
            n = matrix[0].size();
        } else if (type == "Simulate"){
            n = 0;
        }
  
        std::vector<long double> vector_for_type(n);
  
        for (int i = 0; i < n; i++){
            file >> vector_for_type[i];
        }
  
        if (type == "Encode"){
            std::vector<int> encode_result = encode(vector_for_type, matrix);
            for (int i = 0; i < encode_result.size(); i++){
                std::cout << encode_result[i] << " ";
            }
            std::cout << std::endl;
        } else if (type == "Decode"){
            decode(vector_for_type, grid, decode_result);
            for (int i = 0; i < decode_result.size(); i++){
                std::cout << decode_result[i] << " ";
            }
            std::cout << std::endl;
        } else if (type == "Simulate"){
            double noise_level;
            int num_of_iteration, max_errors;
            file >> noise_level >> num_of_iteration >> max_errors;
            long double frequency = simulate(noise_level, num_of_iteration, max_errors, matrix);
            std::cout << frequency << std::endl;
        }       
    } 
}