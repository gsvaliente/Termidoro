# Installing the Man Page

To make the man page available system-wide, you can install it using these commands:

## Install System-wide

```bash
# Copy the man page to the standard location
sudo cp termidoro.1 /usr/local/share/man/man1/

# Update the man database
sudo mandb

# Test the man page
man termidoro
```

## Install in User Directory (Alternative)

If you prefer not to use sudo:

```bash
# Create user man directory if it doesn't exist
mkdir -p ~/.local/share/man/man1/

# Copy the man page
cp termidoro.1 ~/.local/share/man/man1/

# Add to your shell profile if not already present
echo 'export MANPATH="$HOME/.local/share/man:$MANPATH"' >> ~/.bashrc
# or for zsh:
echo 'export MANPATH="$HOME/.local/share/man:$MANPATH"' >> ~/.zshrc

# Reload shell or source the profile
source ~/.bashrc  # or source ~/.zshrc

# Test
man termidoro
```

## Verify Installation

After installation, you should be able to run:

```bash
man termidoro
```

This will display the complete manual page with all usage instructions and options.