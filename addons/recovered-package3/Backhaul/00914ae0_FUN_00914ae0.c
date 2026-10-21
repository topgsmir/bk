/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

undefined8 FUN_00914ae0(undefined8 param_1,int param_2,undefined8 param_3,int param_4)

{
  char cVar1;
  undefined8 extraout_RAX;
  undefined8 uVar2;
  undefined8 extraout_RAX_00;
  undefined8 extraout_RAX_01;
  undefined8 *puVar3;
  int iVar4;
  undefined8 extraout_RAX_02;
  undefined8 uVar5;
  undefined8 uVar6;
  int extraout_RCX;
  int extraout_RCX_00;
  undefined8 extraout_RCX_01;
  undefined8 extraout_RCX_02;
  undefined8 uVar7;
  int extraout_RBX;
  int extraout_RBX_00;
  undefined8 extraout_RBX_01;
  undefined8 extraout_RBX_02;
  undefined8 extraout_RBX_03;
  undefined8 uVar8;
  undefined8 extraout_RDI;
  undefined8 extraout_RDI_00;
  uint uVar9;
  int iVar10;
  int unaff_R14;
  int iStack0000000000000010;
  undefined8 uStack0000000000000018;
  int iStack0000000000000020;
  undefined1 local_f0 [8];
  undefined1 local_e8 [8];
  int local_e0;
  int local_d8;
  undefined8 local_d0;
  undefined8 local_c8;
  undefined8 local_c0;
  undefined8 local_b8;
  undefined8 local_b0;
  undefined8 local_a8;
  undefined8 local_a0;
  undefined8 local_98;
  undefined8 local_90;
  undefined8 *local_88;
  undefined8 local_80;
  undefined8 local_78 [2];
  undefined8 local_68;
  undefined8 local_58 [2];
  undefined8 local_48;
  internal_abi_Type *local_38;
  undefined8 uStack_30;
  internal_abi_Type *local_28;
  undefined8 uStack_20;
  undefined8 local_18;
  undefined8 uStack_10;

  uStack0000000000000018 = param_3;
  iStack0000000000000010 = param_2;
  iStack0000000000000020 = param_4;
  while (&local_a0 <= *(undefined8 **)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  FUN_008f7ea0();
  if (extraout_RCX != 0) {
    local_18 = *(undefined8 *)(extraout_RCX + 8);
    uStack_10 = extraout_RDI;
    uVar2 = FUN_004e7c80("failed to get interface ip: %s",0x1e,&local_18,1,1);
    return uVar2;
  }
  local_e0 = extraout_RBX;
  local_a8 = extraout_RAX;
  FUN_008f8040(iStack0000000000000010,uStack0000000000000018);
  if (extraout_RCX_00 != 0) {
    local_18 = *(undefined8 *)(extraout_RCX_00 + 8);
    uStack_10 = extraout_RDI_00;
    uVar2 = FUN_004e7c80("failed to get outband ip: %s",0x1c,&local_18,1,1);
    return uVar2;
  }
  local_d8 = extraout_RBX_00;
  local_a0 = extraout_RAX_00;
  if ((extraout_RBX_00 != local_e0) ||
     (cVar1 = FUN_00404b60(local_a8,extraout_RAX_00), cVar1 == '\0')) {
    uStack_30 = FUN_0046ebe0(local_a8,local_e0);
    local_38 = &string___String_type;
    uStack_20 = FUN_0046ebe0(local_a0,local_d8);
    local_28 = &string___String_type;
    uVar2 = FUN_004e7c80("interface ip %s and outband ip %s mismatch",0x2a,&local_38,2,2);
    return uVar2;
  }
  LOCK();
  uVar2 = *(undefined8 *)(iStack0000000000000020 + 0x30);
  *(undefined8 *)(iStack0000000000000020 + 0x30) = 10;
  UNLOCK();
  FUN_007e5760(1,extraout_RBX_01,uVar2);
  FUN_0045a700(0,local_a8,local_e0);
  iVar4 = *(int *)(iStack0000000000000020 + 0x18);
  uVar2 = *(undefined8 *)(iStack0000000000000020 + 0x20);
  if (iVar4 != 0) {
    uVar9 = (uint)*(dword *)(iVar4 + 0x10);
    do {
      iVar10 = (uVar9 & *(uint *)PTR_DAT_00e06fc0) * 0x10;
      if (*(int *)(PTR_DAT_00e06fc0 + iVar10 + 8) == *(int *)(iVar4 + 8)) {
        uVar6 = *(undefined8 *)(PTR_DAT_00e06fc0 + iVar10 + 0x10);
        uVar5 = extraout_RAX_01;
        uVar7 = extraout_RCX_01;
        uVar8 = extraout_RBX_02;
        goto LAB_00914f7c;
      }
      uVar9 = uVar9 + 1;
    } while (*(int *)(PTR_DAT_00e06fc0 + iVar10 + 8) != 0);
    local_c8 = extraout_RBX_02;
    local_b8 = extraout_RCX_01;
    local_90 = extraout_RAX_01;
    local_80 = uVar2;
    uVar6 = FUN_00416b20(&PTR_DAT_00e06fc0,*(int *)(iVar4 + 8));
    uVar5 = local_90;
    uVar7 = local_b8;
    uVar8 = local_c8;
    uVar2 = local_80;
LAB_00914f7c:
    FUN_007337c0(uVar6,uVar2,0x2f,uVar5,uVar8,uVar7);
  }
  puVar3 = (undefined8 *)FUN_0049a3c0(3000000000);
  local_b0 = *puVar3;
  local_58[0] = (**(code **)(iStack0000000000000010 + 0x20))(uStack0000000000000018);
  local_48 = local_b0;
  iVar4 = FUN_00452960(local_58,local_e8,0,0,2,1);
  if (iVar4 == 0) {
    return 0;
  }
  do {
    do {
      local_88 = (undefined8 *)FUN_0049a8a0(*(int *)(iStack0000000000000020 + 0x30) * 1000000000);
      local_68 = (**(code **)(iStack0000000000000010 + 0x20))(uStack0000000000000018);
      local_78[0] = *local_88;
      iVar4 = FUN_00452960(local_78,local_f0,0,0,2,1);
      if (iVar4 != 0) {
        if (*(char *)(local_88 + 1) != '\0') {
          FUN_00475180();
        }
        return 0;
      }
      if (*(char *)(local_88 + 1) != '\0') {
        FUN_00475180();
      }
      FUN_0045a700(0,local_a8,local_e0);
      iVar4 = *(int *)(iStack0000000000000020 + 0x18);
      uVar2 = *(undefined8 *)(iStack0000000000000020 + 0x20);
    } while (iVar4 == 0);
    uVar9 = (uint)*(dword *)(iVar4 + 0x10);
    do {
      iVar10 = (uVar9 & *(uint *)PTR_DAT_00e06fe0) * 0x10;
      if (*(int *)(PTR_DAT_00e06fe0 + iVar10 + 8) == *(int *)(iVar4 + 8)) {
        uVar6 = *(undefined8 *)(PTR_DAT_00e06fe0 + iVar10 + 0x10);
        uVar5 = extraout_RAX_02;
        uVar7 = extraout_RCX_02;
        uVar8 = extraout_RBX_03;
        goto LAB_00914ee9;
      }
      uVar9 = uVar9 + 1;
    } while (*(int *)(PTR_DAT_00e06fe0 + iVar10 + 8) != 0);
    local_d0 = extraout_RBX_03;
    local_c0 = extraout_RCX_02;
    local_98 = extraout_RAX_02;
    local_80 = uVar2;
    uVar6 = FUN_00416b20(&PTR_DAT_00e06fe0,*(int *)(iVar4 + 8));
    uVar5 = local_98;
    uVar7 = local_c0;
    uVar8 = local_d0;
    uVar2 = local_80;
LAB_00914ee9:
    FUN_007337c0(uVar6,uVar2,0x2f,uVar5,uVar8,uVar7);
  } while( true );
}
