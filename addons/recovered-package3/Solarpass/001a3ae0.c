/* Recovered diagnostic pseudocode. Not original source, not directly compilable. */

/* WARNING: Removing unreachable block (ram,0x001a5f9b) */
/* WARNING: Globals starting with '_' overlap smaller symbols at the same address */

undefined8 recovered_001a3ae0(long *param_1,char *param_2)

{
  undefined1 (*pauVar1) [16];
  undefined2 uVar2;
  short sVar3;
  code *pcVar4;
  uint uVar5;
  long lVar6;
  long lVar7;
  char **ppcVar8;
  char **ppcVar9;
  char **ppcVar10;
  char **ppcVar11;
  char **ppcVar12;
  undefined1 auVar13 [16];
  undefined4 uVar14;
  undefined4 uVar15;
  undefined4 uVar16;
  undefined4 uVar17;
  undefined4 uVar18;
  char cVar19;
  int iVar21;
  uint uVar22;
  uint uVar23;
  byte bVar20;
  undefined8 *puVar24;
  long lVar25;
  char *pcVar26;
  undefined8 uVar27;
  char *pcVar28;
  ulong uVar29;
  char *pcVar30;
  ulong uVar31;
  undefined1 uVar32;
  char *pcVar33;
  byte bVar34;
  byte extraout_DL;
  undefined4 uVar35;
  char **ppcVar36;
  char *pcVar37;
  long lVar38;
  long *plVar39;
  ushort uVar40;
  long lVar41;
  ulong uVar42;
  char *unaff_R12;
  char *pcVar43;
  char *unaff_R13;
  char *unaff_R14;
  long lVar44;
  char *unaff_R15;
  long *in_FS_OFFSET;
  bool bVar45;
  long *plVar46;
  undefined1 auVar47 [16];
  undefined1 auVar48 [16];
  undefined1 auVar49 [16];
  byte bVar50;
  undefined1 auVar51 [16];
  undefined1 auVar52 [16];
  undefined1 auVar53 [16];
  undefined1 auVar54 [12];
  undefined1 auStack_681 [9];
  char *pcStack_678;
  char *pcStack_670;
  char *pcStack_668;
  undefined8 uStack_660;
  undefined8 uStack_658;
  undefined8 uStack_650;
  undefined4 uStack_648;
  undefined4 uStack_644;
  int iStack_640;
  char *pcStack_638;
  char *pcStack_630;
  char *pcStack_628;
  undefined *puStack_620;
  char *pcStack_618;
  char *pcStack_610;
  char *pcStack_608;
  undefined8 uStack_600;
  char *pcStack_5f0;
  char *pcStack_5e8;
  undefined1 *puStack_5e0;
  long lStack_5d8;
  long lStack_5d0;
  char *pcStack_5c8;
  undefined8 *puStack_5c0;
  undefined8 uStack_5b8;
  undefined *puStack_5b0;
  undefined4 uStack_5a8;
  char *pcStack_598;
  code *pcStack_590;
  long lStack_588;
  long lStack_580;
  char *pcStack_578;
  long lStack_570;
  undefined8 uStack_568;
  undefined7 uStack_560;
  undefined1 uStack_559;
  long *plStack_558;
  char **ppcStack_550;
  undefined8 uStack_548;
  char cStack_540;
  undefined1 uStack_53f;
  undefined2 uStack_53e;
  undefined4 uStack_53c;
  undefined2 uStack_538;
  undefined2 uStack_536;
  undefined4 uStack_534;
  undefined2 uStack_530;
  undefined2 uStack_52e;
  undefined4 uStack_52c;
  undefined8 uStack_528;
  char *pcStack_520;
  char *pcStack_518;
  char *pcStack_510;
  undefined8 uStack_508;
  undefined8 uStack_2e0;
  undefined8 uStack_2d8;
  char *pcStack_2d0;
  undefined *puStack_2c8;
  char *pcStack_2c0;
  char *pcStack_2b8;
  char *pcStack_2b0;
  char *pcStack_2a8;
  char *pcStack_2a0;
  char *pcStack_298;
  char *pcStack_290;
  undefined1 uStack_288;
  undefined1 uStack_278;
  undefined1 uStack_230;
  undefined8 uStack_78;
  undefined6 uStack_70;
  undefined2 uStack_6a;
  undefined6 uStack_68;
  char **ppcStack_62;
  long lStack_58;
  long lStack_50;
  undefined8 uStack_48;
  undefined8 uStack_40;

  uVar32 = 0x50;
  pcStack_670 = param_2;
  uVar35 = uStack_648;
  uVar18 = uStack_644;
  iVar21 = iStack_640;
  ppcVar36 = (char **)CONCAT44(uStack_2d8._4_4_,(undefined4)uStack_2d8);
  ppcVar8 = uStack_658;
  ppcVar10 = uStack_650;
  ppcVar12 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
                    /* WARNING: Switch is manually overridden */
  switch(*(undefined1 *)((long)param_1 + 0x6a)) {
  case 0:
    *(undefined1 *)((long)param_1 + 0x6b) = 1;
    lVar44 = param_1[0xc];
    auVar51 = (**(code **)(lVar44 + 0x18))
                        (param_1[0xb] + (*(long *)(lVar44 + 0x10) - 1U & 0xfffffffffffffff0) + 0x10)
    ;
    uVar32 = (undefined1)lVar44;
    *(undefined1 (*) [16])(param_1 + 0xe) = auVar51;
    auVar51 = (**(code **)(auVar51._8_8_ + 0x18))(auVar51._0_8_,param_2);
    break;
  case 1:
    recovered_0005e9f0(&UNK_00a589e8);
    goto LAB_001a53b7;
  case 2:
override_jmp_001a3b0e_case_2:
    iStack_640 = iVar21;
    uStack_644 = uVar18;
    uStack_648 = uVar35;
                    /* WARNING: Does not return */
    pcVar4 = (code *)invalidInstructionException();
    uStack_2d8 = ppcVar36;
    uStack_658 = ppcVar8;
    uStack_650 = ppcVar10;
    uStack_528 = ppcVar12;
    (*pcVar4)();
  case 3:
    auVar51 = (**(code **)(param_1[0xf] + 0x18))(param_1[0xe],param_2);
    break;
  case 4:
    puStack_5e0 = (undefined1 *)((long)param_1 + 0x10c);
    plStack_558 = param_1 + 0xe;
    uVar32 = 100;
    pcVar33 = param_2;
    pcVar28 = unaff_R14;
    ppcVar36 = (char **)CONCAT44(uStack_2d8._4_4_,(undefined4)uStack_2d8);
    ppcVar12 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
                    /* WARNING: Switch is manually overridden */
    switch(*(undefined1 *)((long)param_1 + 0x10c)) {
    case 0:
      puVar24 = (undefined8 *)param_1[0xe];
      uVar2 = (undefined2)param_1[0x21];
      goto LAB_001a3e37;
    case 1:
      goto override_jmp_001a3d0e_case_1;
    case 2:
      goto override_jmp_001a3b0e_case_2;
    case 3:
      goto override_jmp_001a3d0e_case_3;
    case 4:
      goto override_jmp_001a3d0e_case_4;
    case 5:
      goto override_jmp_001a3d0e_case_5;
    case 6:
      goto override_jmp_001a3d0e_case_6;
    case 7:
      pcVar28 = (char *)(param_1 + 0x3a);
      pcVar43 = (char *)(param_1 + 0x23);
      uVar32 = 0x84;
      ppcVar36 = (char **)CONCAT44(uStack_2d8._4_4_,(undefined4)uStack_2d8);
      ppcVar12 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
                    /* WARNING: Switch is manually overridden */
      switch((char)param_1[0x3a]) {
      case '\0':
        lVar25 = param_1[0x29];
        lVar44 = param_1[0x2c];
        lVar38 = param_1[0x23];
        uVar35 = (undefined4)param_1[0x24];
        plVar39 = (long *)param_1[0x2a];
        plVar46 = (long *)param_1[0x2b];
        goto LAB_001a5b25;
      case '\x01':
        goto override_jmp_001a411d_case_1;
      case '\x02':
        goto override_jmp_001a3b0e_case_2;
      case '\x03':
        goto override_jmp_001a411d_case_3;
      case '\x04':
        goto override_jmp_001a411d_case_4;
      case '\x05':
        plVar39 = (long *)param_1[0x3b];
        uVar32 = (undefined1)param_1[0x3c];
        goto LAB_001a5d91;
      case '\x06':
        uVar32 = 0xac;
        ppcVar36 = uStack_2d8;
        ppcVar12 = uStack_528;
                    /* WARNING: Switch is manually overridden */
        switch((char)param_1[0x42]) {
        case '\0':
          plVar39 = (long *)param_1[0x3b];
          lVar44 = param_1[0x3c];
          uVar29 = param_1[0x3e];
          lVar25 = param_1[0x3d];
          goto LAB_001a5e41;
        case '\x02':
          goto override_jmp_001a3b0e_case_2;
        case '\x03':
          bVar20 = *(byte *)(param_1 + 0x45);
          if (7 < bVar20) goto LAB_001a5ed7;
          goto LAB_001a5e7e;
        case '\x04':
          goto override_jmp_001a4256_case_4;
        }
        break;
      case '\a':
        plVar39 = (long *)param_1[0x3b];
        goto LAB_001a5f57;
      case '\b':
        goto override_jmp_001a411d_case_8;
      case '\t':
        auStack_681._1_8_ = pcVar28;
        goto LAB_001a652f;
      }
    }
    goto override_jmp_001a4256_case_1;
  }
  pcVar33 = auVar51._8_8_;
  if (auVar51._0_8_ == 2) {
    uVar27 = 1;
    uVar32 = 3;
    goto LAB_001a5431;
  }
  lVar44 = param_1[0xe];
  puVar24 = (undefined8 *)param_1[0xf];
  pcVar4 = (code *)*puVar24;
  if (pcVar4 != (code *)0x0) {
    (*pcVar4)(lVar44);
  }
  if (puVar24[1] != 0) {
    (*(code *)PTR_DAT_00a8dfe8)(lVar44);
  }
  ppcVar36 = (char **)CONCAT44(uStack_660._4_4_,(undefined4)uStack_660);
  if ((auVar51._0_8_ & 1) == 0) {
    unaff_R12 = (char *)param_1[2];
    if (unaff_R12 == (char *)0x0) {
      unaff_R13 = Elf64_Ehdr_00000000.e_ident_magic_str;
    }
    else {
      lVar44 = param_1[1];
      unaff_R13 = (char *)(*(code *)PTR_DAT_00a8daa8)(unaff_R12);
      if (unaff_R13 == (char *)0x0) goto LAB_001a5400;
      (*(code *)PTR_DAT_00a8da58)(unaff_R13,lVar44,unaff_R12);
    }
    *(undefined1 *)((long)param_1 + 0x6b) = 0;
    lVar44 = param_1[6];
    lVar41 = param_1[7];
    lVar25 = param_1[10];
    uStack_534 = (undefined4)lVar25;
    uStack_530 = (undefined2)((ulong)lVar25 >> 0x20);
    uStack_52e = (undefined2)((ulong)lVar25 >> 0x30);
    lVar25 = param_1[8];
    lVar38 = param_1[9];
    uStack_548._4_4_ = (undefined4)lVar25;
    cStack_540 = (char)((ulong)lVar25 >> 0x20);
    uStack_53f = (undefined1)((ulong)lVar25 >> 0x28);
    uStack_53e = (undefined2)((ulong)lVar25 >> 0x30);
    uStack_53c = (undefined4)lVar38;
    uStack_538 = (undefined2)((ulong)lVar38 >> 0x20);
    uStack_536 = (undefined2)((ulong)lVar38 >> 0x30);
    puVar24 = (undefined8 *)(*(code *)PTR_DAT_00a8daa8)(0x70);
    if (puVar24 == (undefined8 *)0x0) {
      recovered_0005d9a2(8,0x70);
    }
    else {
      *puVar24 = 1;
      puVar24[1] = 1;
      puVar24[2] = unaff_R12;
      puVar24[3] = unaff_R13;
      puVar24[4] = unaff_R12;
      lVar25 = param_1[4];
      puVar24[5] = param_1[3];
      puVar24[6] = lVar25;
      puVar24[7] = param_1[5];
      puVar24[8] = lVar44;
      *(int *)(puVar24 + 9) = (int)lVar41;
      *(ulong *)((long)puVar24 + 0x4c) = CONCAT44(uStack_548._4_4_,(int)uStack_548);
      *(ulong *)((long)puVar24 + 0x54) =
           CONCAT44(uStack_53c,CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540)));
      puVar24[0xb] = CONCAT26(uStack_536,CONCAT24(uStack_538,uStack_53c));
      puVar24[0xc] = CONCAT26(uStack_52e,CONCAT24(uStack_530,uStack_534));
      puVar24[0xd] = pcVar33;
      uVar2 = (undefined2)param_1[0xd];
      param_1[0xe] = (long)puVar24;
      *(undefined2 *)(param_1 + 0x21) = uVar2;
      *(undefined1 *)((long)param_1 + 0x10c) = 0;
LAB_001a3e37:
      param_2 = pcStack_670;
      plStack_558 = param_1 + 0xe;
      puStack_5e0 = (undefined1 *)((long)param_1 + 0x10c);
      param_1[0xf] = (long)puVar24;
      uStack_2e0 = (char *)((long)param_1 + 0x10a);
      *(undefined2 *)((long)param_1 + 0x10a) = uVar2;
      uStack_2d8._0_4_ = 0x31fda0;
      uStack_2d8._4_4_ = 0;
      recovered_002db180(&uStack_548,&UNK_00726372,&uStack_2e0);
      param_1[0x12] = CONCAT44(uStack_534,CONCAT22(uStack_536,uStack_538));
      param_1[0x10] = CONCAT44(uStack_548._4_4_,(int)uStack_548);
      param_1[0x11] = CONCAT44(uStack_53c,CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540)));
      param_1[0x22] = (long)(param_1 + 0x10);
      *(undefined1 *)(param_1 + 0x27) = 0;
override_jmp_001a3d0e_case_3:
      unaff_R12 = (char *)&uStack_548;
      recovered_0017dab0(unaff_R12,param_1 + 0x22,*(long *)param_2);
      iVar21 = (int)uStack_548;
      lVar44 = CONCAT44(uStack_548._4_4_,(int)uStack_548);
      if (lVar44 == -1) {
        *puStack_5e0 = 3;
        pcVar28 = (char *)0xffffffffffffffff;
        uStack_548 = (char *)0xffffffffffffffff;
        goto LAB_001a719a;
      }
      pcVar33 = (char *)CONCAT44(uStack_53c,CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540)));
      uStack_48 = CONCAT44(uStack_534,CONCAT22(uStack_536,uStack_538));
      uStack_40 = CONCAT44(uStack_52c,CONCAT22(uStack_52e,uStack_530));
      if (((char)param_1[0x27] == '\x03') && ((short)param_1[0x23] == -1)) {
        plVar39 = (long *)param_1[0x24];
        LOCK();
        lVar25 = *plVar39;
        if (lVar25 == 0xcc) {
          *plVar39 = 0x84;
        }
        UNLOCK();
        if (lVar25 != 0xcc) {
          (**(code **)(plVar39[2] + 0x20))();
        }
      }
      if (iVar21 == 2) {
        pcVar28 = (char *)recovered_0005d710(pcVar33);
        if (param_1[0x10] != 0) {
          (*(code *)PTR_DAT_00a8dfe8)(param_1[0x11]);
        }
        plVar39 = (long *)param_1[0xf];
        LOCK();
        *plVar39 = *plVar39 + -1;
        UNLOCK();
        if (*plVar39 == 0) {
          recovered_005a83b0(param_1[0xf]);
        }
        *puStack_5e0 = 1;
        bVar45 = false;
        goto LAB_001a719f;
      }
      puVar24 = (undefined8 *)(*(code *)PTR_DAT_00a8daa8)(0x30);
      if (puVar24 != (undefined8 *)0x0) {
        *puVar24 = 1;
        puVar24[1] = 1;
        puVar24[2] = lVar44;
        puVar24[3] = pcVar33;
        puVar24[4] = uStack_48;
        puVar24[5] = uStack_40;
        param_1[0x13] = (long)puVar24;
        recovered_0059dad0(*(undefined4 *)(puVar24 + 5),param_1[0xf] + 0x50);
        ppcVar8 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
        pcVar28 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
        pcVar33 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
        ppcVar36 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
        pcVar43 = uStack_2e0;
        if (3 < _DAT_00a91d28 - 2) {
          bVar20 = DAT_00a8ffb0;
          if (2 < DAT_00a8ffb0) {
            bVar20 = recovered_00086870(&DAT_00a8ffa0);
            ppcVar8 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
            pcVar28 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
          }
          pcVar33 = pcVar28;
          ppcVar36 = ppcVar8;
          pcVar43 = uStack_2e0;
          if (bVar20 != 0) {
            ppcVar36 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
            pcVar33 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
            if (bVar20 != 2) {
              if (_DAT_00a92030 != 2) goto LAB_001a436c;
              lVar44 = _DAT_00a91d18;
              if (_DAT_00a91d10 == 1) {
                lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0)
                         + 0x10;
              }
              cVar19 = (**(code **)(_DAT_00a91d20 + 0x28))(lVar44);
              ppcVar36 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
              pcVar33 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
              pcVar43 = uStack_2e0;
              if (cVar19 == '\0') goto LAB_001a436c;
            }
            pcVar28 = _DAT_00a91d20;
            puStack_620 = _DAT_00a8ffa0 + 0x30;
            uStack_660._0_4_ = 0x73b3b4;
            uStack_660._4_4_ = 0;
            uStack_658._0_4_ = 0x4d;
            uStack_658._4_4_ = 0;
            pcStack_598 = (char *)(param_1 + 0x10);
            pcStack_5c8 = (char *)(param_1[0xf] + 0x10);
            uStack_548 = (char *)&uStack_660;
            cStack_540 = -0x18;
            uStack_53f = 0xa7;
            uStack_53e = 0xa8;
            uStack_53c = 0;
            uStack_538 = SUB82(&pcStack_598,0);
            uStack_536 = (undefined2)((ulong)&pcStack_598 >> 0x10);
            uStack_534 = (undefined4)((ulong)&pcStack_598 >> 0x20);
            uStack_530 = 0x9a70;
            uStack_52e = 0xa5;
            uStack_52c = 0;
            uStack_528 = &pcStack_5c8;
            pcStack_520 = "";
            pcStack_638 = Elf64_Ehdr_00000000.e_ident_magic_str;
            pcStack_628 = Elf64_Ehdr_00000000.e_ident_magic_str + 2;
            pcStack_2d0 = (char *)&pcStack_638;
            puStack_2c8 = _DAT_00a8ffa0;
            uStack_2e0._0_4_ = 1;
            uStack_2e0._4_4_ = 0;
            pcStack_630 = unaff_R12;
            pcVar33 = uStack_548;
            ppcVar36 = uStack_528;
            pcVar43 = Elf64_Ehdr_00000000.e_ident_magic_str;
            if (_DAT_00a92030 == 2) {
              lVar44 = _DAT_00a91d18;
              if (_DAT_00a91d10 == 1) {
                lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0)
                         + 0x10;
              }
              cVar19 = (**(code **)(_DAT_00a91d20 + 0x50))(lVar44,&uStack_2e0);
              pcVar33 = uStack_548;
              ppcVar36 = uStack_528;
              pcVar43 = (char *)CONCAT44(uStack_2e0._4_4_,(int)uStack_2e0);
              if (cVar19 != '\0') {
                (**(code **)(pcVar28 + 0x58))(lVar44,&uStack_2e0);
                pcVar33 = uStack_548;
                ppcVar36 = uStack_528;
                pcVar43 = (char *)CONCAT44(uStack_2e0._4_4_,(int)uStack_2e0);
              }
            }
          }
        }
