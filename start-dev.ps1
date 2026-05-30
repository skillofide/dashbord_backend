# start-dev.ps1
# Automates the startup of the entire SkillofIDE backend.
# Run this script from the root directory: .\start-dev.ps1

# Set Error Action
$ErrorActionPreference = "Stop"

# Clear Screen and Print Header
Clear-Host
Write-Host '==================================================' -ForegroundColor Magenta
Write-Host '          SkillofIDE Backend Startup Script       ' -ForegroundColor Magenta
Write-Host '==================================================' -ForegroundColor Magenta

$ScriptDir = $PSScriptRoot
if (-not $ScriptDir) {
    $ScriptDir = $PWD
}
cd $ScriptDir

# 1. Verify Docker is running
Write-Host ''
Write-Host '[1/5] Checking if Docker is running...' -ForegroundColor Yellow
& docker info >$null 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host 'Error: Docker is not running.' -ForegroundColor Red
    Write-Host 'Please start Docker Desktop and try again!' -ForegroundColor Yellow
    exit 1
}
Write-Host 'OK: Docker is running.' -ForegroundColor Green

# 2. Check and Build Runner Images
Write-Host ''
Write-Host '[2/5] Checking required Docker runner images...' -ForegroundColor Yellow
$Runners = [ordered]@{
    'skillofide/runner-python:latest'     = 'services/execution-service/runners/python'
    'skillofide/runner-javascript:latest' = 'services/execution-service/runners/javascript'
    'skillofide/runner-java:latest'       = 'services/execution-service/runners/java'
    'skillofide/runner-cpp:latest'        = 'services/execution-service/runners/cpp'
}

foreach ($Img in $Runners.Keys) {
    $Check = docker images -q $Img
    if (-not $Check) {
        Write-Host ('Image ' + $Img + ' not found. Building...') -ForegroundColor Cyan
        $Path = Join-Path $ScriptDir $Runners[$Img]
        docker build -t $Img $Path
        if ($LASTEXITCODE -ne 0) {
            Write-Host ('Error: Failed to build ' + $Img) -ForegroundColor Red
            exit 1
        }
        Write-Host ('Built ' + $Img + ' successfully.') -ForegroundColor Green
    } else {
        Write-Host ('Image ' + $Img + ' already exists.') -ForegroundColor Green
    }
}

# 3. Spin up Infrastructure (PostgreSQL, Redis, NATS)
Write-Host ''
Write-Host '[3/5] Starting databases and messaging queues...' -ForegroundColor Yellow
docker compose up -d postgres redis nats
if ($LASTEXITCODE -ne 0) {
    Write-Host 'Error: Failed to start docker-compose services' -ForegroundColor Red
    exit 1
}
Write-Host 'OK: Databases and queues started.' -ForegroundColor Green

# 4. Run Migrations
Write-Host ''
Write-Host '[4/5] Running database migrations...' -ForegroundColor Yellow
docker compose run --rm db-migrate
if ($LASTEXITCODE -ne 0) {
    Write-Host 'Error: Migrations failed!' -ForegroundColor Red
    exit 1
}
Write-Host 'OK: Database migrations completed successfully.' -ForegroundColor Green

# 5. Launch Backend Services in separate windows
Write-Host ''
Write-Host '[5/5] Launching backend microservices...' -ForegroundColor Yellow

$Services = [ordered]@{
    'Problem Service'      = 'services/problem-service'
    'Execution Service'    = 'services/execution-service'
    'Submission Service'   = 'services/submission-service'
    'Progress Service'     = 'services/progress-service'
    'Notification Service' = 'services/notification-service'
    'API Gateway'          = 'services/api-gateway'
}

foreach ($SvcName in $Services.Keys) {
    $SvcPath = $Services[$SvcName]
    Write-Host ('Launching ' + $SvcName + ' in a new terminal window...') -ForegroundColor Cyan
    
    # Configure API Gateway environment variables
    $EnvSetup = ''
    if ($SvcName -eq 'API Gateway') {
        $EnvSetup = '$env:POSTGRES_DSN="postgres://skillofide:password@localhost:5432/skillofide?sslmode=disable";'
    }
    
    # Construct flat command using format operator to avoid nesting quote issues
    $CommandStr = '$Host.UI.RawUI.WindowTitle = "{0}"; cd "{1}"; {2} go run ./{3}/cmd/main.go' -f $SvcName, $ScriptDir, $EnvSetup, $SvcPath
    
    # Open new window, set title, run service
    Start-Process powershell -ArgumentList '-NoExit', '-Command', $CommandStr
}

Write-Host ''
Write-Host '==================================================' -ForegroundColor Magenta
Write-Host '  OK: All services launched! Check their windows. ' -ForegroundColor Green
Write-Host '==================================================' -ForegroundColor Magenta
