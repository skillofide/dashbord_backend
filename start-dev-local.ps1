# start-dev-local.ps1
# Starts all Go microservices locally WITHOUT Docker.
# Prerequisites: PostgreSQL, Redis, and NATS must be running natively on your machine.
# Run from the project root: .\start-dev-local.ps1

$ErrorActionPreference = "Stop"

Clear-Host
Write-Host '==================================================' -ForegroundColor Magenta
Write-Host '  SkillofIDE Backend - Local Launcher (No Docker) ' -ForegroundColor Magenta
Write-Host '==================================================' -ForegroundColor Magenta
Write-Host ''
Write-Host 'Assumptions:' -ForegroundColor Cyan
Write-Host '  PostgreSQL running on localhost:5432 (user: skillofide, pass: password, db: skillofide)' -ForegroundColor Gray
Write-Host '  Redis running on localhost:6379' -ForegroundColor Gray
Write-Host '  NATS running on localhost:4222' -ForegroundColor Gray
Write-Host ''

$ScriptDir = $PSScriptRoot
if (-not $ScriptDir) {
    $ScriptDir = $PWD
}
cd $ScriptDir

# Build service launch table
$Services = [ordered]@{
    'Problem Service'      = 'services/problem-service'
    'Execution Service'    = 'services/execution-service'
    'Submission Service'   = 'services/submission-service'
    'Progress Service'     = 'services/progress-service'
    'Notification Service' = 'services/notification-service'
    'API Gateway'          = 'services/api-gateway'
}

Write-Host 'Launching backend microservices...' -ForegroundColor Yellow

foreach ($SvcName in $Services.Keys) {
    $SvcPath = $Services[$SvcName]
    Write-Host ('  Launching ' + $SvcName + '...') -ForegroundColor Cyan
    
    # API Gateway also needs POSTGRES_DSN for user auth
    $EnvSetup = ''
    if ($SvcName -eq 'API Gateway') {
        $EnvSetup = '$env:POSTGRES_DSN="postgres://skillofide:password@localhost:5432/skillofide?sslmode=disable";'
    }
    
    $CommandStr = '$Host.UI.RawUI.WindowTitle = "{0}"; cd "{1}"; {2} go run ./{3}/cmd/main.go' -f $SvcName, $ScriptDir, $EnvSetup, $SvcPath
    Start-Process powershell -ArgumentList '-NoExit', '-Command', $CommandStr
}

Write-Host ''
Write-Host '==================================================' -ForegroundColor Magenta
Write-Host '  OK: All services launched! Check their windows. ' -ForegroundColor Green
Write-Host '==================================================' -ForegroundColor Magenta