LAB_001a436c:
        unaff_R14 = *(char **)(param_1[0xf] + 0x68);
        unaff_R15 = unaff_R14 + 0x10;
        LOCK();
        bVar45 = (int)*(long *)(unaff_R14 + 0x10) == 0;
        if (bVar45) {
          *(int *)(unaff_R14 + 0x10) = 1;
        }
        UNLOCK();
        uStack_548 = pcVar33;
        uStack_528 = ppcVar36;
        uStack_2e0 = pcVar43;
        if (!bVar45) {
          recovered_0007d7c0(unaff_R15);
        }
        if ((_DAT_00a92028 & 0x7fffffffffffffff) == 0) {
          uVar29 = 0;
          cVar19 = (char)*(int *)(unaff_R14 + 0x14);
        }
        else {
          uVar23 = recovered_0007d5e0();
          uVar29 = (ulong)uVar23 ^ 1;
          cVar19 = (char)*(int *)(unaff_R14 + 0x14);
        }
        auVar51._8_8_ = uVar29;
        auVar51._0_8_ = unaff_R13;
        auVar13._8_8_ = uVar29;
        auVar13._0_8_ = unaff_R13;
        auVar53._8_8_ = uVar29;
        auVar53._0_8_ = unaff_R13;
        if (cVar19 != '\0') {
          cStack_540 = (char)uVar29;
          uStack_548 = unaff_R15;
          auVar52 = recovered_0005e8b0(&UNK_0078d009,0x2b,&uStack_548,&UNK_00a66208);
          goto LAB_001a457d;
        }
        auVar52 = recovered_0039add0(unaff_R14 + 0x20);
        auVar51 = auVar53;
        if (((char)uVar29 == '\0') && (auVar51 = auVar13, (_DAT_00a92028 & 0x7fffffffffffffff) != 0)
           ) goto LAB_001a5243;
LAB_001a43dc:
        ppcVar36 = auVar52._8_8_;
        LOCK();
        iVar21 = *(int *)unaff_R15;
        *(int *)unaff_R15 = 0;
        UNLOCK();
        if (iVar21 == 2) {
LAB_001a457d:
          ppcVar36 = auVar52._8_8_;
          (*(code *)PTR_DAT_00a8d9a8)(0xca,unaff_R15,0x81,1);
          if ((auVar52._0_8_ & 1) != 0) goto LAB_001a43f2;
LAB_001a45ab:
          if ((_DAT_00a91d28 & 0xfffffffffffffffe) == 4) goto LAB_001a465b;
          bVar20 = DAT_00a8ffc8;
          if (2 < DAT_00a8ffc8) {
            bVar20 = recovered_00086870(&DAT_00a8ffb8);
          }
          if (bVar20 == 0) goto LAB_001a465b;
          if (bVar20 != 2) {
            if (_DAT_00a92030 != 2) goto LAB_001a465b;
            lVar44 = _DAT_00a91d18;
            if (_DAT_00a91d10 == 1) {
              lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                       0x10;
            }
            cVar19 = (**(code **)(_DAT_00a91d20 + 0x28))(lVar44,_DAT_00a8ffb8);
            if (cVar19 == '\0') goto LAB_001a465b;
          }
          lVar44 = _DAT_00a8ffb8 + 0x30;
          uStack_2e0._0_4_ = 0x73b3da;
          uStack_2d8._0_4_ = 0x8b;
          uStack_2d8._4_4_ = 0;
          pcStack_638 = (char *)&uStack_2e0;
          pcStack_630 = "";
          uStack_548 = Elf64_Ehdr_00000000.e_ident_magic_str;
          cStack_540 = (char)&pcStack_638;
          uStack_53f = (undefined1)((ulong)&pcStack_638 >> 8);
          uStack_53e = (undefined2)((ulong)&pcStack_638 >> 0x10);
          uStack_53c = (undefined4)((ulong)&pcStack_638 >> 0x20);
          uStack_538 = 1;
          uStack_536 = 0;
          uStack_534 = 0;
          uStack_530 = (undefined2)lVar44;
          uStack_52e = (undefined2)((ulong)lVar44 >> 0x10);
          uStack_52c = (undefined4)((ulong)lVar44 >> 0x20);
          puVar24 = &uStack_548;
          lVar44 = _DAT_00a8ffb8;
        }
        else {
          if ((auVar52._0_8_ & 1) == 0) goto LAB_001a45ab;
LAB_001a43f2:
          uStack_660 = ppcVar36;
          if (_DAT_00a91d28 - 2 < 4) goto LAB_001a465b;
          bVar20 = DAT_00a90028;
          if (2 < DAT_00a90028) {
            bVar20 = recovered_00086870(&DAT_00a90018);
          }
          if (bVar20 == 0) goto LAB_001a465b;
          if (bVar20 != 2) {
            if (_DAT_00a92030 != 2) goto LAB_001a465b;
            lVar44 = _DAT_00a91d18;
            if (_DAT_00a91d10 == 1) {
              lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                       0x10;
            }
            cVar19 = (**(code **)(_DAT_00a91d20 + 0x28))(lVar44,_DAT_00a90018);
            if (cVar19 == '\0') goto LAB_001a465b;
          }
          puStack_2c8 = (undefined *)(_DAT_00a90018 + 0x30);
          pcStack_638 = &UNK_0073b41f;
          pcStack_630 = (char *)((long)&Elf64_Phdr_ARRAY_00000040[0].p_filesz + 3);
          uStack_548 = (char *)&pcStack_638;
          cStack_540 = -0x18;
          uStack_53f = 0xa7;
          uStack_53e = 0xa8;
          uStack_53c = 0;
          uStack_538 = SUB82(&uStack_660,0);
          uStack_536 = (undefined2)((ulong)&uStack_660 >> 0x10);
          uStack_534 = (undefined4)((ulong)&uStack_660 >> 0x20);
          uStack_530 = 0xb878;
          uStack_52e = 0xa7;
          uStack_52c = 0;
          uStack_2e0._0_4_ = 1;
          uStack_2d8._0_4_ = SUB84(unaff_R12,0);
          uStack_2d8._4_4_ = (undefined4)((ulong)unaff_R12 >> 0x20);
          pcStack_2d0 = Elf64_Ehdr_00000000.e_ident_magic_str + 1;
          puVar24 = &uStack_2e0;
          lVar44 = _DAT_00a90018;
        }
        uStack_2e0._4_4_ = 0;
        recovered_00679790(lVar44,puVar24);
LAB_001a465b:
        auVar54 = recovered_0062c2a0(1);
        *(undefined1 (*) [12])(param_1 + 0x14) = auVar54;
        if ((char)in_FS_OFFSET[-2] != '\x01') goto LAB_001a5186;
        auVar53 = *(undefined1 (*) [16])(in_FS_OFFSET + -4);
        do {
          pcVar28 = auVar51._0_8_;
          unaff_R15 = auVar53._8_8_;
          unaff_R14 = (char *)(auVar53._0_8_ + 1);
          in_FS_OFFSET[-4] = (long)unaff_R14;
          puVar24 = (undefined8 *)(*(code *)PTR_DAT_00a8daa8)(0x68);
          if (puVar24 == (undefined8 *)0x0) {
LAB_001a5168:
            recovered_0005d9a2(8,0x68);
          }
          else {
            *puVar24 = 1;
            puVar24[1] = 1;
            puVar24[2] = 0;
            puVar24[3] = 0;
            *(undefined8 *)((long)puVar24 + 0x19) = 0;
            *(undefined8 *)((long)puVar24 + 0x21) = 0;
            puVar24[6] = 2;
            uVar27 = _UNK_00a8ac28;
            puVar24[7] = _DAT_00a8ac20;
            puVar24[8] = uVar27;
            uVar27 = _UNK_00a8ac38;
            puVar24[9] = _DAT_00a8ac30;
            puVar24[10] = uVar27;
            *(undefined1 (*) [16])(puVar24 + 0xb) = auVar53;
            param_1[0x16] = (long)puVar24;
            pcVar43 = (char *)(auVar53._0_8_ + 2);
            in_FS_OFFSET[-4] = (long)pcVar43;
            puVar24 = (undefined8 *)(*(code *)PTR_DAT_00a8daa8)(0x68);
            if (puVar24 == (undefined8 *)0x0) goto LAB_001a5168;
            *puVar24 = 1;
            puVar24[1] = 1;
            puVar24[2] = 0;
            puVar24[3] = 0;
            *(undefined8 *)((long)puVar24 + 0x19) = 0;
            *(undefined8 *)((long)puVar24 + 0x21) = 0;
            puVar24[6] = 2;
            uVar27 = _UNK_00a8ac28;
            puVar24[7] = _DAT_00a8ac20;
            puVar24[8] = uVar27;
            uVar27 = _UNK_00a8ac38;
            puVar24[9] = _DAT_00a8ac30;
            puVar24[10] = uVar27;
            puVar24[0xb] = unaff_R14;
            puVar24[0xc] = unaff_R15;
            param_1[0x17] = (long)puVar24;
            puVar24 = (undefined8 *)(*(code *)PTR_DAT_00a8daa8)(0x18);
            if (puVar24 == (undefined8 *)0x0) goto LAB_001a73db;
            *puVar24 = 1;
            puVar24[1] = 1;
            puVar24[2] = 1;
            param_1[0x18] = (long)puVar24;
            unaff_R15 = *(char **)(param_1[0xf] + 0x68);
            recovered_00352720(unaff_R15);
            ppcVar11 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
            ppcVar9 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
            unaff_R14 = (char *)param_1[0x13];
            LOCK();
            lVar44 = *(long *)unaff_R14;
            *(long *)unaff_R14 = *(long *)unaff_R14 + 1;
            UNLOCK();
            uVar35 = uStack_648;
            uVar18 = uStack_644;
            iVar21 = iStack_640;
            ppcVar36 = uStack_2d8;
            ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
            ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
            ppcVar12 = uStack_528;
            if (*(long *)unaff_R14 == 0 || SCARRY8(lVar44,1) != *(long *)unaff_R14 < 0)
            goto override_jmp_001a3b0e_case_2;
            unaff_R12 = (char *)param_1[0x17];
            LOCK();
            lVar44 = *(long *)unaff_R12;
            *(long *)unaff_R12 = *(long *)unaff_R12 + 1;
            UNLOCK();
            ppcVar8 = ppcVar9;
            ppcVar10 = ppcVar11;
            if (*(long *)unaff_R12 == 0 || SCARRY8(lVar44,1) != *(long *)unaff_R12 < 0)
            goto override_jmp_001a3b0e_case_2;
            pcStack_678 = (char *)param_1[0x14];
            ppcVar36 = _DAT_00a91850;
            do {
              _DAT_00a91850 = ppcVar36;
              LOCK();
              ppcVar36 = (char **)((long)_DAT_00a91850 + 1);
              UNLOCK();
            } while (_DAT_00a91850 == (char **)0x0);
            uStack_2e0 = (char *)&uStack_660;
            uStack_2d8 = (char **)auStack_681;
            auStack_681._1_4_ = (int)param_1[0x15];
            puStack_2c8 = (undefined *)CONCAT44(puStack_2c8._4_4_,(int)param_1[0x15]);
            uStack_230 = 0;
            param_2 = (char *)(*in_FS_OFFSET + -0x1c0);
            uStack_660 = _DAT_00a91850;
            pcStack_2d0 = pcStack_678;
            pcStack_2c0 = unaff_R15;
            pcStack_2b8 = unaff_R12;
            pcStack_2b0 = unaff_R14;
            if ((char)in_FS_OFFSET[-0x2f] == '\0') {
              uVar29 = in_FS_OFFSET[-0x38];
            }
            else {
              _DAT_00a91850 = ppcVar36;
              if ((char)in_FS_OFFSET[-0x2f] != '\x01') {
                recovered_000a4da0(&pcStack_2d0);
                uVar32 = 1;
                pcVar33 = uStack_2e0;
                goto LAB_001a5277;
              }
              recovered_0062acc0(param_2,recovered_00656300);
              ppcVar11 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
              ppcVar9 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
              *(undefined1 *)(in_FS_OFFSET + -0x2f) = 0;
              uVar29 = in_FS_OFFSET[-0x38];
              ppcVar36 = _DAT_00a91850;
            }
            _DAT_00a91850 = ppcVar36;
            if (0x7ffffffffffffffe < uVar29) goto LAB_001a75eb;
            in_FS_OFFSET[-0x38] = uVar29 + 1;
            lVar44 = in_FS_OFFSET[-0x37];
            (*(code *)PTR_DAT_00a8da58)(&uStack_548,&uStack_2e0,0x268);
            if (lVar44 == 2) goto LAB_001a53b7;
            pcStack_5e8 = *(char **)uStack_548;
            uVar35 = uStack_648;
            uVar18 = uStack_644;
            iVar21 = iStack_640;
            pcStack_5f0 = unaff_R14;
            pcStack_578 = param_2;
            ppcVar36 = uStack_2d8;
            ppcVar12 = uStack_528;
            if ((int)lVar44 == 1) {
              pcVar28 = (char *)in_FS_OFFSET[-0x36];
              LOCK();
              lVar44 = *(long *)pcVar28;
              *(long *)pcVar28 = *(long *)pcVar28 + 1;
              UNLOCK();
              ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
              ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
              if (*(long *)pcVar28 == 0 || SCARRY8(lVar44,1) != *(long *)pcVar28 < 0)
              goto override_jmp_001a3b0e_case_2;
              plVar39 = *(long **)(pcVar28 + 0x210);
              if (plVar39 == (long *)0x0) {
                lStack_5d8 = 0;
                lStack_5d0 = 0;
              }
              else {
                LOCK();
                lVar44 = *plVar39;
                *plVar39 = *plVar39 + 1;
                UNLOCK();
                ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
                ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
                if (*plVar39 == 0 || SCARRY8(lVar44,1) != *plVar39 < 0)
                goto override_jmp_001a3b0e_case_2;
                lStack_5d8 = *(long *)(pcVar28 + 0x210);
                lStack_5d0 = *(long *)(pcVar28 + 0x218);
              }
              pcStack_638 = (char *)0x0;
              iVar21 = (*(code *)PTR_DAT_00a8dfb8)(&pcStack_638,0x80,0x300);
              pcVar33 = pcStack_5e8;
              pcVar43 = (char *)0x0;
              if (iVar21 == 0) {
                pcVar43 = pcStack_638;
              }
              auVar51._8_8_ = 0;
              auVar51._0_8_ = pcVar28;
              if (pcVar43 != (char *)0x0) goto code_r0x001a499c;
            }
            else {
              pcVar28 = (char *)in_FS_OFFSET[-0x36];
              LOCK();
              lVar44 = *(long *)pcVar28;
              *(long *)pcVar28 = *(long *)pcVar28 + 1;
              UNLOCK();
              ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
              ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
              if (*(long *)pcVar28 == 0 || SCARRY8(lVar44,1) != *(long *)pcVar28 < 0)
              goto override_jmp_001a3b0e_case_2;
              plVar39 = *(long **)(pcVar28 + 0x210);
              if (plVar39 == (long *)0x0) {
                lStack_5d8 = 0;
                lStack_5d0 = 0;
              }
              else {
                LOCK();
                lVar44 = *plVar39;
                *plVar39 = *plVar39 + 1;
                UNLOCK();
                ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
                ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
                if (*plVar39 == 0 || SCARRY8(lVar44,1) != *plVar39 < 0)
                goto override_jmp_001a3b0e_case_2;
                lStack_5d8 = *(long *)(pcVar28 + 0x210);
                lStack_5d0 = *(long *)(pcVar28 + 0x218);
              }
              pcStack_638 = (char *)0x0;
              iVar21 = (*(code *)PTR_DAT_00a8dfb8)(&pcStack_638,0x80,0x300);
              pcVar33 = pcStack_5e8;
              pcVar43 = (char *)0x0;
              if (iVar21 == 0) {
                pcVar43 = pcStack_638;
              }
              auVar51._8_8_ = 0;
              auVar51._0_8_ = pcVar28;
              if (pcVar43 != (char *)0x0) {
                pcStack_638[0] = -0x34;
                pcStack_638[1] = '\0';
                pcStack_638[2] = '\0';
                pcStack_638[3] = '\0';
                pcStack_638[4] = '\0';
                pcStack_638[5] = '\0';
                pcStack_638[6] = '\0';
                pcStack_638[7] = '\0';
                *(long *)(pcStack_638 + 8) = 0;
                *(undefined **)(pcStack_638 + 0x10) = &UNK_00a5cd68;
                *(long *)(pcStack_638 + 0x18) = 0;
                *(char **)(pcStack_638 + 0x20) = pcVar28;
                *(char **)(pcStack_638 + 0x28) = pcStack_5e8;
                *(int *)(pcStack_638 + 0x30) = 0;
                *(char **)(pcStack_638 + 0x38) = pcStack_678;
                *(undefined4 *)(pcStack_638 + 0x40) = auStack_681._1_4_;
                *(char **)(pcStack_638 + 0x48) = unaff_R15;
                *(char **)(pcStack_638 + 0x50) = unaff_R12;
                *(char **)(pcStack_638 + 0x58) = pcStack_5f0;
                pcStack_638[0xd8] = '\0';
                *(undefined1 (*) [16])(pcStack_638 + 0x290) = (undefined1  [16])0x0;
                *(long *)(pcStack_638 + 0x2a0) = 0;
                *(long *)(pcStack_638 + 0x2b0) = lStack_5d8;
                *(long *)(pcStack_638 + 0x2b8) = lStack_5d0;
                lVar44 = recovered_00543360(pcVar28 + 0x160,pcStack_638);
                pcStack_638 = pcVar33;
                if (*(long *)(pcVar28 + 0x200) != 0) {
                  (**(code **)(*(long *)(pcVar28 + 0x208) + 0x28))
                            (*(long *)(pcVar28 + 0x200) +
                             (*(long *)(*(long *)(pcVar28 + 0x208) + 0x10) - 1U & 0xfffffffffffffff0
                             ) + 0x10,&pcStack_638);
                }
                if (lVar44 != 0) {
                  recovered_006679b0(in_FS_OFFSET[-0x36],lVar44);
                }
                goto LAB_001a4c0b;
              }
            }
          }
          recovered_0005d9a2(0x80,0x300);
LAB_001a5186:
          auVar53 = recovered_0063fd60();
          in_FS_OFFSET[-3] = auVar53._8_8_;
          *(undefined1 *)(in_FS_OFFSET + -2) = 1;
        } while( true );
      }
    }
    recovered_0005d9a2(8,0x30);
LAB_001a42ed:
    bVar20 = recovered_00086870(&DAT_00a8f790);
    ppcVar36 = (char **)CONCAT44(uStack_660._4_4_,(undefined4)uStack_660);
