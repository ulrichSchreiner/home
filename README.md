Settings for my $HOME

Put
```
fn=`readlink -f $1`
dn=`dirname $fn`
export GOPATH=$dn
```

in a file named `.env`. When entering this directory, your GOPATH will be set automatically.
