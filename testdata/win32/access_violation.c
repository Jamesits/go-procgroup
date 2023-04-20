// gcc access_violation.c -o access_violation.exe

int main(void) {
    int *p = 0;
    *p = 114514;
    return 0;
}