joined_r0x001a42fb:
    ppcVar8 = ppcVar36;
    if (bVar20 != 0) {
      if (bVar20 == 2) {
LAB_001a3bc1:
        pcVar28 = _DAT_00a91d20;
        puStack_620 = (undefined *)(_DAT_00a8f790 + 0x30);
        uStack_660 = &pcStack_598;
        uStack_658._0_4_ = 0x31fcd0;
        uStack_658._4_4_ = 0;
        pcStack_5c8 = &UNK_00725960;
        puStack_5c0 = &uStack_660;
        pcStack_2d0 = (char *)(param_1 + 0xd);
        uStack_2e0 = (char *)&pcStack_5c8;
        uStack_2d8._0_4_ = 0xa8a7e8;
        uStack_2d8._4_4_ = 0;
        puStack_2c8 = &UNK_00a67980;
        pcStack_638 = Elf64_Ehdr_00000000.e_ident_magic_str;
        pcStack_630 = (char *)&uStack_2e0;
        pcStack_628 = Elf64_Ehdr_00000000.e_ident_magic_str + 1;
        uStack_538 = SUB82(&pcStack_638,0);
        uStack_536 = (undefined2)((ulong)&pcStack_638 >> 0x10);
        uStack_534 = (undefined4)((ulong)&pcStack_638 >> 0x20);
        uStack_530 = (undefined2)_DAT_00a8f790;
        uStack_52e = (undefined2)((ulong)_DAT_00a8f790 >> 0x10);
        uStack_52c = (undefined4)((ulong)_DAT_00a8f790 >> 0x20);
        uStack_548._0_4_ = 1;
        uStack_548._4_4_ = 0;
        pcVar33 = pcStack_598;
        ppcVar8 = uStack_660;
        if (_DAT_00a92030 == 2) {
          lVar44 = _DAT_00a91d18;
          if (_DAT_00a91d10 == 1) {
            lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                     0x10;
          }
          cVar19 = (**(code **)(_DAT_00a91d20 + 0x50))(lVar44,&uStack_548);
          pcVar33 = pcStack_598;
          ppcVar8 = uStack_660;
          if (cVar19 != '\0') {
            (**(code **)(pcVar28 + 0x58))(lVar44,&uStack_548);
            pcVar33 = pcStack_598;
            ppcVar8 = uStack_660;
          }
        }
      }
      else {
        ppcVar8 = (char **)CONCAT44(uStack_660._4_4_,(undefined4)uStack_660);
        if (_DAT_00a92030 == 2) {
          lVar44 = _DAT_00a91d18;
          if (_DAT_00a91d10 == 1) {
            lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                     0x10;
          }
          cVar19 = (**(code **)(_DAT_00a91d20 + 0x28))(lVar44);
          ppcVar8 = (char **)CONCAT44(uStack_660._4_4_,(undefined4)uStack_660);
          if (cVar19 != '\0') goto LAB_001a3bc1;
        }
      }
    }
  }
  else {
    pcStack_598 = pcVar33;
    ppcVar8 = (char **)CONCAT44(uStack_660._4_4_,(undefined4)uStack_660);
    if (_DAT_00a91d28 != 5) {
      uStack_2e0 = (char *)CONCAT44(uStack_2e0._4_4_,(int)uStack_2e0);
      bVar20 = DAT_00a8f7a0;
      if (2 < DAT_00a8f7a0) goto LAB_001a42ed;
      goto joined_r0x001a42fb;
    }
  }
  uStack_660 = ppcVar8;
  if ((1 < ((uint)pcVar33 & 3) - 2) && (((ulong)pcVar33 & 3) != 0)) {
    (**(code **)(pcVar33 + 0x17))((Elf64_Ehdr *)(pcVar33 + -1));
  }
  plVar39 = (long *)param_1[0xb];
  LOCK();
  *plVar39 = *plVar39 + -1;
  UNLOCK();
  if (*plVar39 == 0) {
    recovered_00655a20(param_1[0xb],param_1[0xc]);
  }
  if (*param_1 != 0) {
    (*(code *)PTR_DAT_00a8dfe8)(param_1[1]);
  }
  lVar25 = 0x20;
  lVar44 = param_1[3];
  goto joined_r0x001a7733;
LAB_001a53b7:
  recovered_000a4da0(&uStack_538);
  in_FS_OFFSET[-0x38] = in_FS_OFFSET[-0x38] + -1;
  uVar32 = 0;
  pcVar33 = uStack_2e0;
LAB_001a5277:
  uStack_2e0._4_4_ = (undefined4)((ulong)pcVar33 >> 0x20);
  uStack_2e0._1_3_ = (undefined3)((ulong)pcVar33 >> 8);
  uStack_2e0._0_4_ = CONCAT31(uStack_2e0._1_3_,uVar32);
  uStack_548 = (char *)&uStack_2e0;
  cStack_540 = -0x80;
  uStack_53f = 0x85;
  uStack_53e = 0x65;
  uStack_53c = 0;
  puVar24 = &uStack_548;
  recovered_0005e620(&UNK_00729475,puVar24,&UNK_00a59d50);
  ppcVar36 = _DAT_00a91850;
LAB_001a52b8:
  _DAT_00a91850 = ppcVar36;
  pcVar28 = (char *)((ulong)puVar24 & 0xffffffff);
  uVar23 = (uint)puVar24;
  recovered_0062acc0(pcStack_578,recovered_00656300);
  *(undefined1 *)(in_FS_OFFSET + -0x2f) = 0;
  ppcVar36 = _DAT_00a91850;
LAB_001a4cf9:
  _DAT_00a91850 = ppcVar36;
  auStack_681._1_4_ = uVar23;
  ppcVar9 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
  ppcVar11 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
  if (0x7ffffffffffffffe < (ulong)in_FS_OFFSET[-0x38]) goto LAB_001a75eb;
  in_FS_OFFSET[-0x38] = in_FS_OFFSET[-0x38] + 1;
  lVar44 = in_FS_OFFSET[-0x37];
  (*(code *)PTR_DAT_00a8da58)(&uStack_548,&uStack_2e0,0xe8);
  if (lVar44 != 2) goto code_r0x001a4d53;
  recovered_000a5060(&uStack_538);
  in_FS_OFFSET[-0x38] = in_FS_OFFSET[-0x38] + -1;
  uVar32 = 0;
  pcVar33 = uStack_2e0;
LAB_001a5202:
  uStack_2e0._4_4_ = (undefined4)((ulong)pcVar33 >> 0x20);
  uStack_2e0._1_3_ = (undefined3)((ulong)pcVar33 >> 8);
  uStack_2e0._0_4_ = CONCAT31(uStack_2e0._1_3_,uVar32);
  uStack_548 = (char *)&uStack_2e0;
  cStack_540 = -0x80;
  uStack_53f = 0x85;
  uStack_53e = 0x65;
  uStack_53c = 0;
  auVar52 = recovered_0005e620(&UNK_00729475,&uStack_548,&UNK_00a59d68);
LAB_001a5243:
  cVar19 = recovered_0007d5e0();
  auVar51 = auVar52;
  if (cVar19 == '\0') {
    unaff_R14[0x14] = '\x01';
  }
  goto LAB_001a43dc;
code_r0x001a499c:
  pcStack_638[0] = -0x34;
  pcStack_638[1] = '\0';
  pcStack_638[2] = '\0';
  pcStack_638[3] = '\0';
  pcStack_638[4] = '\0';
  pcStack_638[5] = '\0';
  pcStack_638[6] = '\0';
  pcStack_638[7] = '\0';
  *(long *)(pcStack_638 + 8) = 0;
  *(undefined **)(pcStack_638 + 0x10) = &UNK_00a5cdb8;
  *(long *)(pcStack_638 + 0x18) = 0;
  *(char **)(pcStack_638 + 0x20) = pcVar28;
  *(char **)(pcStack_638 + 0x28) = pcStack_5e8;
  *(int *)(pcStack_638 + 0x30) = 0;
  *(char **)(pcStack_638 + 0x38) = pcStack_678;
  *(undefined4 *)(pcStack_638 + 0x40) = auStack_681._1_4_;
  *(char **)(pcStack_638 + 0x48) = unaff_R15;
  *(char **)(pcStack_638 + 0x50) = unaff_R12;
  *(char **)(pcStack_638 + 0x58) = pcStack_5f0;
  pcStack_638[0xd8] = '\0';
  *(undefined1 (*) [16])(pcStack_638 + 0x290) = (undefined1  [16])0x0;
  *(long *)(pcStack_638 + 0x2a0) = 0;
  *(long *)(pcStack_638 + 0x2b0) = lStack_5d8;
  *(long *)(pcStack_638 + 0x2b8) = lStack_5d0;
  pcVar26 = (char *)recovered_00543360(pcVar28 + 0x98,pcStack_638);
  pcStack_638 = pcVar33;
  if (*(long *)(pcVar28 + 0x200) != 0) {
    (**(code **)(*(long *)(pcVar28 + 0x208) + 0x28))
              (*(long *)(pcVar28 + 0x200) +
               (*(long *)(*(long *)(pcVar28 + 0x208) + 0x10) - 1U & 0xfffffffffffffff0) + 0x10,
               &pcStack_638);
  }
  if (pcVar26 != (char *)0x0) {
    pcVar28 = pcVar28 + 0x10;
    pcStack_5c8 = (char *)((ulong)pcStack_5c8 & 0xffffffffffffff00);
    if ((char)in_FS_OFFSET[-0x2f] == '\0') {
LAB_001a4a93:
      if ((*(char *)((long)in_FS_OFFSET + -0x17a) == '\x02') ||
         (pcVar30 = (char *)in_FS_OFFSET[-0x33], pcVar30 == (char *)0x0)) {
        pcVar33 = (char *)0x0;
      }
      else {
        pcVar33 = (char *)0x0;
        if (*pcVar30 != '\0') {
          pcVar33 = pcVar30 + 8;
        }
      }
    }
    else {
      if ((char)in_FS_OFFSET[-0x2f] != '\x02') {
        recovered_0062acc0(pcStack_578,recovered_00656300);
        *(undefined1 *)(in_FS_OFFSET + -0x2f) = 0;
        goto LAB_001a4a93;
      }
      pcVar33 = (char *)0x0;
    }
    pcStack_628 = (char *)&pcStack_5c8;
    pcStack_638 = pcVar28;
    pcStack_630 = pcVar26;
    recovered_00664f00(&pcStack_638,pcVar33);
  }
LAB_001a4c0b:
  in_FS_OFFSET[-0x38] = in_FS_OFFSET[-0x38] + -1;
  LOCK();
  lVar44 = *(long *)pcVar43;
  if (lVar44 == 0xcc) {
    *(long *)pcVar43 = 0x84;
  }
  UNLOCK();
  if (lVar44 != 0xcc) {
    (**(code **)(*(long *)(pcVar43 + 0x10) + 0x20))(pcVar43);
  }
  unaff_R14 = (char *)param_1[0x16];
  LOCK();
  lVar44 = *(long *)unaff_R14;
  *(long *)unaff_R14 = *(long *)unaff_R14 + 1;
  UNLOCK();
  uVar35 = uStack_648;
  uVar18 = uStack_644;
  iVar21 = iStack_640;
  ppcVar36 = uStack_2d8;
  ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
  ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
  ppcVar12 = uStack_528;
  if (*(long *)unaff_R14 == 0 || SCARRY8(lVar44,1) != *(long *)unaff_R14 < 0)
  goto override_jmp_001a3b0e_case_2;
  unaff_R15 = (char *)param_1[0x17];
  LOCK();
  lVar44 = *(long *)unaff_R15;
  *(long *)unaff_R15 = *(long *)unaff_R15 + 1;
  UNLOCK();
  ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
  ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
  if (*(long *)unaff_R15 == 0 || SCARRY8(lVar44,1) != *(long *)unaff_R15 < 0)
  goto override_jmp_001a3b0e_case_2;
  param_2 = (char *)param_1[0x14];
  uVar23 = *(uint *)(param_1 + 0x15);
  puVar24 = (undefined8 *)(ulong)uVar23;
  ppcVar36 = _DAT_00a91850;
  do {
    _DAT_00a91850 = ppcVar36;
    LOCK();
    ppcVar36 = (char **)((long)_DAT_00a91850 + 1);
    UNLOCK();
  } while (_DAT_00a91850 == (char **)0x0);
  unaff_R12 = (char *)((ulong)*(uint *)(param_1[0xf] + 0x48) / 1000000 +
                      *(long *)(param_1[0xf] + 0x40) * 1000);
  uStack_2e0 = (char *)&uStack_660;
  uStack_2d8 = (char **)auStack_681;
  puStack_2c8 = (undefined *)CONCAT44(puStack_2c8._4_4_,uVar23);
  uStack_278 = 0;
  pcStack_2d0 = param_2;
  pcStack_2c0 = unaff_R15;
  pcStack_2b8 = unaff_R14;
  pcStack_2b0 = unaff_R12;
  uStack_660 = _DAT_00a91850;
  if ((char)in_FS_OFFSET[-0x2f] != '\0') {
    if ((char)in_FS_OFFSET[-0x2f] == '\x02') {
      _DAT_00a91850 = ppcVar36;
      recovered_000a5060(&pcStack_2d0);
      uVar32 = 1;
      pcVar33 = uStack_2e0;
      goto LAB_001a5202;
    }
    goto LAB_001a52b8;
  }
  goto LAB_001a4cf9;
code_r0x001a4d53:
  pcStack_5e8 = *(char **)uStack_548;
  uVar35 = uStack_648;
  uVar18 = uStack_644;
  iVar21 = iStack_640;
  pcStack_678 = param_2;
  pcStack_5f0 = unaff_R14;
  ppcVar36 = uStack_2d8;
  ppcVar12 = uStack_528;
  if ((int)lVar44 == 1) {
    pcVar33 = (char *)in_FS_OFFSET[-0x36];
    LOCK();
    lVar44 = *(long *)pcVar33;
    *(long *)pcVar33 = *(long *)pcVar33 + 1;
    UNLOCK();
    ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
    ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
    if (*(long *)pcVar33 == 0 || SCARRY8(lVar44,1) != *(long *)pcVar33 < 0)
    goto override_jmp_001a3b0e_case_2;
    plVar39 = *(long **)(pcVar33 + 0x210);
    if (plVar39 == (long *)0x0) {
      lStack_5d8 = 0;
      lStack_5d0 = 0;
    }
    else {
      LOCK();
      lVar44 = *plVar39;
      *plVar39 = *plVar39 + 1;
      UNLOCK();
      ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
      ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
      if (*plVar39 == 0 || SCARRY8(lVar44,1) != *plVar39 < 0) goto override_jmp_001a3b0e_case_2;
      lStack_5d8 = *(long *)(pcVar33 + 0x210);
      lStack_5d0 = *(long *)(pcVar33 + 0x218);
    }
    pcStack_638 = (char *)0x0;
    iVar21 = (*(code *)PTR_DAT_00a8dfb8)(&pcStack_638,0x80,0x180);
    unaff_R14 = pcStack_5e8;
    pcVar28 = (char *)0x0;
    if (iVar21 == 0) {
      pcVar28 = pcStack_638;
    }
    if (pcVar28 == (char *)0x0) goto LAB_001a73bf;
    pcStack_638[0] = -0x34;
    pcStack_638[1] = '\0';
    pcStack_638[2] = '\0';
    pcStack_638[3] = '\0';
    pcStack_638[4] = '\0';
    pcStack_638[5] = '\0';
    pcStack_638[6] = '\0';
    pcStack_638[7] = '\0';
    *(long *)(pcStack_638 + 8) = 0;
    *(undefined **)(pcStack_638 + 0x10) = &UNK_00a5ce58;
    *(long *)(pcStack_638 + 0x18) = 0;
    *(char **)(pcStack_638 + 0x20) = pcVar33;
    *(char **)(pcStack_638 + 0x28) = pcStack_5e8;
    *(int *)(pcStack_638 + 0x30) = 0;
    *(char **)(pcStack_638 + 0x38) = pcStack_678;
    *(undefined4 *)(pcStack_638 + 0x40) = auStack_681._1_4_;
    *(char **)(pcStack_638 + 0x48) = unaff_R15;
    *(char **)(pcStack_638 + 0x50) = pcStack_5f0;
    *(char **)(pcStack_638 + 0x58) = unaff_R12;
    pcStack_638[0x90] = '\0';
    *(undefined1 (*) [16])(pcStack_638 + 0x110) = (undefined1  [16])0x0;
    *(long *)(pcStack_638 + 0x120) = 0;
    *(long *)(pcStack_638 + 0x130) = lStack_5d8;
    *(long *)(pcStack_638 + 0x138) = lStack_5d0;
    pcVar26 = (char *)recovered_00543360(pcVar33 + 0x98,pcStack_638);
    pcStack_638 = unaff_R14;
    if (*(long *)(pcVar33 + 0x200) != 0) {
      (**(code **)(*(long *)(pcVar33 + 0x208) + 0x28))
                (*(long *)(pcVar33 + 0x200) +
                 (*(long *)(*(long *)(pcVar33 + 0x208) + 0x10) - 1U & 0xfffffffffffffff0) + 0x10,
                 &pcStack_638);
    }
    pcVar43 = pcStack_670;
    if (pcVar26 != (char *)0x0) {
      pcStack_5c8 = (char *)((ulong)pcStack_5c8 & 0xffffffffffffff00);
      if ((char)in_FS_OFFSET[-0x2f] != '\0') {
        if ((char)in_FS_OFFSET[-0x2f] == '\x02') {
          pcStack_628 = (char *)&pcStack_5c8;
          pcStack_638 = pcVar33 + 0x10;
          pcStack_630 = pcVar26;
          recovered_00664f00(&pcStack_638,0);
          pcVar43 = pcStack_670;
          goto LAB_001a5108;
        }
        recovered_0062acc0(pcStack_578,recovered_00656300);
        *(undefined1 *)(in_FS_OFFSET + -0x2f) = 0;
      }
      pcVar43 = pcStack_670;
      pcStack_628 = (char *)&pcStack_5c8;
      if (*(char *)((long)in_FS_OFFSET + -0x17a) == '\x02') {
        pcVar37 = (char *)0x0;
      }
      else {
        pcVar30 = (char *)in_FS_OFFSET[-0x33];
        if (pcVar30 == (char *)0x0) {
          pcVar37 = (char *)0x0;
        }
        else {
          pcVar37 = (char *)0x0;
          if (*pcVar30 != '\0') {
            pcVar37 = pcVar30 + 8;
          }
        }
      }
      pcStack_638 = pcVar33 + 0x10;
      pcStack_630 = pcVar26;
      recovered_00664f00(&pcStack_638,pcVar37);
    }
  }
  else {
    pcVar33 = (char *)in_FS_OFFSET[-0x36];
    LOCK();
    lVar44 = *(long *)pcVar33;
    *(long *)pcVar33 = *(long *)pcVar33 + 1;
    UNLOCK();
    ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
    ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
    if (*(long *)pcVar33 == 0 || SCARRY8(lVar44,1) != *(long *)pcVar33 < 0)
    goto override_jmp_001a3b0e_case_2;
    plVar39 = *(long **)(pcVar33 + 0x210);
    if (plVar39 == (long *)0x0) {
      pcStack_578 = (char *)0x0;
      lStack_570 = 0;
    }
    else {
      LOCK();
      lVar44 = *plVar39;
      *plVar39 = *plVar39 + 1;
      UNLOCK();
      ppcVar8 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
      ppcVar10 = (char **)CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
      if (*plVar39 == 0 || SCARRY8(lVar44,1) != *plVar39 < 0) goto override_jmp_001a3b0e_case_2;
      pcStack_578 = *(char **)(pcVar33 + 0x210);
      lStack_570 = *(long *)(pcVar33 + 0x218);
    }
    pcStack_638 = (char *)0x0;
    iVar21 = (*(code *)PTR_DAT_00a8dfb8)(&pcStack_638,0x80,0x180);
    unaff_R14 = pcStack_5e8;
    pcVar28 = (char *)0x0;
    if (iVar21 == 0) {
      pcVar28 = pcStack_638;
    }
    if (pcVar28 == (char *)0x0) goto LAB_001a73bf;
    pcStack_638[0] = -0x34;
    pcStack_638[1] = '\0';
    pcStack_638[2] = '\0';
    pcStack_638[3] = '\0';
    pcStack_638[4] = '\0';
    pcStack_638[5] = '\0';
    pcStack_638[6] = '\0';
    pcStack_638[7] = '\0';
    *(long *)(pcStack_638 + 8) = 0;
    *(undefined **)(pcStack_638 + 0x10) = &UNK_00a5ce08;
    *(long *)(pcStack_638 + 0x18) = 0;
    *(char **)(pcStack_638 + 0x20) = pcVar33;
    *(char **)(pcStack_638 + 0x28) = pcStack_5e8;
    *(int *)(pcStack_638 + 0x30) = 0;
    *(char **)(pcStack_638 + 0x38) = pcStack_678;
    *(undefined4 *)(pcStack_638 + 0x40) = auStack_681._1_4_;
    *(char **)(pcStack_638 + 0x48) = unaff_R15;
    *(char **)(pcStack_638 + 0x50) = pcStack_5f0;
    *(char **)(pcStack_638 + 0x58) = unaff_R12;
    pcStack_638[0x90] = '\0';
    *(undefined1 (*) [16])(pcStack_638 + 0x110) = (undefined1  [16])0x0;
    *(long *)(pcStack_638 + 0x120) = 0;
    *(char **)(pcStack_638 + 0x130) = pcStack_578;
    *(long *)(pcStack_638 + 0x138) = lStack_570;
    lVar44 = recovered_00543360(pcVar33 + 0x160,pcStack_638);
    pcStack_638 = unaff_R14;
    if (*(long *)(pcVar33 + 0x200) != 0) {
      (**(code **)(*(long *)(pcVar33 + 0x208) + 0x28))
                (*(long *)(pcVar33 + 0x200) +
                 (*(long *)(*(long *)(pcVar33 + 0x208) + 0x10) - 1U & 0xfffffffffffffff0) + 0x10,
                 &pcStack_638);
    }
    pcVar43 = pcStack_670;
    if (lVar44 != 0) {
      recovered_006679b0(in_FS_OFFSET[-0x36],lVar44);
      pcVar43 = pcStack_670;
    }
  }
