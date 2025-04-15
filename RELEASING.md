# Release Process

This repository follows a standardized release process using automated tools. Here's how releases are managed:

1. **Development Setup**:
   - Install dependencies: `npm install`
   - This sets up:
     - Husky for Git hooks
     - Commitlint for commit message validation
     - Release-please for automated releases

2. **Commit Guidelines**:
   Commits must follow conventional commit format with these types:
   - `feat`: New features
   - `fix`: Bug fixes
   - `docs`: Documentation changes
   - `style`: Code style changes
   - `refactor`: Code refactoring
   - `test`: Test changes
   - `chore`: Build process changes
   - `ci`: CI configuration changes
   - `perf`: Performance improvements
   - `revert`: Reverting changes

3. **Branch Protection**:
   The repository enforces:
   - Branch naming conventions:
     - `feat/*` - New features
     - `fix/*` - Bug fixes
     - `docs/*` - Documentation changes
     - `style/*` - Code style changes
     - `refactor/*` - Code refactoring
     - `test/*` - Test changes
     - `chore/*` - Build process changes
     - `ci/*` - CI configuration changes
     - `perf/*` - Performance improvements
     - `revert/*` - Reverting changes
   - Conventional commits
   - No direct commits to release branch
   - Required pull requests for all changes

4. **Creating a Release**:
   You can create a release by using the Makefile:
   
   ```bash
   # Using the Makefile (recommended)
   make release
   ```
   
   This will:
   - Create a release PR
   - Generate a changelog
   - Update version numbers
   - Apply branch protection rules

5. **Release Automation**:
   GitHub Actions handle:
   - Automatic PR creation and updates
   - Changelog generation
   - Version management
   - Release creation
   - GPG signing of releases

For more details about the release process, see:
- `.github/workflows/release.yml` - Release automation
- `.github/workflows/release-please.yml` - Release PR management
- `.github/branch-protection/` - Branch protection rules
- `contrib/release.sh` - Release script
- `package.json` - Development dependencies and scripts
- `Makefile` - Build and release commands 