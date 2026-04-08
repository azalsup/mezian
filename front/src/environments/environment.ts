export const environment = {
  production: false,
  apiBaseUrl: '/api/v1',

  /**
   * SDK log level. Controls what the ApiClient logs to the browser console.
   * 'debug' → log every request/response.
   * 'warn'  → log only warnings and errors (default for production).
   * 'none'  → completely silent.
   */
  logLevel: 'warn' as 'debug' | 'info' | 'warn' | 'error' | 'none',

  auth: {
    /** Set to false to hide OTP option and force password-only auth */
    otpEnabled: true,
    /** 'phone' | 'email' — which identifier is shown first on login/register */
    primaryIdentifier: 'phone' as 'phone' | 'email',
    /** Phone number is mandatory on registration */
    phoneRequired: true,
    /** Email is mandatory on registration */
    emailRequired: false,
    /** Pre-filled country on the address form (ISO 3166-1 alpha-2) */
    defaultCountry: 'MA',
  },
};