LAB_001a5108:
  in_FS_OFFSET[-0x38] = in_FS_OFFSET[-0x38] + -1;
  uVar32 = 0x84;
  LOCK();
  lVar44 = *(long *)pcVar28;
  if (lVar44 == 0xcc) {
    *(long *)pcVar28 = 0x84;
  }
  UNLOCK();
  if (lVar44 != 0xcc) {
    uVar32 = 0x84;
    (**(code **)(*(long *)(pcVar28 + 0x10) + 0x20))(pcVar28);
  }
  lVar44 = (*(code *)PTR_DAT_00a8dbd8)(0xffff,1);
  if (lVar44 != 0) {
    param_1[0x19] = 0xffff;
    param_1[0x1a] = lVar44;
    param_1[0x1b] = 0xffff;
    pcVar33 = pcVar43;
LAB_001a5446:
    do {
      param_1[0x22] = param_1[0x13] + 0x10;
      param_1[0x23] = param_1[0x1a];
      param_1[0x24] = param_1[0x1b];
      *(undefined1 *)(param_1 + 0x3e) = 0;
      pcVar28 = unaff_R14;
override_jmp_001a3d0e_case_4:
      recovered_001f8ad0(&uStack_548,param_1 + 0x22,pcVar33);
      unaff_R14 = uStack_548;
      sVar3 = CONCAT11(uStack_53f,cStack_540);
      if (sVar3 == -1) {
        *puStack_5e0 = 4;
        goto LAB_001a719a;
      }
      pcVar28 = uStack_548;
      uStack_78 = CONCAT26(uStack_538,CONCAT42(uStack_53c,uStack_53e));
      uStack_70 = CONCAT42(uStack_534,uStack_536);
      ppcStack_62 = uStack_528;
      uStack_6a = uStack_530;
      uStack_68 = (undefined6)(CONCAT44(uStack_52c,CONCAT22(uStack_52e,uStack_530)) >> 0x10);
      if (((((char)param_1[0x3e] == '\x03') && ((char)param_1[0x2f] == '\x03')) &&
          ((char)param_1[0x3d] == '\x03')) && ((char)param_1[0x3c] == '\x03')) {
        lVar44 = param_1[0x34];
        LOCK();
        bVar45 = *(char *)(lVar44 + 0x18) == '\0';
        if (bVar45) {
          *(char *)(lVar44 + 0x18) = '\x01';
        }
        UNLOCK();
        if (bVar45) {
          lVar25 = param_1[0x35];
          if (lVar25 == 0) goto LAB_001a5574;
LAB_001a5528:
          plVar39 = (long *)param_1[0x36];
          *(long **)(lVar25 + 8) = plVar39;
          if (plVar39 == (long *)0x0) goto LAB_001a558a;
LAB_001a5538:
          *plVar39 = param_1[0x35];
LAB_001a553e:
          *(undefined1 (*) [16])(param_1 + 0x35) = (undefined1  [16])0x0;
        }
        else {
          recovered_000618e0(lVar44 + 0x18);
          lVar25 = param_1[0x35];
          if (lVar25 != 0) goto LAB_001a5528;
LAB_001a5574:
          if (*(long **)(lVar44 + 0x20) == param_1 + 0x35) {
            plVar39 = (long *)param_1[0x36];
            *(long **)(lVar44 + 0x20) = plVar39;
            if (plVar39 != (long *)0x0) goto LAB_001a5538;
LAB_001a558a:
            if (*(long **)(lVar44 + 0x28) != param_1 + 0x35) goto LAB_001a5594;
            *(long *)(lVar44 + 0x28) = param_1[0x35];
            goto LAB_001a553e;
          }
        }
LAB_001a5594:
        LOCK();
        bVar45 = *(char *)(lVar44 + 0x18) == '\x01';
        if (bVar45) {
          *(char *)(lVar44 + 0x18) = '\0';
        }
        UNLOCK();
        if (bVar45) {
          lVar44 = param_1[0x37];
        }
        else {
          recovered_00061630(lVar44 + 0x18);
          lVar44 = param_1[0x37];
        }
        if (lVar44 != 0) {
          (**(code **)(lVar44 + 0x18))(param_1[0x38]);
        }
      }
      pcVar33 = pcStack_670;
      if (sVar3 != 2) {
        param_1[0x1c] = (long)unaff_R14;
        *(short *)(param_1 + 0x1d) = sVar3;
        *(undefined8 *)((long)param_1 + 0xea) = uStack_78;
        *(ulong *)((long)param_1 + 0xf2) = CONCAT26(uStack_6a,uStack_70);
        param_1[0x1f] = CONCAT62(uStack_68,uStack_6a);
        param_1[0x20] = (long)ppcStack_62;
        param_1[0x22] = param_1[0x16] + 0x10;
        *(undefined1 *)(param_1 + 0x30) = 0;
override_jmp_001a3d0e_case_5:
        param_2 = pcStack_670;
        unaff_R15 = (char *)recovered_002001b0(param_1 + 0x22,*(long *)pcStack_670);
        if (unaff_R15 != (char *)0x0) {
          if ((((char)param_1[0x30] == '\x03') && ((char)param_1[0x2f] == '\x03')) &&
             ((char)param_1[0x26] == '\x04')) {
            if ((char)param_1[0x2e] == '\x01') {
              pcVar33 = (char *)param_1[0x27];
              LOCK();
              cVar19 = *pcVar33;
              if (cVar19 == '\0') {
                *pcVar33 = '\x01';
              }
              UNLOCK();
              if (cVar19 != '\0') {
                recovered_000618e0(pcVar33);
              }
              pauVar1 = (undefined1 (*) [16])(param_1 + 0x2a);
              if (param_1[0x2a] == 0) {
                if (*(long **)(pcVar33 + 8) == param_1 + 0x28) {
                  lVar44 = param_1[0x2b];
                  *(long *)(pcVar33 + 8) = lVar44;
                  goto joined_r0x001a5880;
                }
              }
              else {
                lVar44 = param_1[0x2b];
                *(long *)(param_1[0x2a] + 0x18) = lVar44;
joined_r0x001a5880:
                if (lVar44 == 0) {
                  if (*(long **)(pcVar33 + 0x10) != param_1 + 0x28) goto LAB_001a5890;
                  *(undefined8 *)(pcVar33 + 0x10) = *(undefined8 *)*pauVar1;
                }
                else {
                  *(undefined8 *)(lVar44 + 0x10) = *(undefined8 *)*pauVar1;
                }
                *pauVar1 = (undefined1  [16])0x0;
              }
LAB_001a5890:
              if (param_1[0x2d] == param_1[0x2c]) {
                LOCK();
                cVar19 = *pcVar33;
                if (cVar19 == '\x01') {
                  *pcVar33 = '\0';
                }
                UNLOCK();
                if (cVar19 != '\x01') {
                  recovered_00061630(pcVar33);
                }
              }
              else {
                recovered_00659940(param_1[0x27]);
              }
            }
            if (param_1[0x28] != 0) {
              (**(code **)(param_1[0x28] + 0x18))(param_1[0x29]);
            }
          }
          plVar39 = param_1 + 0x1d;
          if (*(long *)(unaff_R15 + 0x40) == 0) {
LAB_001a5a05:
            LOCK();
            cVar19 = *unaff_R15;
            if (cVar19 == '\0') {
              *unaff_R15 = '\x01';
            }
            UNLOCK();
            if (cVar19 != '\0') {
              recovered_000618e0(unaff_R15);
            }
            recovered_00659940(unaff_R15,1,unaff_R15);
            LOCK();
            plVar46 = (long *)(param_1[0x18] + 0x10);
            lVar44 = *plVar46;
            *plVar46 = *plVar46 + 1;
            UNLOCK();
            param_1[0x22] = lVar44;
            lVar25 = param_1[0xf] + 0x10;
            lVar41 = *plVar39;
            lVar6 = param_1[0x1e];
            lVar7 = param_1[0x1f];
            uStack_534 = (undefined4)lVar7;
            uStack_530 = (undefined2)((ulong)lVar7 >> 0x20);
            uStack_52e = (undefined2)((ulong)lVar7 >> 0x30);
            uStack_52c = (undefined4)param_1[0x20];
            uStack_528._0_4_ = (undefined4)((ulong)param_1[0x20] >> 0x20);
            uStack_548._4_4_ = (undefined4)lVar41;
            cStack_540 = (char)((ulong)lVar41 >> 0x20);
            uStack_53f = (undefined1)((ulong)lVar41 >> 0x28);
            uStack_53e = (undefined2)((ulong)lVar41 >> 0x30);
            uStack_53c = (undefined4)lVar6;
            uStack_538 = (undefined2)((ulong)lVar6 >> 0x20);
            uStack_536 = (undefined2)((ulong)lVar6 >> 0x30);
            lVar38 = param_1[0x14];
            uVar35 = (undefined4)param_1[0x15];
            plVar39 = param_1 + 0x16;
            plVar46 = param_1 + 0x17;
            uStack_5a8 = (undefined4)uStack_528;
            uStack_5b8 = CONCAT44(uStack_534,(int)((ulong)lVar6 >> 0x20));
            puStack_5b0 = (undefined *)CONCAT44(uStack_52c,(int)((ulong)lVar7 >> 0x20));
            pcStack_5c8 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
            puStack_5c0 = (undefined8 *)CONCAT44(uStack_53c,(int)((ulong)lVar41 >> 0x20));
            param_1[0x23] = lVar38;
            *(undefined4 *)(param_1 + 0x24) = uVar35;
            *(char **)((long)param_1 + 0x124) = pcStack_5c8;
            *(undefined8 **)((long)param_1 + 300) = puStack_5c0;
            *(undefined8 *)((long)param_1 + 0x134) = uStack_5b8;
            *(undefined **)((long)param_1 + 0x13c) = puStack_5b0;
            *(undefined4 *)((long)param_1 + 0x144) = (undefined4)uStack_528;
            param_1[0x29] = lVar25;
            param_1[0x2a] = (long)plVar39;
            param_1[0x2b] = (long)plVar46;
            param_1[0x2c] = lVar44;
            *(undefined1 *)(param_1 + 0x3a) = 0;
LAB_001a5b25:
            param_1[0x2d] = lVar25;
            param_1[0x2e] = lVar44;
            param_1[0x2f] = param_1[0x25];
            param_1[0x30] = param_1[0x26];
            param_1[0x31] = param_1[0x27];
            param_1[0x32] = param_1[0x28];
            param_1[0x33] = lVar38;
            *(undefined4 *)(param_1 + 0x34) = uVar35;
            param_1[0x35] = (long)plVar39;
            param_1[0x36] = (long)plVar46;
            lVar44 = *(long *)(lVar25 + 0x58);
            uVar29 = *(ulong *)((long)&__DT_RELA[0xf7].r_addend + lVar44);
            param_1[0x3b] = lVar25 + 0x58;
            param_1[0x3c] = lVar44 + 0x1950;
            param_1[0x3d] = uVar29 >> 2;
            *(undefined1 (*) [16])(param_1 + 0x3e) = (undefined1  [16])0x0;
            *(undefined1 (*) [16])(param_1 + 0x40) = (undefined1  [16])0x0;
            *(undefined1 (*) [16])((long)param_1 + 0x209) = (undefined1  [16])0x0;
override_jmp_001a411d_case_3:
            pcVar28 = (char *)(param_1 + 0x3a);
            auVar51._8_8_ = param_2;
            auVar51._0_8_ = pcVar28;
            pcVar43 = (char *)(param_1 + 0x23);
            recovered_00371980(&uStack_548,param_1 + 0x3b,*(long *)param_2);
            iVar21 = (int)uStack_548;
            unaff_R14 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
            if (unaff_R14 == Elf64_Ehdr_00000000.e_ident_magic_str + 1) {
              *pcVar28 = '\x03';
              bVar45 = true;
              pcVar33 = param_2;
              pcVar26 = unaff_R14;
              goto LAB_001a6d6d;
            }
            pcStack_638 = (char *)CONCAT44(uStack_53c,
                                           CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540)));
            pcStack_630 = (char *)CONCAT44(uStack_534,CONCAT22(uStack_536,uStack_538));
            pcStack_628 = (char *)CONCAT44(uStack_52c,CONCAT22(uStack_52e,uStack_530));
            puStack_620 = (undefined *)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
            pcStack_608 = pcStack_510;
            uStack_600 = uStack_508;
            pcStack_618 = pcStack_520;
            pcStack_610 = pcStack_518;
            recovered_006680f0(param_1 + 0x3c);
            if (param_1[0x40] != 0) {
              (**(code **)(param_1[0x40] + 0x18))(param_1[0x41]);
            }
            if (iVar21 == 1) {
              auStack_681._1_8_ = pcVar28;
              pcVar28 = (char *)recovered_0005d6e0(&pcStack_638);
              goto LAB_001a6d63;
            }
            pcStack_510 = pcStack_628;
            pcStack_520 = pcStack_638;
            pcStack_518 = pcStack_630;
            uStack_548._0_4_ = (int)puStack_620;
            uStack_548._4_4_ = (undefined4)((ulong)puStack_620 >> 0x20);
            cStack_540 = (char)pcStack_618;
            uStack_53f = (undefined1)((ulong)pcStack_618 >> 8);
            uStack_53e = (undefined2)((ulong)pcStack_618 >> 0x10);
            uStack_53c = (undefined4)((ulong)pcStack_618 >> 0x20);
            uStack_538 = SUB82(pcStack_610,0);
            uStack_536 = (undefined2)((ulong)pcStack_610 >> 0x10);
            uStack_534 = (undefined4)((ulong)pcStack_610 >> 0x20);
            uStack_530 = SUB82(pcStack_608,0);
            uStack_52e = (undefined2)((ulong)pcStack_608 >> 0x10);
            uStack_52c = (undefined4)((ulong)pcStack_608 >> 0x20);
            uStack_528._0_4_ = (undefined4)uStack_600;
            uStack_528._4_4_ = (undefined4)((ulong)uStack_600 >> 0x20);
            puVar24 = (undefined8 *)(*(code *)PTR_DAT_00a8daa8)(0x40);
            if (puVar24 != (undefined8 *)0x0) {
              puVar24[6] = pcStack_518;
              puVar24[7] = pcStack_510;
              puVar24[4] = CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
              puVar24[5] = pcStack_520;
              puVar24[2] = CONCAT44(uStack_534,CONCAT22(uStack_536,uStack_538));
              puVar24[3] = CONCAT44(uStack_52c,CONCAT22(uStack_52e,uStack_530));
              *puVar24 = CONCAT44(uStack_548._4_4_,(int)uStack_548);
              puVar24[1] = CONCAT44(uStack_53c,CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540)))
              ;
              param_1[0x37] = (long)puVar24;
              param_1[0x38] = (long)&UNK_00a58e30;
              lVar44 = *(long *)(param_1[0x2d] + 0x28);
              param_1[0x3d] = *(long *)(param_1[0x2d] + 0x20);
              param_1[0x3e] = lVar44;
              param_1[0x40] = (long)(param_1 + 0x37);
              *(undefined2 *)((long)param_1 + 0x209) = 0x200;
override_jmp_001a411d_case_4:
              pcVar28 = (char *)(param_1 + 0x3a);
              pcVar43 = (char *)(param_1 + 0x23);
              recovered_00189ec0(&uStack_548,param_1 + 0x3b,param_2);
              pcVar26 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
              if ((char)uStack_548 == -2) {
                *pcVar28 = '\x04';
                bVar45 = true;
                pcVar33 = param_2;
                goto LAB_001a6d6d;
              }
              uStack_568._0_7_ = CONCAT43(uStack_548._4_4_,uStack_548._1_3_);
              uStack_560 = (undefined7)
                           (CONCAT44(uStack_53c,CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540))
                                    ) >> 8);
              uStack_568._7_1_ = cStack_540;
              if ((char)uStack_548 != -1) {
                auStack_681._1_8_ = pcVar28;
                pcVar28 = (char *)recovered_0005d5d0(&uStack_548);
LAB_001a6207:
                lVar44 = param_1[0x37];
                unaff_R14 = (char *)param_1[0x38];
                if (*(code **)unaff_R14 != (code *)0x0) {
                  (**(code **)unaff_R14)(lVar44);
                }
                if (*(long *)(unaff_R14 + 8) != 0) {
                  (*(code *)PTR_DAT_00a8dfe8)(lVar44);
                }
                goto LAB_001a6d63;
              }
              plVar39 = param_1 + 0x37;
              param_1[0x3b] = (long)plVar39;
              *(undefined1 *)(param_1 + 0x3c) = 3;
              uVar32 = 3;
