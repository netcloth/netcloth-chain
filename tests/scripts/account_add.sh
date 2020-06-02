#!/usr/bin/expect

# before use this script, you should have expect installed, if not:apt-get install expect
spawn nchcli keys add [lindex $argv 0]
expect {
    "Enter a passphrase to encrypt your key to disk" {
        send "testtest\n";
        exp_continue;
    }
    "Repeat the passphrase" {
        send "testtest\n";
        exp_continue;
    }
}
