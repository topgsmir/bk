/* DIAGNOSTIC PSEUDOCODE. Not original source; not directly compilable. */

undefined8 FUN_008f7ea0(void)

{
  byte bVar1;
  char cVar2;
  int *piVar3;
  undefined8 uVar4;
  int iVar5;
  int *piVar6;
  int extraout_RDI;
  int iVar7;
  int iVar8;
  int unaff_R14;
  undefined1 auVar9 [16];

  while (&stack0x00000000 <= *(undefined1 **)(unaff_R14 + 0x10)) {
    FUN_00477040();
  }
  auVar9 = FUN_0055ebc0();
  if (extraout_RDI != 0) {
    return 0;
  }
  do {
    piVar3 = auVar9._0_8_;
    if (auVar9._8_8_ < 1) {
      FUN_004e7c80("no public IP address found",0x1a,0,0,0);
      return 0;
    }
    if ((undefined **)*piVar3 == &PTR_gFfCXg_NzAGozd___Interface_type_00b42d60) {
      piVar6 = (int *)piVar3[1];
      bVar1 = FUN_00560780(*piVar6,piVar6[1],piVar6[2]);
      bVar1 = bVar1 ^ 1;
    }
    else {
      bVar1 = 0;
      piVar6 = (int *)0x0;
    }
    if (bVar1 != 0) {
      iVar5 = piVar6[1];
      iVar7 = *piVar6;
      iVar8 = piVar6[2];
      if (iVar5 != 4) {
        if (iVar5 == 0x10) {
          for (iVar5 = 0; iVar5 < 10; iVar5 = iVar5 + 1) {
            if (*(char *)(iVar7 + iVar5) != '\0') goto LAB_008f7f75;
          }
          if ((*(char *)(iVar7 + 10) == -1) && (*(char *)(iVar7 + 0xb) == -1)) {
            iVar7 = iVar7 + 0xc;
            iVar8 = iVar8 + -0xc;
            iVar5 = 4;
            goto LAB_008f7f80;
          }
        }
LAB_008f7f75:
        iVar7 = 0;
        iVar8 = 0;
        iVar5 = 0;
      }
LAB_008f7f80:
      if ((iVar7 != 0) && (cVar2 = FUN_008f7a20(iVar7,iVar5,iVar8), cVar2 == '\0')) {
        uVar4 = FUN_00560ac0(iVar7,iVar5,iVar8);
        return uVar4;
      }
    }
    auVar9._8_8_ = auVar9._8_8_ + -1;
    auVar9._0_8_ = piVar3 + 2;
  } while( true );
}
