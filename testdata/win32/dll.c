// gcc dll.c -o dll.dll -shared

__declspec(dllexport) int __cdecl test()
{
  return 0;
}
