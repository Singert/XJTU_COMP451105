int square(int x) {
    return x * x;
}

int compute(int a, int b) {
    int val = square(a + b);
    if (val > 100 && a != b) {
        return val;
    } else {
        return 0;
    }
}
