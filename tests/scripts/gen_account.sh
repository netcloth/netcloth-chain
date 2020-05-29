bash batch_add_account.sh acc 1000 |grep "\"address\"" |awk -F '"' '{print $4}' >n1_accs
