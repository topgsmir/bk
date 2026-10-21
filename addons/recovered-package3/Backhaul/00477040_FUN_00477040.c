/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

void FUN_00477040(void)

{
  int *extraout_RBX;
  int *extraout_RBX_00;
  int *piVar1;
  undefined8 unaff_RBP;
  int extraout_RDI;
  int extraout_RDI_00;
  int iVar2;
  int in_FS_OFFSET;
  undefined8 unaff_retaddr;
  int in_stack_00000008;

  iVar2 = *(int *)(in_FS_OFFSET + -8);
  piVar1 = *(int **)(iVar2 + 0x30);
  *(undefined8 *)(iVar2 + 0x40) = unaff_retaddr;
  *(int **)(iVar2 + 0x38) = &stack0x00000008;
  *(undefined8 *)(iVar2 + 0x60) = unaff_RBP;
  *(undefined8 *)(iVar2 + 0x50) = 0;
  if (iVar2 == *piVar1) {
    FUN_0047b580();
    FUN_00478d40();
    piVar1 = extraout_RBX;
    iVar2 = extraout_RDI;
  }
  if (iVar2 == piVar1[9]) {
    FUN_0047b5a0();
    FUN_00478d40();
    piVar1 = extraout_RBX_00;
    iVar2 = extraout_RDI_00;
  }
  piVar1[2] = in_stack_00000008;
  piVar1[1] = (int)&stack0x00000010;
  piVar1[3] = iVar2;
  iVar2 = *piVar1;
  *(int *)(in_FS_OFFSET + -8) = iVar2;
  iVar2 = *(int *)(iVar2 + 0x38);
  *(undefined8 *)(iVar2 + -8) = 0x47701d;
  FUN_0047b6e0();
  *(undefined8 *)(iVar2 + -8) = 0x477022;
  FUN_00478d40();
  return;
}
