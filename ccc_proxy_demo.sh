# Address of TES proxy
PROXY='http://localhost:10001'
export PROXY

# Describe sites: Central, PDX, Other

# Show DTS file mock data

# Route to file on PDX site
funnel run -S $PROXY 'echo $in' --name PDX-task --container alpine --in in='ccc://file_pdx'

# Route to file on Other site
funnel run -S $PROXY 'echo $in' --name Other-task --container alpine --in in='ccc://file_other'

# Route to file on Central site
funnel run -S $PROXY 'echo $in' --name Central-task --container alpine --in in='ccc://file_central'

# Default to PDX. Probably needs improvement based on submission site inference.
funnel run -S $PROXY 'echo $in' --container alpine --in in='ccc://file_pdx_central'

# Two files + Fetch_file strategy.
# Should run on PDX, not Central, because files need to be fetched to PDX.
funnel run -S $PROXY 'echo $in $two' --name PDX-fetch --container alpine --in two='ccc://file_pdx' --in in='ccc://file_central' --tag strategy='fetch_file'

# Two files without a shared site.
# Should fail to be scheduled/routed.
funnel run -S $PROXY 'echo $in $two' --name No-site --container alpine --in two='ccc://file_pdx' --in in='ccc://file_other'

# Same as above, but fetch_file.
funnel run -S $PROXY 'echo $in $two' --name No-site-fetch --container alpine --in two='ccc://file_pdx' --in in='ccc://file_other' --tag strategy='fetch_file'

# Same as above, but push_file.
funnel run -S $PROXY 'echo $in $two' --name No-site-push --container alpine --in two='ccc://file_pdx' --in in='ccc://file_other' --tag strategy='push_file'


# Two files, one on Central, and push_file.
# Should fail because the file needs to be on PDX in order to push to central.
funnel run -S $PROXY 'echo $in $two' --name Cannot-push --container alpine --in two='ccc://file_pdx' --in in='ccc://file_central' --tag strategy='push_file'

# Two files, and push_file.
funnel run -S $PROXY 'echo $in $two' --name PDX-push --container alpine --in two='ccc://file_pdx' --in in='ccc://file_pdx_central' --tag strategy='push_file'





#funnel run -S $PROXY 'echo $in $two' --container alpine --in two='ccc://file_pdx' --in in='ccc://file_central'
