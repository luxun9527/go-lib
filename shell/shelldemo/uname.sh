#!/bin/bash
cmdline=$(cat /proc/cmdline |grep 'hpet'|tr -s '\n' )
if [ ! -z "$cmdline"]
then
     echo 'has set'
     exit 0
fi
path=/tmp/phet
mkdir $path
mount /dev/sda1 $path
echo '#Sample GRUB configuration file

      # Boot automatically after 30 secs.
      set timeout=3
      loadfont tos

      if [ "$grub_platform" = "efi" ]; then
              insmod efi_gop
              insmod efi_uga
              insmod gfxterm
              terminal_output gfxterm

              set gfxmode=1024x768,auto
              set menu_color_normal=white/black
              set menu_color_highlight=black/light-gray
              set uefi=" (UEFI)"
              set gfxpayload=keep
      fi

      set default=gnulinux
      # Fallback to GNU/Hurd.
      set fallback=gnuhurd

      # For booting initialize system
      menuentry "UTOS-X86-S64${uefi}" --id gnulinux {
              insmod mdraid1x # load raid1 support for grub

              if [ -e (md/UTOSCORE-X86-S64)/boot/bzImage ]; then
                      echo "----------------------------"
                      echo "- booting from RAID system -"
                      echo "----------------------------"
                      set root=(md/UTOSCORE-X86-S64)
                      linux /boot/bzImage type=raid hpet=disable
              else
                      echo "----------------------------"
                      echo "- booting form INIT system -"
                      echo "----------------------------"

                      insmod search_label
                      search --no-floppy --label UTOSBOOT --set=tmproot
                      if [ -z "${tmproot}" ]; then
                              search --no-floppy --label UTOSBOOT-X86-S64 --set=tmproot
                      fi
                      set from=${tmproot}
                      if [ -z "${tmproot}" ]; then
                              search --no-floppy --label UTOSDISK --set=tmproot
                              set from=${tmproot}
                      fi
                      if [ -z "${tmproot}" ]; then
                              search --no-floppy --label UTOSDISK-X86-S64 --set=tmproot
                              set from=${tmproot}
                      fi
                      set root=${tmproot}
                      linux /boot/bzImage type=usb from=$from
              fi
      }' > "$path/boot/grub/grub.cfg"
umount $path
rm -rf $path