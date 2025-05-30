# Define RPC credentials and URL
$rpcUser = "username"
$rpcPassword = "password"
$rpcUrl = "http://127.0.0.1:18443/"

# Prepare the descriptor payload
$payload = @{
    jsonrpc = "1.0"
    id = "powershell"
    method = "importdescriptors"
    params = @(
        @(
            @{
                desc = "addr(bcrt1q4qps0ar0q293un44ah05dwukezsp8h8g9e7hyr)"
                timestamp = "now"
                watchonly = $true
                label = "MyWatchOnlyAddress"
                active = $false
            }
        )
    )
}

# Convert payload to JSON
$jsonPayload = $payload | ConvertTo-Json -Depth 10 # -Compress

$jsonPayload

# Prepare authentication
$base64AuthInfo = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes("$rpcUser`:$rpcPassword"))

# Send the RPC call
$response = Invoke-RestMethod -Uri $rpcUrl `
                              -Method Post `
                              -Body $jsonPayload `
                              -Headers @{Authorization=("Basic {0}" -f $base64AuthInfo); "Content-Type"="text/plain"}

# Output the response
$response | ConvertTo-Json -Depth 10