LAB_001a5d91:
              pcVar30 = (char *)(param_1 + 0x3a);
              pcVar43 = (char *)(param_1 + 0x23);
              uStack_548._0_4_ = CONCAT31(uStack_548._1_3_,uVar32);
              pcVar33 = Elf64_Ehdr_00000000.e_ident_magic_str;
              auVar51 = (**(code **)(plVar39[1] + 0x20))(*plVar39,param_2,&uStack_548);
              lVar44 = auVar51._8_8_;
              pcVar26 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
              auStack_681._1_8_ = pcVar30;
              if (auVar51._0_8_ != 0) {
                lVar25 = lVar44;
                if (auVar51._0_8_ == 2) {
                  *pcVar30 = '\x05';
                  bVar45 = true;
                  pcVar33 = param_2;
                  pcVar28 = pcVar30;
                  goto LAB_001a6d6d;
                }
LAB_001a61ff:
                auStack_681._1_8_ = pcVar30;
                pcVar28 = (char *)recovered_0005d710(lVar25);
                goto LAB_001a6207;
              }
              lVar25 = 0x1700000003;
              if (lVar44 == 0) goto LAB_001a61ff;
              if (lVar44 == 1) {
                plVar39 = param_1 + 0x37;
                uVar29 = param_1[0x2e];
                lVar44 = *(long *)(param_1[0x2d] + 8);
                lVar25 = *(long *)(param_1[0x2d] + 0x10);
                param_1[0x3b] = (long)plVar39;
                param_1[0x3c] = lVar44;
                param_1[0x3d] = lVar25;
                param_1[0x3e] = uVar29;
                *(undefined1 *)(param_1 + 0x42) = 0;
LAB_001a5e41:
                param_1[0x3f] = (long)plVar39;
                param_1[0x40] = lVar44;
                param_1[0x41] = lVar25;
                param_1[0x43] = (long)plVar39;
                param_1[0x44] =
                     uVar29 >> 0x38 | (uVar29 & 0xff000000000000) >> 0x28 |
                     (uVar29 & 0xff0000000000) >> 0x18 | (uVar29 & 0xff00000000) >> 8 |
                     (uVar29 & 0xff000000) << 8 | (uVar29 & 0xff0000) << 0x18 |
                     (uVar29 & 0xff00) << 0x28 | uVar29 << 0x38;
                *(undefined1 *)(param_1 + 0x45) = 0;
                bVar20 = 0;
LAB_001a5e7e:
                pcStack_678 = (char *)(param_1 + 0x23);
                auStack_681._1_8_ = param_1 + 0x3a;
                unaff_R14 = (char *)(param_1 + 0x43);
                do {
                  auVar51 = (**(code **)((*(undefined8 **)unaff_R14)[1] + 0x20))
                                      (**(undefined8 **)unaff_R14,param_2,
                                       (char *)((long)(param_1 + 0x44) + (ulong)bVar20),
                                       8 - (ulong)bVar20);
                  lVar44 = auVar51._8_8_;
                  pcVar43 = pcStack_678;
                  if (auVar51._0_8_ != 0) {
                    if (auVar51._0_8_ == 2) {
                      *(undefined1 *)(param_1 + 0x42) = 3;
                      pcVar33 = param_2;
                      pcVar28 = (char *)(param_1 + 0x44);
                      goto LAB_001a61be;
                    }
LAB_001a619e:
                    pcVar28 = (char *)recovered_0005d710(lVar44);
                    pcVar43 = pcStack_678;
                    goto LAB_001a61ae;
                  }
                  if (lVar44 == 0) {
                    lVar44 = 0x1700000003;
                    goto LAB_001a619e;
                  }
                  bVar20 = auVar51[8] + *(byte *)(param_1 + 0x45);
                  *(byte *)(param_1 + 0x45) = bVar20;
                  pcVar28 = (char *)auStack_681._1_8_;
                } while (bVar20 < 8);
LAB_001a5ed7:
                param_1[0x43] = param_1[0x3f];
                param_1[0x44] = param_1[0x40];
                param_1[0x45] = param_1[0x41];
                *(undefined1 *)(param_1 + 0x49) = 0;
override_jmp_001a4256_case_4:
                recovered_0018ac50(&uStack_548,param_1 + 0x43,param_2);
                ppcVar36 = uStack_2d8;
                auStack_681._1_8_ = pcVar28;
                if ((char)uStack_548 != -2) {
                  uStack_2d8._0_4_ =
                       (undefined4)
                       (CONCAT44(uStack_53c,CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540))) >>
                       8);
                  uStack_2d8._7_1_ = SUB81(ppcVar36,7);
                  uStack_2d8._4_3_ = (undefined3)((uint)uStack_53c >> 8);
                  uStack_2e0._0_4_ = (int)CONCAT43(uStack_548._4_4_,uStack_548._1_3_);
                  uStack_2e0._4_4_ =
                       (undefined4)
                       (CONCAT17(cStack_540,CONCAT43(uStack_548._4_4_,uStack_548._1_3_)) >> 0x20);
                  pcVar28 = (char *)0x0;
                  if ((char)uStack_548 != -1) {
                    uStack_53c = (undefined4)
                                 (CONCAT34(uStack_2d8._4_3_,(undefined4)uStack_2d8) >> 0x18);
                    uStack_548._4_4_ =
                         (undefined4)(CONCAT44(uStack_2e0._4_4_,(int)uStack_2e0) >> 0x18);
                    pcVar28 = (char *)recovered_0005d5d0(&uStack_548);
                  }
LAB_001a61ae:
                  *(undefined1 *)(param_1 + 0x42) = 1;
                  if (pcVar28 != (char *)0x0) goto LAB_001a6207;
                  plVar39 = param_1 + 0x37;
                  param_1[0x3b] = (long)plVar39;
                  pcVar28 = (char *)auStack_681._1_8_;
LAB_001a5f57:
                  auVar51 = (**(code **)(plVar39[1] + 0x28))(*plVar39,param_2);
                  lVar25 = auVar51._8_8_;
                  pcVar26 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
                  if (auVar51._0_8_ == 1) {
                    *pcVar28 = '\a';
                    bVar45 = true;
                    pcVar33 = param_2;
                    goto LAB_001a6d6d;
                  }
                  pcVar30 = pcVar28;
                  if (lVar25 != 0) goto LAB_001a61ff;
                  lVar44 = param_1[0x33];
                  lVar25 = param_1[0x34];
                  auVar54 = recovered_0062c2a0(1);
                  uStack_660._0_4_ = auVar54._0_4_;
                  uStack_660._4_4_ = auVar54._4_4_;
                  uStack_2e0._0_4_ = (int)lVar44;
                  uStack_2e0._4_4_ = (undefined4)((ulong)lVar44 >> 0x20);
                  uStack_658._0_4_ = auVar54._8_4_;
                  uStack_2d8._0_4_ = (int)lVar25;
                  recovered_0062c1d0(&uStack_548,&uStack_660,&uStack_2e0);
                  unaff_R14 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
                  unaff_R15 = (char *)CONCAT44(uStack_53c,
                                               CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540)))
                  ;
                  uVar23 = CONCAT22(uStack_536,uStack_538);
                  auVar51._8_4_ = uVar23;
                  auVar51._0_8_ = pcVar28;
                  auVar51._12_4_ = 0;
                  puVar24 = (undefined8 *)(*(code *)PTR_DAT_00a8daa8)(0x18);
                  if (puVar24 == (undefined8 *)0x0) {
LAB_001a73db:
                    recovered_0005d9a2(8,0x18);
                    goto LAB_001a73ea;
                  }
                  lVar44 = 0;
                  if (unaff_R14 == (char *)0x0) {
                    lVar44 = (long)unaff_R15 * 1000 + (ulong)uVar23 / 1000000;
                  }
                  *puVar24 = 1;
                  puVar24[1] = 1;
                  puVar24[2] = lVar44;
                  param_1[0x39] = (long)puVar24;
                  param_1[0x3b] = *(long *)param_1[0x35] + 0x10;
                  *(undefined1 *)(param_1 + 0x49) = 0;
override_jmp_001a411d_case_8:
                  pcVar33 = pcStack_670;
                  auVar51._8_8_ = pcStack_670;
                  auVar51._0_8_ = pcVar28;
                  unaff_R15 = (char *)recovered_002001b0(param_1 + 0x3b,*(long *)pcStack_670);
                  pcVar26 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
                  auStack_681._1_8_ = pcVar28;
                  if (unaff_R15 == (char *)0x0) {
                    *pcVar28 = '\b';
                    bVar45 = true;
                    goto LAB_001a6d6d;
                  }
                  if ((((char)param_1[0x49] == '\x03') && ((char)param_1[0x48] == '\x03')) &&
                     ((char)param_1[0x3f] == '\x04')) {
                    if ((char)param_1[0x47] == '\x01') {
                      pcVar33 = (char *)param_1[0x40];
                      LOCK();
                      cVar19 = *pcVar33;
                      if (cVar19 == '\0') {
                        *pcVar33 = '\x01';
                      }
                      UNLOCK();
                      if (cVar19 != '\0') {
                        recovered_000618e0(pcVar33);
                        unaff_R14 = pcVar33;
                      }
                      pauVar1 = (undefined1 (*) [16])(param_1 + 0x43);
                      if (param_1[0x43] == 0) {
                        if (*(long **)(pcVar33 + 8) == param_1 + 0x41) {
                          lVar44 = param_1[0x44];
                          *(long *)(pcVar33 + 8) = lVar44;
                          goto joined_r0x001a645e;
                        }
                      }
                      else {
                        lVar44 = param_1[0x44];
                        *(long *)(param_1[0x43] + 0x18) = lVar44;
joined_r0x001a645e:
                        if (lVar44 == 0) {
                          if (*(long **)(pcVar33 + 0x10) != param_1 + 0x41) goto LAB_001a646e;
                          *(long *)(pcVar33 + 0x10) = *(long *)*pauVar1;
                        }
                        else {
                          *(undefined8 *)(lVar44 + 0x10) = *(undefined8 *)*pauVar1;
                        }
                        *pauVar1 = (undefined1  [16])0x0;
                      }
LAB_001a646e:
                      if (param_1[0x46] == param_1[0x45]) {
                        LOCK();
                        cVar19 = *pcVar33;
                        if (cVar19 == '\x01') {
                          *pcVar33 = '\0';
                        }
                        UNLOCK();
                        if (cVar19 != '\x01') {
                          recovered_00061630(pcVar33);
                        }
                      }
                      else {
                        recovered_00659940(param_1[0x40]);
                      }
                    }
                    if (param_1[0x41] != 0) {
                      (**(code **)(param_1[0x41] + 0x18))(param_1[0x42]);
                    }
                  }
                  lVar44 = param_1[0x30];
                  lVar25 = param_1[0x31];
                  lVar38 = param_1[0x32];
                  uStack_538 = (undefined2)lVar25;
                  uStack_536 = (undefined2)((ulong)lVar25 >> 0x10);
                  uStack_534 = (undefined4)((ulong)lVar25 >> 0x20);
                  uStack_530 = (undefined2)lVar38;
                  uStack_52e = (undefined2)((ulong)lVar38 >> 0x10);
                  uStack_52c = (undefined4)((ulong)lVar38 >> 0x20);
                  uStack_548._0_4_ = (int)param_1[0x2f];
                  uStack_548._4_4_ = (undefined4)((ulong)param_1[0x2f] >> 0x20);
                  cStack_540 = (char)lVar44;
                  uStack_53f = (undefined1)((ulong)lVar44 >> 8);
                  uStack_53e = (undefined2)((ulong)lVar44 >> 0x10);
                  uStack_53c = (undefined4)((ulong)lVar44 >> 0x20);
                  recovered_00275a80(unaff_R15 + 0x28,&uStack_548,param_1[0x2e]);
                  LOCK();
                  cVar19 = *unaff_R15;
                  if (cVar19 == '\0') {
                    *unaff_R15 = '\x01';
                  }
                  UNLOCK();
                  if (cVar19 != '\0') goto LAB_001a73f9;
                  goto LAB_001a6503;
                }
                *(undefined1 *)(param_1 + 0x42) = 4;
                pcVar33 = param_2;
LAB_001a61be:
                pcVar26 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
                *(char *)auStack_681._1_8_ = '\x06';
                bVar45 = true;
                goto LAB_001a6d6d;
              }
              uVar29 = 0x28;
              recovered_0005e7c0(&UNK_0078d074,0x28,&UNK_00a862b0);
              bVar20 = extraout_DL;
              goto LAB_001a759d;
            }
LAB_001a73ea:
            recovered_0005d9a2(8,0x40);
LAB_001a73f9:
            recovered_000618e0(unaff_R15);
LAB_001a6503:
            param_2 = auVar51._8_8_;
            pcVar28 = auVar51._0_8_;
            recovered_00659940(unaff_R15,1,unaff_R15);
            param_1[0x3b] = *(long *)param_1[0x36] + 0x10;
            *(undefined1 *)(param_1 + 0x49) = 0;
