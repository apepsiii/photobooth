@echo off
title Photobooth Sekolah - Kiosk Launcher
cd /d "%~dp0"

echo ===================================================
echo     MEMULAI APLIKASI PHOTOBOOTH SEKOLAH
echo ===================================================
echo.

:: 1. Cek apakah binary photobooth.exe sudah ada
if not exist "photobooth.exe" (
    echo Mengompilasi aplikasi Go...
    go build -o photobooth.exe .
)

:: 2. Jalankan Photobooth Server di background
echo Menjalankan Photobooth Server di port 8090...
start "" /b "photobooth.exe"

:: Tunggu 2 detik agar server siap
timeout /t 2 /nobreak >nul

:: 3. Cari browser (Google Chrome atau Microsoft Edge)
set BROWSER_EXE=""

if exist "C:\Program Files\Google\Chrome\Application\chrome.exe" (
    set BROWSER_EXE="C:\Program Files\Google\Chrome\Application\chrome.exe"
) else if exist "C:\Program Files (x86)\Google\Chrome\Application\chrome.exe" (
    set BROWSER_EXE="C:\Program Files (x86)\Google\Chrome\Application\chrome.exe"
) else if exist "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" (
    set BROWSER_EXE="C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe"
) else if exist "C:\Program Files\Microsoft\Edge\Application\msedge.exe" (
    set BROWSER_EXE="C:\Program Files\Microsoft\Edge\Application\msedge.exe"
)

:: 4. Buka di Fullscreen Kiosk Mode
if not %BROWSER_EXE%=="" (
    echo Membuka antarmuka Kiosk Fullscreen...
    start "" %BROWSER_EXE% --kiosk http://localhost:8090/kiosk --disable-pinch --overscroll-history-navigation=0 --no-first-run
) else (
    echo Browser Chrome/Edge tidak ditemukan di folder standar.
    echo Membuka browser default...
    start http://localhost:8090/kiosk
)

echo.
echo ===================================================
echo  Aplikasi Photobooth berjalan!
echo  - Layar Kiosk: http://localhost:8090/kiosk
echo  - Panel Admin: http://localhost:8090/admin/sessions
echo ===================================================

