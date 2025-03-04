@echo off

echo Start building Koolo
echo Cleaning up previous artifacts...
if exist build rmdir /s /q build > NUL || goto :error

echo Building Koolo binary...
set PATH=%PATH%;C:\Program Files (x86)\CMake\bin;C:\mingw-w64\x86_64-8.1.0-posix-seh-rt_v6-rev0\mingw64\bin;C:\opencv\build\install\x64\mingw\bin
go build -trimpath -tags static --ldflags -extldflags="-static" -o build/koolo.exe ./cmd/koolo/main.go > NUL || goto :error

echo Copying assets...
mkdir build\config > NUL || goto :error
copy config\config.yaml.dist build\config\config.yaml  > NUL || goto :error
xcopy /q /E /I /y config\pickit build\config\pickit  > NUL || goto :error
xcopy /q /E /I /y config\pickit_leveling build\config\pickit_leveling  > NUL || goto :error
xcopy /q /E /I /y assets build\assets  > NUL || goto :error
xcopy /q /y koolo-map.exe build > NUL || goto :error
xcopy /q /y README.md build > NUL || goto :error

echo Done! Artifacts are in build directory.

:error
if %errorlevel% neq 0 (
    echo Error occurred #%errorlevel%.
    exit /b %errorlevel%
)
