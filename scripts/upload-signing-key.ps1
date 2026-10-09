# Run only when publication has been authorized and GitHub CLI is signed in.
$ErrorActionPreference = 'Stop'
$backpackRepoRoot = Split-Path -Parent $PSScriptRoot
$backpackKeyPath = Join-Path $backpackRepoRoot '.publisher/release-signing-key.dpapi'
if (-not (Test-Path -LiteralPath $backpackKeyPath)) {
    throw 'Local encrypted signing key is missing. Generate a new key with go run ./tools/releasekey and update the pinned public key before publishing.'
}
gh auth status
if ($LASTEXITCODE -ne 0) { throw 'Sign into GitHub first: gh auth login' }
$backpackEncrypted = Get-Content -LiteralPath $backpackKeyPath -Raw
$backpackSecure = ConvertTo-SecureString -String $backpackEncrypted.Trim()
$backpackSecret = [Net.NetworkCredential]::new('', $backpackSecure).Password
try {
    $backpackSecret | gh secret set RELEASE_SIGNING_KEY --repo topgsmir/BackPack
    if ($LASTEXITCODE -ne 0) { throw 'GitHub did not accept the signing key.' }
} finally {
    $backpackSecret = $null
    $backpackSecure.Dispose()
}
Write-Output 'Signing key saved as a GitHub Actions secret; no release has been published.'
