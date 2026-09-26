@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul

set user_text=Atualizacao
if not "%~1"=="" set user_text=%~1

git add .
git commit -m "%user_text%"
git push origin main
