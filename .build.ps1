$Version = GitVersion.exe | ConvertFrom-Json
goversioninfo -64 -product-version="$($Version.SemVer)" -ver-major="$($Version.Major)" -ver-minor="$($Version.Minor)" -ver-patch="$($Version.Patch)" -product-ver-major="$($Version.Major)" -product-ver-minor="$($Version.Minor)" -product-ver-patch="$($Version.Patch)"
go build -ldflags "-H=windowsgui -X main.version=$($Version.SemVer)"
