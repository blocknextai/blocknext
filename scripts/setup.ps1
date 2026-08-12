#Requires -Version 5.1

$ErrorActionPreference = 'Stop'

Set-Location (Split-Path -Parent $PSScriptRoot)

if (Test-Path .env) {
    Write-Host '.env already exists - leaving it untouched'
    exit 0
}

Copy-Item .env.example .env

$text = Get-Content .env -Raw
$tokens = [regex]::Matches($text, 'REPLACE_ME_OPENSSL_[A-Z_]+') | ForEach-Object { $_.Value } | Sort-Object -Unique

foreach ($token in $tokens) {
    $bytes = [byte[]]::new(32)
    [System.Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($bytes)
    $text = $text -replace $token, ([System.BitConverter]::ToString($bytes) -replace '-', '').ToLower()
}

Set-Content .env $text -NoNewline

if ($text -match 'REPLACE_ME_OPENSSL') {
    Write-Host 'warning: unreplaced secret placeholders remain in .env'
    exit 1
}

Write-Host ".env created with generated secrets - run 'docker compose -f docker-compose.prod.yml up -d' to start the stack"
