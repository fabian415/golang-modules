@echo off
REM Sample Build Script for BSP Image Building
REM This script simulates the build process with progress updates
REM Outputs one line per second

echo Starting BSP Image Build Process...

REM Check if we should simulate an error
if "%1"=="error" (
    echo [INFO] Starting build with error simulation...
    >nul 2>&1 ping 127.0.0.1 -n 2
    echo BSP Build Progress: 30%%
    >nul 2>&1 ping 127.0.0.1 -n 2
    echo BSP Build Progress: 60%%
    >nul 2>&1 ping 127.0.0.1 -n 2
    echo [ERROR] Build failed at step 6: Compilation error in kernel module
    exit /b 1
)

REM Simulate progress - output one line per second
REM Each echo is followed by a 1-second delay using ping
echo [INFO] Initializing build environment...
>nul 2>&1 ping 127.0.0.1 -n 2
echo BSP Build Progress: 0%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [INFO] Downloading dependencies...
>nul 2>&1 ping 127.0.0.1 -n 2
echo BSP Build Progress: 10%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [INFO] Configuring build parameters...
>nul 2>&1 ping 127.0.0.1 -n 2
echo BSP Build Progress: 20%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [INFO] Compiling kernel modules...
>nul 2>&1 ping 127.0.0.1 -n 2
echo BSP Build Progress: 30%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [INFO] Building device tree...
>nul 2>&1 ping 127.0.0.1 -n 2
echo BSP Build Progress: 40%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [INFO] Compiling user space applications...
>nul 2>&1 ping 127.0.0.1 -n 2
echo BSP Build Progress: 50%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [INFO] Creating root filesystem...
>nul 2>&1 ping 127.0.0.1 -n 2
echo BSP Build Progress: 60%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [INFO] Packaging BSP image...
>nul 2>&1 ping 127.0.0.1 -n 2
echo BSP Build Progress: 70%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [INFO] Generating checksums...
>nul 2>&1 ping 127.0.0.1 -n 2
echo BSP Build Progress: 80%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [INFO] Finalizing build...
>nul 2>&1 ping 127.0.0.1 -n 2
echo BSP Build Progress: 90%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [SUCCESS] BSP Build Progress: 100%%
>nul 2>&1 ping 127.0.0.1 -n 2
echo [SUCCESS] Script build completed successfully!
>nul 2>&1 ping 127.0.0.1 -n 2
echo Final image generated at: /tmp/sample_bsp_image.img

:end
echo Build script completed.
