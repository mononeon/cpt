$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot
$CPT = "..\..\cpt.exe"

Write-Host "--- Testing Communication (Twice-Run) Problem ---" -ForegroundColor Cyan

# Create test cases
mkdir -Force .cpt\sol_ac
Set-Content -Path .cpt\sol_ac\in_1.txt -Value "4`n0000000000000000000000000000000000000000000000000000000000000000000`n1111111111111111111111111111111111111111111111111111111111111111111`n2222222222222222222222222222222222222222222222222222222222222222222`n3333333333333333333333333333333333333333333333333333333333333333333`n"
Set-Content -Path .cpt\sol_ac\out_1.txt -Value ""

mkdir -Force .cpt\sol_wa
Set-Content -Path .cpt\sol_wa\in_1.txt -Value "1`n1234567890123456789012345678901234567890123456789012345678901234567`n"
Set-Content -Path .cpt\sol_wa\out_1.txt -Value ""


Write-Host "`nTesting correct communication solution with verbose logging..." -ForegroundColor Yellow
& $CPT communicate -v --interactor interactor.cpp --test 1 sol_ac

Write-Host "`nTesting WRONG communication solution with verbose logging..." -ForegroundColor Yellow
& $CPT communicate -v --interactor interactor.cpp --test 1 sol_wa

Write-Host "`nAll communication tests completed!" -ForegroundColor Green
