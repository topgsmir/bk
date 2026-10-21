/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

void FUN_00920800(undefined8 param_1,undefined8 param_2,int param_3,undefined8 param_4)

{
  int *piVar1;
  undefined *puVar2;
  undefined *puVar3;
  undefined *puVar4;
  undefined *puVar5;
  undefined *puVar6;
  int extraout_RAX;
  undefined8 extraout_RAX_00;
  int iVar7;
  undefined8 uVar8;
  undefined8 extraout_RAX_01;
  undefined8 *puVar9;
  undefined8 uVar10;
  undefined **extraout_RCX;
  int extraout_RCX_00;
  undefined8 extraout_RCX_01;
  code **ppcVar11;
  int iVar12;
  int iVar13;
  int extraout_RBX;
  undefined8 extraout_RBX_00;
  int extraout_RBX_01;
  char *pcVar14;
  undefined8 extraout_RSI;
  undefined8 extraout_RSI_00;
  undefined8 extraout_RSI_01;
  undefined8 extraout_RSI_02;
  undefined8 extraout_RSI_03;
  undefined8 extraout_RSI_04;
  undefined **extraout_RDI;
  undefined8 *extraout_R11;
  undefined8 *extraout_R11_00;
  undefined8 *extraout_R11_01;
  undefined8 *extraout_R11_02;
  undefined8 *extraout_R11_03;
  undefined8 *extraout_R11_04;
  undefined8 *extraout_R11_05;
  undefined8 *extraout_R11_06;
  undefined8 *extraout_R11_07;
  int *extraout_R11_08;
  int unaff_R14;
  undefined1 auVar15 [16];
  undefined8 uStack0000000000000008;
  int iStack0000000000000018;
  undefined8 uStack0000000000000020;
  code *local_48;
  undefined8 local_40;
  code *local_38;
  undefined8 local_30;
  code *local_28;
  undefined8 local_20;
  internal_abi_Type *local_18;
  undefined **ppuStack_10;

  uStack0000000000000008 = param_1;
  iStack0000000000000018 = param_3;
  uStack0000000000000020 = param_4;
  while (&local_48 <= *(code ***)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  FUN_00921260();
  if (extraout_RBX != 0) {
    local_18 = *(internal_abi_Type **)(extraout_RBX + 8);
    ppuStack_10 = extraout_RCX;
    FUN_0053d8c0(DAT_00e92328,"failed to load configuration: %v",0x20,&local_18,1,1);
  }
  auVar15 = FUN_00540dc0(extraout_RAX);
  if (auVar15._0_8_ != 0) {
    local_18 = *(internal_abi_Type **)(auVar15._0_8_ + 8);
    ppuStack_10 = auVar15._8_8_;
    FUN_0053d8c0(DAT_00e92328,"%v",2,&local_18,1,1);
  }
  FUN_008f7ea0();
  if (extraout_RCX_00 != 0) {
    local_18 = *(internal_abi_Type **)(extraout_RCX_00 + 8);
    ppuStack_10 = extraout_RDI;
    FUN_0053d8c0(DAT_00e92328,"%v",2,&local_18,1,1);
  }
  if ((((*(int *)(extraout_RAX + 0xe0) == 3) && (**(short **)(extraout_RAX + 0xd8) == 0x7574)) &&
      ((char)(*(short **)(extraout_RAX + 0xd8))[1] == 'n')) &&
     (((*(int *)(extraout_RAX + 0x148) == 3 && (**(short **)(extraout_RAX + 0x140) == 0x7069)) &&
      ((char)(*(short **)(extraout_RAX + 0x140))[1] == 'x')))) {
    pcVar14 = (char *)&DAT_00a53410;
    iVar12 = 3;
  }
  else if ((*(int *)(extraout_RAX + 8) == 0) && (*(int *)(extraout_RAX + 0x18) == 0)) {
    if ((*(int *)(extraout_RAX + 0x30) == 0) && (*(int *)(extraout_RAX + 0x40) == 0)) {
      local_18 = &string___String_type;
      ppuStack_10 = &PTR_DAT_00b3c6e0;
      FUN_0053da40(DAT_00e92328,&local_18,1,1);
      iVar12 = 0;
      pcVar14 = (char *)0x0;
    }
    else {
      iVar12 = 6;
      pcVar14 = "client";
    }
  }
  else {
    iVar12 = 6;
    pcVar14 = "server";
  }
  FUN_007d13c0(pcVar14,iVar12,extraout_RAX,extraout_RAX_00,extraout_RBX_00);
  puVar6 = DAT_00e001e8;
  puVar5 = PTR_DAT_00e001e0;
  puVar4 = PTR_DAT_00e001d8;
  puVar3 = PTR_DAT_00e001d0;
  puVar2 = PTR_DAT_00e001c8;
  iVar13 = *(int *)(extraout_RAX + 0x208);
  piVar1 = *(int **)(extraout_RAX + 0x200);
  iVar7 = extraout_RAX;
  if (iVar13 < 9) {
    if (iVar13 == 7) {
      if ((((sdword)*piVar1 == 0x5f776f6c) && (*(short *)((int)piVar1 + 4) == 0x7063)) &&
         (*(char *)((int)piVar1 + 6) == 'u')) {
        if (DAT_013e2350 != 0) {
          iVar7 = FUN_00478f20();
          *extraout_R11 = puVar4;
          extraout_R11[1] = extraout_RSI;
        }
        DAT_00e92358 = puVar4;
      }
      else {
LAB_00920c35:
        if (DAT_013e2350 != 0) {
          iVar7 = FUN_00478f20();
          *extraout_R11_04 = puVar3;
          extraout_R11_04[1] = extraout_RSI_04;
        }
        DAT_00e92358 = puVar3;
      }
    }
    else {
      if ((iVar13 != 8) || (*piVar1 != 0x6465636e616c6162)) goto LAB_00920c35;
      if (DAT_013e2350 != 0) {
        iVar7 = FUN_00478f20();
        *extraout_R11_00 = puVar3;
        extraout_R11_00[1] = extraout_RSI_00;
      }
      DAT_00e92358 = puVar3;
    }
  }
  else if (iVar13 == 10) {
    if ((*piVar1 != 0x6f6d656d5f776f6c) || ((short)piVar1[1] != 0x7972)) goto LAB_00920c35;
    if (DAT_013e2350 != 0) {
      iVar7 = FUN_00478f20();
      *extraout_R11_01 = puVar2;
      extraout_R11_01[1] = extraout_RSI_01;
    }
    DAT_00e92358 = puVar2;
  }
  else if (iVar13 == 0xd) {
    if (((*piVar1 != 0x6f6c5f6172746c75) || ((sdword)piVar1[1] != 0x70635f77)) ||
       (*(char *)((int)piVar1 + 0xc) != 'u')) goto LAB_00920c35;
    if (DAT_013e2350 != 0) {
      iVar7 = FUN_00478f20();
      *extraout_R11_02 = puVar5;
      extraout_R11_02[1] = extraout_RSI_02;
    }
    DAT_00e92358 = puVar5;
  }
  else {
    if ((((iVar13 != 0xf) || (*piVar1 != 0x5f656d6572747865)) || ((sdword)piVar1[1] != 0x5f776f6c))
       || ((*(short *)((int)piVar1 + 0xc) != 0x7063 || (*(char *)((int)piVar1 + 0xe) != 'u'))))
    goto LAB_00920c35;
    if (DAT_013e2350 != 0) {
      iVar7 = FUN_00478f20();
      *extraout_R11_03 = puVar6;
      extraout_R11_03[1] = extraout_RSI_03;
    }
    DAT_00e92358 = puVar6;
  }
  iVar13 = *(int *)(iVar7 + 0x230);
  if (iVar13 < 1) {
    iVar13 = 0x78;
  }
  DAT_00e27c60 = iVar13 * 1000000000;
  if (((iVar12 == 3) && ((short)*(sdword *)pcVar14 == 0x7069)) &&
     (*(char *)((int)pcVar14 + 2) == 'x')) {
    FUN_007e5760(0);
    iVar7 = extraout_RAX;
  }
  if (*(char *)(iVar7 + 0x1c0) != '\0') {
    FUN_00923720(*(undefined8 *)(iVar7 + 0x1c8),*(undefined8 *)(iVar7 + 0x1d0));
    iVar7 = extraout_RAX;
  }
  if (iVar12 == 3) {
    if (((short)*(sdword *)pcVar14 == 0x7069) && (*(char *)((int)pcVar14 + 2) == 'x')) {
      uVar10 = FUN_00922a80();
      FUN_007edce0(*(undefined8 *)(extraout_RAX + 0x238),*(undefined8 *)(extraout_RAX + 0x240));
      uVar8 = FUN_007ed7e0();
      FUN_007e74e0(uVar10,uVar8);
      if (extraout_RBX_01 != 0) {
        local_18 = *(internal_abi_Type **)(extraout_RBX_01 + 8);
        ppuStack_10 = (undefined **)extraout_RCX_01;
        FUN_0053da40(DAT_00e92328,&local_18,1,1);
        return;
      }
      puVar9 = (undefined8 *)FUN_0041a160(&DAT_009e0d60);
      *puVar9 = &LAB_00921140;
      if (DAT_013e2350 != 0) {
        puVar9 = (undefined8 *)FUN_00478f20();
        *extraout_R11_05 = extraout_RAX_01;
        extraout_R11_05[1] = uStack0000000000000020;
      }
      puVar9[1] = extraout_RAX_01;
      puVar9[2] = iStack0000000000000018;
      puVar9[3] = uStack0000000000000020;
      FUN_0044b100();
      local_48 = FUN_00925560;
      ppcVar11 = &local_48;
      iVar7 = extraout_RAX;
      local_40 = extraout_RAX_01;
      goto LAB_00920fc1;
    }
  }
  else if (iVar12 == 6) {
    if ((*(sdword *)pcVar14 == 0x65696c63) && ((short)*(sdword *)((int)pcVar14 + 4) == 0x746e)) {
      uVar10 = FUN_008f6360(iVar7,iStack0000000000000018,uStack0000000000000020);
      puVar9 = (undefined8 *)FUN_0041a160(&DAT_009bdd00);
      *puVar9 = &LAB_009211a0;
      if (DAT_013e2350 != 0) {
        puVar9 = (undefined8 *)FUN_00478f00();
        *extraout_R11_06 = uVar10;
      }
      puVar9[1] = uVar10;
      FUN_0044b100();
      local_38 = (code *)&LAB_00925500;
      ppcVar11 = &local_38;
      iVar7 = extraout_RAX;
      local_30 = uVar10;
      goto LAB_00920fc1;
    }
    if ((*(sdword *)pcVar14 == 0x76726573) && ((short)*(sdword *)((int)pcVar14 + 4) == 0x7265)) {
      uVar10 = FUN_0091bb60(iVar7,iStack0000000000000018,uStack0000000000000020);
      puVar9 = (undefined8 *)FUN_0041a160(&DAT_009bde00);
      *puVar9 = &LAB_00921200;
      if (DAT_013e2350 != 0) {
        puVar9 = (undefined8 *)FUN_00478f00();
        *extraout_R11_07 = uVar10;
      }
      puVar9[1] = uVar10;
      FUN_0044b100();
      local_28 = (code *)&LAB_009254a0;
      ppcVar11 = &local_28;
      iVar7 = extraout_RAX;
      local_20 = uVar10;
      goto LAB_00920fc1;
    }
  }
  ppcVar11 = (code **)0x0;
LAB_00920fc1:
  if ((((((iVar12 == 3) && ((short)*(sdword *)pcVar14 == 0x7069)) &&
        (*(char *)((int)pcVar14 + 2) == 'x')) &&
       ((*(int *)(iVar7 + 0x298) == 6 && (**(sdword **)(iVar7 + 0x290) == 0x76726573)))) &&
      ((short)(*(sdword **)(iVar7 + 0x290))[1] == 0x7265)) ||
     (((iVar12 == 6 && (*(sdword *)pcVar14 == 0x76726573)) &&
      (((short)*(sdword *)((int)pcVar14 + 4) == 0x7265 &&
       (((*(int *)(iVar7 + 0xe0) == 3 && (**(short **)(iVar7 + 0xd8) == 0x7574)) &&
        ((char)(*(short **)(iVar7 + 0xd8))[1] == 'n')))))))) {
    puVar9 = (undefined8 *)FUN_0041a160(&DAT_009bdc80);
    *puVar9 = &LAB_009210e0;
    if (DAT_013e2350 != 0) {
      puVar9 = (undefined8 *)FUN_00478f00();
      *extraout_R11_08 = extraout_RAX;
    }
    puVar9[1] = extraout_RAX;
    FUN_0044b100();
  }
  uVar10 = (**(code **)(iStack0000000000000018 + 0x20))(uStack0000000000000020);
  FUN_004127e0(uVar10,0);
  if (ppcVar11 != (code **)0x0) {
    (**ppcVar11)();
  }
  return;
}
