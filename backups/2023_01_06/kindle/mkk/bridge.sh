#!/bin/sh
#
# Kindle Touch/PaperWhite JailBreak Bridge
#
# $Id: bridge.sh 18869 2022-02-13 19:01:52Z NiLuJe $
#
##

BRIDGE_REV="$( echo '$Revision: 18869 $' | cut -d ' ' -f 2 )"

ROOT=""
ROOTPART="mmcblk0p1"
IS_TOUCH="false"
IS_PW="false"
IS_PW2="false"
IS_KV="false"
IS_KT2="false"
IS_PW3="false"
IS_KOA="false"
IS_KT3="false"
IS_KOA2="false"
IS_PW4="false"
IS_KT4="false"
IS_KOA3="false"
IS_PW5="false"
K5_ATLEAST_54="false"
K5_ATLEAST_512="false"
MKK_PERSISTENT_STORAGE="/var/local/mkk"
RP_PERSISTENT_STORAGE="/var/local/rp"
VARLOCAL_OOS="false"

# Pull some helper functions for logging
source /etc/upstart/functions

LOG_DOMAIN="jb_bridge"

logmsg()
{
	f_log "${1}" "${LOG_DOMAIN}" "${2}" "${3}" "${4}"
}

RW=""
mount_rw() {
	if [ -z "${RW}" ] ; then
		RW="yes"
		mount -o rw,remount /
	fi
}

mount_ro() {
	if [ -n "${RW}" ] ; then
		RW=""
		mount -o ro,remount /
	fi
}

make_mutable() {
	local my_path="${1}"
	# NOTE: Can't do that on symlinks, hence the hoop-jumping...
	if [ -d "${my_path}" ] ; then
		find "${my_path}" -type d -exec chattr -i '{}' \;
		find "${my_path}" -type f -exec chattr -i '{}' \;
	elif [ -f "${my_path}" ] ; then
		chattr -i "${my_path}"
	fi
}

make_immutable() {
	local my_path="${1}"
	if [ -d "${my_path}" ] ; then
		find "${my_path}" -type d -exec chattr +i '{}' \;
		find "${my_path}" -type f -exec chattr +i '{}' \;
	elif [ -f "${my_path}" ] ; then
		chattr +i "${my_path}"
	fi
}

find_rootfs()
{
	# We need to know our model first...
	check_model

	# Handle the fact that the KOA2 switched to an all-new partition layout (àla Android)
	# NOTE: For some mysterious reason, rdev's output doesn't match what we can get from the elektra db via:
	#       echo "$(kdb get system/driver/filesystem/DEV_ROOT)$(kdb get system/driver/filesystem/DEV_PART_ROOTFS)"
	#       while that was true on previous models...
	#       FWIW, kdb returns /dev/mmcblk2p4 instead of rdev's p5... :?
	# NOTE: Assume the KOA3 behaves the same, as differences between the two appear to be minimal...
	#       At least the kdb output matches the KOA2, still.
	# NOTE: c.f., /etc/default/layout on recent FW versions
	if [ "${IS_KOA2}" == "true" ] || [ "${IS_KOA3}" == "true" ] ; then
		ROOTPART="mmcblk2p5"
	elif [ "${IS_PW4}" == "true" ] || [ "${IS_KT4}" == "true" ] ; then
		# NOTE: At least this time, rdev agrees with kdb...
		ROOTPART="mmcblk1p8"
	elif [ "${IS_PW5}" == "true" ] ; then
		# NOTE: TBC!
		ROOTPART="mmcblk0p8"
	fi

	# Check if we're running from main...
	DEV="$(rdev | awk '{ print $1 }')"
	# Don't do something stupid if rdev is missing for some reason...
	if [ "${DEV}" != "/dev/${ROOTPART}" ] && [ -n "${DEV}" ] ; then
		ROOT="/var/tmp/rootfs"
	fi
}