LAB_001a652f:
            pcVar30 = (char *)recovered_002001b0(param_1 + 0x3b,*(long *)param_2);
            if (pcVar30 != (char *)0x0) {
              if ((((char)param_1[0x49] != '\x03') || ((char)param_1[0x48] != '\x03')) ||
                 ((char)param_1[0x3f] != '\x04')) goto LAB_001a664a;
              if ((char)param_1[0x47] == '\x01') {
                pcVar33 = (char *)param_1[0x40];
                LOCK();
                cVar19 = *pcVar33;
                if (cVar19 == '\0') {
                  *pcVar33 = '\x01';
                }
                UNLOCK();
                if (cVar19 != '\0') {
                  recovered_000618e0(pcVar33);
                }
                pauVar1 = (undefined1 (*) [16])(param_1 + 0x43);
                if (param_1[0x43] == 0) {
                  if (*(long **)(pcVar33 + 8) == param_1 + 0x41) {
                    lVar44 = param_1[0x44];
                    *(long *)(pcVar33 + 8) = lVar44;
                    goto joined_r0x001a65f5;
                  }
                }
                else {
                  lVar44 = param_1[0x44];
                  *(long *)(param_1[0x43] + 0x18) = lVar44;
joined_r0x001a65f5:
                  if (lVar44 == 0) {
                    if (*(long **)(pcVar33 + 0x10) != param_1 + 0x41) goto LAB_001a6601;
                    *(undefined8 *)(pcVar33 + 0x10) = *(undefined8 *)*pauVar1;
                  }
                  else {
                    *(undefined8 *)(lVar44 + 0x10) = *(undefined8 *)*pauVar1;
                  }
                  *pauVar1 = (undefined1  [16])0x0;
                }
LAB_001a6601:
                if (param_1[0x46] == param_1[0x45]) {
                  LOCK();
                  cVar19 = *pcVar33;
                  if (cVar19 == '\x01') {
                    *pcVar33 = '\0';
                  }
                  UNLOCK();
                  if (cVar19 != '\x01') {
                    recovered_00061630(pcVar33);
                  }
                }
                else {
                  recovered_00659940(param_1[0x40]);
                }
              }
              if (param_1[0x41] != 0) {
                (**(code **)(param_1[0x41] + 0x18))(param_1[0x42]);
              }
LAB_001a664a:
              unaff_R15 = (char *)param_1[0x2e];
              lVar38 = param_1[0x30];
              lVar44 = param_1[0x31];
              lVar25 = param_1[0x32];
              uStack_538 = (undefined2)lVar44;
              uStack_536 = (undefined2)((ulong)lVar44 >> 0x10);
              uStack_534 = (undefined4)((ulong)lVar44 >> 0x20);
              uStack_530 = (undefined2)lVar25;
              uStack_52e = (undefined2)((ulong)lVar25 >> 0x10);
              uStack_52c = (undefined4)((ulong)lVar25 >> 0x20);
              uStack_548._0_4_ = (int)param_1[0x2f];
              uStack_548._4_4_ = (undefined4)((ulong)param_1[0x2f] >> 0x20);
              cStack_540 = (char)lVar38;
              uStack_53f = (undefined1)((ulong)lVar38 >> 8);
              uStack_53e = (undefined2)((ulong)lVar38 >> 0x10);
              uStack_53c = (undefined4)((ulong)lVar38 >> 0x20);
              unaff_R14 = (char *)param_1[0x39];
              plVar39 = *(long **)(pcVar30 + 0x48);
              uVar31 = recovered_00643940(plVar39,*(long *)(pcVar30 + 0x50),unaff_R15);
              if (*(long *)(pcVar30 + 0x38) == 0) {
                plVar39 = (long *)(pcVar30 + 0x28);
                recovered_0007b280(plVar39,pcVar30 + 0x48);
              }
              pcVar33 = *(char **)(pcVar30 + 0x28);
              uVar29 = *(ulong *)(pcVar30 + 0x30);
              bVar20 = (byte)(uVar31 >> 0x38);
              bVar34 = bVar20 >> 1;
              auVar51 = ZEXT216(CONCAT11(bVar20 >> 1,bVar20 >> 1));
              auVar51 = pshuflw(auVar51,auVar51,0);
              bVar45 = false;
              lVar44 = 0;
              do {
                uVar31 = uVar31 & uVar29;
                auVar53 = *(undefined1 (*) [16])(pcVar33 + uVar31);
                auVar52[0] = -(auVar53[0] == auVar51[0]);
                auVar52[1] = -(auVar53[1] == auVar51[1]);
                auVar52[2] = -(auVar53[2] == auVar51[2]);
                auVar52[3] = -(auVar53[3] == auVar51[3]);
                auVar52[4] = -(auVar53[4] == auVar51[4]);
                auVar52[5] = -(auVar53[5] == auVar51[5]);
                auVar52[6] = -(auVar53[6] == auVar51[6]);
                auVar52[7] = -(auVar53[7] == auVar51[7]);
                auVar52[8] = -(auVar53[8] == auVar51[0]);
                auVar52[9] = -(auVar53[9] == auVar51[1]);
                auVar52[10] = -(auVar53[10] == auVar51[2]);
                auVar52[0xb] = -(auVar53[0xb] == auVar51[3]);
                auVar52[0xc] = -(auVar53[0xc] == auVar51[4]);
                auVar52[0xd] = -(auVar53[0xd] == auVar51[5]);
                bVar50 = auVar53[0xf];
                auVar52[0xe] = -(auVar53[0xe] == auVar51[6]);
                auVar52[0xf] = -(bVar50 == auVar51[7]);
                uVar40 = (ushort)(SUB161(auVar52 >> 7,0) & 1) |
                         (ushort)(SUB161(auVar52 >> 0xf,0) & 1) << 1 |
                         (ushort)(SUB161(auVar52 >> 0x17,0) & 1) << 2 |
                         (ushort)(SUB161(auVar52 >> 0x1f,0) & 1) << 3 |
                         (ushort)(SUB161(auVar52 >> 0x27,0) & 1) << 4 |
                         (ushort)(SUB161(auVar52 >> 0x2f,0) & 1) << 5 |
                         (ushort)(SUB161(auVar52 >> 0x37,0) & 1) << 6 |
                         (ushort)(SUB161(auVar52 >> 0x3f,0) & 1) << 7 |
                         (ushort)(SUB161(auVar52 >> 0x47,0) & 1) << 8 |
                         (ushort)(SUB161(auVar52 >> 0x4f,0) & 1) << 9 |
                         (ushort)(SUB161(auVar52 >> 0x57,0) & 1) << 10 |
                         (ushort)(SUB161(auVar52 >> 0x5f,0) & 1) << 0xb |
                         (ushort)(SUB161(auVar52 >> 0x67,0) & 1) << 0xc |
                         (ushort)(SUB161(auVar52 >> 0x6f,0) & 1) << 0xd |
                         (ushort)(SUB161(auVar52 >> 0x77,0) & 1) << 0xe |
                         (ushort)(auVar52[0xf] >> 7) << 0xf;
                uVar23 = (uint)uVar40;
                while (uVar40 != 0) {
                  uVar22 = 0;
                  for (uVar5 = uVar23; (uVar5 & 1) == 0; uVar5 = uVar5 >> 1 | 0x80000000) {
                    uVar22 = uVar22 + 1;
                  }
                  uVar42 = uVar22 + uVar31 & uVar29;
                  if ((long *)unaff_R15 == *(long **)(pcVar33 + (uVar42 * -3 + -3) * 0x10)) {
                    sVar3 = *(short *)(pcVar33 + (uVar42 * -3 + -3) * 0x10 + 8);
                    plVar39 = *(long **)(pcVar33 + (uVar42 * -3 + -1) * 0x10 + 8);
                    *(ulong *)(pcVar33 + (uVar42 * -3 + -2) * 0x10 + 8) =
                         CONCAT44(uStack_534,CONCAT22(uStack_536,uStack_538));
                    *(ulong *)(pcVar33 + (uVar42 * -3 + -1) * 0x10) =
                         CONCAT44(uStack_52c,CONCAT22(uStack_52e,uStack_530));
                    *(ulong *)(pcVar33 + (uVar42 * -3 + -3) * 0x10 + 8) =
                         CONCAT44(uStack_548._4_4_,(int)uStack_548);
                    *(ulong *)(pcVar33 + (uVar42 * -3 + -2) * 0x10) =
                         CONCAT44(uStack_53c,CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540)));
                    *(char **)(pcVar33 + (uVar42 * -3 + -1) * 0x10 + 8) = unaff_R14;
                    if (sVar3 != 2) {
                      LOCK();
                      *plVar39 = *plVar39 + -1;
                      UNLOCK();
                      if (*plVar39 == 0) {
                        recovered_00604330();
                      }
                    }
                    goto LAB_001a67f7;
                  }
                  uVar40 = (ushort)(uVar23 - 1) & (ushort)uVar23;
                  uVar23 = CONCAT22((short)(uVar23 - 1 >> 0x10),uVar40);
                }
                if (bVar45) {
LAB_001a6728:
                  auVar48[0] = -(auVar53[0] == -1);
                  auVar48[1] = -(auVar53[1] == -1);
                  auVar48[2] = -(auVar53[2] == -1);
                  auVar48[3] = -(auVar53[3] == -1);
                  auVar48[4] = -(auVar53[4] == -1);
                  auVar48[5] = -(auVar53[5] == -1);
                  auVar48[6] = -(auVar53[6] == -1);
                  auVar48[7] = -(auVar53[7] == 0xff);
                  auVar48[8] = -(auVar53[8] == -1);
                  auVar48[9] = -(auVar53[9] == -1);
                  auVar48[10] = -(auVar53[10] == -1);
                  auVar48[0xb] = -(auVar53[0xb] == -1);
                  auVar48[0xc] = -(auVar53[0xc] == -1);
                  auVar48[0xd] = -(auVar53[0xd] == -1);
                  auVar48[0xe] = -(auVar53[0xe] == -1);
                  auVar48[0xf] = -(bVar50 == 0xff);
                  if ((((((((((((((((SUB161(auVar48 >> 7,0) & 1) != 0 ||
                                   (SUB161(auVar48 >> 0xf,0) & 1) != 0) ||
                                  (SUB161(auVar48 >> 0x17,0) & 1) != 0) ||
                                 (SUB161(auVar48 >> 0x1f,0) & 1) != 0) ||
                                (SUB161(auVar48 >> 0x27,0) & 1) != 0) ||
                               (SUB161(auVar48 >> 0x2f,0) & 1) != 0) ||
                              (SUB161(auVar48 >> 0x37,0) & 1) != 0) ||
                             (SUB161(auVar48 >> 0x3f,0) & 1) != 0) ||
                            (SUB161(auVar48 >> 0x47,0) & 1) != 0) ||
                           (SUB161(auVar48 >> 0x4f,0) & 1) != 0) ||
                          (SUB161(auVar48 >> 0x57,0) & 1) != 0) ||
                         (SUB161(auVar48 >> 0x5f,0) & 1) != 0) ||
                        (SUB161(auVar48 >> 0x67,0) & 1) != 0) ||
                       (SUB161(auVar48 >> 0x6f,0) & 1) != 0) || (SUB161(auVar48 >> 0x77,0) & 1) != 0
                      ) || auVar48[0xf] < '\0') goto LAB_001a678b;
                  bVar45 = true;
                }
                else {
                  uVar40 = (ushort)(SUB161(auVar53 >> 7,0) & 1) |
                           (ushort)(SUB161(auVar53 >> 0xf,0) & 1) << 1 |
                           (ushort)(SUB161(auVar53 >> 0x17,0) & 1) << 2 |
                           (ushort)(SUB161(auVar53 >> 0x1f,0) & 1) << 3 |
                           (ushort)(SUB161(auVar53 >> 0x27,0) & 1) << 4 |
                           (ushort)(SUB161(auVar53 >> 0x2f,0) & 1) << 5 |
                           (ushort)(SUB161(auVar53 >> 0x37,0) & 1) << 6 |
                           (ushort)(SUB161(auVar53 >> 0x3f,0) & 1) << 7 |
                           (ushort)(SUB161(auVar53 >> 0x47,0) & 1) << 8 |
                           (ushort)(SUB161(auVar53 >> 0x4f,0) & 1) << 9 |
                           (ushort)(SUB161(auVar53 >> 0x57,0) & 1) << 10 |
                           (ushort)(SUB161(auVar53 >> 0x5f,0) & 1) << 0xb |
                           (ushort)(SUB161(auVar53 >> 0x67,0) & 1) << 0xc |
                           (ushort)(SUB161(auVar53 >> 0x6f,0) & 1) << 0xd |
                           (ushort)(SUB161(auVar53 >> 0x77,0) & 1) << 0xe |
                           (ushort)(bVar50 >> 7) << 0xf;
                  if (uVar40 != 0) {
                    uVar23 = 0;
                    for (uVar22 = (uint)uVar40; (uVar22 & 1) == 0; uVar22 = uVar22 >> 1 | 0x80000000
                        ) {
                      uVar23 = uVar23 + 1;
                    }
                    plVar39 = (long *)(uVar23 + uVar31 & uVar29);
                    goto LAB_001a6728;
                  }
                  bVar45 = false;
                  plVar39 = (long *)0x0;
                }
                uVar31 = uVar31 + lVar44 + 0x10;
                lVar44 = lVar44 + 0x10;
              } while( true );
            }
            *(char *)auStack_681._1_8_ = '\t';
            bVar45 = true;
            pcVar33 = param_2;
            pcVar26 = uStack_548;
            goto LAB_001a6d6d;
          }
          uVar27 = recovered_0053f070(*(long *)(unaff_R15 + 0x48),*(long *)(unaff_R15 + 0x50),
                                      plVar39);
          lVar44 = recovered_00093c40(*(long *)(unaff_R15 + 0x28),*(long *)(unaff_R15 + 0x30),uVar27
                                      ,plVar39);
          if (lVar44 == 0) goto LAB_001a5a05;
          unaff_R14 = *(char **)(lVar44 + -8);
          LOCK();
          cVar19 = *unaff_R15;
          if (cVar19 == '\0') {
            *unaff_R15 = '\x01';
          }
          UNLOCK();
          if (cVar19 != '\0') goto LAB_001a74bd;
          goto LAB_001a592a;
        }
        *puStack_5e0 = 5;
        pcVar28 = unaff_R14;
        goto LAB_001a719a;
      }
      pcStack_598 = pcVar28;
      ppcVar36 = uStack_660;
      pcVar28 = uStack_2e0;
      pcVar43 = uStack_548;
      if ((_DAT_00a91d28 & 0xfffffffffffffffe) != 4) {
        bVar20 = DAT_00a8ffe0;
        if (2 < DAT_00a8ffe0) {
          bVar20 = recovered_00086870(&DAT_00a8ffd0);
        }
        ppcVar36 = uStack_660;
        pcVar28 = uStack_2e0;
        pcVar43 = uStack_548;
        if (bVar20 != 0) {
          if (bVar20 == 2) {
LAB_001a560c:
            pcVar26 = _DAT_00a91d20;
            puStack_2c8 = (undefined *)(_DAT_00a8ffd0 + 0x30);
            pcStack_638 = (char *)&pcStack_598;
            pcStack_630 = (char *)recovered_0031fcd0;
            uStack_660._0_4_ = 0x72649b;
            uStack_660._4_4_ = 0;
            uStack_658 = &pcStack_638;
            pcStack_5c8 = (char *)&uStack_660;
            puStack_5c0 = (undefined8 *)&UNK_00a8a7e8;
            uStack_2e0._0_4_ = 1;
            uStack_2e0._4_4_ = 0;
            uStack_2d8 = &pcStack_5c8;
            pcStack_2d0 = Elf64_Ehdr_00000000.e_ident_magic_str;
            uStack_538 = SUB82(&uStack_2e0,0);
            uStack_536 = (undefined2)((ulong)&uStack_2e0 >> 0x10);
            uStack_534 = (undefined4)((ulong)&uStack_2e0 >> 0x20);
            uStack_530 = (undefined2)_DAT_00a8ffd0;
            uStack_52e = (undefined2)((ulong)_DAT_00a8ffd0 >> 0x10);
            uStack_52c = (undefined4)((ulong)_DAT_00a8ffd0 >> 0x20);
            uStack_548._0_4_ = 1;
            uStack_548._4_4_ = 0;
            ppcVar36 = (char **)&UNK_0072649b;
            pcVar28 = Elf64_Ehdr_00000000.e_ident_magic_str;
            pcVar43 = Elf64_Ehdr_00000000.e_ident_magic_str;
            if (_DAT_00a92030 == 2) {
              lVar44 = _DAT_00a91d18;
              if (_DAT_00a91d10 == 1) {
                lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0)
                         + 0x10;
              }
              cVar19 = (**(code **)(_DAT_00a91d20 + 0x50))(lVar44,&uStack_548);
              unaff_R14 = pcVar26;
              ppcVar36 = (char **)CONCAT44(uStack_660._4_4_,(undefined4)uStack_660);
              pcVar28 = (char *)CONCAT44(uStack_2e0._4_4_,(int)uStack_2e0);
              pcVar43 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
              if (cVar19 != '\0') {
                (**(code **)(pcVar26 + 0x58))(lVar44,&uStack_548);
                ppcVar36 = (char **)CONCAT44(uStack_660._4_4_,(undefined4)uStack_660);
                pcVar28 = (char *)CONCAT44(uStack_2e0._4_4_,(int)uStack_2e0);
                pcVar43 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
              }
            }
          }
          else if (_DAT_00a92030 == 2) {
            lVar44 = _DAT_00a91d18;
            if (_DAT_00a91d10 == 1) {
              lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                       0x10;
            }
            cVar19 = (**(code **)(_DAT_00a91d20 + 0x28))(lVar44);
            ppcVar36 = uStack_660;
            pcVar28 = uStack_2e0;
            pcVar43 = uStack_548;
            if (cVar19 != '\0') goto LAB_001a560c;
          }
        }
      }
      uStack_548 = pcVar43;
      uStack_2e0 = pcVar28;
      uStack_660 = ppcVar36;
      if ((1 < ((uint)pcStack_598 & 3) - 2) && (((ulong)pcStack_598 & 3) != 0)) {
        (**(code **)(pcStack_598 + 0x17))((Elf64_Ehdr *)(pcStack_598 + -1));
      }
    } while( true );
  }
  recovered_0005d891(1,0xffff);
LAB_001a5400:
  recovered_0005d891(1,unaff_R12);
override_jmp_001a3d0e_case_1:
  recovered_0005e9f0(&UNK_00a59d38);
override_jmp_001a411d_case_1:
  recovered_0005e9f0(&UNK_00a59d08);
override_jmp_001a4256_case_1:
  uVar27 = recovered_0005e9f0(&UNK_00a586d0);
  goto LAB_001a5431;
LAB_001a678b:
  bVar50 = pcVar33[(long)plVar39];
  bVar20 = bVar20 >> 1;
  if (-1 < (char)bVar50) {
LAB_001a759d:
    bVar34 = bVar20;
    auVar51 = *(undefined1 (*) [16])pcVar33;
    uVar23 = 0;
    for (uVar22 = (uint)(ushort)((ushort)(SUB161(auVar51 >> 7,0) & 1) |
                                 (ushort)(SUB161(auVar51 >> 0xf,0) & 1) << 1 |
                                 (ushort)(SUB161(auVar51 >> 0x17,0) & 1) << 2 |
                                 (ushort)(SUB161(auVar51 >> 0x1f,0) & 1) << 3 |
                                 (ushort)(SUB161(auVar51 >> 0x27,0) & 1) << 4 |
                                 (ushort)(SUB161(auVar51 >> 0x2f,0) & 1) << 5 |
                                 (ushort)(SUB161(auVar51 >> 0x37,0) & 1) << 6 |
                                 (ushort)(SUB161(auVar51 >> 0x3f,0) & 1) << 7 |
                                 (ushort)(SUB161(auVar51 >> 0x47,0) & 1) << 8 |
                                 (ushort)(SUB161(auVar51 >> 0x4f,0) & 1) << 9 |
                                 (ushort)(SUB161(auVar51 >> 0x57,0) & 1) << 10 |
                                 (ushort)(SUB161(auVar51 >> 0x5f,0) & 1) << 0xb |
                                 (ushort)(SUB161(auVar51 >> 0x67,0) & 1) << 0xc |
                                 (ushort)(SUB161(auVar51 >> 0x6f,0) & 1) << 0xd |
                                 (ushort)(SUB161(auVar51 >> 0x77,0) & 1) << 0xe |
                                (ushort)(byte)(auVar51[0xf] >> 7) << 0xf); (uVar22 & 1) == 0;
        uVar22 = uVar22 >> 1 | 0x80000000) {
      uVar23 = uVar23 + 1;
    }
    plVar39 = (long *)(ulong)uVar23;
    bVar50 = pcVar33[(long)plVar39];
  }
  pcVar33[(long)plVar39] = bVar34;
  (*(undefined1 (*) [16])(pcVar33 + 0x10))[(ulong)(plVar39 + -2) & uVar29] = bVar34;
  lVar44 = *(long *)(pcVar30 + 0x40) - _UNK_0070d1e8;
  *(ulong *)(pcVar30 + 0x38) = *(long *)(pcVar30 + 0x38) - (ulong)(bVar50 & 1);
  *(long *)(pcVar30 + 0x40) = lVar44;
  *(char **)(pcVar33 + ((long)plVar39 * -3 + -3) * 0x10) = unaff_R15;
  lVar44 = (long)plVar39 * -3 + -3;
  *(ulong *)(pcVar33 + lVar44 * 0x10 + 8) = CONCAT44(uStack_548._4_4_,(int)uStack_548);
  *(ulong *)(pcVar33 + lVar44 * 0x10 + 0x10) =
       CONCAT44(uStack_53c,CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540)));
  lVar44 = (long)plVar39 * -3 + -2;
  *(ulong *)(pcVar33 + lVar44 * 0x10 + 8) = CONCAT44(uStack_534,CONCAT22(uStack_536,uStack_538));
  *(ulong *)(pcVar33 + lVar44 * 0x10 + 0x10) = CONCAT44(uStack_52c,CONCAT22(uStack_52e,uStack_530));
  *(char **)(pcVar33 + ((long)plVar39 * -3 + -1) * 0x10 + 8) = unaff_R14;
