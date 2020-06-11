accs_file=$1

output_file=account_not_exist
rm -rf $output_file

num=0
for acc in `cat $accs_file`
do
num=$((num+1))
echo $num:$acc
nchcli q account $acc |grep amount
if [ $? -eq 1 ]; then
echo $acc >> $output_file
fi
done
