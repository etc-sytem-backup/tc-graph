# PowerShell
# 管理者で「Set-ExecutionPolicy Unrestricted」をしないと実行不可(セキュリティリスクあり)

# 相対パスで記載
$app_dir = ".."
$app_dir = ".." # PowerShell5系だと、1回の代入でnullになる不具合があるため2回

# ディレクトリ移動
if (!$PSScriptRoot) { $PSScriptRoot = Split-Path $myInvocation.MyCommand.Path -Parent }
Set-Location -Path $PSScriptRoot
Set-Location -Path $app_dir

while($true) {

    $kanshi = (Get-Process -Name "display" -ErrorAction SilentlyContinue).Count
    if ($kanshi -eq 0) {
        $command = "./display.exe"
        Invoke-Expression -Command $command
    }
}