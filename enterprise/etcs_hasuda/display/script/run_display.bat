@REM display.ps1と同じディレクトリに配置してください。
@echo off
where /q pwsh.exe
if %ERRORLEVEL% EQU 0 (
    pwsh.exe -ExecutionPolicy Bypass -File "%~dp0display.ps1"
) else (
    PowerShell.exe -ExecutionPolicy Bypass -File "%~dp0display.ps1"
)