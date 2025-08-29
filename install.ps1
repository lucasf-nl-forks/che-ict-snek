# Set up TLS security protocol
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

# Function to get latest release version
function Get-LatestReleaseVersion {
    param (
        [string]$RepositoryOwner,
        [string]$RepositoryName
    )

    $releasesUrl = "https://api.github.com/repos/$RepositoryOwner/$RepositoryName/releases/latest"
    try {
        $response = Invoke-WebRequest -Uri $releasesUrl -UseBasicParsing
        $releaseData = ConvertFrom-Json $response.Content
        return $releaseData.tag_name
    }
    catch {
        Write-Error "Failed to fetch latest release version: $_"
        exit 1
    }
}

# Function to download and extract zip
function DownloadAndExtract-Zip {
    param (
        [string]$RepositoryOwner,
        [string]$RepositoryName,
        [string]$Version,
        [string]$FileName
    )

    # Create temp directory
    $tempDir = Join-Path $env:TEMP ([Guid]::NewGuid().ToString())
    New-Item -ItemType Directory -Force -Path $tempDir

    try {
        # Download zip
        $downloadUrl = "https://github.com/$RepositoryOwner/$RepositoryName/releases/download/$Version/$FileName"
        $zipPath = Join-Path $tempDir $FileName

        Write-Host "Downloading $FileName..."
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath

        # Extract contents
        Write-Host "Extracting files..."
        Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force

        return $tempDir
    }
    catch {
        Write-Error "Failed to download or extract zip: $_"
        Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
        exit 1
    }
}

# Function to add directory to PATH
function Add-ToUserPath {
    param (
        [string]$Directory
    )

    $pathKey = "HKCU:\Environment"
    $pathValue = "Path"

    try {
        $currentPaths = (Get-ItemProperty -Path $pathKey -Name $pathValue).$pathValue -split ";"

        if ($currentPaths -contains $Directory) {
            Write-Host "Directory is already in PATH"
            return
        }

        $newPaths = @($currentPaths + $Directory)
        Set-ItemProperty -Path $pathKey -Name $pathValue -Type ExpandString -Value ($newPaths -join ";")

        Write-Host "Added $Directory to user PATH"
    }
    catch {
        Write-Error "Failed to modify PATH: $_"
        exit 1
    }
}

# Main script
try {
    # Repository configuration
    $repoOwner = "che-ict"
    $repoName = "snek"
    $fileName = "snek-{tag_name}-windows-{architecture}.zip"
    $exeName = "snek.exe"

    # Determine architecture
    $architecture = if ($ENV:PROCESSOR_ARCHITECTURE -eq "ARM64") {
        "arm64"
    } elseif ($ENV:PROCESSOR_ARCHITECTURE -eq "AMD64") {
        "amd64"
    } elseif ($ENV:PROCESSOR_ARCHITECTURE -eq "386") {
        "386"
    } else {
        throw "Unsupported architecture: $($ENV:PROCESSOR_ARCHITECTURE)"
    }

    # Construct filename based on architecture
    $fileName = $fileName -replace "{architecture}", $architecture


    # Get latest version
    Write-Host "Getting latest release version..."
    $latestVersion = Get-LatestReleaseVersion -RepositoryOwner $repoOwner -RepositoryName $repoName

    $fileName = $fileName -replace "{tag_name}", $latestVersion


    # Download and extract
    Write-Host "Downloading and extracting release..."
    $extractedPath = DownloadAndExtract-Zip -RepositoryOwner $repoOwner -RepositoryName $repoName -Version $latestVersion -FileName $fileName

    # Find and copy exe
    $exePath = Get-ChildItem -Path $extractedPath -Filter $exeName -Recurse
    if (-not $exePath) {
        throw "Could not find executable $exeName in extracted files"
    }

    # Create bin directory if it doesn't exist
    $binPath = Join-Path $env:LOCALAPPDATA "bin"
    New-Item -ItemType Directory -Force -Path $binPath

    # Copy executable
    Write-Host "Copying executable to bin directory..."
    Copy-Item -Path $exePath.FullName -Destination $binPath -Force

    # Add to PATH if necessary
    Write-Host "Adding bin directory to PATH..."
    Add-ToUserPath -Directory $binPath

    Write-Host "Installation complete!"
}
catch {
    Write-Error "An error occurred: $_"
    exit 1
}