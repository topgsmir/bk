/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

void FUN_00916ec0(int param_1)

{
  undefined8 extraout_RAX;
  undefined8 *puVar1;
  int iVar2;
  int extraout_RCX;
  undefined8 extraout_RBX;
  undefined8 extraout_RDI;
  int unaff_R14;
  int iStack0000000000000008;
  undefined1 local_60 [8];
  undefined8 local_58;
  undefined8 local_50;
  undefined8 local_48;
  undefined8 local_40;
  undefined8 local_38 [2];
  undefined8 local_28;
  undefined8 local_18;
  undefined8 uStack_10;

  iStack0000000000000008 = param_1;
  while (&uStack_10 <= *(undefined8 **)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  FUN_008f7ea0();
  if (extraout_RCX != 0) {
    local_18 = *(undefined8 *)(extraout_RCX + 8);
    uStack_10 = extraout_RDI;
    FUN_0053d8c0(*(undefined8 *)(iStack0000000000000008 + 0x38),"failed to get server ip: %v",0x1b,
                 &local_18,1,1);
    return;
  }
  local_58 = 5000000000;
  local_50 = extraout_RBX;
  local_48 = extraout_RAX;
  while( true ) {
    local_40 = (**(code **)(*(int *)(iStack0000000000000008 + 0x20) + 0x20))
                         (*(undefined8 *)(iStack0000000000000008 + 0x28));
    puVar1 = (undefined8 *)FUN_0049a3c0(local_58);
    local_38[0] = *puVar1;
    local_28 = local_40;
    iVar2 = FUN_00452960(local_38,local_60,0,0,2,1);
    if (iVar2 != 0) break;
    FUN_00917000(iStack0000000000000008,local_48,local_50,&local_58);
  }
  return;
}
