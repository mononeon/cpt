$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot
$CPT = "..\..\cpt.exe"

Write-Host "--- Testing Normal Problem Types ---" -ForegroundColor Cyan

# Create dummy tests
mkdir -Force .cpt\sol_ac
Set-Content -Path .cpt\sol_ac\in_1.txt -Value "1`n5 3`n"
Set-Content -Path .cpt\sol_ac\out_1.txt -Value "8`n"

mkdir -Force .cpt\sol_wa
Set-Content -Path .cpt\sol_wa\in_1.txt -Value "1`n5 3`n"
Set-Content -Path .cpt\sol_wa\out_1.txt -Value "8`n"

# 1. AC Test
Write-Host "`nTesting AC Solution..." -ForegroundColor Yellow
& $CPT test sol_ac --tests 1

# 2. WA Test
Write-Host "`nTesting WA Solution..." -ForegroundColor Yellow
& $CPT test sol_wa --tests 1

# 3. AC Stress Test (Verbose)
Write-Host "`nTesting AC Stress (Verbose, 3 iters)..." -ForegroundColor Yellow
& $CPT stress -v --gen gen.cpp --brute sol_ac.cpp --iters 3 sol_ac

# 4. TLE Test
Write-Host "`nTesting TLE Solution (with 1s limit)..." -ForegroundColor Yellow
& $CPT stress --gen gen.cpp --brute sol_ac.cpp --iters 100 --tl 1s sol_tle

# 5. RE Test
Write-Host "`nTesting RE Solution (100 iters)..." -ForegroundColor Yellow
& $CPT stress --gen gen.cpp --brute sol_ac.cpp --iters 100 sol_re

# 5. CLI Case Management
Write-Host "`nAdding custom test case via CLI..." -ForegroundColor Yellow
& $CPT add sol_ac

Write-Host "`nRunning sol_ac with newly added test case..." -ForegroundColor Yellow
& $CPT test -tests 2 sol_ac

Write-Host "`nRemoving custom test case via CLI..." -ForegroundColor Yellow
& $CPT rm sol_ac 2

Write-Host "`nAll normal tests completed!" -ForegroundColor Green