LAB_001a67f7:
  LOCK();
  cVar19 = *pcVar30;
  if (cVar19 == '\0') {
    *pcVar30 = '\x01';
  }
  UNLOCK();
  if (cVar19 != '\0') {
    recovered_000618e0(pcVar30);
  }
  recovered_00659940(pcVar30,1,pcVar30);
  param_2 = *(char **)param_1[0x35];
  LOCK();
  lVar44 = *(long *)param_2;
  *(long *)param_2 = *(long *)param_2 + 1;
  UNLOCK();
  uVar35 = uStack_648;
  uVar18 = uStack_644;
  iVar21 = iStack_640;
  ppcVar36 = uStack_2d8;
  ppcVar8 = uStack_658;
  ppcVar10 = uStack_650;
  ppcVar12 = uStack_528;
  if (*(long *)param_2 == 0 || SCARRY8(lVar44,1) != *(long *)param_2 < 0)
  goto override_jmp_001a3b0e_case_2;
  unaff_R15 = *(char **)param_1[0x36];
  LOCK();
  lVar44 = *(long *)unaff_R15;
  *(long *)unaff_R15 = *(long *)unaff_R15 + 1;
  UNLOCK();
  pcStack_678 = pcVar43;
  if (*(long *)unaff_R15 == 0 || SCARRY8(lVar44,1) != *(long *)unaff_R15 < 0)
  goto override_jmp_001a3b0e_case_2;
  pcVar28 = (char *)param_1[0x2e];
  unaff_R14 = (char *)param_1[0x37];
  unaff_R12 = (char *)param_1[0x38];
  pcStack_598 = (char *)param_1[0x2f];
  pcStack_590 = (code *)param_1[0x30];
  lStack_588 = param_1[0x31];
  lStack_580 = param_1[0x32];
  ppcVar36 = _DAT_00a91850;
  do {
    _DAT_00a91850 = ppcVar36;
    LOCK();
    ppcVar36 = (char **)((long)_DAT_00a91850 + 1);
    UNLOCK();
  } while (_DAT_00a91850 == (char **)0x0);
  ppcStack_550 = _DAT_00a91850;
  uStack_2e0 = (char *)&ppcStack_550;
  uStack_2d8 = (char **)auStack_681;
  pcStack_2d0 = (char *)param_1[0x2f];
  puStack_2c8 = (undefined *)param_1[0x30];
  pcStack_2c0 = (char *)param_1[0x31];
  pcStack_2b8 = (char *)param_1[0x32];
  uStack_288 = 0;
  lStack_5d8 = *in_FS_OFFSET + -0x1c0;
  _DAT_00a91850 = ppcVar36;
  pcStack_2b0 = unaff_R14;
  pcStack_2a8 = unaff_R12;
  pcStack_2a0 = unaff_R15;
  pcStack_298 = param_2;
  pcStack_290 = pcVar28;
  if ((char)in_FS_OFFSET[-0x2f] == '\0') {
    uVar29 = in_FS_OFFSET[-0x38];
    ppcVar9 = uStack_658;
    ppcVar11 = uStack_650;
    if (uVar29 < 0x7fffffffffffffff) goto LAB_001a6948;
LAB_001a75eb:
    uStack_658 = ppcVar9;
    uStack_650 = ppcVar11;
    uVar23 = recovered_0005e840(&UNK_00a89bc8);
LAB_001a75f7:
    if (uVar23 == 2) {
      uStack_660._0_4_ = SUB84(unaff_R14,0);
      uStack_660._4_4_ = (undefined4)((ulong)unaff_R14 >> 0x20);
      uStack_658._0_4_ = SUB84(unaff_R12,0);
      uStack_658._4_4_ = (undefined4)((ulong)unaff_R12 >> 0x20);
      pcVar33 = (char *)0x0;
    }
    else {
      recovered_0062acc0(lStack_5d8,recovered_00656300);
      *(undefined1 *)(in_FS_OFFSET + -0x2f) = 0;
LAB_001a6b4e:
      uStack_660._0_4_ = SUB84(unaff_R14,0);
      uStack_660._4_4_ = (undefined4)((ulong)unaff_R14 >> 0x20);
      uStack_658._0_4_ = SUB84(unaff_R12,0);
      uStack_658._4_4_ = (undefined4)((ulong)unaff_R12 >> 0x20);
      param_2 = pcStack_670;
      pcVar33 = (char *)0x0;
      if (*(char *)((long)in_FS_OFFSET + -0x17a) != '\x02') {
        pcVar43 = (char *)in_FS_OFFSET[-0x33];
        pcVar33 = (char *)0x0;
        if (pcVar43 != (char *)0x0) {
          pcVar33 = (char *)0x0;
          if (*pcVar43 != '\0') {
            pcVar33 = pcVar43 + 8;
          }
        }
      }
    }
    uStack_650 = &pcStack_668;
    recovered_00664f00(&uStack_660,pcVar33);
    pcVar43 = pcStack_678;
LAB_001a6d3b:
    in_FS_OFFSET[-0x38] = in_FS_OFFSET[-0x38] + -1;
    LOCK();
    lVar44 = *(long *)pcVar28;
    if (lVar44 == 0xcc) {
      *(long *)pcVar28 = 0x84;
    }
    UNLOCK();
    if (lVar44 != 0xcc) {
      (**(code **)(*(long *)(pcVar28 + 0x10) + 0x20))(pcVar28);
    }
    pcVar28 = (char *)0x0;
LAB_001a6d63:
    *(char *)auStack_681._1_8_ = '\x01';
    bVar45 = false;
    pcVar33 = param_2;
    pcVar26 = uStack_548;
LAB_001a6d6d:
    uStack_548 = pcVar26;
    if (bVar45) goto LAB_001a718f;
    recovered_000aa6b0(pcVar43);
    unaff_R12 = pcVar43;
    if (pcVar28 == (char *)0x0) goto LAB_001a6ef5;
    pcStack_668 = pcVar28;
    ppcVar36 = uStack_658;
    ppcVar8 = uStack_2d8;
    pcVar28 = uStack_548;
    if ((_DAT_00a91d28 & 0xfffffffffffffffe) != 4) {
      bVar20 = DAT_00a8fff8;
      if (2 < DAT_00a8fff8) {
        bVar20 = recovered_00086870(&DAT_00a8ffe8);
      }
      ppcVar36 = uStack_658;
      ppcVar8 = uStack_2d8;
      pcVar28 = uStack_548;
      if (bVar20 != 0) {
        if (bVar20 != 2) {
          if (_DAT_00a92030 != 2) goto LAB_001a7455;
          lVar44 = _DAT_00a91d18;
          if (_DAT_00a91d10 == 1) {
            lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                     0x10;
          }
          cVar19 = (**(code **)(_DAT_00a91d20 + 0x28))(lVar44);
          ppcVar36 = uStack_658;
          ppcVar8 = uStack_2d8;
          pcVar28 = uStack_548;
          if (cVar19 == '\0') goto LAB_001a7455;
        }
        pcVar43 = _DAT_00a91d20;
        puStack_620 = (undefined *)(_DAT_00a8ffe8 + 0x30);
        uStack_660 = &pcStack_668;
        uStack_658._0_4_ = 0x2ddb10;
        uStack_658._4_4_ = 0;
        pcStack_598 = &UNK_00726459;
        pcStack_590 = (code *)&uStack_660;
        uStack_568._0_7_ = SUB87(param_1 + 0x1d,0);
        uStack_568._7_1_ = (char)((ulong)(param_1 + 0x1d) >> 0x38);
        uStack_2e0 = (char *)&pcStack_598;
        uStack_2d8._0_4_ = 0xa8a7e8;
        uStack_2d8._4_4_ = 0;
        pcStack_2d0 = (char *)&uStack_568;
        puStack_2c8 = &UNK_00a58960;
        pcStack_638 = Elf64_Ehdr_00000000.e_ident_magic_str;
        pcStack_630 = (char *)&uStack_2e0;
        pcStack_628 = Elf64_Ehdr_00000000.e_ident_magic_str + 1;
        uStack_538 = SUB82(&pcStack_638,0);
        uStack_536 = (undefined2)((ulong)&pcStack_638 >> 0x10);
        uStack_534 = (undefined4)((ulong)&pcStack_638 >> 0x20);
        uStack_530 = (undefined2)_DAT_00a8ffe8;
        uStack_52e = (undefined2)((ulong)_DAT_00a8ffe8 >> 0x10);
        uStack_52c = (undefined4)((ulong)_DAT_00a8ffe8 >> 0x20);
        uStack_548._0_4_ = 1;
        uStack_548._4_4_ = 0;
        ppcVar36 = (char **)recovered_002ddb10;
        ppcVar8 = (char **)&UNK_00a8a7e8;
        pcVar28 = Elf64_Ehdr_00000000.e_ident_magic_str;
        if (_DAT_00a92030 == 2) {
          lVar44 = _DAT_00a91d18;
          if (_DAT_00a91d10 == 1) {
            lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                     0x10;
          }
          cVar19 = (**(code **)(_DAT_00a91d20 + 0x50))(lVar44,&uStack_548);
          unaff_R14 = pcVar43;
          ppcVar36 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
          ppcVar8 = (char **)CONCAT44(uStack_2d8._4_4_,(undefined4)uStack_2d8);
          pcVar28 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
          if (cVar19 != '\0') {
            (**(code **)(pcVar43 + 0x58))(lVar44,&uStack_548);
            ppcVar36 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
            ppcVar8 = (char **)CONCAT44(uStack_2d8._4_4_,(undefined4)uStack_2d8);
            pcVar28 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
          }
        }
      }
    }
LAB_001a7455:
    uStack_658 = ppcVar36;
    uStack_2d8 = ppcVar8;
    uStack_548 = pcVar28;
    (*(code *)**(undefined8 **)pcStack_668)();
    goto LAB_001a5446;
  }
  if ((char)in_FS_OFFSET[-0x2f] == '\x02') {
    recovered_000a4c00(&pcStack_2d0);
    uVar32 = 1;
    pcVar33 = uStack_2e0;
  }
  else {
    recovered_0062acc0(lStack_5d8,recovered_00656300);
    *(undefined1 *)(in_FS_OFFSET + -0x2f) = 0;
    uVar29 = in_FS_OFFSET[-0x38];
    ppcVar9 = uStack_658;
    ppcVar11 = uStack_650;
    if (0x7ffffffffffffffe < uVar29) goto LAB_001a75eb;
LAB_001a6948:
    in_FS_OFFSET[-0x38] = uVar29 + 1;
    lVar44 = in_FS_OFFSET[-0x37];
    pcStack_5e8 = unaff_R14;
    (*(code *)PTR_DAT_00a8da58)(&uStack_548,&uStack_2e0,0xd8);
    iVar21 = iStack_640;
    uVar18 = uStack_644;
    uVar35 = uStack_648;
    ppcVar10 = uStack_650;
    ppcVar8 = uStack_658;
    ppcVar36 = uStack_660;
    uVar14 = (undefined4)uStack_660;
    if (lVar44 != 2) {
      pcVar33 = *(char **)CONCAT44(uStack_548._4_4_,(int)uStack_548);
      uStack_650._4_4_ = (undefined4)lStack_588;
      uStack_648 = (undefined4)((ulong)lStack_588 >> 0x20);
      uStack_644 = (undefined4)lStack_580;
      iStack_640 = (int)((ulong)lStack_580 >> 0x20);
      uStack_660._4_4_ = SUB84(pcStack_598,0);
      uVar16 = uStack_660._4_4_;
      uStack_658._0_4_ = (undefined4)((ulong)pcStack_598 >> 0x20);
      uStack_658._4_4_ = SUB84(pcStack_590,0);
      uStack_650._0_4_ = (undefined4)((ulong)pcStack_590 >> 0x20);
      uStack_660._0_4_ = SUB84(pcVar33,0);
      uVar15 = (undefined4)uStack_660;
      uStack_660._4_4_ = (undefined4)((ulong)pcVar33 >> 0x20);
      uVar17 = uStack_660._4_4_;
      uStack_660 = ppcVar36;
      uStack_660._0_4_ = uVar14;
      uStack_660._4_4_ = uVar16;
      pcStack_5f0 = pcVar28;
      pcStack_578 = param_2;
      ppcVar36 = uStack_2d8;
      ppcVar12 = uStack_528;
      if ((int)lVar44 == 1) {
        unaff_R14 = (char *)in_FS_OFFSET[-0x36];
        LOCK();
        lVar44 = *(long *)unaff_R14;
        *(long *)unaff_R14 = *(long *)unaff_R14 + 1;
        UNLOCK();
        uStack_548 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
        if (*(long *)unaff_R14 == 0 || SCARRY8(lVar44,1) != *(long *)unaff_R14 < 0)
        goto override_jmp_001a3b0e_case_2;
        plVar39 = *(long **)(unaff_R14 + 0x210);
        lVar44 = 0;
        lVar25 = 0;
        if (plVar39 != (long *)0x0) {
          LOCK();
          lVar44 = *plVar39;
          *plVar39 = *plVar39 + 1;
          UNLOCK();
          uStack_548 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
          if (*plVar39 == 0 || SCARRY8(lVar44,1) != *plVar39 < 0) goto override_jmp_001a3b0e_case_2;
          lVar44 = *(long *)(unaff_R14 + 0x210);
          lVar25 = *(long *)(unaff_R14 + 0x218);
        }
        pcStack_668 = (char *)0x0;
        lStack_58 = lVar44;
        lStack_50 = lVar25;
        iVar21 = (*(code *)PTR_DAT_00a8dfb8)(&pcStack_668,0x80,0x180);
        pcVar28 = (char *)0x0;
        if (iVar21 == 0) {
          pcVar28 = pcStack_668;
        }
        if (pcVar28 == (char *)0x0) goto LAB_001a73bf;
        pcStack_668[0] = -0x34;
        pcStack_668[1] = '\0';
        pcStack_668[2] = '\0';
        pcStack_668[3] = '\0';
        pcStack_668[4] = '\0';
        pcStack_668[5] = '\0';
        pcStack_668[6] = '\0';
        pcStack_668[7] = '\0';
        *(long *)(pcStack_668 + 8) = 0;
        *(undefined **)(pcStack_668 + 0x10) = &UNK_00a5cd18;
        *(long *)(pcStack_668 + 0x18) = 0;
        *(char **)(pcStack_668 + 0x20) = unaff_R14;
        *(char **)(pcStack_668 + 0x28) = pcVar33;
        *(int *)(pcStack_668 + 0x30) = 0;
        *(ulong *)(pcStack_668 + 0x34) = CONCAT44(uStack_660._4_4_,(undefined4)uStack_660);
        *(ulong *)(pcStack_668 + 0x3c) = CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
        *(ulong *)(pcStack_668 + 0x44) = CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
        *(ulong *)(pcStack_668 + 0x4c) = CONCAT44(uStack_644,uStack_648);
        *(int *)(pcStack_668 + 0x54) = iStack_640;
        *(char **)(pcStack_668 + 0x58) = pcStack_5e8;
        *(char **)(pcStack_668 + 0x60) = unaff_R12;
        *(char **)(pcStack_668 + 0x68) = unaff_R15;
        *(char **)(pcStack_668 + 0x70) = pcStack_578;
        *(char **)(pcStack_668 + 0x78) = pcStack_5f0;
        pcStack_668[0x80] = '\0';
        *(long *)(pcStack_668 + 0x100) = 0;
        *(long *)(pcStack_668 + 0x108) = 0;
        *(long *)(pcStack_668 + 0x110) = 0;
        *(long *)(pcStack_668 + 0x120) = lStack_58;
        *(long *)(pcStack_668 + 0x128) = lStack_50;
        unaff_R12 = (char *)recovered_00543360(unaff_R14 + 0x98,pcStack_668);
        uStack_660._0_4_ = uVar15;
        uStack_660._4_4_ = uVar17;
        if (*(long *)(unaff_R14 + 0x200) != 0) {
          (**(code **)(*(long *)(unaff_R14 + 0x208) + 0x28))
                    (*(long *)(unaff_R14 + 0x200) +
                     (*(long *)(*(long *)(unaff_R14 + 0x208) + 0x10) - 1U & 0xfffffffffffffff0) +
                     0x10,&uStack_660);
        }
        param_2 = pcStack_670;
        pcVar43 = pcStack_678;
        if (unaff_R12 != (char *)0x0) {
          unaff_R14 = unaff_R14 + 0x10;
          pcStack_668 = (char *)((ulong)pcStack_668 & 0xffffffffffffff00);
          uVar23 = (uint)*(byte *)(in_FS_OFFSET + -0x2f);
          if (*(byte *)(in_FS_OFFSET + -0x2f) != 0) goto LAB_001a75f7;
          goto LAB_001a6b4e;
        }
      }
      else {
        unaff_R14 = (char *)in_FS_OFFSET[-0x36];
        LOCK();
        lVar44 = *(long *)unaff_R14;
        *(long *)unaff_R14 = *(long *)unaff_R14 + 1;
        UNLOCK();
        uStack_548 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
        if (*(long *)unaff_R14 == 0 || SCARRY8(lVar44,1) != *(long *)unaff_R14 < 0)
        goto override_jmp_001a3b0e_case_2;
        plVar39 = *(long **)(unaff_R14 + 0x210);
        lVar44 = 0;
        lVar25 = 0;
        if (plVar39 != (long *)0x0) {
          LOCK();
          lVar44 = *plVar39;
          *plVar39 = *plVar39 + 1;
          UNLOCK();
          uStack_548 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
          if (*plVar39 == 0 || SCARRY8(lVar44,1) != *plVar39 < 0) goto override_jmp_001a3b0e_case_2;
          lVar44 = *(long *)(unaff_R14 + 0x210);
          lVar25 = *(long *)(unaff_R14 + 0x218);
        }
        pcStack_668 = (char *)0x0;
        lStack_5d8 = lVar44;
        lStack_5d0 = lVar25;
        iVar21 = (*(code *)PTR_DAT_00a8dfb8)(&pcStack_668,0x80,0x180);
        pcVar28 = (char *)0x0;
        if (iVar21 == 0) {
          pcVar28 = pcStack_668;
        }
        if (pcVar28 == (char *)0x0) goto LAB_001a73bf;
        pcStack_668[0] = -0x34;
        pcStack_668[1] = '\0';
        pcStack_668[2] = '\0';
        pcStack_668[3] = '\0';
        pcStack_668[4] = '\0';
        pcStack_668[5] = '\0';
        pcStack_668[6] = '\0';
        pcStack_668[7] = '\0';
        *(long *)(pcStack_668 + 8) = 0;
        *(undefined **)(pcStack_668 + 0x10) = &UNK_00a5ccc8;
        *(long *)(pcStack_668 + 0x18) = 0;
        *(char **)(pcStack_668 + 0x20) = unaff_R14;
        *(char **)(pcStack_668 + 0x28) = pcVar33;
        *(int *)(pcStack_668 + 0x30) = 0;
        *(ulong *)(pcStack_668 + 0x34) = CONCAT44(uStack_660._4_4_,(undefined4)uStack_660);
        *(ulong *)(pcStack_668 + 0x3c) = CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
        *(ulong *)(pcStack_668 + 0x44) = CONCAT44(uStack_650._4_4_,(undefined4)uStack_650);
        *(ulong *)(pcStack_668 + 0x4c) = CONCAT44(uStack_644,uStack_648);
        *(int *)(pcStack_668 + 0x54) = iStack_640;
        *(char **)(pcStack_668 + 0x58) = pcStack_5e8;
        *(char **)(pcStack_668 + 0x60) = unaff_R12;
        *(char **)(pcStack_668 + 0x68) = unaff_R15;
        *(char **)(pcStack_668 + 0x70) = pcStack_578;
        *(char **)(pcStack_668 + 0x78) = pcStack_5f0;
        pcStack_668[0x80] = '\0';
        *(long *)(pcStack_668 + 0x100) = 0;
        *(long *)(pcStack_668 + 0x108) = 0;
        *(long *)(pcStack_668 + 0x110) = 0;
        *(long *)(pcStack_668 + 0x120) = lStack_5d8;
        *(long *)(pcStack_668 + 0x128) = lStack_5d0;
        lVar44 = recovered_00543360(unaff_R14 + 0x160,pcStack_668);
        uStack_660._0_4_ = uVar15;
        uStack_660._4_4_ = uVar17;
        if (*(long *)(unaff_R14 + 0x200) != 0) {
          (**(code **)(*(long *)(unaff_R14 + 0x208) + 0x28))
                    (*(long *)(unaff_R14 + 0x200) +
                     (*(long *)(*(long *)(unaff_R14 + 0x208) + 0x10) - 1U & 0xfffffffffffffff0) +
                     0x10,&uStack_660);
        }
        param_2 = pcStack_670;
        pcVar43 = pcStack_678;
        if (lVar44 != 0) {
          recovered_006679b0(in_FS_OFFSET[-0x36],lVar44);
        }
      }
      goto LAB_001a6d3b;
    }
    recovered_000a4c00(&uStack_538);
    in_FS_OFFSET[-0x38] = in_FS_OFFSET[-0x38] + -1;
    uVar32 = 0;
    unaff_R14 = Elf64_Ehdr_00000000.e_ident_magic_str + 1;
    pcVar33 = uStack_2e0;
  }
  uStack_2e0._4_4_ = (undefined4)((ulong)pcVar33 >> 0x20);
  uStack_2e0._1_3_ = (undefined3)((ulong)pcVar33 >> 8);
  uStack_2e0._0_4_ = CONCAT31(uStack_2e0._1_3_,uVar32);
  uStack_548 = (char *)&uStack_2e0;
  cStack_540 = -0x80;
  uStack_53f = 0x85;
  uStack_53e = 0x65;
  uStack_53c = 0;
  recovered_0005e620(&UNK_00729475,&uStack_548,&UNK_00a59d20);
LAB_001a74bd:
  recovered_000618e0(unaff_R15);
LAB_001a592a:
  recovered_00659940(unaff_R15,1,unaff_R15);
  param_1[0x22] = (long)unaff_R14;
  param_1[0x23] = param_1[0x17] + 0x10;
  *(undefined1 *)(param_1 + 0x31) = 0;
  pcVar33 = param_2;
