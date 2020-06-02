accs_file=$1
num=0

for acc in `cat $accs_file`
do
num=$((num+1))
nchcli q account $acc >/dev/null
echo $num:$acc
if [ $? -eq 1 ]; then
echo $acc >>hehe
fi
done
