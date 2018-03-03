set -e

go install github.com/barnex/godoc2ghmd
out=godoc.md
pre=github.com/barnex/bruteray/
cmd='godoc2ghmd -ex'

$cmd $pre > $out
echo >> $out
echo >> $out
$cmd $pre/shape >> $out

echo >> $out
echo >> $out
echo "- - -" >> $out
echo 'Generated by a modified [godoc2ghmd](https://github.com/GandalfUK/godoc2ghmd)' >> $out
