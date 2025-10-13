# Корень проекта — папка, где лежит этот скрипт
$root = Split-Path -Parent $MyInvocation.MyCommand.Path

# Ищем go.exe в PATH и берём полный путь
try {
    $go = (Get-Command go -ErrorAction Stop).Source
} catch {
    Write-Error "Go не найден в PATH. Установите Go и/или добавьте его в PATH."
    exit 1
}

# Если у сервисов нет go.mod, замените 'go run .' на 'go run main.go' ниже
$args = @(
  'new-tab','--title','Orchestrator_service','-d',"$root\orchestrator_service",
  'powershell','-NoExit','-Command', ("& '{0}' run ." -f $go),

  ';','split-pane','-V','--title','Document_service','-d',"$root\document_service",
  'powershell','-NoExit','-Command', ("& '{0}' run ." -f $go),

  ';','focus-pane','--target','0',
  ';','split-pane','-H','--title','BD_service','-d',"$root\bd_service",
  'powershell','-NoExit','-Command', ("& '{0}' run ." -f $go),

  ';','focus-pane','--target','1',
  ';','split-pane','-H','--title','Garant_service','-d',"$root\garant_service",
  'powershell','-NoExit','-Command', ("& '{0}' run ." -f $go)
)


wt --window myProject @args
