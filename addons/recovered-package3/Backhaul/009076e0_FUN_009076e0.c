/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

void FUN_009076e0(int param_1,int param_2,undefined8 param_3,undefined8 *param_4)

{
  undefined8 uVar1;
  char cVar2;
  undefined8 extraout_RAX;
  undefined8 extraout_RAX_00;
  undefined8 extraout_RAX_01;
  undefined8 *puVar3;
  int iVar4;
  undefined8 extraout_RAX_02;
  int extraout_RCX;
  int extraout_RCX_00;
  undefined8 extraout_RCX_01;
  undefined8 extraout_RCX_02;
  int extraout_RBX;
  int extraout_RBX_00;
  undefined8 extraout_RBX_01;
  undefined8 extraout_RBX_02;
  undefined8 extraout_RBX_03;
  undefined8 extraout_RDI;
  undefined8 extraout_RDI_00;
  int unaff_R14;
  int iStack0000000000000008;
  int iStack0000000000000010;
  undefined8 uStack0000000000000018;
  undefined8 *puStack0000000000000020;
  undefined1 local_f0 [8];
  undefined1 local_e8 [8];
  int local_e0;
  int local_d8;
  undefined8 local_d0;
  undefined8 local_c8;
  undefined8 local_c0;
  undefined8 *local_b8;
  undefined8 local_b0;
  undefined8 local_a8 [2];
  undefined8 local_98;
  undefined8 local_88 [2];
  undefined8 local_78;
  internal_abi_Type *local_68;
  undefined8 uStack_60;
  internal_abi_Type *local_58;
  undefined8 uStack_50;
  internal_abi_Type *local_48;
  undefined **ppuStack_40;
  undefined8 local_38;
  undefined8 uStack_30;
  internal_abi_Type *local_28;
  undefined **ppuStack_20;
  undefined8 local_18;
  undefined8 uStack_10;

  iStack0000000000000008 = param_1;
  uStack0000000000000018 = param_3;
  iStack0000000000000010 = param_2;
  puStack0000000000000020 = param_4;
  while (local_a8 <= *(undefined8 **)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  FUN_008f7ea0();
  if (extraout_RCX != 0) {
    local_28 = &string___String_type;
    ppuStack_20 = &PTR_DAT_00b3dd40;
    local_18 = *(undefined8 *)(extraout_RCX + 8);
    uStack_10 = extraout_RDI;
    FUN_0053d960(*(undefined8 *)(iStack0000000000000008 + 0x18),2,&local_28,2,2);
    return;
  }
  local_d8 = extraout_RBX;
  local_c0 = extraout_RAX;
  FUN_008f8040(iStack0000000000000010,uStack0000000000000018);
  if (extraout_RCX_00 != 0) {
    local_48 = &string___String_type;
    ppuStack_40 = &PTR_DAT_00b3dd50;
    local_38 = *(undefined8 *)(extraout_RCX_00 + 8);
    uStack_30 = extraout_RDI_00;
    FUN_0053d960(*(undefined8 *)(iStack0000000000000008 + 0x18),2,&local_48,2,2);
    return;
  }
  local_e0 = extraout_RBX_00;
  local_c8 = extraout_RAX_00;
  if ((extraout_RBX_00 == local_d8) &&
     (cVar2 = FUN_00404b60(local_c0,extraout_RAX_00), cVar2 != '\0')) {
    LOCK();
    uVar1 = puStack0000000000000020[10];
    puStack0000000000000020[10] = 10;
    UNLOCK();
    FUN_007e5760(1,extraout_RBX_01,uVar1);
    FUN_0045a700(0,local_c0,local_d8);
    if (*(int *)*puStack0000000000000020 != 0) {
      FUN_007337c0(&PTR_dWhibQI_ES9pBAp6d___Interface_type_00b418e0,(int *)*puStack0000000000000020,
                   0x2f,extraout_RAX_01,extraout_RBX_02,extraout_RCX_01);
    }
    puVar3 = (undefined8 *)FUN_0049a3c0(3000000000);
    local_d0 = *puVar3;
    local_88[0] = (**(code **)(iStack0000000000000010 + 0x20))(uStack0000000000000018);
    local_78 = local_d0;
    iVar4 = FUN_00452960(local_88,local_e8,0,0,2,1);
    if (iVar4 != 0) {
      while( true ) {
        local_b8 = (undefined8 *)FUN_0049a8a0(puStack0000000000000020[10] * 1000000000);
        local_98 = (**(code **)(iStack0000000000000010 + 0x20))(uStack0000000000000018);
        local_a8[0] = *local_b8;
        iVar4 = FUN_00452960(local_a8,local_f0,0,0,2,1);
        if (iVar4 != 0) break;
        if (*(char *)(local_b8 + 1) != '\0') {
          FUN_00475180();
        }
        FUN_0045a700(0,local_c0,local_d8);
        if (*(int *)*puStack0000000000000020 != 0) {
          FUN_007337c0(&PTR_dWhibQI_ES9pBAp6d___Interface_type_00b418e0,
                       (int *)*puStack0000000000000020,0x2f,extraout_RAX_02,extraout_RBX_03,
                       extraout_RCX_02);
        }
      }
      if (*(char *)(local_b8 + 1) != '\0') {
        FUN_00475180();
      }
      return;
    }
    return;
  }
  local_b0 = *(undefined8 *)(iStack0000000000000008 + 0x18);
  uStack_60 = FUN_0046ebe0(local_c0,local_d8);
  local_68 = &string___String_type;
  uStack_50 = FUN_0046ebe0(local_c8,local_e0);
  local_58 = &string___String_type;
  FUN_0053d720(local_b0,2,"interface ip %s and outband ip %s mismatch",0x2a,&local_68,2,2);
  return;
}
