#!/usr/bin/env pwsh
# PowerShell scripts para desarrollo

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Write-ColorText {
    param($Text, $Color = "White")
    Write-Host $Text -ForegroundColor $Color
}

function Show-Help {
    Write-ColorText "`nPlaylist Migrator - Scripts de desarrollo PowerShell`n" "Blue"
    Write-ColorText "Uso: .\scripts\dev.ps1 [comando]`n"
    Write-ColorText "Comandos disponibles:" "Blue"
    Write-ColorText "  setup       - Configuración inicial" "Green"
    Write-ColorText "  dev         - Ejecutar en modo desarrollo" "Green"
    Write-ColorText "  build       - Compilar aplicación" "Green"
    Write-ColorText "  test        - Ejecutar tests" "Green"
    Write-ColorText "  docker-up   - Levantar servicios Docker" "Green"
    Write-ColorText "  docker-down - Detener servicios Docker" "Green"
    Write-ColorText "  migrate     - Ejecutar migraciones" "Green"
    Write-ColorText "  clean       - Limpiar archivos generados" "Green"
    Write-ColorText "  check       - Verificar herramientas instaladas`n" "Green"
}

function Invoke-Setup {
    Write-ColorText "[Setup] Configurando proyecto..." "Yellow"
    
    if (-not (Test-Path ".env")) {
        Write-ColorText "[Setup] Copiando .env.example a .env" "Blue"
        Copy-Item ".env.example" ".env"
    }
    
    Write-ColorText "[Setup] Instalando dependencias Go..." "Blue"
    go mod tidy
    go mod download
    
    Write-ColorText "[Setup] Instalando herramientas..." "Blue"
    go install github.com/air-verse/air@latest
    
    Write-ColorText "[Setup] ✅ Configuración completada" "Green"
    Write-ColorText "[Setup] 💡 Edita .env con tus credenciales y ejecuta: .\scripts\dev.ps1 docker-up" "Blue"
}

function Invoke-Dev {
    Write-ColorText "[Dev] Verificando Air..." "Yellow"
    
    if (Get-Command air -ErrorAction SilentlyContinue) {
        Write-ColorText "[Dev] Ejecutando con Air live reload..." "Blue"
        air --build.cmd "go build -o ./tmp/main.exe cmd/server/main.go" --build.bin "./tmp/main.exe"
    } else {
        Write-ColorText "[Dev] Air no encontrado. Ejecutando sin live reload..." "Yellow"
        go run cmd/server/main.go
    }
}

function Invoke-Build {
    Write-ColorText "[Build] Compilando aplicación..." "Yellow"
    
    go build -o dist/main.exe cmd/server/main.go
    if ($LASTEXITCODE -eq 0) {
        Write-ColorText "[Build] ✅ Compilación exitosa" "Green"
    } else {
        Write-ColorText "[Build] ❌ Error en compilación" "Red"
        exit 1
    }
}

function Invoke-Test {
    Write-ColorText "[Test] Ejecutando tests..." "Yellow"
    
    go test -v ./...
    if ($LASTEXITCODE -eq 0) {
        Write-ColorText "[Test] ✅ Tests completados" "Green"
    } else {
        Write-ColorText "[Test] ❌ Tests fallaron" "Red"
        exit 1
    }
}

function Invoke-DockerUp {
    Write-ColorText "[Docker] Verificando Docker..." "Yellow"
    
    if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-ColorText "[Docker] ❌ Docker no está instalado" "Red"
        exit 1
    }
}
function Invoke-DockerUp {
    Write-ColorText "[Docker] Verificando Docker..." "Yellow"
    
    if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-ColorText "[Docker] ❌ Docker no está instalado" "Red"
        exit 1
    }
    
    Write-ColorText "[Docker] Levantando servicios PostgreSQL y Redis..." "Blue"
    docker-compose up -d postgres redis
    
    if ($LASTEXITCODE -eq 0) {
        Write-ColorText "[Docker] ✅ Servicios levantados" "Green"
        Write-ColorText "[Docker] 💡 Ejecuta: .\scripts\makefile.ps1 migrate" "Blue"
    } else {
        Write-ColorText "[Docker] ❌ Error levantando servicios" "Red"
        exit 1
    }
}

