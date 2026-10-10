/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

undefined1  [16]
FUN_004e7c80(undefined8 param_1,undefined8 param_2,int param_3,uint param_4,undefined8 param_5)

{
  uint uVar1;
  undefined8 *extraout_RAX;
  undefined8 uVar2;
  int extraout_RAX_00;
  int iVar3;
  int extraout_RAX_01;
  undefined8 *puVar4;
  int iVar5;
  int iVar6;
  int extraout_RCX;
  uint extraout_RCX_00;
  undefined8 *extraout_RCX_01;
  undefined **ppuVar7;
  int extraout_RBX;
  uint extraout_RBX_00;
  undefined8 extraout_RBX_01;
  undefined8 uVar8;
  int iVar9;
  int extraout_RSI;
  uint uVar10;
  uint extraout_RDI;
  int extraout_R8;
  int iVar11;
  int iVar12;
  int extraout_R9;
  int iVar13;
  uint uVar14;
  uint extraout_R10;
  int iVar15;
  undefined8 *extraout_R11;
  undefined8 *extraout_R11_00;
  int *extraout_R11_01;
  undefined8 *extraout_R11_02;
  undefined8 *extraout_R11_03;
  uint uVar16;
  undefined8 uVar17;
  int unaff_R14;
  undefined1 auVar18 [16];
  undefined8 uStack0000000000000008;
  undefined8 uStack0000000000000010;
  int iStack0000000000000018;
  uint uStack0000000000000020;
  undefined8 uStack0000000000000028;
  undefined8 *local_48;
  int local_40;
  undefined8 local_38;
  undefined8 *local_30;
  undefined8 local_28;
  int local_20;
  int local_18;
  int local_10;

  uStack0000000000000008 = param_1;
  iStack0000000000000018 = param_3;
  uStack0000000000000010 = param_2;
  uStack0000000000000028 = param_5;
  uStack0000000000000020 = param_4;
  while (&local_48 <= *(undefined8 ***)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  local_48 = (undefined8 *)FUN_004ea7c0();
  *(undefined1 *)((int)local_48 + 0xb4) = 1;
  FUN_004f0620(local_48,uStack0000000000000008,uStack0000000000000010,iStack0000000000000018,
               uStack0000000000000020,uStack0000000000000028);
  auVar18 = FUN_0045a540(0,*local_48,local_48[1]);
  uVar2 = auVar18._8_8_;
  local_28 = auVar18._0_8_;
  uVar10 = local_48[0x18];
  if (uVar10 == 0) {
    puVar4 = (undefined8 *)FUN_0041a160(&DAT_009b2460);
    puVar4[1] = uVar2;
    if (DAT_013e2350 != 0) {
      puVar4 = (undefined8 *)FUN_00478f00();
      *extraout_R11_00 = local_28;
    }
    *puVar4 = local_28;
    ppuVar7 = &PTR_error___Interface_type_00b40520;
  }
  else if (uVar10 == 1) {
    auVar18 = FUN_0041a160(&DAT_009d09e0);
    *(undefined8 *)(auVar18._0_8_ + 8) = uVar2;
    if (DAT_013e2350 != 0) {
      auVar18 = FUN_00478f00();
      *extraout_R11 = local_28;
    }
    puVar4 = auVar18._0_8_;
    *puVar4 = local_28;
    if (local_48[0x18] == 0) {
                    /* WARNING: Subroutine does not return */
      FUN_004792c0(0,auVar18._8_8_,0);
    }
    uVar10 = *(uint *)local_48[0x17];
    if (uStack0000000000000020 <= uVar10) {
                    /* WARNING: Subroutine does not return */
      FUN_004792c0(uVar10);
    }
    iVar12 = *(int *)(iStack0000000000000018 + uVar10 * 0x10);
    iVar6 = *(int *)(iStack0000000000000018 + 8 + uVar10 * 0x10);
    uVar2 = 0;
    if (iVar12 != 0) {
      uVar10 = (uint)*(dword *)(iVar12 + 0x10);
      do {
        iVar9 = (uVar10 & *(uint *)PTR_DAT_00e03820) * 0x10;
        if (iVar12 == *(int *)(PTR_DAT_00e03820 + iVar9 + 8)) {
          uVar2 = *(undefined8 *)(PTR_DAT_00e03820 + iVar9 + 0x10);
          goto LAB_004e7eef;
        }
        uVar10 = uVar10 + 1;
      } while (*(int *)(PTR_DAT_00e03820 + iVar9 + 8) != 0);
      local_18 = iVar6;
      uVar2 = FUN_00416b20(&PTR_DAT_00e03820);
      iVar6 = local_18;
    }
LAB_004e7eef:
    puVar4[2] = uVar2;
    if (DAT_013e2350 != 0) {
      FUN_00478f20(puVar4,puVar4[3]);
      *extraout_R11_01 = extraout_RCX;
      extraout_R11_01[1] = extraout_RBX;
      puVar4 = extraout_RAX;
      iVar6 = extraout_RCX;
    }
    puVar4[3] = iVar6;
    ppuVar7 = &PTR_error___Interface_type_00b407e0;
  }
  else {
    if (*(char *)(local_48 + 0x16) != '\0') {
      iVar12 = 0x3f;
      if (uVar10 != 0) {
        for (; uVar10 >> iVar12 == 0; iVar12 = iVar12 + -1) {
        }
      }
      if (uVar10 == 0) {
        iVar12 = -1;
      }
      FUN_004f1c20(&PTR_PTR_00b4a0e0,local_48[0x17],uVar10,local_48[0x19],0,uVar10,iVar12 + 1);
    }
    iVar12 = local_48[0x17];
    iVar6 = local_48[0x18];
    uVar10 = 0;
    iVar15 = 0;
    uVar16 = 0;
    puVar4 = local_48;
    uVar8 = uVar2;
    iVar11 = iVar6;
    iVar13 = iStack0000000000000018;
    uVar14 = uStack0000000000000020;
    local_10 = iVar12;
    for (iVar9 = 0; local_40 = iVar15, iVar9 < iVar11; iVar9 = iVar9 + 1) {
      uVar1 = *(uint *)(iVar12 + iVar9 * 8);
      if (iVar9 < 1) {
LAB_004e7fec:
        if (uVar14 <= uVar1) {
                    /* WARNING: Subroutine does not return */
          FUN_004792c0(uVar1,uVar8,uVar14);
        }
        iVar15 = *(int *)(iVar13 + uVar1 * 0x10);
        local_38 = *(undefined8 *)(iVar13 + 8 + uVar1 * 0x10);
        iVar5 = 0;
        if (iVar15 != 0) {
          uVar14 = (uint)*(dword *)(iVar15 + 0x10);
          do {
            iVar13 = (uVar14 & *(uint *)PTR_DAT_00e03840) * 0x10;
            if (iVar15 == *(int *)(PTR_DAT_00e03840 + iVar13 + 8)) {
              iVar5 = *(int *)(PTR_DAT_00e03840 + iVar13 + 0x10);
              uVar8 = uVar2;
              iVar13 = iStack0000000000000018;
              uVar14 = uStack0000000000000020;
              goto LAB_004e8027;
            }
            uVar14 = uVar14 + 1;
          } while (*(int *)(PTR_DAT_00e03840 + iVar13 + 8) != 0);
          iVar5 = FUN_00416b20(&PTR_DAT_00e03840,iVar15);
          puVar4 = local_48;
          iVar12 = local_10;
          uVar8 = uVar2;
          iVar11 = iVar6;
          iVar13 = iStack0000000000000018;
          uVar14 = uStack0000000000000020;
        }
LAB_004e8027:
        iVar15 = local_40;
        if (iVar5 != 0) {
          uVar16 = uVar16 + 1;
          if (uVar10 < uVar16) {
            local_18 = iVar5;
            FUN_00473640(local_40,uVar16,uVar10,1,&error___Interface_type);
            puVar4 = local_48;
            iVar12 = local_10;
            uVar8 = uVar2;
            uVar10 = extraout_RCX_00;
            iVar11 = iVar6;
            iVar13 = iStack0000000000000018;
            uVar14 = uStack0000000000000020;
            iVar15 = extraout_RAX_00;
            uVar16 = extraout_RBX_00;
            iVar5 = local_18;
          }
          iVar3 = (uVar16 - 1) * 0x10;
          *(int *)(iVar15 + iVar3) = iVar5;
          uVar17 = local_38;
          if (DAT_013e2350 != 0) {
            uVar8 = *(undefined8 *)(iVar15 + 8 + iVar3);
            local_20 = iVar15;
            FUN_00478f20();
            *extraout_R11_02 = uVar17;
            extraout_R11_02[1] = uVar8;
            iVar3 = extraout_RAX_01;
            puVar4 = extraout_RCX_01;
            uVar8 = extraout_RBX_01;
            iVar9 = extraout_RSI;
            uVar10 = extraout_RDI;
            iVar11 = extraout_R8;
            iVar13 = extraout_R9;
            uVar14 = extraout_R10;
            iVar15 = local_20;
          }
          *(undefined8 *)(iVar15 + 8 + iVar3) = uVar17;
        }
      }
      else {
        if ((uint)puVar4[0x18] <= iVar9 - 1U) {
                    /* WARNING: Subroutine does not return */
          FUN_004792c0(iVar9 - 1U,uVar8,puVar4[0x18]);
        }
        if (*(uint *)(puVar4[0x17] + -8 + iVar9 * 8) != uVar1) goto LAB_004e7fec;
      }
    }
    puVar4 = (undefined8 *)FUN_0041a160(&DAT_009d0a80,uVar8);
    puVar4[1] = uVar2;
    if (DAT_013e2350 != 0) {
      puVar4 = (undefined8 *)FUN_00478f20();
      *extraout_R11_03 = local_28;
      extraout_R11_03[1] = local_40;
    }
    *puVar4 = local_28;
    puVar4[3] = uVar16;
    puVar4[4] = uVar10;
    puVar4[2] = local_40;
    ppuVar7 = &PTR_error___Interface_type_00b40800;
  }
  local_30 = puVar4;
  FUN_004ea860(local_48);
  auVar18._8_8_ = local_30;
  auVar18._0_8_ = ppuVar7;
  return auVar18;
}
