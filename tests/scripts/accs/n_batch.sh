n=$1
send_book_file=$2

for((i=0;i<$n;i++))
do
	bash send_by_send_book.sh $send_book_file 1pnch testtest $((i+1)) $n
done
