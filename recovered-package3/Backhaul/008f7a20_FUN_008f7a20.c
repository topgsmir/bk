/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

undefined8 FUN_008f7a20(undefined8 param_1,undefined8 param_2,undefined8 param_3)

{
  undefined8 uVar1;
  undefined8 uVar2;
  char cVar3;
  undefined8 extraout_RAX;
  undefined8 uVar4;
  undefined8 extraout_RCX;
  undefined8 uVar5;
  int iVar6;
  undefined8 extraout_RBX;
  undefined8 uVar7;
  undefined8 extraout_RSI;
  uint uVar8;
  undefined8 extraout_R8;
  undefined8 extraout_R9;
  int unaff_R14;
  undefined8 in_XMM15_Qb;
  undefined8 uStack0000000000000008;
  undefined8 uStack0000000000000010;
  undefined8 uStack0000000000000018;
  undefined local_198 [8];
  undefined4 uStack_190;
  undefined4 uStack_18c;
  undefined local_188 [8];
  undefined4 uStack_180;
  undefined4 uStack_17c;
  undefined local_178 [8];
  undefined4 uStack_170;
  undefined4 uStack_16c;
  byte local_168 [32];
  byte local_148 [16];
  undefined1 local_138 [16];
  byte local_128 [32];
  int local_108;
  undefined8 local_100;
  undefined8 uStack_f8;
  undefined8 local_f0;
  undefined8 uStack_e8;
  undefined4 local_e0;
  undefined4 uStack_dc;
  undefined4 uStack_d8;
  undefined4 uStack_d4;
  undefined8 local_d0;
  undefined8 uStack_c8;
  undefined8 local_c0;
  byte *pbStack_b8;
  undefined8 local_b0;
  undefined8 uStack_a8;
  undefined8 local_a0;
  undefined8 local_98;
  undefined8 local_90;
  byte *local_88;
  undefined8 local_80;
  undefined8 local_78;
  undefined8 local_70;
  undefined8 local_68;
  undefined8 local_60;
  byte *local_58;
  undefined8 local_50;
  undefined8 local_48;
  undefined8 local_40;
  undefined8 uStack_38;
  undefined8 local_30;
  undefined8 uStack_28;
  undefined8 local_20;
  undefined8 uStack_18;
  undefined8 *local_10;

  uStack0000000000000008 = param_1;
  uStack0000000000000018 = param_3;
  uStack0000000000000010 = param_2;
  while (uStack_170 = (undefined4)in_XMM15_Qb, local_138 <= *(undefined1 **)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  iVar6 = 0x10;
  if (DAT_00e038e8 < 0x10) {
    iVar6 = DAT_00e038e8;
  }
  uStack_180 = uStack_170;
  if (PTR_DAT_00e038e0 != local_178) {
    FUN_00479ca0(local_178,PTR_DAT_00e038e0,iVar6);
  }
  _uStack_170 = CONCAT44(10,uStack_170);
  uVar8 = 8;
  for (iVar6 = 0; iVar6 < 4; iVar6 = iVar6 + 1) {
    if (uVar8 < 8) {
      local_128[iVar6] = ~(0xffU >> ((byte)uVar8 & 0x1f));
      uVar8 = 0;
    }
    else {
      local_128[iVar6] = 0xff;
      uVar8 = uVar8 - 8;
    }
  }
  iVar6 = 0x10;
  if (DAT_00e038e8 < 0x10) {
    iVar6 = DAT_00e038e8;
  }
  uStack_190 = uStack_180;
  if (PTR_DAT_00e038e0 != local_188) {
    FUN_00479ca0(local_188,PTR_DAT_00e038e0,iVar6,uVar8,local_178);
  }
  _uStack_180 = CONCAT44(0x10ac,uStack_180);
  uVar8 = 0xc;
  for (iVar6 = 0; iVar6 < 4; iVar6 = iVar6 + 1) {
    if (uVar8 < 8) {
      local_148[iVar6] = ~(0xffU >> ((byte)uVar8 & 0x1f));
      uVar8 = 0;
    }
    else {
      local_148[iVar6] = 0xff;
      uVar8 = uVar8 - 8;
    }
  }
  iVar6 = 0x10;
  if (DAT_00e038e8 < 0x10) {
    iVar6 = DAT_00e038e8;
  }
  if (PTR_DAT_00e038e0 != local_198) {
    FUN_00479ca0(local_198,PTR_DAT_00e038e0,iVar6,uVar8,local_178,local_188);
  }
  _uStack_190 = CONCAT44(0xa8c0,uStack_190);
  uVar8 = 0x10;
  for (iVar6 = 0; iVar6 < 4; iVar6 = iVar6 + 1) {
    if (uVar8 < 8) {
      local_168[iVar6] = ~(0xffU >> ((byte)uVar8 & 0x1f));
      uVar8 = 0;
    }
    else {
      local_168[iVar6] = 0xff;
      uVar8 = uVar8 - 8;
    }
  }
  FUN_004795b9(uStack0000000000000008,uStack0000000000000010,uStack0000000000000018,&local_100,
               local_178,local_188,local_198);
  uStack_c8 = 0x10;
  local_c0 = 0x10;
  local_d0 = extraout_RSI;
  local_b0 = 4;
  uStack_a8 = 4;
  pbStack_b8 = local_128;
  local_98 = 0x10;
  local_90 = 0x10;
  local_a0 = extraout_R8;
  local_80 = 4;
  local_78 = 4;
  local_88 = local_148;
  local_68 = 0x10;
  local_60 = 0x10;
  local_70 = extraout_R9;
  local_50 = 4;
  local_48 = 4;
  local_58 = local_168;
  local_10 = &local_d0;
  local_108 = 0;
  uVar4 = extraout_RAX;
  uVar5 = extraout_RCX;
  uVar7 = extraout_RBX;
  while( true ) {
    if (2 < local_108) {
      return 0;
    }
    local_100 = *local_10;
    uStack_f8 = local_10[1];
    local_40 = local_100;
    uStack_38 = uStack_f8;
    local_f0 = local_10[2];
    uStack_e8 = local_10[3];
    local_30 = local_f0;
    uStack_28 = uStack_e8;
    uVar1 = local_10[4];
    uVar2 = local_10[5];
    local_20._0_4_ = (undefined4)uVar1;
    local_20._4_4_ = (undefined4)((uint)uVar1 >> 0x20);
    uStack_18._0_4_ = (undefined4)uVar2;
    uStack_18._4_4_ = (undefined4)((uint)uVar2 >> 0x20);
    local_e0 = (undefined4)local_20;
    uStack_dc = local_20._4_4_;
    uStack_d8 = (undefined4)uStack_18;
    uStack_d4 = uStack_18._4_4_;
    local_20 = uVar1;
    uStack_18 = uVar2;
    cVar3 = FUN_00561900(&local_100,uVar4,uVar7,uVar5);
    if (cVar3 != '\0') break;
    local_10 = local_10 + 6;
    local_108 = local_108 + 1;
    uVar4 = uStack0000000000000008;
    uVar5 = uStack0000000000000018;
    uVar7 = uStack0000000000000010;
  }
  return 1;
}
