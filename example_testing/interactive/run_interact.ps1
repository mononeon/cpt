$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot
$CPT = "..\..\cpt.exe"

Write-Host "--- Testing Interactive Problem ---" -ForegroundColor Cyan

# Create test case directory and in_1.txt
mkdir -Force .cpt\sol
Set-Content -Path .cpt\sol\in_1.txt -Value "2`n60 1 123456789012345678`n60 3 987654321098765432`n"
Set-Content -Path .cpt\sol\out_1.txt -Value ""

Write-Host "`nTesting correct interactive solution with verbose logging..." -ForegroundColor Yellow
& $CPT interact -v --interactor interactor.cpp --test 1 sol

Write-Host "`nTesting WRONG interactive solution with verbose logging..." -ForegroundColor Yellow
& $CPT interact -v --interactor interactor.cpp --test 1 sol_wa

Write-Host "`nAll interactive tests completed!" -ForegroundColor Green
