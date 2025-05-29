int foo(int a){
    int b = 0;
    if (a > 0) {
        b = a * 2;
    } else {
        b = a + 5;
    }
    return b;
}
foo(10);