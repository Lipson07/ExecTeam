param(
    [ValidateSet("init","new","up","down","status","version")]
    [string]$Command,
    [string]$Name
)

$DB_HOST = "localhost"
$DB_PORT = "5432"
$DB_USER = "postgres"
$DB_NAME = "exec_team"
$env:PGPASSWORD = "1234"
$MIGRATION_DIR = "$PSScriptRoot/migration"

switch ($Command) {
    "init" {
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMP DEFAULT NOW());"
        Write-Host "schema_migrations created" -ForegroundColor Green
    }
    
    "new" {
        if (-not $Name) { Write-Host "Usage: .\migrate.ps1 new -Name create_users" -ForegroundColor Yellow; return }
        if (-not (Test-Path $MIGRATION_DIR)) { New-Item -ItemType Directory -Path $MIGRATION_DIR | Out-Null }
        $ts = Get-Date -Format "yyyyMMddHHmmss"
        $up = "$MIGRATION_DIR/$ts`_$Name.up.sql"
        $down = "$MIGRATION_DIR/$ts`_$Name.down.sql"
        
        # Создать БЕЗ BOM
        $utf8 = New-Object System.Text.UTF8Encoding $false
        [System.IO.File]::WriteAllText($up, "", $utf8)
        [System.IO.File]::WriteAllText($down, "", $utf8)
        
        Write-Host "Created:" -ForegroundColor Green
        Write-Host "  $up"
        Write-Host "  $down"
    }
    
    "up" {
        Get-ChildItem "$MIGRATION_DIR/*.up.sql" | Sort-Object Name | ForEach-Object {
            $ver = ($_.BaseName -split '_')[0]
            $check = psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT 1 FROM schema_migrations WHERE version = '$ver';" 2>$null
            if ($check -notmatch '1') {
                Write-Host "APPLY: $($_.Name)" -ForegroundColor Green
                # Убрать BOM из файла перед применением
                $content = [System.IO.File]::ReadAllText($_.FullName)
                $utf8 = New-Object System.Text.UTF8Encoding $false
                $tmpFile = "$env:TEMP/migration_tmp.sql"
                [System.IO.File]::WriteAllText($tmpFile, $content, $utf8)
                psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $tmpFile
                
                if ($LASTEXITCODE -eq 0) {
                    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "INSERT INTO schema_migrations (version) VALUES ('$ver');" | Out-Null
                    Write-Host "  OK" -ForegroundColor Green
                } else {
                    Write-Host "  ERROR!" -ForegroundColor Red
                    break
                }
            } else {
                Write-Host "SKIP: $($_.Name)" -ForegroundColor DarkGray
            }
        }
        Write-Host "Done." -ForegroundColor Green
    }
    
    "down" {
        $ver = psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;" 2>$null
        $ver = ($ver -join '').Trim()
        if ($ver) {
            $df = Get-ChildItem "$MIGRATION_DIR/$ver*.down.sql" | Select-Object -First 1
            if ($df) {
                Write-Host "ROLLBACK: $($df.Name)" -ForegroundColor Yellow
                $content = [System.IO.File]::ReadAllText($df.FullName)
                $utf8 = New-Object System.Text.UTF8Encoding $false
                $tmpFile = "$env:TEMP/migration_tmp.sql"
                [System.IO.File]::WriteAllText($tmpFile, $content, $utf8)
                psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $tmpFile
                
                if ($LASTEXITCODE -eq 0) {
                    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "DELETE FROM schema_migrations WHERE version = '$ver';" | Out-Null
                    Write-Host "  OK" -ForegroundColor Green
                } else {
                    Write-Host "  ERROR!" -ForegroundColor Red
                }
            } else {
                Write-Host "DOWN file not found for: $ver" -ForegroundColor Red
            }
        } else {
            Write-Host "No migrations to rollback" -ForegroundColor Yellow
        }
    }
    
    "status" {
        Write-Host ""
        Write-Host "=== APPLIED ===" -ForegroundColor Cyan
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT version, applied_at FROM schema_migrations ORDER BY version;"
        Write-Host ""
        Write-Host "=== FILES ===" -ForegroundColor Cyan
        Get-ChildItem "$MIGRATION_DIR/*.sql" | Sort-Object Name | ForEach-Object { Write-Host "  $($_.Name)" }
        Write-Host ""
    }
    
    "version" {
        $v = psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;" 2>$null
        $v = ($v -join '').Trim()
        if ($v) { Write-Host "Version: $v" -ForegroundColor Cyan }
        else { Write-Host "No migrations" -ForegroundColor Yellow }
    }
}