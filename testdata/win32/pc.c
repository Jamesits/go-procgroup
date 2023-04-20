// gcc pc.c -o pc.exe -mwindows -municode -mwin32
#include <windows.h>

void (*p)() = 0;

int wWinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPWSTR lpCmdLine, int nShowCmd) {
    SetErrorMode(0);
    SetThreadErrorMode(0, 0);
    p();
    return 0;
}
