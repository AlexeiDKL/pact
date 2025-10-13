# restart-pact-services.ps1

# 1. Завершаем все процессы go
Get-Process -Name "go" -ErrorAction SilentlyContinue | Stop-Process -Force

# 2. Список сервисов (те же, что в start-pact-services.ps1)
$serviceTitles = @(
    "Orchestrator_service",
    "Document_service",
    "BD_service",
    "Garant_service"
)

# 3. Находим процессы WindowsTerminal.exe, где в командной строке есть хотя бы один из --title
$procs = Get-CimInstance Win32_Process -Filter "Name = 'WindowsTerminal.exe'" |
    Where-Object {
        foreach ($title in $serviceTitles) {
            if ($_.CommandLine -match "--title $title") { return $true }
        }
        return $false
    }

if ($procs) {
    $procs | ForEach-Object { Stop-Process -Id $_.ProcessId -Force }
    Write-Host "Close: $($serviceTitles -join ', ')"
} else {
    Write-Host "Services are not running."
}

# 4. Подождать секунду
Start-Sleep -Seconds 4

# 5. Запустить сервисы снова
.\start-pact-services.ps1
