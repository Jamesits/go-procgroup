// gcc missing_dll.c -L. -ldll -o missing_dll.exe

#include <stdio.h>
__declspec(dllimport) int __cdecl test();

int main(int argc, char** argv)
{
    printf("%d\n", test());
    return 0;
}
