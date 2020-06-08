accs_file=$1
num=0

for acc in `cat $accs_file`
do
num=$((num+1))
echo $num:$acc
nchcli q account $acc |grep amount
if [ $? -eq 1 ]; then
echo $acc >>hehe
fi
done
