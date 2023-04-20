// gcc messagebox.c -o messagebox.exe

#include <windows.h>

int main(void) {
    MessageBox(NULL, TEXT("Hello"), TEXT("UwU"), MB_ICONERROR | MB_OK);

    return 0;
}
