srcAccountName='goblobcopysrc'
srcAccountKey='rVPlzjcgyggGWFY1y/VbvaUv7Wk/iUZ7Qvo/a9KXBnvVsbOe/IU2rLiF/jtgSRw47DoYyQmdf5uq+AStdJuNmA=='
srcContainerName='src'
dstAccountName='goblobcopydest'
dstContainerName='dst'
dstAccountKey='X28GbCkFfcNZLcLMCkCFPxYlAvUwjc4leLc9eH5SjWOcBgZLCglozvbQHcsuWdNmKEl7tZwLpGXY+AStOB9lZw=='
blobName='test.dat'

./async-blob-copy \
-sa $srcAccountName \
-sk $srcAccountKey \
-sc $srcContainerName \
-da $dstAccountName \
-dk $dstAccountKey \
-dc $dstContainerName \
-b $blobName
