// gcc harderror.c -lntdll -o harderror.exe -mwindows -municode -mwin32
#include <windows.h>
NTSTATUS NtRaiseHardError(NTSTATUS ErrorStatus, ULONG NumberOfParameters, ULONG UnicodeStringParameterMask OPTIONAL, PULONG_PTR Parameters, ULONG ResponseOption, PULONG Response);

int wWinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPWSTR lpCmdLine, int nShowCmd) {
    SetErrorMode(0);
    SetThreadErrorMode(0, 0);
    ULONG resp;
    NtRaiseHardError(0xC0000005, 0, 0, 0, 2, &resp);
    return 0;
}
