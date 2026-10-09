/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

bool FUN_00560780(char *param_1,int param_2,undefined8 param_3)

{
  bool bVar1;
  undefined1 uVar2;
  char *pcVar3;
  int iVar4;
  int unaff_R14;
  char *pcStack0000000000000008;

  pcStack0000000000000008 = param_1;
  while (&stack0x00000000 <= *(undefined1 **)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  if (param_2 == 4) {
    bVar1 = true;
    pcVar3 = pcStack0000000000000008;
  }
  else {
    if (param_2 == 0x10) {
      for (iVar4 = 0; iVar4 < 10; iVar4 = iVar4 + 1) {
        if (pcStack0000000000000008[iVar4] != '\0') goto LAB_005607f9;
      }
      if ((pcStack0000000000000008[10] == -1) && (pcStack0000000000000008[0xb] == -1)) {
        bVar1 = true;
        pcVar3 = pcStack0000000000000008 + 0xc;
        goto LAB_005607b0;
      }
    }
LAB_005607f9:
    bVar1 = false;
    pcVar3 = (char *)0x0;
  }
LAB_005607b0:
  if (pcVar3 == (char *)0x0) {
    uVar2 = FUN_00561280(pcStack0000000000000008,param_2,param_3,PTR_DAT_00e03940,DAT_00e03948,
                         DAT_00e03950);
    return (bool)uVar2;
  }
  if (!bVar1) {
                    /* WARNING: Subroutine does not return */
    FUN_004792c0(0,0,0);
  }
  return *pcVar3 == '\x7f';
}
