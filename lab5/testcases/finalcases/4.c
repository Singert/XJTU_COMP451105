void process(int multiplier , int size) {
    int i;
    int data[5] = {1, 2, 3, 4, 5};
    for (i = 0; i < size; i = i + 1) {
        data[i] = data[i] * multiplier;
    }
}

int main() {
    int arr[5] = {1, 2, 3, 4, 5};
    process(arr, 5);
    process(arr, 5, 3);
    return 0;
}