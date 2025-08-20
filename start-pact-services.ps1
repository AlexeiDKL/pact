<#
.SYNOPSIS
  Открывает одно окно Windows Terminal, разбитое на 4 панели, и в каждой панели:
    1. переходит в C:\project\golang\pact
    2. очищает экран
    3. запускает go run для своего сервиса
#>

# 1. Путь к корню проекта
$basePath = 'C:\project\golang\pact'

# 2. Формируем список аргументов для wt.exe
$wtArgs = @(
	newTab powershell;
	splitPane -H --title "orchestrator_service"  powershell;
	splitPane -V powershell ; 
	splitPane -H powershell ;
)

# 3. Вызываем Windows Terminal напрямую, передавая массив аргументов
#    Оператор & гарантирует запуск внешнего исполняемого файла wt.exe
& wt @wtArgs
