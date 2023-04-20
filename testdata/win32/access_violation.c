// gcc access_violation.c -o access_violation.exe
#include <windows.h>

int main(void) {
    SetErrorMode(0);
    SetThreadErrorMode(0, 0);
    volatile int *p = 0;
    *p = 114514;
    return 0;
}
