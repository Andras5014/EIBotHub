$ErrorActionPreference = 'Stop'

$AppDir = Resolve-Path (Join-Path $PSScriptRoot '..')

Set-Location $AppDir

& (Join-Path $AppDir 'eibothub.exe')
