meta:
  id: archive_header
  endian: le

seq:
  - id: signature
    type: u1
    repeat: expr
    repeat-expr: 8

enums:
  archive_type:
    0: unknown
    1: zip
    2: rar4
    3: rar5
    4: seven_z
    5: gzip
    6: bzip2
    7: xz
    8: cab
    9: arj
    10: tar

instances:
  archive:
    enum: archive_type
    value: |
      (signature[0] == 0x50 and signature[1] == 0x4b and signature[2] == 0x03 and signature[3] == 0x04) ? 1 :
      (signature[0] == 0x37 and signature[1] == 0x7a and signature[2] == 0xbc and signature[3] == 0xaf and signature[4] == 0x27 and signature[5] == 0x1c) ? 4 :
      (signature[0] == 0x52 and signature[1] == 0x61 and signature[2] == 0x72 and signature[3] == 0x21 and signature[4] == 0x1a and signature[5] == 0x07 and signature[6] == 0x00) ? 2 :
      (signature[0] == 0x52 and signature[1] == 0x61 and signature[2] == 0x72 and signature[3] == 0x21 and signature[4] == 0x1a and signature[5] == 0x07 and signature[6] == 0x01 and signature[7] == 0x00) ? 3 :
      (signature[0] == 0x1F and signature[1] == 0x8B and signature[2] == 0x08) ? 5 :
      (signature[0] == 0x42 and signature[1] == 0x5A and signature[2] == 0x68) ? 6 :
      (signature[0] == 0xFD and signature[1] == 0x37 and signature[2] == 0x7A and signature[3] == 0x58 and signature[4] == 0x5A and signature[5] == 0x00) ? 7 :
      (signature[0] == 0x4D and signature[1] == 0x53 and signature[2] == 0x43 and signature[3] == 0x46) ? 8 :
      (signature[0] == 0x60 and signature[1] == 0xEA) ? 9 :
      0

