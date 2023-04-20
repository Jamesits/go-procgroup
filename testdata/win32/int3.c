// gcc int3.c -o int3.exe

int main(void) {
    asm("int3");
    return 0;
}
