./build.sh

srcAccountName='goblobcopysrc'
srcAccountKey='2Ex1p6i7UKBPxOKmjLL/s83Ro6fbOHLwNbo1OyIYI8NEWiU5IT52TVTdqhlR1+Z577X2nrg5BRYh+AStQrf22w=='
srcContainerName='src'
dstAccountName='goblobcopydest'
dstContainerName='dst'
dstAccountKey='xxed5stXgMfzgJiBCtY5AQ9KIxnk/kjUTVCsrTQxbO3NIIT9x/9nc1fowr8Yur+GrrEejRljWueQ+ASt3AZCGQ=='
srcBlobName='test.dat'
dstBlobName='test-copy.dat'

./async-blob-copy \
-sa $srcAccountName \
-sk $srcAccountKey \
-sc $srcContainerName \
-da $dstAccountName \
-dk $dstAccountKey \
-dc $dstContainerName \
-sb $srcBlobName \
-db $dstBlobName
