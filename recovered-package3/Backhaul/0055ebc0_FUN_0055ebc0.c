/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

undefined8 FUN_0055ebc0(void)

{
  undefined8 uVar1;
  undefined8 *puVar2;
  undefined8 extraout_RSI;
  int extraout_RDI;
  undefined8 *extraout_R11;
  int unaff_R14;
  undefined8 in_XMM15_Qa;
  undefined8 in_XMM15_Qb;

  while (&stack0x00000000 <= *(undefined1 **)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  uVar1 = FUN_0055ff20(0);
  if (extraout_RDI != 0) {
    puVar2 = (undefined8 *)FUN_0041a160(&DAT_00a093e0);
    puVar2[1] = 5;
    *puVar2 = &DAT_00a54d0b;
    puVar2[3] = 6;
    puVar2[2] = &DAT_00a56979;
    puVar2[4] = in_XMM15_Qa;
    puVar2[5] = in_XMM15_Qb;
    puVar2[6] = in_XMM15_Qa;
    puVar2[7] = in_XMM15_Qb;
    puVar2[8] = extraout_RDI;
    if (DAT_013e2350 != 0) {
      puVar2 = (undefined8 *)FUN_00478f00();
      *extraout_R11 = extraout_RSI;
    }
    puVar2[9] = extraout_RSI;
  }
  return uVar1;
}