mount_root_rw()
{
	# Make sure we use the right rootfs ;p
	find_rootfs

	if [ -n "${ROOT}" ] ; then
		logmsg "I" "mount_root_rw" "" "Running from diags, mounting main rootfs ${ROOTPART} on ${ROOT}"
		mkdir -p "${ROOT}"
		# NOTE: We enforce an initial mount to be able to run the checks,
		#       so account for the fact that we can effectively call mount twice, in which case,
		#       handle the error by issuing a remount instead.
		mount -o rw "/dev/${ROOTPART}" "${ROOT}" || mount -o rw,remount "/dev/${ROOTPART}" "${ROOT}"
	else
		logmsg "I" "mount_root_rw" "" "Mounting rootfs rw"
		mount_rw
	fi
}

check_model()
{
	# Do the S/N dance...
	kmfc="$(cut -c1 /proc/usid)"
	if [ "${kmfc}" == "B" ] || [ "${kmfc}" == "9" ] ; then
		# Older device ID scheme
		kmodel="$(cut -c3-4 /proc/usid)"
		case "${kmodel}" in
			"24" | "1B" | "1D" | "1F" | "1C" | "20" )
				# PaperWhite 1 (2012)
				IS_PW="true"
			;;
			"D4" | "5A" | "D5" | "D6" | "D7" | "D8" | "F2" | "17" | "60" | "F4" | "F9" | "62" | "61" | "5F" )
				# PaperWhite 2 (2013)
				IS_PW="true"
				IS_PW2="true"
			;;
			"13" | "54" | "2A" | "4F" | "52" | "53" )
				# Voyage...
				IS_KV="true"
			;;
			"C6" | "DD" )
				# KT2...
				IS_TOUCH="true"
				IS_KT2="true"
			;;
			"0F" | "11" | "10" | "12" )
				# Touch
				IS_TOUCH="true"
			;;
			* )
				# Fallback... We shouldn't ever hit that.
				IS_TOUCH="true"
			;;
		esac
	else
		# Try the new device ID scheme...
		kmodel="$(cut -c4-6 /proc/usid)"
		case "${kmodel}" in
			"0G1" | "0G2" | "0G4" | "0G5" | "0G6" | "0G7" | "0KB" | "0KC" | "0KD" | "0KE" | "0KF" | "0KG" | "0LK" | "0LL" )
				# PW3...
				IS_PW3="true"
			;;
			"0GC" | "0GD" | "0GR" | "0GS" | "0GT" | "0GU" )
				# Oasis...
				IS_KOA="true"
			;;
			"0DU" | "0K9" | "0KA" )
				# KT3...
				IS_KT3="true"
			;;
			"0LM" | "0LN" | "0LP" | "0LQ" | "0P1" | "0P2" | "0P6" | "0P7" | "0P8" | "0S1" | "0S2" | "0S3" | "0S4" | "0S7" | "0SA" )
				# KOA2...
				IS_KOA2="true"
			;;
			"0PP" | "0T1" | "0T2" | "0T3" | "0T4" | "0T5" | "0T6" | "0T7" | "0TJ" | "0TK" | "0TL" | "0TM" | "0TN" | "102" | "103" | "16Q" | "16R" | "16S" | "16T" | "16U" | "16V" )
				# PW4...
				IS_PW4="true"
			;;
			"10L" | "0WF" | "0WG" | "0WH" | "0WJ" | "0VB" )
				# KT4...
				IS_KT4="true"
			;;
			"11L" | "0WQ" | "0WP" | "0WN" | "0WM" | "0WL" )
				# KOA3...
				IS_KOA3="true"
			;;
			"1LG" | "1Q0" | "1PX" | "1VD" | "219" | "21A" | "2BH" | "2BJ" )
				# PW5...
				IS_PW5="true"
			;;
			* )
				# Fallback... We shouldn't ever hit that.
				IS_TOUCH="true"
			;;
		esac
	fi

	# Use the proper constants for our screen...
	if [ "${IS_KV}" == "true" ] || [ "${IS_PW3}" == "true" ] || [ "${IS_KOA}" == "true" ] || [ "${IS_PW4}" == "true" ] ; then
		SCREEN_X_RES=1088
		SCREEN_Y_RES=1448
		EIPS_X_RES=16
		EIPS_Y_RES=24
	elif [ "${IS_PW}" == "true" ] ; then
		SCREEN_X_RES=768
		SCREEN_Y_RES=1024
		EIPS_X_RES=16
		EIPS_Y_RES=24
	elif [ "${IS_KT2}" == "true" ] || [ "${IS_KT3}" == "true" ] || [ "${IS_KT4}" == "true" ] ; then
		SCREEN_X_RES=608
		SCREEN_Y_RES=800
		EIPS_X_RES=16
		EIPS_Y_RES=24
	elif [ "${IS_KOA2}" == "true" ] || [ "${IS_KOA3}" == "true" ] ; then
		SCREEN_X_RES=1280
		SCREEN_Y_RES=1680
		EIPS_X_RES=16
		EIPS_Y_RES=24
	elif [ "${IS_PW5}" == "true" ] ; then
		SCREEN_X_RES=1236
		SCREEN_Y_RES=1648
		EIPS_X_RES=16
		EIPS_Y_RES=24
	else
		SCREEN_X_RES=600
		SCREEN_Y_RES=800
		EIPS_X_RES=12
		EIPS_Y_RES=20
	fi
	EIPS_MAXCHARS="$((${SCREEN_X_RES} / ${EIPS_X_RES}))"
	EIPS_MAXLINES="$((${SCREEN_Y_RES} / ${EIPS_Y_RES}))"
}

