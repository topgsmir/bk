# Run only when publication has been authorized and GitHub CLI is signed in.
$ErrorActionPreference = 'Stop'
$bkRepoRoot = Split-Path -Parent $PSScriptRoot
$bkKeyPath = Join-Path $bkRepoRoot '.publisher/release-signing-key.dpapi'
if (-not (Test-Path -LiteralPath $bkKeyPath)) {
    throw 'Local encrypted signing key is missing. Generate a new key with go run ./tools/releasekey and update the pinned public key before publishing.'
}
gh auth status
if ($LASTEXITCODE -ne 0) { throw 'Sign into GitHub first: gh auth login' }
$bkEncrypted = Get-Content -LiteralPath $bkKeyPath -Raw
$bkSecure = ConvertTo-SecureString -String $bkEncrypted.Trim()
$bkSecret = [Net.NetworkCredential]::new('', $bkSecure).Password
try {
    $bkSecret | gh secret set RELEASE_SIGNING_KEY --repo topgsmir/bk
    if ($LASTEXITCODE -ne 0) { throw 'GitHub did not accept the signing key.' }
} finally {
    $bkSecret = $null
    $bkSecure.Dispose()
}
Write-Output 'Signing key saved as a GitHub Actions secret; no release has been published.'
