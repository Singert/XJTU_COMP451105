void initMatrix(int matrix[2][3]) {
    int i;
    int j;
    for (i = 0; i < 2; i = i + 1) {
        for (j = 0; j < 3; j = j + 1) {
            matrix[i][j] = i * 3 + j;
        }
    }
}

int sumMatrix(int matrix[2][3]) {
    int i;
    int j;
    int sum = 0;
    for (i = 0; i < 2; i = i + 1) {
        for (j = 0; j < 3; j = j + 1) {
            sum = sum + matrix[i][j];
        }
    }
    return sum;
}

int main() {
    int m[2][3];
    initMatrix(m);
    int total = sumMatrix(m);
    printf("Sum is %d\n", total);
    return 0;
}
