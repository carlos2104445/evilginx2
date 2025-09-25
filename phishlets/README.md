# Evilginx2 Phishlets Collection

This directory contains comprehensive phishlet configurations for popular sites, sourced from the well-maintained [An0nUD4Y/Evilginx2-Phishlets](https://github.com/An0nUD4Y/Evilginx2-Phishlets) repository.

## Available Phishlets

### Google (`google.yaml`)
- **Author**: @hash3liZer
- **Features**: 
  - JavaScript injection for botguard bypass
  - Extensive domain filtering (googleapis.com, gstatic.com, youtube.com)
  - Advanced authentication flow handling
  - Force_post configuration for multi-step authentication
- **Domains**: accounts.google.com, apis.google.com, ssl.gstatic.com, content.googleapis.com

### Microsoft (`microsoft.yaml`)
- **Author**: @white_fi
- **Features**:
  - Force_post configuration for "Keep me signed in" functionality
  - Comprehensive live.com and microsoft.com domain coverage
  - Multi-domain authentication support
- **Domains**: login.microsoftonline.com, login.live.com, account.microsoft.com

### LinkedIn (`linkedin.yaml`)
- **Author**: @white_fi
- **Features**:
  - Arkose Labs CAPTCHA integration
  - Extensive licdn.com domain filtering
  - JavaScript injection for form automation
  - Advanced session token capture
- **Domains**: linkedin.com, licdn.com

### Facebook (`facebook.yaml`)
- **Author**: @white_fi
- **Features**:
  - Standard Facebook authentication flow
  - Session token capture
  - Mobile and desktop compatibility
- **Domains**: facebook.com, fbcdn.net

### Instagram (`instagram.yaml`)
- **Author**: @white_fi
- **Features**:
  - Facebook OAuth integration support
  - Mobile-optimized authentication flow
  - Session token capture
- **Domains**: instagram.com, cdninstagram.com

### Twitter/X (`twitter.yaml`)
- **Author**: @white_fi (updated for X.com rebranding)
- **Features**:
  - Updated for X.com rebranding with backward compatibility
  - Supports both twitter.com and x.com domains
  - Cross-domain filtering for seamless transitions
  - Session token capture for both platforms
- **Domains**: twitter.com, x.com, twimg.com

## Usage

1. **Enable a phishlet**:
   ```
   phishlets enable <phishlet_name>
   ```

2. **Set hostname**:
   ```
   phishlets hostname <phishlet_name> <your_domain>
   ```

3. **Create lure**:
   ```
   lures create <phishlet_name>
   ```

## Advanced Features

### JavaScript Injection
Many phishlets include JavaScript injection capabilities to bypass modern anti-bot measures:
- **Google**: Botguard bypass scripts
- **LinkedIn**: Form automation and CAPTCHA handling
- **Microsoft**: Multi-factor authentication support

### Domain Filtering
Comprehensive sub_filters ensure seamless proxying of:
- Static resources (CSS, JS, images)
- API endpoints
- Authentication redirects
- Cross-domain references

### Session Token Capture
All phishlets are configured to capture relevant session tokens:
- **Google**: OAuth tokens, session cookies
- **Microsoft**: Authentication tokens, refresh tokens
- **LinkedIn**: Session cookies, CSRF tokens
- **Facebook/Instagram**: Session cookies, authentication tokens
- **Twitter/X**: Session cookies, authentication tokens

## Compatibility

- **Minimum Version**: 2.3.0+
- **Certificate Management**: Compatible with updated certificate management system
- **Security Features**: Works with modernized proxy system and security improvements

## Security Considerations

These phishlets are designed for authorized security testing and educational purposes only. Always ensure you have proper authorization before conducting any security assessments.

## Contributing

When updating phishlets:
1. Test YAML syntax with `./build/evilginx -developer -debug -c ./test-config -p ./phishlets`
2. Verify all phishlets show as "visible" without parsing errors
3. Test certificate generation in developer mode
4. Update documentation for any new features or bypass techniques
