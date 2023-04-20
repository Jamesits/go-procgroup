// gcc fatalappexit.c -o fatalappexit.exe -mwindows -municode -mwin32
#include <windows.h>

int wWinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPWSTR lpCmdLine, int nShowCmd) {
    SetErrorMode(0);
    SetThreadErrorMode(0, 0);
    FatalAppExitW(0, L"114514");
    return 0;
}
