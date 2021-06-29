$ErrorActionPreference = 'Stop'

$Version = GitVersion.exe | ConvertFrom-Json
Write-Output ('Version: {0}' -f $Version.SemVer)
goversioninfo -64 -product-version="$($Version.SemVer)" -ver-major="$($Version.Major)" -ver-minor="$($Version.Minor)" -ver-patch="$($Version.Patch)" -product-ver-major="$($Version.Major)" -product-ver-minor="$($Version.Minor)" -product-ver-patch="$($Version.Patch)"
Write-Output 'Updated program build info'