check_version()
{
	# The great version check!
	kpver="$(grep '^Kindle 5' ${ROOT}/etc/prettyversion.txt 2>&1)"
	if [ $? -ne 0 ] ; then
		logmsg "W" "check_version" "" "couldn't detect the kindle major version!"
		# We're in a bit of a pickle... Make an educated guess...
		if [ "${IS_PW2}" == "true" ] ; then
			# The PW2 shipped on 5.4.0 ;)
			logmsg "I" "check_version" "" "PW2 detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
		elif [ "${IS_KV}" == "true" ] ; then
			# The KV shipped on 5.5.0 ;)
			logmsg "I" "check_version" "" "KV detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
		elif [ "${IS_KT2}" == "true" ] ; then
			# The KT2 shipped on 5.6.0 ;)
			logmsg "I" "check_version" "" "KT2 detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
		elif [ "${IS_PW3}" == "true" ] ; then
			# The PW3 shipped on 5.6.1 ;)
			logmsg "I" "check_version" "" "PW3 detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
		elif [ "${IS_KOA}" == "true" ] ; then
			# The Oasis shipped on 5.7.1.1 ;)
			logmsg "I" "check_version" "" "Oasis detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
		elif [ "${IS_KT3}" == "true" ] ; then
			# The KT3 shipped on >= 5.7.x ;)
			logmsg "I" "check_version" "" "KT3 detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
		elif [ "${IS_KOA2}" == "true" ] ; then
			# The KOA2 shipped on >= 5.9.0.x ;)
			logmsg "I" "check_version" "" "KOA2 detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
		elif [ "${IS_PW4}" == "true" ] ; then
			# The PW4 shipped on >= 5.10.0.x ;)
			logmsg "I" "check_version" "" "PW4 detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
		elif [ "${IS_KT4}" == "true" ] ; then
			# The KT4 shipped on >= 5.1x.y ;)
			logmsg "I" "check_version" "" "KT4 detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
		elif [ "${IS_KOA3}" == "true" ] ; then
			# The KOA3 shipped on >= 5.12.y ;)
			logmsg "I" "check_version" "" "KOA3 detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
			K5_ATLEAST_512="true"
		elif [ "${IS_PW5}" == "true" ] ; then
			# The PW5 shipped on >= 5.14.y ;)
			logmsg "I" "check_version" "" "PW5 detected, assuming >= 5.4"
			K5_ATLEAST_54="true"
			K5_ATLEAST_512="true"
		else
			# Poor man's last resort trick. See if we can find a new feature of FW 5.4 on the FS...
			if [ -f "${ROOT}/etc/upstart/contentpackd.conf" ] ; then
				logmsg "I" "check_version" "" "found a fw >= 5.4 feature"
				K5_ATLEAST_54="true"
			fi
			# NOTE: Alternative checks:
			# -x ${ROOT}/usr/bin/contentpackd
			# -f ${ROOT}/opt/amazon/ebook/lib/VocabBuilderSDK.jar
			# -f ${ROOT}/opt/amazon/ebook/booklet/VocabBuilderBooklet.jar
		fi
	else
		# Weeee, the great case switch!
		khver="$(echo ${kpver} | sed -n -r 's/^(Kindle)([[:blank:]]*)([[:digit:].]*)(.*?)$/\3/p')"
		case "${khver}" in
			5.0.* | 5.0 )
				K5_ATLEAST_54="false"
			;;
			5.1.* | 5.1 )
				K5_ATLEAST_54="false"
			;;
			5.2.* | 5.2 )
				K5_ATLEAST_54="false"
			;;
			5.3.* | 5.3 )
				K5_ATLEAST_54="false"
			;;
			5.4.* | 5.4 )
				K5_ATLEAST_54="true"
			;;
			5.5.* | 5.5 )
				K5_ATLEAST_54="true"
			;;
			5.6.* | 5.6 )
				K5_ATLEAST_54="true"
			;;
			5.7.* | 5.7 )
				K5_ATLEAST_54="true"
			;;
			5.8.* | 5.8 )
				K5_ATLEAST_54="true"
			;;
			5.9.* | 5.9 )
				K5_ATLEAST_54="true"
			;;
			5.10.* | 5.10 )
				K5_ATLEAST_54="true"
			;;
			5.11.* | 5.11 )
				K5_ATLEAST_54="true"
			;;
			5.12.* | 5.12 )
				K5_ATLEAST_54="true"
				K5_ATLEAST_512="true"
			;;
			5.13.* | 5.13 )
				K5_ATLEAST_54="true"
				K5_ATLEAST_512="true"
			;;
			5.14.* | 5.14 )
				K5_ATLEAST_54="true"
				K5_ATLEAST_512="true"
			;;
			5.* )
				# Assume newer, just to be safe ;)
				K5_ATLEAST_54="true"
				K5_ATLEAST_512="true"
			;;
			* )
				# Given the previous checks, this shouldn't be reachable, but cover all bases anyway...
				logmsg "W" "check_version" "" "couldn't detect the kindle version!"
				# Poor man's last resort trick. See if we can find a new feature of FW 5.4 on the FS...
				if [ -f "${ROOT}/etc/upstart/contentpackd.conf" ] ; then
					logmsg "I" "check_version" "" "found a fw >= 5.4 feature"
					K5_ATLEAST_54="true"
				fi
			;;
		esac
	fi
}

