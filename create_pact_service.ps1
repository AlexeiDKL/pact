# ...existing code...
function Build-GoProject {
    param(
        [Parameter(Mandatory = $true)]
        [string]$MainGoPath,      # Путь к main.go

        [Parameter(Mandatory = $true)]
        [string]$OutputName,      # Имя итогового бинарника

        [Parameter(Mandatory = $true)]
        [string]$OutputDir        # Папка, куда класть результат
    )

    $Goos = $env:GOOS
    $Goarch = $env:GOARCH

    Write-Host "Сборка для $Goos/$Goarch"

    # Проверяем, что main.go существует
    if (-not (Test-Path $MainGoPath)) {
        Write-Error "Файл $MainGoPath не найден."
        return
    }

    # Определяем рабочую директорию (где лежит main.go)
    $projectDir = Split-Path -Parent $MainGoPath

    # Создаём выходную папку, если её нет
    if (-not (Test-Path $OutputDir)) {
        New-Item -ItemType Directory -Path $OutputDir | Out-Null
    }

    # Путь к итоговому бинарнику
    $outputBinary = Join-Path $OutputDir $OutputName

    Write-Host "Сборка проекта $MainGoPath → $outputBinary"

    # Выполняем сборку в каталоге проекта (чтобы корректно срабатывать модули)
    Push-Location $projectDir
    try {
        #явно вызываем go с текущими env
        & go build -o $outputBinary .
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Ошибка при сборке Go проекта (exit $LASTEXITCODE)."
            Pop-Location
            return
        }
    } finally {
        Pop-Location
    }

    # Проверки результата: существование, размер и ELF magic (чтобы не получилось текстовое "Hello World")
    if (-not (Test-Path $outputBinary)) {
        Write-Error "Файл результата не найден: $outputBinary"
        return
    }

    $fi = Get-Item $outputBinary
    Write-Host "Результат: $($fi.Length) байт"

    if ($fi.Length -lt 1024) {
        Write-Warning "Размер бинарника <$([int]1024) байт — подозрительно маленький ($($fi.Length) байт)."
    }

    try {
        $fs = [System.IO.File]::Open($outputBinary, [System.IO.FileMode]::Open, [System.IO.FileAccess]::Read)
        $hdr = New-Object byte[] 4
        $fs.Read($hdr,0,4) | Out-Null
        $fs.Close()
        if (!($hdr[0] -eq 0x7F -and $hdr[1] -eq 0x45 -and $hdr[2] -eq 0x4C -and $hdr[3] -eq 0x46)) {
            Write-Error "Сгенерированный файл не является ELF‑бинарником (magic mismatch). Удаляю $outputBinary"
            Remove-Item $outputBinary -Force -ErrorAction SilentlyContinue
            return
        }
    } catch {
        Write-Error "Не удалось проверить magic файла: $_"
        return
    }

    # Копируем config/config.yaml, если он существует
    $configPath = Join-Path $projectDir "config\config.yaml"
    if (Test-Path $configPath) {
        $destConfigDir = Join-Path $OutputDir "config"
        if (-not (Test-Path $destConfigDir)) {
            New-Item -ItemType Directory -Path $destConfigDir | Out-Null
        }
        Copy-Item $configPath -Destination $destConfigDir -Force
        Write-Host "Скопирован config.yaml → $destConfigDir"
    }
    else {
        Write-Warning "Файл config/config.yaml не найден рядом с main.go"
    }

    Write-Host "✅ Сборка завершена успешно!"
}

# Устанавливаем новые
$env:GOOS  = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

$projects = @(
    @{
        MainGoPath = "C:\project\golang\pact\bd_service\main.go"
        OutputName = "bd_service"
        OutputDir  = "C:\project\golang\pact_build_linux\bd_service"
    },
     @{
        MainGoPath = "C:\project\golang\pact\document_service\main.go"
        OutputName = "document_service"
        OutputDir  = "C:\project\golang\pact_build_linux\document_service"
    },
     @{
        MainGoPath = "C:\project\golang\pact\garant_service\main.go"
        OutputName = "garant_service"
        OutputDir  = "C:\project\golang\pact_build_linux\garant_service"
    },
     @{
        MainGoPath = "C:\project\golang\pact\orchestrator_service\main.go"
        OutputName = "orchestrator_service"
        OutputDir  = "C:\project\golang\pact_build_linux\orchestrator_service"
    },
     @{
        MainGoPath = "C:\project\golang\qu\tester_service\main.go"
        OutputName = "tester_service"
        OutputDir  = "C:\project\golang\pact_build_linux\tester_service"
    }
)

# Сборка
foreach ($proj in $projects) {
    Build-GoProject -MainGoPath $proj.MainGoPath `
                    -OutputName $proj.OutputName `
                    -OutputDir  $proj.OutputDir
}





# # Возвращаем обратно
# $env:GOOS  = $oldGOOS
# $env:GOARCH = $oldGOARCH
# $env:CGO_ENABLED = $oldCGO
