int double_if_positive(int x) {
    if (x > 0) {
        return x * 2;
    }
    return x;
}

int main() {
    int a = -3;
    int b = 4;
    int result = double_if_positive(a + b);
    return result;
}