function Invoke-DockerDown {
    Write-ColorText "[Docker] Deteniendo servicios..." "Yellow"
    docker-compose down
    Write-ColorText "[Docker] ✅ Servicios detenidos" "Green"
}

function Invoke-Migrate {
    Write-ColorText "[Migrate] Ejecutando migraciones..." "Yellow"
    
    go run cmd/migrate/main.go up
    if ($LASTEXITCODE -eq 0) {
        Write-ColorText "[Migrate] ✅ Migraciones completadas" "Green"
    } else {
        Write-ColorText "[Migrate] ❌ Error en migraciones" "Red"
        exit 1
    }
}

function Invoke-Clean {
    Write-ColorText "[Clean] Limpiando archivos..." "Yellow"

    $filesToRemove = @("dist/main.exe", "coverage.out", "coverage.html")
    foreach ($file in $filesToRemove) {
        if (Test-Path $file) {
            Remove-Item $file -Force
            Write-ColorText "[Clean] Eliminado: $file" "Blue"
        }
    }
    
    if (Test-Path "tmp") {
        Remove-Item "tmp" -Recurse -Force
        Write-ColorText "[Clean] Eliminado directorio: tmp" "Blue"
    }
    
    Write-ColorText "[Clean] ✅ Limpieza completada" "Green"
}

function Invoke-Check {
    Write-ColorText "🔍 Verificando herramientas instaladas:" "Blue"
    
    Write-ColorText "Go:" "Blue"
    if (Get-Command go -ErrorAction SilentlyContinue) {
        $goVersion = go version
        Write-ColorText "  ✅ $goVersion" "Green"
    } else {
        Write-ColorText "  ❌ Go no encontrado" "Red"
    }
    
    Write-ColorText "Docker:" "Blue"
    if (Get-Command docker -ErrorAction SilentlyContinue) {
        $dockerVersion = docker --version
        Write-ColorText "  ✅ $dockerVersion" "Green"
    } else {
        Write-ColorText "  ❌ Docker no encontrado" "Red"
    }
    
    Write-ColorText "Docker Compose:" "Blue"
    if (Get-Command docker-compose -ErrorAction SilentlyContinue) {
        $composeVersion = docker-compose --version
        Write-ColorText "  ✅ $composeVersion" "Green"
    } else {
        Write-ColorText "  ❌ Docker Compose no encontrado" "Red"
    }
    
    Write-ColorText "Air (live reload):" "Blue"
    if (Get-Command air -ErrorAction SilentlyContinue) {
        Write-ColorText "  ✅ Air encontrado" "Green"
    } else {
        Write-ColorText "  ❌ Air no encontrado (ejecuta: go install github.com/air-verse/air@latest)" "Yellow"
    }
    
    Write-ColorText "Golangci-lint:" "Blue"
    if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
        Write-ColorText "  ✅ Golangci-lint encontrado" "Green"
    } else {
        Write-ColorText "  ❌ Golangci-lint no encontrado" "Yellow"
        Write-ColorText "    Descarga desde: https://golangci-lint.run/docs/welcome/install/" "Blue"
    }
}

# Ejecutar comando
switch ($Command) {
    "help" { Show-Help }
    "setup" { Invoke-Setup }
    "dev" { Invoke-Dev }
    "build" { Invoke-Build }
    "test" { Invoke-Test }
    "docker-up" { Invoke-DockerUp }
    "docker-down" { Invoke-DockerDown }
    "migrate" { Invoke-Migrate }
    "clean" { Invoke-Clean }
    "check" { Invoke-Check }
    default { 
        Write-ColorText "❌ Comando no reconocido: $Command" "Red"
        Show-Help 
        exit 1
    }
}