override_jmp_001a3d0e_case_6:
  unaff_R12 = (char *)recovered_002001b0(param_1 + 0x23,*(long *)pcVar33);
  if (unaff_R12 == (char *)0x0) {
    *puStack_5e0 = 6;
    pcVar28 = unaff_R14;
    goto LAB_001a719a;
  }
  if ((((char)param_1[0x31] == '\x03') && ((char)param_1[0x30] == '\x03')) &&
     ((char)param_1[0x27] == '\x04')) {
    if ((char)param_1[0x2f] == '\x01') {
      pcVar28 = (char *)param_1[0x28];
      LOCK();
      cVar19 = *pcVar28;
      if (cVar19 == '\0') {
        *pcVar28 = '\x01';
      }
      UNLOCK();
      if (cVar19 != '\0') {
        recovered_000618e0(pcVar28);
      }
      pauVar1 = (undefined1 (*) [16])(param_1 + 0x2b);
      if (param_1[0x2b] == 0) {
        if (*(long **)(pcVar28 + 8) == param_1 + 0x29) {
          lVar44 = param_1[0x2c];
          *(long *)(pcVar28 + 8) = lVar44;
          goto joined_r0x001a628a;
        }
      }
      else {
        lVar44 = param_1[0x2c];
        *(long *)(param_1[0x2b] + 0x18) = lVar44;
joined_r0x001a628a:
        if (lVar44 == 0) {
          if (*(long **)(pcVar28 + 0x10) != param_1 + 0x29) goto LAB_001a629a;
          *(undefined8 *)(pcVar28 + 0x10) = *(undefined8 *)*pauVar1;
        }
        else {
          *(undefined8 *)(lVar44 + 0x10) = *(undefined8 *)*pauVar1;
        }
        *pauVar1 = (undefined1  [16])0x0;
      }
LAB_001a629a:
      if (param_1[0x2e] == param_1[0x2d]) {
        LOCK();
        cVar19 = *pcVar28;
        if (cVar19 == '\x01') {
          *pcVar28 = '\0';
        }
        UNLOCK();
        if (cVar19 != '\x01') {
          recovered_00061630(pcVar28);
        }
      }
      else {
        recovered_00659940(param_1[0x28]);
      }
    }
    if (param_1[0x29] != 0) {
      (**(code **)(param_1[0x29] + 0x18))(param_1[0x2a]);
    }
  }
  if (*(long *)(unaff_R12 + 0x40) != 0) {
    lVar44 = param_1[0x22];
    uVar29 = recovered_00643940(*(long *)(unaff_R12 + 0x48),*(long *)(unaff_R12 + 0x50),lVar44);
    lVar25 = *(long *)(unaff_R12 + 0x28);
    bVar20 = (byte)(uVar29 >> 0x39);
    auVar51 = ZEXT216(CONCAT11(bVar20,bVar20));
    auVar51 = pshuflw(auVar51,auVar51,0);
    lVar38 = 0;
    while( true ) {
      uVar29 = uVar29 & *(ulong *)(unaff_R12 + 0x30);
      auVar53 = *(undefined1 (*) [16])(lVar25 + uVar29);
      auVar49[0] = -(auVar53[0] == auVar51[0]);
      auVar49[1] = -(auVar53[1] == auVar51[1]);
      auVar49[2] = -(auVar53[2] == auVar51[2]);
      auVar49[3] = -(auVar53[3] == auVar51[3]);
      auVar49[4] = -(auVar53[4] == auVar51[4]);
      auVar49[5] = -(auVar53[5] == auVar51[5]);
      auVar49[6] = -(auVar53[6] == auVar51[6]);
      auVar49[7] = -(auVar53[7] == auVar51[7]);
      auVar49[8] = -(auVar53[8] == auVar51[0]);
      auVar49[9] = -(auVar53[9] == auVar51[1]);
      auVar49[10] = -(auVar53[10] == auVar51[2]);
      auVar49[0xb] = -(auVar53[0xb] == auVar51[3]);
      auVar49[0xc] = -(auVar53[0xc] == auVar51[4]);
      auVar49[0xd] = -(auVar53[0xd] == auVar51[5]);
      auVar49[0xe] = -(auVar53[0xe] == auVar51[6]);
      auVar49[0xf] = -(auVar53[0xf] == auVar51[7]);
      uVar40 = (ushort)(SUB161(auVar49 >> 7,0) & 1) | (ushort)(SUB161(auVar49 >> 0xf,0) & 1) << 1 |
               (ushort)(SUB161(auVar49 >> 0x17,0) & 1) << 2 |
               (ushort)(SUB161(auVar49 >> 0x1f,0) & 1) << 3 |
               (ushort)(SUB161(auVar49 >> 0x27,0) & 1) << 4 |
               (ushort)(SUB161(auVar49 >> 0x2f,0) & 1) << 5 |
               (ushort)(SUB161(auVar49 >> 0x37,0) & 1) << 6 |
               (ushort)(SUB161(auVar49 >> 0x3f,0) & 1) << 7 |
               (ushort)(SUB161(auVar49 >> 0x47,0) & 1) << 8 |
               (ushort)(SUB161(auVar49 >> 0x4f,0) & 1) << 9 |
               (ushort)(SUB161(auVar49 >> 0x57,0) & 1) << 10 |
               (ushort)(SUB161(auVar49 >> 0x5f,0) & 1) << 0xb |
               (ushort)(SUB161(auVar49 >> 0x67,0) & 1) << 0xc |
               (ushort)(SUB161(auVar49 >> 0x6f,0) & 1) << 0xd |
               (ushort)(SUB161(auVar49 >> 0x77,0) & 1) << 0xe | (ushort)(auVar49[0xf] >> 7) << 0xf;
      uVar23 = (uint)uVar40;
      while (uVar40 != 0) {
        uVar22 = 0;
        for (uVar5 = uVar23; (uVar5 & 1) == 0; uVar5 = uVar5 >> 1 | 0x80000000) {
          uVar22 = uVar22 + 1;
        }
        lVar41 = (uVar22 + uVar29 & *(ulong *)(unaff_R12 + 0x30)) * -0x30;
        if (lVar44 == *(long *)(lVar25 + -0x30 + lVar41)) {
          lVar44 = *(long *)(lVar25 + lVar41 + -8);
          lVar25 = param_1[0x14];
          lVar38 = param_1[0x15];
          auVar54 = recovered_0062c2a0(1);
          pcVar33 = pcStack_670;
          pcStack_638 = auVar54._0_8_;
          pcStack_630 = (char *)CONCAT44(pcStack_630._4_4_,auVar54._8_4_);
          uStack_2e0._0_4_ = (int)lVar25;
          uStack_2e0._4_4_ = (undefined4)((ulong)lVar25 >> 0x20);
          uStack_2d8._0_4_ = (int)lVar38;
          recovered_0062c1d0(&uStack_548,&pcStack_638,&uStack_2e0);
          lVar25 = 0;
          if ((char)uStack_548 == '\0') {
            lVar25 = (ulong)CONCAT22(uStack_536,uStack_538) / 1000000 +
                     CONCAT44(uStack_53c,CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540))) *
                     1000;
          }
          *(long *)(lVar44 + 0x10) = lVar25;
          goto LAB_001a6425;
        }
        uVar40 = (ushort)(uVar23 - 1) & (ushort)uVar23;
        uVar23 = CONCAT22((short)(uVar23 - 1 >> 0x10),uVar40);
      }
      auVar47[0] = -(auVar53[0] == -1);
      auVar47[1] = -(auVar53[1] == -1);
      auVar47[2] = -(auVar53[2] == -1);
      auVar47[3] = -(auVar53[3] == -1);
      auVar47[4] = -(auVar53[4] == -1);
      auVar47[5] = -(auVar53[5] == -1);
      auVar47[6] = -(auVar53[6] == -1);
      auVar47[7] = -(auVar53[7] == -1);
      auVar47[8] = -(auVar53[8] == -1);
      auVar47[9] = -(auVar53[9] == -1);
      auVar47[10] = -(auVar53[10] == -1);
      auVar47[0xb] = -(auVar53[0xb] == -1);
      auVar47[0xc] = -(auVar53[0xc] == -1);
      auVar47[0xd] = -(auVar53[0xd] == -1);
      auVar47[0xe] = -(auVar53[0xe] == -1);
      auVar47[0xf] = -(auVar53[0xf] == -1);
      if ((((((((((((((((SUB161(auVar47 >> 7,0) & 1) != 0 || (SUB161(auVar47 >> 0xf,0) & 1) != 0) ||
                      (SUB161(auVar47 >> 0x17,0) & 1) != 0) || (SUB161(auVar47 >> 0x1f,0) & 1) != 0)
                    || (SUB161(auVar47 >> 0x27,0) & 1) != 0) || (SUB161(auVar47 >> 0x2f,0) & 1) != 0
                   ) || (SUB161(auVar47 >> 0x37,0) & 1) != 0) ||
                 (SUB161(auVar47 >> 0x3f,0) & 1) != 0) || (SUB161(auVar47 >> 0x47,0) & 1) != 0) ||
               (SUB161(auVar47 >> 0x4f,0) & 1) != 0) || (SUB161(auVar47 >> 0x57,0) & 1) != 0) ||
             (SUB161(auVar47 >> 0x5f,0) & 1) != 0) || (SUB161(auVar47 >> 0x67,0) & 1) != 0) ||
           (SUB161(auVar47 >> 0x6f,0) & 1) != 0) || (SUB161(auVar47 >> 0x77,0) & 1) != 0) ||
          auVar47[0xf] < '\0') break;
      uVar29 = uVar29 + lVar38 + 0x10;
      lVar38 = lVar38 + 0x10;
    }
  }
LAB_001a6425:
  LOCK();
  cVar19 = *unaff_R12;
  if (cVar19 == '\0') {
    *unaff_R12 = '\x01';
  }
  UNLOCK();
  if (cVar19 != '\0') goto LAB_001a73ce;
  while( true ) {
    recovered_00659940(unaff_R12,1,unaff_R12);
LAB_001a6ef5:
    if ((ulong)param_1[0x1c] <= (ulong)param_1[0x1b]) break;
    recovered_0005e6f0(0,param_1[0x1c],param_1[0x1b],&UNK_00a59d80);
LAB_001a73bf:
    recovered_0005d9a2(0x80,0x180);
LAB_001a73ce:
    recovered_000618e0(unaff_R12);
  }
  unaff_R14 = (char *)param_1[0xf];
  recovered_005b1080(&uStack_548,param_1[0x22],param_1[0x1a]);
  recovered_0036bd70(&uStack_2e0,*(long *)(unaff_R14 + 0x68),&uStack_548);
  if ((int)uStack_2e0 == -1) goto LAB_001a5446;
  pcStack_518 = pcStack_2b0;
  uStack_528._0_4_ = SUB84(pcStack_2c0,0);
  uStack_528._4_4_ = (undefined4)((ulong)pcStack_2c0 >> 0x20);
  pcStack_520 = pcStack_2b8;
  uStack_538 = SUB82(pcStack_2d0,0);
  uStack_536 = (undefined2)((ulong)pcStack_2d0 >> 0x10);
  uStack_534 = (undefined4)((ulong)pcStack_2d0 >> 0x20);
  uStack_530 = SUB82(puStack_2c8,0);
  uStack_52e = (undefined2)((ulong)puStack_2c8 >> 0x10);
  uStack_52c = (undefined4)((ulong)puStack_2c8 >> 0x20);
  cStack_540 = (char)uStack_2d8;
  uStack_53f = (undefined1)((ulong)uStack_2d8 >> 8);
  uStack_53e = (undefined2)((ulong)uStack_2d8 >> 0x10);
  uStack_53c = uStack_2d8._4_4_;
  ppcVar36 = uStack_658;
  uStack_548 = uStack_2e0;
  if (3 < _DAT_00a91d28 - 2) {
    bVar20 = DAT_00a90010;
    if (2 < DAT_00a90010) {
      bVar20 = recovered_00086870(&DAT_00a90000);
    }
    ppcVar36 = uStack_658;
    if (bVar20 != 0) {
      if (bVar20 != 2) {
        if (_DAT_00a92030 != 2) goto LAB_001a74de;
        lVar44 = _DAT_00a91d18;
        if (_DAT_00a91d10 == 1) {
          lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                   0x10;
        }
        cVar19 = (**(code **)(_DAT_00a91d20 + 0x28))(lVar44);
        ppcVar36 = uStack_658;
        if (cVar19 == '\0') goto LAB_001a74de;
      }
      pcVar28 = _DAT_00a91d20;
      puStack_5b0 = _DAT_00a90000 + 0x30;
      pcStack_590 = recovered_003740a0;
      uStack_568._0_7_ = 0x726481;
      uStack_568._7_1_ = '\0';
      uStack_560 = SUB87(&pcStack_598,0);
      uStack_559 = (undefined1)((ulong)&pcStack_598 >> 0x38);
      pcStack_668 = (char *)(param_1 + 0x1d);
      uStack_660 = (char **)&uStack_568;
      uStack_658._0_4_ = 0xa8a7e8;
      uStack_658._4_4_ = 0;
      uStack_650 = &pcStack_668;
      uStack_648 = 0xa58960;
      uStack_644 = 0;
      pcStack_5c8 = Elf64_Ehdr_00000000.e_ident_magic_str;
      puStack_5c0 = &uStack_660;
      uStack_5b8 = 2;
      pcStack_628 = (char *)&pcStack_5c8;
      puStack_620 = _DAT_00a90000;
      pcStack_638 = Elf64_Ehdr_00000000.e_ident_magic_str;
      pcStack_598 = (char *)&uStack_548;
      ppcVar36 = (char **)&UNK_00a8a7e8;
      if (_DAT_00a92030 == 2) {
        lVar44 = _DAT_00a91d18;
        if (_DAT_00a91d10 == 1) {
          lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                   0x10;
        }
        cVar19 = (**(code **)(_DAT_00a91d20 + 0x50))(lVar44,&pcStack_638);
        unaff_R14 = pcVar28;
        ppcVar36 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
        if (cVar19 != '\0') {
          (**(code **)(pcVar28 + 0x58))(lVar44,&pcStack_638);
          ppcVar36 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
        }
      }
    }
  }
LAB_001a74de:
  uStack_658 = ppcVar36;
  uStack_528 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
  if (uStack_548 < Elf64_Ehdr_00000000.e_ident_pad + 1) {
    pcVar28 = Elf64_Ehdr_00000000.e_ident_magic_str + 1;
    if (Elf64_Ehdr_00000000.e_ident_magic_str < uStack_548) {
      pcVar28 = uStack_548 + -2;
    }
    if (pcVar28 == Elf64_Ehdr_00000000.e_ident_magic_str + 2) {
      (**(code **)(CONCAT44(uStack_53c,CONCAT22(uStack_53e,CONCAT11(uStack_53f,cStack_540))) + 0x20)
      )(CONCAT44(uStack_528._4_4_,(undefined4)uStack_528),
        CONCAT44(uStack_534,CONCAT22(uStack_536,uStack_538)),
        CONCAT44(uStack_52c,CONCAT22(uStack_52e,uStack_530)));
    }
    else if (pcVar28 == Elf64_Ehdr_00000000.e_ident_magic_str + 1) {
      (**(code **)(CONCAT44(uStack_534,CONCAT22(uStack_536,uStack_538)) + 0x20))
                (pcStack_520,CONCAT44(uStack_52c,CONCAT22(uStack_52e,uStack_530)),
                 CONCAT44(uStack_528._4_4_,(undefined4)uStack_528));
    }
    else {
      uStack_528 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528);
      if ((pcVar28 == Elf64_Ehdr_00000000.e_ident_magic_str) &&
         (uStack_528 = (char **)CONCAT44(uStack_528._4_4_,(undefined4)uStack_528),
         CONCAT44(uStack_52c,CONCAT22(uStack_52e,uStack_530)) != 0)) {
        (*(code *)PTR_DAT_00a8dfe8)(CONCAT44(uStack_528._4_4_,(undefined4)uStack_528));
      }
    }
  }
  goto LAB_001a5446;
LAB_001a718f:
  *puStack_5e0 = 7;
  pcVar28 = unaff_R14;
LAB_001a719a:
  bVar45 = true;
LAB_001a719f:
  if (bVar45) {
    uVar27 = 1;
    uVar32 = 4;
    goto LAB_001a5431;
  }
  recovered_000aa7c0(plStack_558);
  if (pcVar28 != (char *)0x0) {
    pcStack_598 = pcVar28;
    ppcVar36 = uStack_658;
    ppcVar8 = uStack_2d8;
    pcVar33 = uStack_548;
    if (_DAT_00a91d28 != 5) {
      bVar20 = DAT_00a8f7b8;
      if (2 < DAT_00a8f7b8) {
        bVar20 = recovered_00086870(&DAT_00a8f7a8);
      }
      ppcVar36 = uStack_658;
      ppcVar8 = uStack_2d8;
      pcVar33 = uStack_548;
      if (bVar20 != 0) {
        if (bVar20 == 2) {
LAB_001a7201:
          pcVar28 = _DAT_00a91d20;
          puStack_620 = (undefined *)(_DAT_00a8f7a8 + 0x30);
          uStack_660 = &pcStack_598;
          uStack_658._0_4_ = 0x2ddb10;
          uStack_658._4_4_ = 0;
          pcStack_5c8 = &UNK_0072597c;
          puStack_5c0 = &uStack_660;
          pcStack_2d0 = (char *)(param_1 + 0xd);
          uStack_2e0 = (char *)&pcStack_5c8;
          uStack_2d8._0_4_ = 0xa8a7e8;
          uStack_2d8._4_4_ = 0;
          puStack_2c8 = &UNK_00a67980;
          pcStack_638 = Elf64_Ehdr_00000000.e_ident_magic_str;
          pcStack_630 = (char *)&uStack_2e0;
          pcStack_628 = Elf64_Ehdr_00000000.e_ident_magic_str + 1;
          uStack_538 = SUB82(&pcStack_638,0);
          uStack_536 = (undefined2)((ulong)&pcStack_638 >> 0x10);
          uStack_534 = (undefined4)((ulong)&pcStack_638 >> 0x20);
          uStack_530 = (undefined2)_DAT_00a8f7a8;
          uStack_52e = (undefined2)((ulong)_DAT_00a8f7a8 >> 0x10);
          uStack_52c = (undefined4)((ulong)_DAT_00a8f7a8 >> 0x20);
          uStack_548._0_4_ = 1;
          uStack_548._4_4_ = 0;
          ppcVar36 = (char **)recovered_002ddb10;
          ppcVar8 = (char **)&UNK_00a8a7e8;
          pcVar33 = Elf64_Ehdr_00000000.e_ident_magic_str;
          if (_DAT_00a92030 == 2) {
            lVar44 = _DAT_00a91d18;
            if (_DAT_00a91d10 == 1) {
              lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                       0x10;
            }
            cVar19 = (**(code **)(_DAT_00a91d20 + 0x50))(lVar44,&uStack_548);
            ppcVar36 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
            ppcVar8 = (char **)CONCAT44(uStack_2d8._4_4_,(undefined4)uStack_2d8);
            pcVar33 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
            if (cVar19 != '\0') {
              (**(code **)(pcVar28 + 0x58))(lVar44,&uStack_548);
              ppcVar36 = (char **)CONCAT44(uStack_658._4_4_,(undefined4)uStack_658);
              ppcVar8 = (char **)CONCAT44(uStack_2d8._4_4_,(undefined4)uStack_2d8);
              pcVar33 = (char *)CONCAT44(uStack_548._4_4_,(int)uStack_548);
            }
          }
        }
        else if (_DAT_00a92030 == 2) {
          lVar44 = _DAT_00a91d18;
          if (_DAT_00a91d10 == 1) {
            lVar44 = _DAT_00a91d18 + (*(long *)(_DAT_00a91d20 + 0x10) - 1U & 0xfffffffffffffff0) +
                     0x10;
          }
          cVar19 = (**(code **)(_DAT_00a91d20 + 0x28))(lVar44);
          ppcVar36 = uStack_658;
          ppcVar8 = uStack_2d8;
          pcVar33 = uStack_548;
          if (cVar19 != '\0') goto LAB_001a7201;
        }
      }
    }
    uStack_658 = ppcVar36;
    uStack_2d8 = ppcVar8;
    uStack_548 = pcVar33;
    (*(code *)**(undefined8 **)pcStack_598)();
  }
  plVar39 = (long *)param_1[0xb];
  LOCK();
  *plVar39 = *plVar39 + -1;
  UNLOCK();
  if (*plVar39 == 0) {
    recovered_00655a20(param_1[0xb],param_1[0xc]);
  }
  lVar25 = 8;
  lVar44 = *param_1;
joined_r0x001a7733:
  if (lVar44 != 0) {
    (*(code *)PTR_DAT_00a8dfe8)(*(undefined8 *)((long)param_1 + lVar25));
  }
  uVar32 = 1;
  uVar27 = 0;
LAB_001a5431:
  *(undefined1 *)((long)param_1 + 0x6a) = uVar32;
  return uVar27;
}
