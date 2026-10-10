/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

undefined1  [32] FUN_008f8040(int param_1,undefined8 param_2)

{
  int extraout_RAX;
  undefined8 extraout_RAX_00;
  int extraout_RAX_01;
  undefined8 extraout_RAX_02;
  undefined8 extraout_RAX_03;
  int iVar1;
  int iVar12;
  undefined8 extraout_RAX_04;
  undefined8 extraout_RAX_05;
  undefined8 *puVar2;
  undefined8 *puVar3;
  int iVar4;
  undefined8 *extraout_RCX;
  undefined8 extraout_RCX_00;
  uint extraout_RCX_01;
  uint extraout_RCX_02;
  undefined8 extraout_RCX_03;
  uint extraout_RCX_04;
  uint extraout_RCX_05;
  uint extraout_RCX_06;
  uint extraout_RCX_07;
  undefined8 extraout_RCX_08;
  uint extraout_RCX_09;
  uint extraout_RCX_10;
  uint extraout_RCX_14;
  uint extraout_RCX_15;
  uint uVar9;
  byte bVar10;
  byte bVar11;
  undefined8 extraout_RBX;
  int extraout_RBX_00;
  int extraout_RBX_01;
  undefined8 extraout_RBX_02;
  undefined8 extraout_RBX_03;
  undefined1 auVar5 [32];
  uint extraout_RCX_11;
  uint extraout_RCX_12;
  int extraout_RCX_13;
  undefined8 extraout_RBX_04;
  int extraout_RBX_05;
  undefined1 auVar6 [32];
  undefined8 extraout_RBX_06;
  undefined8 uVar13;
  undefined1 auVar7 [32];
  undefined1 auVar8 [32];
  int extraout_RSI;
  undefined8 extraout_RSI_00;
  undefined8 extraout_RSI_01;
  byte extraout_DIL;
  undefined8 extraout_RDI;
  undefined8 extraout_R8;
  undefined8 *extraout_R11;
  undefined8 *extraout_R11_00;
  undefined8 *extraout_R11_01;
  undefined8 *extraout_R11_02;
  undefined8 *extraout_R11_03;
  undefined8 *extraout_R11_04;
  undefined8 *extraout_R11_05;
  undefined8 *extraout_R11_06;
  undefined8 *extraout_R11_07;
  undefined8 *extraout_R11_08;
  undefined8 *extraout_R11_09;
  undefined8 *extraout_R11_10;
  undefined8 *extraout_R11_11;
  int unaff_R14;
  undefined8 *in_XMM15_Qa;
  code **in_XMM15_Qb;
  undefined1 auVar14 [16];
  int iStack0000000000000008;
  undefined8 uStack0000000000000010;
  undefined1 local_b0 [16];
  undefined8 *local_a0;
  code **ppcStack_98;
  undefined8 *local_90;
  int local_88;
  int local_80;
  undefined8 local_78;
  undefined8 local_70;
  code *local_68;
  undefined8 local_60;
  undefined8 local_58;
  internal_abi_Type *local_50;
  undefined8 uStack_48;
  undefined8 *local_40;
  undefined8 local_38;
  undefined8 local_30;
  undefined8 *local_28;
  undefined8 *local_20;
  undefined8 *local_18;
  code **ppcStack_10;

  iStack0000000000000008 = param_1;
  uStack0000000000000010 = param_2;
  while (local_b0 <= *(undefined1 **)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  local_18 = in_XMM15_Qa;
  ppcStack_10 = in_XMM15_Qb;
  (**(code **)(iStack0000000000000008 + 0x18))(uStack0000000000000010);
  iVar1 = iStack0000000000000008;
  local_30 = uStack0000000000000010;
  if (extraout_DIL == 0) {
    FUN_004a72a0(iStack0000000000000008,uStack0000000000000010,3000000000);
    iVar1 = extraout_RAX;
    local_30 = extraout_RBX;
    local_18 = extraout_RCX;
  }
  local_88 = FUN_0041a160(&DAT_009fed20);
  *(undefined8 *)(local_88 + 0x28) = 3000000000;
  FUN_00799060(iVar1,local_30,"GET",3,"https://api.ipify.org?format=json",0x21,0,0);
  bVar10 = extraout_DIL ^ 1;
  bVar11 = bVar10;
  if (extraout_RBX_00 == 0) {
    FUN_00762800(local_88,extraout_RAX_00);
    if (extraout_RBX_01 == 0) {
      bVar11 = bVar10 | 2;
      local_60 = *(undefined8 *)(extraout_RAX_01 + 0x40);
      local_58 = *(undefined8 *)(extraout_RAX_01 + 0x48);
      local_68 = FUN_008f8b80;
      ppcStack_10 = &local_68;
      if (*(int *)(extraout_RAX_01 + 0x10) == 200) {
        local_80 = extraout_RAX_01;
        local_28 = (undefined8 *)FUN_0041a160(&DAT_009a7080);
        iVar4 = *(int *)(local_80 + 0x40);
        local_78 = *(undefined8 *)(local_80 + 0x48);
        local_38 = 0;
        if (iVar4 != 0) {
          uVar9 = (uint)*(dword *)(iVar4 + 0x10);
          do {
            iVar12 = (uVar9 & *(uint *)PTR_DAT_00e03800) * 0x10;
            if (*(int *)(PTR_DAT_00e03800 + iVar12 + 8) == *(int *)(iVar4 + 8)) {
              local_38 = *(undefined8 *)(PTR_DAT_00e03800 + iVar12 + 0x10);
              goto LAB_008f8975;
            }
            uVar9 = uVar9 + 1;
          } while (*(int *)(PTR_DAT_00e03800 + iVar12 + 8) != 0);
          local_38 = FUN_00416b20(&PTR_DAT_00e03800);
        }
LAB_008f8975:
        puVar3 = (undefined8 *)FUN_0041a160(&DAT_00a24ba0);
        *puVar3 = local_38;
        if (DAT_013e2350 != 0) {
          puVar3 = (undefined8 *)FUN_00478f00();
          *extraout_R11_09 = local_78;
        }
        puVar3[1] = local_78;
        iVar4 = FUN_0050aba0(puVar3,&DAT_009697e0,local_28);
        if ((iVar4 == 0) && (ppcStack_98 = (code **)local_28[1], ppcStack_98 != (code **)0x0)) {
          local_a0 = (undefined8 *)*local_28;
          local_b0._0_8_ = in_XMM15_Qa;
          local_b0._8_8_ = in_XMM15_Qb;
          (**ppcStack_10)();
          if ((bVar10 & 1) != 0) {
            (*(code *)*local_18)();
          }
          auVar8._8_8_ = ppcStack_98;
          auVar8._0_8_ = local_a0;
          auVar8._16_8_ = local_b0._0_8_;
          auVar8._24_8_ = local_b0._8_8_;
          return auVar8;
        }
        auVar14 = FUN_004e7c80("ipify decode failed",0x13,0,0,0);
        local_38 = auVar14._8_8_;
        local_90 = (undefined8 *)FUN_00473640(0,1,0,1,&error___Interface_type);
        *local_90 = auVar14._0_8_;
        uVar9 = extraout_RCX_14;
        puVar3 = extraout_R11_10;
        if (DAT_013e2350 != 0) {
          uVar13 = local_90[1];
          local_90 = (undefined8 *)FUN_00478f20();
          *extraout_R11_11 = local_38;
          extraout_R11_11[1] = uVar13;
          uVar9 = extraout_RCX_15;
          puVar3 = extraout_R11_11;
        }
        local_90[1] = local_38;
      }
      else {
        uStack_48 = FUN_0046eb60(*(undefined8 *)(extraout_RAX_01 + 0x10));
        local_50 = &int___Int_type;
        auVar14 = FUN_004e7c80("ipify bad status: %d",0x14,&local_50,1,1);
        local_38 = auVar14._8_8_;
        local_90 = (undefined8 *)FUN_00473640(0,1,0,1,&error___Interface_type);
        *local_90 = auVar14._0_8_;
        uVar9 = extraout_RCX_04;
        puVar3 = extraout_R11_01;
        if (DAT_013e2350 != 0) {
          uVar13 = local_90[1];
          local_90 = (undefined8 *)FUN_00478f20();
          *extraout_R11_02 = local_38;
          extraout_R11_02[1] = uVar13;
          uVar9 = extraout_RCX_05;
          puVar3 = extraout_R11_02;
        }
        local_90[1] = local_38;
      }
    }
    else {
      local_50 = *(internal_abi_Type **)(extraout_RBX_01 + 8);
      uStack_48 = extraout_RCX_03;
      auVar14 = FUN_004e7c80("ipify request failed: %w",0x18,&local_50,1,1);
      local_38 = auVar14._8_8_;
      local_90 = (undefined8 *)FUN_00473640(0,1,0,1,&error___Interface_type);
      *local_90 = auVar14._0_8_;
      uVar9 = extraout_RCX_06;
      puVar3 = extraout_R11_03;
      if (DAT_013e2350 != 0) {
        uVar13 = local_90[1];
        local_90 = (undefined8 *)FUN_00478f20();
        *extraout_R11_04 = local_38;
        extraout_R11_04[1] = uVar13;
        uVar9 = extraout_RCX_07;
        puVar3 = extraout_R11_04;
      }
      local_90[1] = local_38;
    }
  }
  else {
    local_50 = *(internal_abi_Type **)(extraout_RBX_00 + 8);
    uStack_48 = extraout_RCX_00;
    auVar14 = FUN_004e7c80("ipify request build failed: %w",0x1e,&local_50,1,1);
    local_38 = auVar14._8_8_;
    local_90 = (undefined8 *)FUN_00473640(0,1,0,1,&error___Interface_type);
    *local_90 = auVar14._0_8_;
    uVar9 = extraout_RCX_01;
    puVar3 = extraout_R11;
    if (DAT_013e2350 != 0) {
      uVar13 = local_90[1];
      local_90 = (undefined8 *)FUN_00478f20();
      *extraout_R11_00 = local_38;
      extraout_R11_00[1] = uVar13;
      uVar9 = extraout_RCX_02;
      puVar3 = extraout_R11_00;
    }
    local_90[1] = local_38;
  }
  FUN_007d0700(iVar1,local_30,"GET",3,"/getip",6,0,0,puVar3,in_XMM15_Qa,in_XMM15_Qb);
  if (extraout_RSI == 0) {
    local_70 = extraout_RAX_02;
    local_20 = (undefined8 *)FUN_0041a160(&dwXnm9Hr_Q1::_IpOXsKfLWZg_T_1thXQ___Struct_type);
    iVar1 = FUN_004fa0e0(local_70,extraout_RBX_02,extraout_RCX_08,
                         &dwXnm9Hr_Q1::__IpOXsKfLWZg_T_1thXQ___Pointer_type,local_20);
    if ((iVar1 == 0) && (ppcStack_98 = (code **)local_20[1], ppcStack_98 != (code **)0x0)) {
      local_a0 = (undefined8 *)*local_20;
      local_b0._0_8_ = in_XMM15_Qa;
      local_b0._8_8_ = in_XMM15_Qb;
      if ((bVar11 & 2) != 0) {
        bVar11 = bVar11 & 0xfd;
        (**ppcStack_10)();
      }
      if ((bVar11 & 1) != 0) {
        (*(code *)*local_18)();
      }
      auVar5._8_8_ = ppcStack_98;
      auVar5._0_8_ = local_a0;
      auVar5._16_8_ = local_b0._0_8_;
      auVar5._24_8_ = local_b0._8_8_;
      return auVar5;
    }
    auVar14 = FUN_004e7c80("internal API invalid response",0x1d,0,0,0);
    uVar13 = auVar14._8_8_;
    puVar3 = local_90;
    if (uVar9 < 2) {
      local_30 = uVar13;
      puVar3 = (undefined8 *)FUN_00473640(local_90,2,uVar9,1,&error___Interface_type);
      uVar9 = extraout_RCX_11;
      uVar13 = local_30;
    }
    puVar3[2] = auVar14._0_8_;
    if (DAT_013e2350 != 0) {
      FUN_00478f20(puVar3[3]);
      *extraout_R11_06 = extraout_RBX_04;
      extraout_R11_06[1] = extraout_RAX_04;
      uVar9 = extraout_RCX_12;
      uVar13 = extraout_RBX_04;
    }
    puVar3[3] = uVar13;
  }
  else {
    local_50 = *(internal_abi_Type **)(extraout_RSI + 8);
    uStack_48 = extraout_R8;
    auVar14 = FUN_004e7c80("internal API failed: %w",0x17,&local_50,1,1);
    uVar13 = auVar14._8_8_;
    puVar3 = local_90;
    if (uVar9 < 2) {
      local_30 = uVar13;
      puVar3 = (undefined8 *)FUN_00473640(local_90,2,uVar9,1,&error___Interface_type);
      uVar9 = extraout_RCX_09;
      uVar13 = local_30;
    }
    puVar3[2] = auVar14._0_8_;
    if (DAT_013e2350 != 0) {
      FUN_00478f20(puVar3[3]);
      *extraout_R11_05 = extraout_RBX_03;
      extraout_R11_05[1] = extraout_RAX_03;
      uVar9 = extraout_RCX_10;
      uVar13 = extraout_RBX_03;
    }
    puVar3[3] = uVar13;
  }
  local_90 = puVar3;
  FUN_008f7ea0();
  if (extraout_RCX_13 == 0) {
    if (extraout_RBX_05 != 0) {
      local_b0._0_8_ = in_XMM15_Qa;
      local_b0._8_8_ = in_XMM15_Qb;
      local_a0 = (undefined8 *)extraout_RAX_05;
      ppcStack_98 = (code **)extraout_RBX_05;
      if ((bVar11 & 2) != 0) {
        bVar11 = bVar11 & 0xfd;
        (**ppcStack_10)();
      }
      if ((bVar11 & 1) != 0) {
        (*(code *)*local_18)();
      }
      auVar6._8_8_ = ppcStack_98;
      auVar6._0_8_ = local_a0;
      auVar6._16_8_ = local_b0._0_8_;
      auVar6._24_8_ = local_b0._8_8_;
      return auVar6;
    }
    puVar3 = (undefined8 *)FUN_0041a160(&DAT_009b2460);
    puVar3[1] = 0x12;
    *puVar3 = &DAT_00a7264e;
    puVar2 = local_90;
    if (uVar9 < 3) {
      local_40 = puVar3;
      puVar2 = (undefined8 *)FUN_00473640(local_90,3,uVar9,1,&error___Interface_type);
      puVar3 = local_40;
    }
    puVar2[4] = &PTR_error___Interface_type_00b40520;
    if (DAT_013e2350 != 0) {
      puVar3 = (undefined8 *)FUN_00478f20();
      *extraout_R11_08 = puVar3;
      extraout_R11_08[1] = extraout_RSI_01;
    }
    puVar2[5] = puVar3;
  }
  else {
    local_50 = (internal_abi_Type *)0x0;
    if (extraout_RCX_13 != 0) {
      local_50 = *(internal_abi_Type **)(extraout_RCX_13 + 8);
    }
    uStack_48 = extraout_RDI;
    auVar14 = FUN_004e7c80("interface IP failed: %w",0x17,&local_50,1,1);
    uVar13 = auVar14._8_8_;
    puVar2 = local_90;
    if (uVar9 < 3) {
      local_30 = uVar13;
      puVar2 = (undefined8 *)FUN_00473640(local_90,3,uVar9,1,&error___Interface_type);
      uVar13 = local_30;
    }
    puVar2[4] = auVar14._0_8_;
    if (DAT_013e2350 != 0) {
      FUN_00478f20();
      *extraout_R11_07 = extraout_RBX_06;
      extraout_R11_07[1] = extraout_RSI_00;
      uVar13 = extraout_RBX_06;
    }
    puVar2[5] = uVar13;
  }
  uStack_48 = FUN_0046ec80(puVar2,3);
  local_50 = (internal_abi_Type *)&DAT_0096abe0;
  local_b0 = FUN_004e7c80("all outbound IP methods failed: %v",0x22,&local_50,1,1);
  local_a0 = in_XMM15_Qa;
  ppcStack_98 = in_XMM15_Qb;
  if ((bVar11 & 2) != 0) {
    bVar11 = bVar11 & 0xfd;
    (**ppcStack_10)();
  }
  if ((bVar11 & 1) != 0) {
    (*(code *)*local_18)();
  }
  auVar7._8_8_ = ppcStack_98;
  auVar7._0_8_ = local_a0;
  auVar7._16_8_ = local_b0._0_8_;
  auVar7._24_8_ = local_b0._8_8_;
  return auVar7;
}
