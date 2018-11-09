#! /bin/bash

usage() {
    echo ""
    echo "./start -u yuyouqi -h hive-prod -dt '2018-11-11'"
    echo " -u   username - 登录hive客户端机器的用户名"
    echo " -h   host - 客户端所在host"
    echo " -dt  date - 快照的日期精确到日"
    echo " -help     - for this usage message"
    echo ""
}

while [ $# -ne 0 ] ; do
    PARAM=$1
    shift
    case $PARAM in
        -help) usage; exit 0;;
        -u) USERNAME=$1; shift;;
        -h) HOST=$1; shift;;
        -dt) DT=$1; shift;;
        *) usage; exit 0;;
    esac
done

if [ -z "$USERNAME" -o -z "$HOST" -o -z "$DT" ]; then
    usage
    exit 1
fi
RED="\033[0;31m"
NC="\033[0m"
GREEN="\033[0;32m"
echo -e "username: ${RED}${USERNAME}${NC} host ${RED}${HOST}${NC} dt ${RED}${DT}${NC} "
DT2=${DT//-/}

HiveCMD="hive -e"
MRSql="set hive.execution.engine=MR"
NsrSql="select c.cert_name, c.cert_no, c.certification_type, u.gender, u.country_code, u.mobile_number from dwd.dwd_putong_pyment_yay_alipay_certs_a_d c join dwd.dwd_putong_yay_users_a_d u on (c.user_id = u.id) where c.dt='${DT}' and u.dt='${DT}' and c.certification_status = 'passed';"
LwsdSql="select c.cert_name, c.cert_no, c.certification_type, m.income from dwd.dwd_putong_pyment_yay_monthly_incomes_a_d m join dwd.dwd_putong_pyment_yay_alipay_certs_a_d c on (m.user_id = c.user_id)
where m.dt = '${DT}' and c.dt = '${DT}' and c.certification_status = 'passed' and deposit_month = '${DT2:0:6}'"
echo -e ${GREEN}${NsrSql}${NC}
echo -e ${GREEN}${LwsdSql}${NC}

read -p "Continue? [y/n]" confirm && [[ $confirm == [yY] ]] || exit 0
NsrCmd="${HiveCMD}  \"${MRSql};${NsrSql}\"  > nsrjcxx.csv"
LwsdCmd="${HiveCMD}  \"${MRSql};${LwsdSql}\"  > lwbcsd.csv"
echo    " ############ begin hive ############## "
ssh -l ${USERNAME} ${HOST} "source /etc/profile;source ~/.bashrc; ${NsrCmd};${LwsdCmd}"
echo    " ############ copy file  ############### "
scp ${USERNAME}@${HOST}:~/nsrjcxx.csv ./conf/
scp ${USERNAME}@${HOST}:~/lwbcsd.csv ./conf/
read -p " ############ begin write to mysql, confirm?[y/n] ###########" confirm && [[ ${confirm} == [yY] ]] || exit 0
build/bin/tax_tool --config="./conf" --nsr="./conf/nsrjcxx.csv" --lwsd="./conf/lwbcsd.csv" 
echo -e " ############ ${GREEN}success${NC} ################# "

