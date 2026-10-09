/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

char * FUN_00560ac0(int param_1,int param_2)

{
  uint uVar1;
  undefined1 uVar2;
  undefined1 *puVar3;
  char *pcVar4;
  uint uVar5;
  int iVar6;
  int unaff_R14;
  undefined1 auVar7 [16];
  int iStack0000000000000008;
  undefined1 local_50 [24];
  undefined1 local_38 [8];
  uint local_30;
  undefined1 local_28 [32];

  iStack0000000000000008 = param_1;
  while (local_38 <= *(undefined1 **)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  if (param_2 == 0) {
    return "<nil>";
  }
  if ((param_2 != 4) && (param_2 != 0x10)) {
    uVar5 = param_2 * 2;
    if (uVar5 < 0x21) {
      puVar3 = local_28;
    }
    else {
      local_30 = uVar5;
      puVar3 = (undefined1 *)FUN_00473560(&uint8___Uint8_type,uVar5);
      uVar5 = local_30;
    }
    iVar6 = 0;
    while( true ) {
      if (param_2 <= iVar6) {
        auVar7 = FUN_0045a540(local_50,puVar3);
        pcVar4 = (char *)FUN_00459f40(0,"?",1,auVar7._0_8_,auVar7._8_8_);
        return pcVar4;
      }
      uVar2 = (&DAT_00a70434)[*(byte *)(iStack0000000000000008 + iVar6) & 0xf];
      if (uVar5 <= (uint)(iVar6 * 2)) break;
      puVar3[iVar6 * 2] = (&DAT_00a70434)[*(byte *)(iStack0000000000000008 + iVar6) >> 4];
      uVar1 = iVar6 * 2 + 1;
      if (uVar5 <= uVar1) {
                    /* WARNING: Subroutine does not return */
        FUN_004792c0(uVar1);
      }
      puVar3[iVar6 * 2 + 1] = uVar2;
      iVar6 = iVar6 + 1;
    }
                    /* WARNING: Subroutine does not return */
    FUN_004792c0(iVar6 * 2);
  }
  auVar7 = FUN_00560cc0();
  pcVar4 = (char *)FUN_0045a540(0,auVar7._0_8_,auVar7._8_8_);
  return pcVar4;
}