print_jb_install_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff...
	kh_eips_string="**** JAILBREAK ****"

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}

print_fw54_exec_install_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff...
	kh_eips_string="**** FW 5.4 JB ****"

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}

print_fw512_debugging_install_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff...
	kh_eips_string="**** 5.12.x :( ****"

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}


print_mkk_dev_keystore_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff... Print an extra warning if we failed to copy the keys...
	if [ "${VARLOCAL_OOS}" == "true" ] ; then
		kh_eips_string="**** WARNING: FAILED TO COPY MKK KEYS ****"
	else
		kh_eips_string="**** MKK KEYS **** "
	fi

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}

print_mkk_kindlet_jb_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff...
	kh_eips_string="**** MKK K JB **** "

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}

print_gandalf_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff...
	kh_eips_string=" **** GANDALF **** "

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}

print_rp_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff...
	kh_eips_string="   **** RP ****    "

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}

print_crp_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff...
	kh_eips_string="   **** CRP ****   "

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}

print_dispatch_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff...
	kh_eips_string="**** DISPATCH **** "

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}

print_log_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff...
	kh_eips_string="   **** LOG ****   "

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}

print_bridge_job_feedback()
{
	# We need to know our model
	check_model
	# Prepare our stuff...
	kh_eips_string=" **** BRIDGE+ **** "

	# And finally, show our message, centered on the bottom of the screen
	eips $(((${EIPS_MAXCHARS} - ${#kh_eips_string}) / 2)) $((${EIPS_MAXLINES} - 2)) "${kh_eips_string}"
}

install_touch_update_key()
{
	mount_root_rw
	logmsg "I" "install_touch_update_key" "" "Copying the jailbreak updater key"
	make_mutable "${ROOT}/etc/uks/pubdevkey01.pem"
	rm -rf "${ROOT}/etc/uks/pubdevkey01.pem"
	cat > "${ROOT}/etc/uks/pubdevkey01.pem" << EOF
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDJn1jWU+xxVv/eRKfCPR9e47lP
WN2rH33z9QbfnqmCxBRLP6mMjGy6APyycQXg3nPi5fcb75alZo+Oh012HpMe9Lnp
eEgloIdm1E4LOsyrz4kttQtGRlzCErmBGt6+cAVEV86y2phOJ3mLk0Ek9UQXbIUf
rvyJnS2MKLG2cczjlQIDAQAB
-----END PUBLIC KEY-----
EOF
	# Harmonize permissions
	chown root:root "${ROOT}/etc/uks/pubdevkey01.pem"
	chmod 0644 "${ROOT}/etc/uks/pubdevkey01.pem"
	make_immutable "${ROOT}/etc/uks/pubdevkey01.pem"
	mount_ro

	# Show some feedback
	print_jb_install_feedback
}

install_fw54_exec_userstore_flag()
{
	# Make sure we're on FW >= 5.4...
	check_model
	check_version

	if [ "${K5_ATLEAST_54}" == "true" ] ; then
		mount_root_rw
		logmsg "I" "install_fw54_exec_userstore_flag" "" "Creating the userstore exec flag file"
		make_mutable "${ROOT}/MNTUS_EXEC"
		rm -rf "${ROOT}/MNTUS_EXEC"
		touch "${ROOT}/MNTUS_EXEC"
		make_immutable "${ROOT}/MNTUS_EXEC"
		mount_ro

		# Show some feedback
		print_fw54_exec_install_feedback
	fi
}

install_fw512_debugging_flag()
{
	# Make sure we're on FW >= 5.12...
	check_model
	check_version

	if [ "${K5_ATLEAST_512}" == "true" ] ; then
		mount_root_rw
		logmsg "I" "install_fw512_debugging_flag" "" "Creating the debugging flag file"
		make_mutable "${ROOT}/PRE_GM_DEBUGGING_FEATURES_ENABLED__REMOVE_AT_GMC"
		rm -rf "${ROOT}/PRE_GM_DEBUGGING_FEATURES_ENABLED__REMOVE_AT_GMC"
		touch "${ROOT}/PRE_GM_DEBUGGING_FEATURES_ENABLED__REMOVE_AT_GMC"
		make_immutable "${ROOT}/PRE_GM_DEBUGGING_FEATURES_ENABLED__REMOVE_AT_GMC"
		mount_ro

		# Show some feedback
		print_fw512_debugging_install_feedback
	fi
}

install_mkk_dev_keystore()
{
	logmsg "I" "install_mkk_dev_keystore" "" "Copying the kindlet keystore"
	# We shouldn't need to do anything specific to read/write /var/local
	if [ "$(df -k /var/local | awk '$3 ~ /[0-9]+/ { print $4 }')" -lt "512" ] ; then
		# Hu ho... Keep track of this...
		VARLOCAL_OOS="true"
		logmsg "W" "install_mkk_dev_keystore" "" "Failed to copy the kindlet keystore: not enough space left on device"
	else
		# NOTE: This might have gone poof on newer devices without Kindlet support, so, create it as needed
		mkdir -p "/var/local/java/keystore"
		cp -f "${MKK_PERSISTENT_STORAGE}/developer.keystore" "/var/local/java/keystore/developer.keystore"
	fi

	# Show some feedback
	print_mkk_dev_keystore_feedback
}

install_mkk_kindlet_jb()
{
	mount_root_rw
	logmsg "I" "install_mkk_kindlet_jb" "" "Copying the kindlet jailbreak"
	cp -f "${MKK_PERSISTENT_STORAGE}/json_simple-1.1.jar" "${ROOT}/opt/amazon/ebook/lib/json_simple-1.1.jar"
	chmod 0664 "${ROOT}/opt/amazon/ebook/lib/json_simple-1.1.jar"
	mount_ro

	# Show some feedback
	print_mkk_kindlet_jb_feedback
}

setup_gandalf()
{
	logmsg "I" "setup_gandalf" "" "Setting up gandalf... you shall not pass!"
	make_mutable "${MKK_PERSISTENT_STORAGE}"
	chown root:root "${MKK_PERSISTENT_STORAGE}/gandalf"
	chmod a+rx "${MKK_PERSISTENT_STORAGE}/gandalf"
	chmod +s "${MKK_PERSISTENT_STORAGE}/gandalf"
	ln -sf "${MKK_PERSISTENT_STORAGE}/gandalf" "${MKK_PERSISTENT_STORAGE}/su"
	make_immutable "${MKK_PERSISTENT_STORAGE}"

	# Show some feedback
	print_gandalf_feedback
}

install_rp()
{
	mount_root_rw
	logmsg "I" "install_rp" "" "Copying the RP"
	make_mutable "${ROOT}/etc/upstart/debrick.conf"
	cp -f "${RP_PERSISTENT_STORAGE}/debrick.conf" "${ROOT}/etc/upstart/debrick.conf"
	chmod 0664 "${ROOT}/etc/upstart/debrick.conf"
	make_immutable "${ROOT}/etc/upstart/debrick.conf"
	make_mutable "${ROOT}/etc/upstart/debrick"
	cp -f "${RP_PERSISTENT_STORAGE}/debrick" "${ROOT}/etc/upstart/debrick"
	chmod 0755 "${ROOT}/etc/upstart/debrick"
	make_immutable "${ROOT}/etc/upstart/debrick"
	mount_ro

	# Show some feedback
	print_rp_feedback
}

install_crp()
{
	mount_root_rw
	logmsg "I" "install_crp" "" "Copying the CRP"
	make_mutable "${ROOT}/etc/upstart/cowardsdebrick.conf"
	cp -f "${RP_PERSISTENT_STORAGE}/cowardsdebrick.conf" "${ROOT}/etc/upstart/cowardsdebrick.conf"
	chmod 0664 "${ROOT}/etc/upstart/cowardsdebrick.conf"
	make_immutable "${ROOT}/etc/upstart/cowardsdebrick.conf"
	# My version of CRP doesn't use a separate script ;)
	if [ -f "${RP_PERSISTENT_STORAGE}/cowardsdebrick" ] ; then
		cp -f "${RP_PERSISTENT_STORAGE}/cowardsdebrick" "${ROOT}/etc/upstart/cowardsdebrick"
		chmod 0755 "${ROOT}/etc/upstart/cowardsdebrick"
	fi
	mount_ro

	# Show some feedback
	print_crp_feedback
}

install_dispatch()
{
	mount_root_rw
	logmsg "I" "install_dispatch" "" "Copying the dispatch script"
	make_mutable "${ROOT}/usr/bin/logThis.sh"
	rm -rf "${ROOT}/usr/bin/logThis.sh"
	cp -f "${MKK_PERSISTENT_STORAGE}/dispatch.sh" "${ROOT}/usr/bin/logThis.sh"
	chmod 0755 "${ROOT}/usr/bin/logThis.sh"
	make_immutable "${ROOT}/usr/bin/logThis.sh"
	mount_ro

	# Show some feedback
	print_dispatch_feedback
}

install_log()
{
	mount_root_rw
	logmsg "I" "install_log" "" "Patching in the dispatch command"
	sed -e '/^{/a\' -e '    ";log" : "/usr/bin/logThis.sh",' -i "${ROOT}/usr/share/webkit-1.0/pillow/debug_cmds.json"
	mount_ro

	# Show some feedback
	print_log_feedback
}

install_bridge_job()
{
	mount_root_rw
	logmsg "I" "install_bridge_job" "" "Copying the bridge job"
	make_mutable "${ROOT}/etc/upstart/bridge.conf"
	rm -rf "${ROOT}/etc/upstart/bridge.conf"
	cp -f "${MKK_PERSISTENT_STORAGE}/bridge.conf" "${ROOT}/etc/upstart/bridge.conf"
	chmod 0664 "${ROOT}/etc/upstart/bridge.conf"
	make_immutable "${ROOT}/etc/upstart/bridge.conf"
	mount_ro

	# Show some feedback
	print_bridge_job_feedback
}

clean_up()
{
	# Unmount main rootfs if we're on diags...
	if [ -n "${ROOT}" ] ; then
		logmsg "I" "clean_up" "" "Unmounting main rootfs"
		umount "${ROOT}"
	fi
}

# Here we go...
logmsg "I" "main" "" "i can fix this (r${BRIDGE_REV})"

# We'll begin by checking where our rootfs is...
find_rootfs

# And if we're in diags, we'll need to mount it first, otherwise all those ${ROOT} checks will be useless ;).
if [ -n "${ROOT}" ] ; then
	mount_root_rw
fi

# Start with the userstore exec flag on FW >= 5.4 (so that the last eips print shown will make sense)
if [ ! -f "${ROOT}/MNTUS_EXEC" ] ; then
	install_fw54_exec_userstore_flag
fi

# Check if we need to do something with the OTA pubkey
if [ ! -f "${ROOT}/etc/uks/pubdevkey01.pem" ] ; then
	# No jailbreak key, install it
	install_touch_update_key
else
	# Jailbreak key found... Check it.
	if [ "$(md5sum "${ROOT}/etc/uks/pubdevkey01.pem" | awk '{ print $1; }')" != "7130ce39bb3596c5067cabb377c7a9ed" ] ; then
		# Unknown (?!) jailbreak key, install it
		install_touch_update_key
	fi
	if [ ! -O "${ROOT}/etc/uks/pubdevkey01.pem" ] || [ ! -G "${ROOT}/etc/uks/pubdevkey01.pem" ] ; then
		# Not our own? Make it so!
		install_touch_update_key
	fi
fi

# Check if we need to do something with the Kindlet developer keystore
if [ -f "${MKK_PERSISTENT_STORAGE}/developer.keystore" ] ; then
	# No developer keystore, install it
	if [ ! -f "/var/local/java/keystore/developer.keystore" ] ; then
		install_mkk_dev_keystore
	else
		# Developer keystore doesn't match, install it
		# NOTE: This *will* mess with real, official developer keystores. Not that we really care about it, but it should be noted ;).
		if [ "$(md5sum "/var/local/java/keystore/developer.keystore" | awk '{ print $1; }')" != "$(md5sum "${MKK_PERSISTENT_STORAGE}/developer.keystore" | awk '{ print $1; }')" ] ; then
			install_mkk_dev_keystore
		fi
	fi
fi

# Check if we need to do something with the Kindlet JB
if [ -f "${MKK_PERSISTENT_STORAGE}/json_simple-1.1.jar" ] ; then
	# Kindlet JB doesn't match, install it
	if [ "$(md5sum "${ROOT}/opt/amazon/ebook/lib/json_simple-1.1.jar" | awk '{ print $1; }')" != "$(md5sum "${MKK_PERSISTENT_STORAGE}/json_simple-1.1.jar" | awk '{ print $1; }')" ] ; then
		install_mkk_kindlet_jb
	fi
fi

# Check if we need to do something with Gandalf
if [ -f "${MKK_PERSISTENT_STORAGE}/gandalf" ] ; then
	# NOTE: The bridge job already does this, too.
	if [ ! -O "${MKK_PERSISTENT_STORAGE}/gandalf" ] || [ ! -G "${MKK_PERSISTENT_STORAGE}/gandalf" ] ; then
		setup_gandalf
	fi
	if [ ! -x "${MKK_PERSISTENT_STORAGE}/gandalf" ] ; then
		setup_gandalf
	fi
	if [ ! -u "${MKK_PERSISTENT_STORAGE}/gandalf" ] ; then
		setup_gandalf
	fi
	if [ ! -h "${MKK_PERSISTENT_STORAGE}/su" ] ; then
		setup_gandalf
	fi
	if [ ! -x "${MKK_PERSISTENT_STORAGE}/su" ] ; then
		setup_gandalf
	fi
	# NOTE: This will actually end up a NOOP, because -O & -G tests don't behave all that well with symlinks...
	if [ ! -O "${MKK_PERSISTENT_STORAGE}/su" ] || [ ! -G "${MKK_PERSISTENT_STORAGE}/su" ] ; then
		setup_gandalf
	fi
fi

# Check if we need to do something with the RP
if [ -f "${RP_PERSISTENT_STORAGE}/debrick.conf" ] ; then
	if [ ! -f "${ROOT}/etc/upstart/debrick.conf" ] ; then
		install_rp
	fi
fi

# Check if we need to do something with the CRP
if [ -f "${RP_PERSISTENT_STORAGE}/cowardsdebrick.conf" ] ; then
	if [ ! -f "${ROOT}/etc/upstart/cowardsdebrick.conf" ] ; then
		install_crp
	fi
fi

# Check if we need to do something with the dispatch script
if [ -f "${MKK_PERSISTENT_STORAGE}/dispatch.sh" ] ; then
	if [ ! -f "${ROOT}/usr/bin/logThis.sh" ] ; then
		install_dispatch
	else
		# If it's not ours, install it
		if ! grep -q "Dispatch" "${ROOT}/usr/bin/logThis.sh" ; then
			install_dispatch
		fi
	fi
fi

# Check if we need to do something about the dispatch command
if [ -f "${ROOT}/usr/share/webkit-1.0/pillow/debug_cmds.json" ] ; then
	if ! grep -q "logThis.sh" "${ROOT}/usr/share/webkit-1.0/pillow/debug_cmds.json" ; then
		install_log
	fi
fi

# Check if we need to do something with the bridge job
if [ -f "${MKK_PERSISTENT_STORAGE}/bridge.conf" ] ; then
	if [ ! -f "${ROOT}/etc/upstart/bridge.conf" ] ; then
		install_bridge_job
	fi
fi

# And finish with the last ditch effort for FW 5.12.x...
# NOTE: This flag has far ranging effects, so, here be dragons.
if [ ! -f "${ROOT}/PRE_GM_DEBUGGING_FEATURES_ENABLED__REMOVE_AT_GMC" ] ; then
	install_fw512_debugging_flag
fi

# Nothing to do or cleanup...
clean_up

logmsg "I" "main" "" "these are not the droids you're looking for"

# And don't try anything fancier, the userstore isn't mounted yet...

return 0
