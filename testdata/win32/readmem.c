// gcc readmem.c -o readmem.exe -mwindows -municode -mwin32
#include <windows.h>

int wWinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPWSTR lpCmdLine, int nShowCmd) {
    SetErrorMode(0);
    SetThreadErrorMode(0, 0);
    asm("xor %rax, %rax");
    asm("movq (%rax), %rax");
    return 0;
}
