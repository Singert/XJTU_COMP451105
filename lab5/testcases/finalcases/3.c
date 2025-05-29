float compute(float x = 1, int y) {
    float result;
    result = (float)(y) * -x + (float)(-y);
    return result;
}

void printResult(float value) {
    printf("Result: %f\n", value);
}

int main() {
    float res = compute(3.5, 4);
    printResult(res);
    return 0;
}
