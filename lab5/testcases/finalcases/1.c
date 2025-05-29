int max(int a, int b) {
    if (a > b) {
        return a;
    } else {
        return b;
    }
}

void printMax(int x, int y) {
    int m = max(x, y);
    if (m > 0) {
        printf("Max is %d\n", m);
    } else {
        printf("No max found\n");
    }
}

int main() {
    int a = 10;
    int b = 20;
    printMax(a, b);
    return 0;
}
