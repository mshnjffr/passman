# 🍺 Homebrew Tap Setup Instructions

## Step 1: Create the Tap Repository

1. Go to GitHub: https://github.com/new
2. Repository name: `homebrew-passman` (the `homebrew-` prefix is required)
3. Description: "Homebrew formula for Passman - A beautiful terminal password manager"
4. Make it **Public**
5. Initialize with README ✅
6. Click "Create repository"

## Step 2: Add the Formula File

1. In your new `homebrew-passman` repository, click "Create new file"
2. Name the file: `passman.rb`
3. Copy and paste this content:

```ruby
class Passman < Formula
  desc "Beautiful terminal password manager built with Go and Bubble Tea"
  homepage "https://github.com/mshnjffr/passman"
  url "https://github.com/mshnjffr/passman/archive/refs/tags/v1.0.2.tar.gz"
  sha256 "4ba1fc2bb5e8a2c03115418778c566d7222d0ff0f9b2a2bb61057f308b6e5698"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", bin/"passman"
  end

  test do
    assert_match "passman", shell_output("#{bin}/passman --version")
  end
end
```

4. Commit the file with message: "Add passman formula"

## Step 3: Test the Installation

Once the tap repository is created, users can install with:

```bash
# Add the tap
brew tap mshnjffr/passman

# Install passman
brew install passman

# Test it works
passman --version
```

## Step 4: Update for New Versions

When you release new versions of passman:

1. Get the new SHA256:
   ```bash
   curl -sL https://github.com/mshnjffr/passman/archive/refs/tags/vX.X.X.tar.gz | shasum -a 256
   ```

2. Update the formula:
   - Change the `url` to point to the new tag
   - Update the `sha256` with the new hash
   - Commit the changes

## Repository Structure

Your `homebrew-passman` repository should look like:
```
homebrew-passman/
├── README.md
└── passman.rb
```

## Usage for End Users

After setup, users can install passman with:

```bash
# Method 1: Add tap first, then install
brew tap mshnjffr/passman
brew install passman

# Method 2: Install directly (auto-adds tap)
brew install mshnjffr/passman/passman

# Updates
brew upgrade passman
```

## Benefits of Homebrew Distribution

✅ **No Go required**: Users don't need Go installed
✅ **Automatic updates**: `brew upgrade` keeps it current  
✅ **Dependency management**: Homebrew handles everything
✅ **Familiar workflow**: Standard `brew install` command
✅ **Cross-platform**: Works on macOS and Linux
✅ **Uninstall support**: `brew uninstall passman`
