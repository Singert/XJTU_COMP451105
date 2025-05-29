int main() {
    int arr[5][5];
    arr[0][2] = 10;
    arr[1][2] = (int)(3.14);
    int x = arr[0][0] + arr[1][2];
    return x;
}
